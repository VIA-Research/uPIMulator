package vm

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"uPIMulator/src/device/core"
	"uPIMulator/src/device/simulator/channel"
	"uPIMulator/src/device/simulator/dpu"
	"uPIMulator/src/encoding"
	"uPIMulator/src/host/abi"
	"uPIMulator/src/host/vm/arena"
	"uPIMulator/src/host/vm/base"
	"uPIMulator/src/host/vm/dram"
	"uPIMulator/src/host/vm/dram/bank"
	"uPIMulator/src/host/vm/frame"
	"uPIMulator/src/host/vm/stack"
	"uPIMulator/src/host/vm/symbol"
	"uPIMulator/src/host/vm/type_system"
	"uPIMulator/src/misc"
	"uPIMulator/src/program"
)

type VirtualMachine struct {
	bin_dirpath string

	num_channels          int
	num_ranks_per_channel int
	num_dpus_per_rank     int

	num_dpus int

	verbose int

	app  *program.App
	task *program.Task

	arena             *arena.Arena
	frame_chain       *frame.FrameChain
	registry          *type_system.Registry
	garbage_collector *arena.GarbageCollector

	cur_skeleton_name *string

	memory_controller *dram.MemoryController
	channels          []*channel.Channel

	prepare_xfer_buf map[*dpu.Dpu]int64
	push_xfer        map[*bank.TransferCommand]bool
}

func (this *VirtualMachine) Init(command_line_parser *misc.CommandLineParser) {
	this.bin_dirpath = command_line_parser.StringParameter("bin_dirpath")

	this.num_channels = int(command_line_parser.IntParameter("num_channels"))
	this.num_ranks_per_channel = int(command_line_parser.IntParameter("num_ranks_per_channel"))
	this.num_dpus_per_rank = int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = this.num_channels * this.num_ranks_per_channel * this.num_dpus_per_rank

	this.verbose = int(command_line_parser.IntParameter("verbose"))

	this.app = nil
	this.task = nil

	this.arena = new(arena.Arena)
	this.arena.Init()

	this.frame_chain = new(frame.FrameChain)
	this.frame_chain.Init()

	this.registry = new(type_system.Registry)
	this.registry.Init()

	this.garbage_collector = new(arena.GarbageCollector)
	this.garbage_collector.Init()
	this.garbage_collector.ConnectArena(this.arena)
	this.garbage_collector.ConnectFrameChain(this.frame_chain)
	this.garbage_collector.ConnectRegistry(this.registry)

	this.cur_skeleton_name = nil

	this.memory_controller = new(dram.MemoryController)
	this.memory_controller.Init(command_line_parser)

	this.channels = make([]*channel.Channel, 0)
	for i := 0; i < this.num_channels; i++ {
		channel_ := new(channel.Channel)
		channel_.Init(i, command_line_parser)

		this.channels = append(this.channels, channel_)
	}

	this.memory_controller.ConnectChannels(this.channels)

	this.prepare_xfer_buf = make(map[*dpu.Dpu]int64)
	this.push_xfer = make(map[*bank.TransferCommand]bool)
}

func (this *VirtualMachine) Fini() {
	for _, channel_ := range this.channels {
		channel_.Fini()
	}

	this.memory_controller.Fini()

	if len(this.prepare_xfer_buf) != 0 {
		err := errors.New("VM's prepare xfer buf is not empty")
		panic(err)
	}
}

func (this *VirtualMachine) Load(app *program.App, task *program.Task) {
	this.app = app
	this.task = task

	bootstrap := this.app.Label("__bootstrap")
	this.frame_chain.Bootstrap(bootstrap)
}

func (this *VirtualMachine) CanAdvance() bool {
	return this.frame_chain.CanAdvance()
}

func (this *VirtualMachine) Advance() {
	if !this.frame_chain.CanAdvance() {
		err := errors.New("frame chain cannot advance")
		panic(err)
	}

	this.garbage_collector.MarkAndSweep()

	bytecode := this.frame_chain.Advance()

	if this.verbose >= 1 {
		fmt.Printf("%s\n", bytecode.Stringify())
	}

	if bytecode.OpCode() == abi.NEW_SCOPE {
		this.frame_chain.LastFrame().FastScopeChain().NewScope()
	} else if bytecode.OpCode() == abi.DELETE_SCOPE {
		this.frame_chain.LastFrame().FastScopeChain().DeleteScope()
	} else if bytecode.OpCode() == abi.PUSH_CHAR {
		value := bytecode.Arg1()
		this.PushChar(value)
	} else if bytecode.OpCode() == abi.PUSH_SHORT {
		value := bytecode.Arg1()
		this.PushShort(value)
	} else if bytecode.OpCode() == abi.PUSH_INT {
		value := bytecode.Arg1()
		this.PushInt(value)
	} else if bytecode.OpCode() == abi.PUSH_LONG {
		value := bytecode.Arg1()
		this.PushLong(value)
	} else if bytecode.OpCode() == abi.PUSH_STRING {
		value := bytecode.Str1()
		this.PushString(value)
	} else if bytecode.OpCode() == abi.POP {
		this.Pop()
	} else if bytecode.OpCode() == abi.BEGIN_STRUCT {
		skeleton_name := bytecode.Str1()
		this.BeginStruct(skeleton_name)
	} else if bytecode.OpCode() == abi.APPEND_VOID {
		num_stars := bytecode.Arg1()
		field_name := bytecode.Str1()
		this.AppendVoid(num_stars, field_name)
	} else if bytecode.OpCode() == abi.APPEND_CHAR {
		num_stars := bytecode.Arg1()
		field_name := bytecode.Str1()
		this.AppendChar(num_stars, field_name)
	} else if bytecode.OpCode() == abi.APPEND_SHORT {
		num_stars := bytecode.Arg1()
		field_name := bytecode.Str1()
		this.AppendShort(num_stars, field_name)
	} else if bytecode.OpCode() == abi.APPEND_INT {
		num_stars := bytecode.Arg1()
		field_name := bytecode.Str1()
		this.AppendInt(num_stars, field_name)
	} else if bytecode.OpCode() == abi.APPEND_LONG {
		num_stars := bytecode.Arg1()
		field_name := bytecode.Str1()
		this.AppendLong(num_stars, field_name)
	} else if bytecode.OpCode() == abi.APPEND_STRUCT {
		num_stars := bytecode.Arg1()
		struct_name := bytecode.Str1()
		field_name := bytecode.Str2()
		this.AppendStruct(num_stars, struct_name, field_name)
	} else if bytecode.OpCode() == abi.END_STRUCT {
		this.EndStruct()
	} else if bytecode.OpCode() == abi.NEW_GLOBAL_VOID {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewGlobalVoid(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_GLOBAL_CHAR {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewGlobalChar(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_GLOBAL_SHORT {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewGlobalShort(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_GLOBAL_INT {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewGlobalInt(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_GLOBAL_LONG {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewGlobalLong(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_FAST_VOID {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewFastVoid(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_FAST_CHAR {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewFastChar(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_FAST_SHORT {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewFastShort(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_FAST_INT {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewFastInt(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_FAST_LONG {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewFastLong(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_FAST_STRUCT {
		num_stars := bytecode.Arg1()
		struct_name := bytecode.Str1()
		symbol_name := bytecode.Str2()
		this.NewFastStruct(num_stars, struct_name, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_ARG_VOID {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewArgVoid(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_ARG_CHAR {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewArgChar(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_ARG_SHORT {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewArgShort(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_ARG_INT {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewArgInt(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_ARG_LONG {
		num_stars := bytecode.Arg1()
		symbol_name := bytecode.Str1()
		this.NewArgLong(num_stars, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_ARG_STRUCT {
		num_stars := bytecode.Arg1()
		struct_name := bytecode.Str1()
		symbol_name := bytecode.Str2()
		this.NewArgStruct(num_stars, struct_name, symbol_name)
	} else if bytecode.OpCode() == abi.NEW_RETURN_VOID {
		num_stars := bytecode.Arg1()
		this.NewReturnVoid(num_stars)
	} else if bytecode.OpCode() == abi.NEW_RETURN_CHAR {
		num_stars := bytecode.Arg1()
		this.NewReturnChar(num_stars)
	} else if bytecode.OpCode() == abi.NEW_RETURN_SHORT {
		num_stars := bytecode.Arg1()
		this.NewReturnShort(num_stars)
	} else if bytecode.OpCode() == abi.NEW_RETURN_INT {
		num_stars := bytecode.Arg1()
		this.NewReturnInt(num_stars)
	} else if bytecode.OpCode() == abi.NEW_RETURN_LONG {
		num_stars := bytecode.Arg1()
		this.NewReturnLong(num_stars)
	} else if bytecode.OpCode() == abi.NEW_RETURN_STRUCT {
		num_stars := bytecode.Arg1()
		struct_name := bytecode.Str1()
		this.NewReturnStruct(num_stars, struct_name)
	} else if bytecode.OpCode() == abi.SIZEOF_VOID {
		num_stars := bytecode.Arg1()
		this.SizeofVoid(num_stars)
	} else if bytecode.OpCode() == abi.SIZEOF_CHAR {
		num_stars := bytecode.Arg1()
		this.SizeofChar(num_stars)
	} else if bytecode.OpCode() == abi.SIZEOF_SHORT {
		num_stars := bytecode.Arg1()
		this.SizeofShort(num_stars)
	} else if bytecode.OpCode() == abi.SIZEOF_INT {
		num_stars := bytecode.Arg1()
		this.SizeofInt(num_stars)
	} else if bytecode.OpCode() == abi.SIZEOF_LONG {
		num_stars := bytecode.Arg1()
		this.SizeofLong(num_stars)
	} else if bytecode.OpCode() == abi.SIZEOF_STRUCT {
		num_stars := bytecode.Arg1()
		struct_name := bytecode.Str1()
		this.SizeofStruct(num_stars, struct_name)
	} else if bytecode.OpCode() == abi.GET_IDENTIFIER {
		symbol_name := bytecode.Str1()
		this.GetIdentifier(symbol_name)
	} else if bytecode.OpCode() == abi.GET_ARG_IDENTIFIER {
		symbol_name := bytecode.Str1()
		this.GetArgIdentifier(symbol_name)
	} else if bytecode.OpCode() == abi.GET_SUBSCRIPT {
		this.GetSubscript()
	} else if bytecode.OpCode() == abi.GET_ACCESS {
		field_name := bytecode.Str1()
		this.GetAccess(field_name)
	} else if bytecode.OpCode() == abi.GET_REFERENCE {
		field_name := bytecode.Str1()
		this.GetReference(field_name)
	} else if bytecode.OpCode() == abi.GET_ADDRESS {
		this.GetAddress()
	} else if bytecode.OpCode() == abi.GET_VALUE {
		this.GetValue()
	} else if bytecode.OpCode() == abi.ALLOC {
		this.Alloc()
	} else if bytecode.OpCode() == abi.FREE {
		this.Free()
	} else if bytecode.OpCode() == abi.ASSERT {
		this.Assert()
	} else if bytecode.OpCode() == abi.ADD {
		this.Add()
	} else if bytecode.OpCode() == abi.SUB {
		this.Sub()
	} else if bytecode.OpCode() == abi.MUL {
		this.Mul()
	} else if bytecode.OpCode() == abi.DIV {
		this.Div()
	} else if bytecode.OpCode() == abi.MOD {
		this.Mod()
	} else if bytecode.OpCode() == abi.LSHIFT {
		this.Lshift()
	} else if bytecode.OpCode() == abi.RSHIFT {
		this.Rshift()
	} else if bytecode.OpCode() == abi.NEGATE {
		this.Negate()
	} else if bytecode.OpCode() == abi.TILDE {
		this.Tilde()
	} else if bytecode.OpCode() == abi.SQRT {
		this.Sqrt()
	} else if bytecode.OpCode() == abi.BITWISE_AND {
		this.BitwiseAnd()
	} else if bytecode.OpCode() == abi.BITWISE_XOR {
		this.BitwiseXor()
	} else if bytecode.OpCode() == abi.BITWISE_OR {
		this.BitwiseOr()
	} else if bytecode.OpCode() == abi.LOGICAL_AND {
		this.LogicalAnd()
	} else if bytecode.OpCode() == abi.LOGICAL_OR {
		this.LogicalOr()
	} else if bytecode.OpCode() == abi.LOGICAL_NOT {
		this.LogicalNot()
	} else if bytecode.OpCode() == abi.EQ {
		this.Eq()
	} else if bytecode.OpCode() == abi.NOT_EQ {
		this.NotEq()
	} else if bytecode.OpCode() == abi.LESS {
		this.Less()
	} else if bytecode.OpCode() == abi.LESS_EQ {
		this.LessEq()
	} else if bytecode.OpCode() == abi.GREATER {
		this.Greater()
	} else if bytecode.OpCode() == abi.GREATER_EQ {
		this.GreaterEq()
	} else if bytecode.OpCode() == abi.CONDITIONAL {
		this.Conditional()
	} else if bytecode.OpCode() == abi.ASSIGN {
		this.Assign()
	} else if bytecode.OpCode() == abi.ASSIGN_STAR {
		this.AssignStar()
	} else if bytecode.OpCode() == abi.ASSIGN_DIV {
		this.AssignDiv()
	} else if bytecode.OpCode() == abi.ASSIGN_MOD {
		this.AssignMod()
	} else if bytecode.OpCode() == abi.ASSIGN_ADD {
		this.AssignAdd()
	} else if bytecode.OpCode() == abi.ASSIGN_SUB {
		this.AssignSub()
	} else if bytecode.OpCode() == abi.ASSIGN_LSHIFT {
		this.AssignLshift()
	} else if bytecode.OpCode() == abi.ASSIGN_RSHIFT {
		this.AssignRshift()
	} else if bytecode.OpCode() == abi.ASSIGN_BITWISE_AND {
		this.AssignBitwiseAnd()
	} else if bytecode.OpCode() == abi.ASSIGN_BITWISE_XOR {
		this.AssignBitwiseXor()
	} else if bytecode.OpCode() == abi.ASSIGN_BITWISE_OR {
		this.AssignBitwiseOr()
	} else if bytecode.OpCode() == abi.ASSIGN_PLUS_PLUS {
		this.AssignPlusPlus()
	} else if bytecode.OpCode() == abi.ASSIGN_MINUS_MINUS {
		this.AssignMinusMinus()
	} else if bytecode.OpCode() == abi.ASSIGN_RETURN {
		this.AssignReturn()
	} else if bytecode.OpCode() == abi.JUMP {
		label_name := bytecode.Str1()
		label := this.app.Label(label_name)

		this.Jump(label)
	} else if bytecode.OpCode() == abi.JUMP_IF_ZERO {
		label_name := bytecode.Str1()
		label := this.app.Label(label_name)

		this.JumpIfZero(label)
	} else if bytecode.OpCode() == abi.JUMP_IF_NONZERO {
		label_name := bytecode.Str1()
		label := this.app.Label(label_name)

		this.JumpIfNonZero(label)
	} else if bytecode.OpCode() == abi.CALL {
		label_name := bytecode.Str1()
		label := this.app.Label(label_name)

		this.Call(label)
	} else if bytecode.OpCode() == abi.RETURN {
		this.Return()
	} else if bytecode.OpCode() == abi.NOP {
		this.Nop()
	} else if bytecode.OpCode() == abi.DPU_ALLOC {
		dpu_id := bytecode.Arg1()

		this.DpuAlloc(dpu_id)
	} else if bytecode.OpCode() == abi.DPU_LOAD {
		this.DpuLoad()
	} else if bytecode.OpCode() == abi.DPU_PREPARE {
		this.DpuPrepare()
	} else if bytecode.OpCode() == abi.DPU_TRANSFER {
		this.DpuTransfer()
	} else if bytecode.OpCode() == abi.DPU_COPY_TO {
		this.DpuCopyTo()
	} else if bytecode.OpCode() == abi.DPU_COPY_FROM {
		this.DpuCopyFrom()
	} else if bytecode.OpCode() == abi.DPU_LAUNCH {
		this.DpuLaunch()
	} else if bytecode.OpCode() == abi.DPU_FREE {
		this.DpuFree()
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	if this.verbose >= 2 {
		fmt.Printf("%s\n", this.Stringify())
	}
}

func (this *VirtualMachine) PushChar(value int64) {
	object := this.arena.NewChar(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.CHAR, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) PushShort(value int64) {
	object := this.arena.NewShort(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.SHORT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) PushInt(value int64) {
	object := this.arena.NewInt(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) PushLong(value int64) {
	object := this.arena.NewLong(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.LONG, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) PushString(value string) {
	object := this.arena.NewString(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.STRING, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Pop() {
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) BeginStruct(skeleton_name string) {
	skeleton := new(type_system.Skeleton)
	skeleton.Init(skeleton_name)

	this.registry.AddSkeleton(skeleton)

	this.cur_skeleton_name = new(string)
	*this.cur_skeleton_name = skeleton_name
}

func (this *VirtualMachine) AppendVoid(num_stars int64, field_name string) {
	skeleton := this.registry.Skeleton(*this.cur_skeleton_name)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.VOID, int(num_stars))

	field := new(type_system.Field)
	field.Init(type_variable, field_name)

	skeleton.Append(field)
}

func (this *VirtualMachine) AppendChar(num_stars int64, field_name string) {
	skeleton := this.registry.Skeleton(*this.cur_skeleton_name)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.CHAR, int(num_stars))

	field := new(type_system.Field)
	field.Init(type_variable, field_name)

	skeleton.Append(field)
}

func (this *VirtualMachine) AppendShort(num_stars int64, field_name string) {
	skeleton := this.registry.Skeleton(*this.cur_skeleton_name)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.SHORT, int(num_stars))

	field := new(type_system.Field)
	field.Init(type_variable, field_name)

	skeleton.Append(field)
}

func (this *VirtualMachine) AppendInt(num_stars int64, field_name string) {
	skeleton := this.registry.Skeleton(*this.cur_skeleton_name)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, int(num_stars))

	field := new(type_system.Field)
	field.Init(type_variable, field_name)

	skeleton.Append(field)
}

func (this *VirtualMachine) AppendLong(num_stars int64, field_name string) {
	skeleton := this.registry.Skeleton(*this.cur_skeleton_name)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.LONG, int(num_stars))

	field := new(type_system.Field)
	field.Init(type_variable, field_name)

	skeleton.Append(field)
}

func (this *VirtualMachine) AppendStruct(num_stars int64, struct_name string, field_name string) {
	skeleton := this.registry.Skeleton(*this.cur_skeleton_name)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitStruct(type_system.STRUCT, struct_name, int(num_stars))

	field := new(type_system.Field)
	field.Init(type_variable, field_name)

	skeleton.Append(field)
}

func (this *VirtualMachine) EndStruct() {
	this.cur_skeleton_name = nil
}

func (this *VirtualMachine) NewGlobalVoid(num_stars int64, symbol_name string) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.VOID, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.GlobalScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewGlobalChar(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewChar(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.CHAR, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.GlobalScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewGlobalShort(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewShort(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.SHORT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.GlobalScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewGlobalInt(num_stars int64, symbol_name string) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.GlobalScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewGlobalLong(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewLong(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.LONG, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.GlobalScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewFastVoid(num_stars int64, symbol_name string) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.VOID, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().FastScopeChain().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewFastChar(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewChar(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.CHAR, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().FastScopeChain().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewFastShort(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewShort(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.SHORT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().FastScopeChain().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewFastInt(num_stars int64, symbol_name string) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().FastScopeChain().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewFastLong(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewLong(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.LONG, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().FastScopeChain().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewFastStruct(num_stars int64, struct_name string, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewStruct(struct_name, this.registry.SkeletonSize(struct_name))
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitStruct(type_system.STRUCT, struct_name, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().FastScopeChain().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewArgVoid(num_stars int64, symbol_name string) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.VOID, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().ArgScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewArgChar(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewChar(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.CHAR, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().ArgScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewArgShort(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewShort(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.SHORT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().ArgScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewArgInt(num_stars int64, symbol_name string) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().ArgScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewArgLong(num_stars int64, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewLong(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.LONG, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().ArgScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewArgStruct(num_stars int64, struct_name string, symbol_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewStruct(struct_name, this.registry.SkeletonSize(struct_name))
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitStruct(type_system.STRUCT, struct_name, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	symbol_ := new(symbol.Symbol)
	symbol_.Init(symbol_name, type_variable, object)

	this.frame_chain.LastFrame().ArgScope().AddSymbol(symbol_)
}

func (this *VirtualMachine) NewReturnVoid(num_stars int64) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.VOID, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().ReturnStack().Push(stack_item)
}

func (this *VirtualMachine) NewReturnChar(num_stars int64) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewChar(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.CHAR, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().ReturnStack().Push(stack_item)
}

func (this *VirtualMachine) NewReturnShort(num_stars int64) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewShort(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.SHORT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().ReturnStack().Push(stack_item)
}

func (this *VirtualMachine) NewReturnInt(num_stars int64) {
	object := this.arena.NewInt(0)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().ReturnStack().Push(stack_item)
}

func (this *VirtualMachine) NewReturnLong(num_stars int64) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewLong(0)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.LONG, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().ReturnStack().Push(stack_item)
}

func (this *VirtualMachine) NewReturnStruct(num_stars int64, struct_name string) {
	var object *base.Object
	if num_stars > 0 {
		object = this.arena.NewInt(0)
	} else {
		object = this.arena.NewStruct(struct_name, this.registry.SkeletonSize(struct_name))
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitStruct(type_system.STRUCT, struct_name, int(num_stars))

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().ReturnStack().Push(stack_item)
}

func (this *VirtualMachine) SizeofVoid(num_stars int64) {
	var value int64
	if num_stars == 0 {
		err := errors.New("num stars == 0")
		panic(err)
	} else {
		value = 4
	}

	object := this.arena.NewInt(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) SizeofChar(num_stars int64) {
	var value int64
	if num_stars == 0 {
		value = 1
	} else {
		value = 4
	}

	object := this.arena.NewInt(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) SizeofShort(num_stars int64) {
	var value int64
	if num_stars == 0 {
		value = 2
	} else {
		value = 4
	}

	object := this.arena.NewInt(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) SizeofInt(num_stars int64) {
	value := int64(4)

	object := this.arena.NewInt(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) SizeofLong(num_stars int64) {
	var value int64
	if num_stars == 0 {
		value = 8
	} else {
		value = 4
	}

	object := this.arena.NewInt(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) SizeofStruct(num_stars int64, struct_name string) {
	var value int64
	if num_stars == 0 {
		value = this.registry.SkeletonSize(struct_name)
	} else {
		value = 4
	}

	object := this.arena.NewInt(value)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GetIdentifier(symbol_name string) {
	var symbol_ *symbol.Symbol
	if this.frame_chain.LastFrame().FastScopeChain().HasSymbol(symbol_name) {
		symbol_ = this.frame_chain.LastFrame().FastScopeChain().Symbol(symbol_name)
	} else if this.frame_chain.GlobalScope().HasSymbol(symbol_name) {
		symbol_ = this.frame_chain.GlobalScope().Symbol(symbol_name)
	} else {
		err_msg := fmt.Sprintf("symbol (%s) is not found", symbol_name)
		err := errors.New(err_msg)
		panic(err)
	}

	type_variable := symbol_.TypeVariable()
	object := symbol_.Object()

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GetArgIdentifier(symbol_name string) {
	var symbol_ *symbol.Symbol
	if this.frame_chain.LastFrame().ArgScope().HasSymbol(symbol_name) {
		symbol_ = this.frame_chain.LastFrame().ArgScope().Symbol(symbol_name)
	} else {
		err_msg := fmt.Sprintf("symbol (%s) is not found", symbol_name)
		err := errors.New(err_msg)
		panic(err)
	}

	type_variable := symbol_.TypeVariable()
	object := symbol_.Object()

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GetSubscript() {
	base_ := this.frame_chain.LastFrame().Stack().Front(1)
	index := this.frame_chain.LastFrame().Stack().Front(0)

	base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()
	index_value := this.arena.Pool().Memory().Read(index.Address(), index.Size()).SignedValue()

	var offset int64
	var size int64
	type_variable := new(type_system.TypeVariable)
	if base_.TypeVariable().NumStars() == 0 {
		err := errors.New("base is not a pointer")
		panic(err)
	} else if base_.TypeVariable().NumStars() == 1 {
		if base_.TypeVariable().TypeVariableType() == type_system.VOID {
			err := errors.New("type variable type is void")
			panic(err)
		} else if base_.TypeVariable().TypeVariableType() == type_system.CHAR {
			offset = index_value

			size = 1

			type_variable.InitPrimitive(type_system.CHAR, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.SHORT {
			offset = index_value * 2

			size = 2

			type_variable.InitPrimitive(type_system.SHORT, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.INT {
			offset = index_value * 4

			size = 4

			type_variable.InitPrimitive(type_system.INT, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.LONG {
			offset = index_value * 8

			size = 8

			type_variable.InitPrimitive(type_system.LONG, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.STRUCT {
			struct_name := base_.TypeVariable().StructName()

			offset = index_value * this.registry.SkeletonSize(struct_name)

			size = this.registry.SkeletonSize(struct_name)

			type_variable.InitStruct(type_system.STRUCT, struct_name, base_.TypeVariable().NumStars()-1)
		} else {
			err := errors.New("type variable type is not valid")
			panic(err)
		}
	} else {
		offset = index_value * 4
		size = 4

		if base_.TypeVariable().TypeVariableType() == type_system.VOID {
			type_variable.InitPrimitive(type_system.VOID, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.CHAR {
			type_variable.InitPrimitive(type_system.CHAR, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.SHORT {
			type_variable.InitPrimitive(type_system.SHORT, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.INT {
			type_variable.InitPrimitive(type_system.INT, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.LONG {
			type_variable.InitPrimitive(type_system.LONG, base_.TypeVariable().NumStars()-1)
		} else if base_.TypeVariable().TypeVariableType() == type_system.STRUCT {
			struct_name := base_.TypeVariable().StructName()

			type_variable.InitStruct(type_system.STRUCT, struct_name, base_.TypeVariable().NumStars()-1)
		} else {
			err := errors.New("type variable type is not valid")
			panic(err)
		}
	}

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, base_value+offset, size)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GetAccess(field_name string) {
	base_ := this.frame_chain.LastFrame().Stack().Front(0)

	if base_.TypeVariable().NumStars() != 0 {
		err := errors.New("base is a pointer")
		panic(err)
	} else if base_.TypeVariable().TypeVariableType() != type_system.STRUCT {
		err := errors.New("base is not a struct")
		panic(err)
	}

	struct_name := base_.TypeVariable().StructName()

	offset := this.registry.FieldOffset(struct_name, field_name)

	field := this.registry.Skeleton(struct_name).Field(field_name)

	type_variable := field.TypeVariable()

	var size int64
	if field.TypeVariable().TypeVariableType() == type_system.VOID {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			err := errors.New("void type has no star")
			panic(err)
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.CHAR {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = 1
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.SHORT {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = 2
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.INT {
		size = 4
	} else if field.TypeVariable().TypeVariableType() == type_system.LONG {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = 8
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.STRUCT {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = this.registry.SkeletonSize(base_.TypeVariable().StructName())
		}
	} else {
		err := errors.New("type variable type is not valid")
		panic(err)
	}

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, base_.Address()+offset, size)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GetReference(field_name string) {
	base_ := this.frame_chain.LastFrame().Stack().Front(0)

	if base_.TypeVariable().NumStars() == 0 {
		err := errors.New("base is not a pointer")
		panic(err)
	} else if base_.TypeVariable().NumStars() > 1 {
		err := errors.New("base is a multi-dimensional pointer")
		panic(err)
	} else if base_.TypeVariable().TypeVariableType() != type_system.STRUCT {
		err := errors.New("base is not a struct")
		panic(err)
	}

	struct_name := base_.TypeVariable().StructName()

	offset := this.registry.FieldOffset(struct_name, field_name)

	field := this.registry.Skeleton(struct_name).Field(field_name)

	type_variable := field.TypeVariable()

	var size int64
	if field.TypeVariable().TypeVariableType() == type_system.VOID {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			err := errors.New("void type has no star")
			panic(err)
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.CHAR {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = 1
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.SHORT {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = 2
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.INT {
		size = 4
	} else if field.TypeVariable().TypeVariableType() == type_system.LONG {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = 8
		}
	} else if field.TypeVariable().TypeVariableType() == type_system.STRUCT {
		if field.TypeVariable().NumStars() > 0 {
			size = 4
		} else {
			size = this.registry.SkeletonSize(base_.TypeVariable().StructName())
		}
	} else {
		err := errors.New("type variable type is not valid")
		panic(err)
	}

	address := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, address+offset, size)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GetAddress() {
	base_ := this.frame_chain.LastFrame().Stack().Front(0)

	object := this.arena.NewInt(base_.Address())

	type_variable := new(type_system.TypeVariable)
	if base_.TypeVariable().TypeVariableType() == type_system.VOID {
		type_variable.InitPrimitive(type_system.VOID, base_.TypeVariable().NumStars()+1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.CHAR {
		type_variable.InitPrimitive(type_system.CHAR, base_.TypeVariable().NumStars()+2)
	} else if base_.TypeVariable().TypeVariableType() == type_system.SHORT {
		type_variable.InitPrimitive(type_system.SHORT, base_.TypeVariable().NumStars()+1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.INT {
		type_variable.InitPrimitive(type_system.INT, base_.TypeVariable().NumStars()+1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.LONG {
		type_variable.InitPrimitive(type_system.LONG, base_.TypeVariable().NumStars()+1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.STRUCT {
		type_variable.InitStruct(type_system.STRUCT, base_.TypeVariable().StructName(), base_.TypeVariable().NumStars()+1)
	} else {
		err := errors.New("type variable type is not valid")
		panic(err)
	}

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GetValue() {
	base_ := this.frame_chain.LastFrame().Stack().Front(0)

	base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

	if base_.TypeVariable().NumStars() < 1 {
		err := errors.New("base is not a pointer")
		panic(err)
	}

	object := this.arena.Pool().Object(base_value)

	type_variable := new(type_system.TypeVariable)
	if base_.TypeVariable().TypeVariableType() == type_system.VOID {
		type_variable.InitPrimitive(type_system.VOID, base_.TypeVariable().NumStars()-1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.CHAR {
		type_variable.InitPrimitive(type_system.CHAR, base_.TypeVariable().NumStars()-2)
	} else if base_.TypeVariable().TypeVariableType() == type_system.SHORT {
		type_variable.InitPrimitive(type_system.SHORT, base_.TypeVariable().NumStars()-1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.INT {
		type_variable.InitPrimitive(type_system.INT, base_.TypeVariable().NumStars()-1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.LONG {
		type_variable.InitPrimitive(type_system.LONG, base_.TypeVariable().NumStars()-1)
	} else if base_.TypeVariable().TypeVariableType() == type_system.STRUCT {
		type_variable.InitStruct(type_system.STRUCT, base_.TypeVariable().StructName(), base_.TypeVariable().NumStars()-1)
	} else {
		err := errors.New("type variable type is not valid")
		panic(err)
	}

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Alloc() {
	base_ := this.frame_chain.LastFrame().Stack().Front(0)

	base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

	object := this.arena.NewPointer(base_value)
	pointer := this.arena.NewInt(object.Address())

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.VOID, 1)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, pointer.Address(), pointer.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Free() {
	base_ := this.frame_chain.LastFrame().Stack().Front(0)

	base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

	this.arena.Free(base_value)

	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) Assert() {
	base_ := this.frame_chain.LastFrame().Stack().Front(0)

	base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

	if base_value == 0 {
		err := errors.New("assert")
		panic(err)
	}

	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) Add() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value + roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be added")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Sub() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value - roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be subtracted")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Mul() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value * roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Div() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value / roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be divided")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Mod() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value % roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the modular operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Lshift() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value << roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be left shifted")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Rshift() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value >> roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be right shifted")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Negate() {
	operand := this.frame_chain.LastFrame().Stack().Front(0)

	if operand.TypeVariable().NumStars() != 0 {
		err := errors.New("operand's num stars != 0")
		panic(err)
	}

	operand_value := this.arena.Pool().
		Memory().
		Read(operand.Address(), operand.Size()).
		SignedValue()

	value := -operand_value

	var object *base.Object
	if operand.TypeVariable().TypeVariableType() == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.INT {
		object = this.arena.NewInt(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be negated")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(operand.TypeVariable().TypeVariableType(), 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Tilde() {
	operand := this.frame_chain.LastFrame().Stack().Front(0)

	if operand.TypeVariable().NumStars() != 0 {
		err := errors.New("operand's num stars != 0")
		panic(err)
	}

	operand_value := this.arena.Pool().
		Memory().
		Read(operand.Address(), operand.Size()).
		SignedValue()

	value := ^operand_value

	var object *base.Object
	if operand.TypeVariable().TypeVariableType() == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.INT {
		object = this.arena.NewInt(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be tilded")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(operand.TypeVariable().TypeVariableType(), 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Sqrt() {
	operand := this.frame_chain.LastFrame().Stack().Front(0)

	if operand.TypeVariable().NumStars() != 0 {
		err := errors.New("operand's num stars != 0")
		panic(err)
	}

	operand_value := this.arena.Pool().
		Memory().
		Read(operand.Address(), operand.Size()).
		SignedValue()

	value := int64(math.Sqrt(float64(operand_value)))

	var object *base.Object
	if operand.TypeVariable().TypeVariableType() == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.INT {
		object = this.arena.NewInt(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be tilded")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(operand.TypeVariable().TypeVariableType(), 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) BitwiseAnd() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value & roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the bitwise and operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) BitwiseXor() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value ^ roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the bitwise xor operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) BitwiseOr() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	value := loperand_value | roperand_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the bitwise or operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) LogicalAnd() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value != 0 && roperand_value != 0 {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the bitwise or operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) LogicalOr() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value != 0 || roperand_value != 0 {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the bitwise or operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) LogicalNot() {
	operand := this.frame_chain.LastFrame().Stack().Front(0)

	if operand.TypeVariable().NumStars() != 0 {
		err := errors.New("operand's num stars != 0")
		panic(err)
	}

	operand_value := this.arena.Pool().
		Memory().
		Read(operand.Address(), operand.Size()).
		SignedValue()

	var value int64
	if operand_value != 0 {
		value = 0
	} else {
		value = 1
	}

	var object *base.Object
	if operand.TypeVariable().TypeVariableType() == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.INT {
		object = this.arena.NewInt(value)
	} else if operand.TypeVariable().TypeVariableType() == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot be tilded")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(operand.TypeVariable().TypeVariableType(), 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Eq() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value == roperand_value {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the eq operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) NotEq() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value != roperand_value {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the not eq operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Less() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value < roperand_value {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the less operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) LessEq() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value <= roperand_value {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the less eq operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Greater() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value > roperand_value {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the greater operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) GreaterEq() {
	loperand := this.frame_chain.LastFrame().Stack().Front(1)
	roperand := this.frame_chain.LastFrame().Stack().Front(0)

	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(loperand, roperand)

	loperand_value := this.arena.Pool().
		Memory().
		Read(loperand.Address(), loperand.Size()).
		SignedValue()
	roperand_value := this.arena.Pool().
		Memory().
		Read(roperand.Address(), roperand.Size()).
		SignedValue()

	var value int64
	if loperand_value >= roperand_value {
		value = 1
	} else {
		value = 0
	}

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("loperand and roperand cannot conduct the greater eq operation")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Push(stack_item)
}

func (this *VirtualMachine) Conditional() {
	condition_stack_item := this.frame_chain.LastFrame().Stack().Front(2)
	true_stack_item := this.frame_chain.LastFrame().Stack().Front(1)
	false_stack_item := this.frame_chain.LastFrame().Stack().Front(0)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()

	condition_value := this.arena.Pool().
		Memory().
		Read(condition_stack_item.Address(), condition_stack_item.Size()).
		SignedValue()

	if condition_value != 0 {
		this.frame_chain.LastFrame().Stack().Push(true_stack_item)
	} else {
		this.frame_chain.LastFrame().Stack().Push(false_stack_item)
	}
}

func (this *VirtualMachine) Assign() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	rvalue_byte_stream := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size())

	if lvalue.Size() > rvalue_byte_stream.Size() {
		for lvalue.Size() != rvalue_byte_stream.Size() {
			if rvalue_byte_stream.Signbit() {
				rvalue_byte_stream.Append(255)
			} else {
				rvalue_byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < rvalue_byte_stream.Size() {
		for lvalue.Size() != rvalue_byte_stream.Size() {
			rvalue_byte_stream.Remove(int(rvalue_byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), rvalue_byte_stream)

	if this.arena.Pool().HasObject(lvalue.Address()) {
		lvalue_object := this.arena.Pool().Object(lvalue.Address())

		if !lvalue_object.HasTypeVariable() {
			lvalue_object.SetTypeVariable(rvalue.TypeVariable())
		}
	}

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignStar() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value * rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignDiv() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value / rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignMod() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value % rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignAdd() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value + rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignSub() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value - rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignLshift() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value << rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignRshift() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value >> rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignBitwiseAnd() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value & rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignBitwiseXor() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value ^ rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignBitwiseOr() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(1)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	} else if rvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("rvalue's num stars != 0")
		panic(err)
	}

	type_variable_type := this.TypeCast(lvalue, rvalue)

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()
	rvalue_value := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size()).SignedValue()

	value := lvalue_value | rvalue_value

	var object *base.Object
	if type_variable_type == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if type_variable_type == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if type_variable_type == type_system.INT {
		object = this.arena.NewInt(value)
	} else if type_variable_type == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_variable_type, 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignPlusPlus() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	}

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()

	value := lvalue_value + 1

	var object *base.Object
	if lvalue.TypeVariable().TypeVariableType() == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if lvalue.TypeVariable().TypeVariableType() == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if lvalue.TypeVariable().TypeVariableType() == type_system.INT {
		object = this.arena.NewInt(value)
	} else if lvalue.TypeVariable().TypeVariableType() == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(lvalue.TypeVariable().TypeVariableType(), 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignMinusMinus() {
	lvalue := this.frame_chain.LastFrame().Stack().Front(0)

	if lvalue.TypeVariable().NumStars() != 0 {
		err := errors.New("lvalue's num stars != 0")
		panic(err)
	}

	lvalue_value := this.arena.Pool().Memory().Read(lvalue.Address(), lvalue.Size()).SignedValue()

	value := lvalue_value - 1

	var object *base.Object
	if lvalue.TypeVariable().TypeVariableType() == type_system.CHAR {
		object = this.arena.NewChar(value)
	} else if lvalue.TypeVariable().TypeVariableType() == type_system.SHORT {
		object = this.arena.NewShort(value)
	} else if lvalue.TypeVariable().TypeVariableType() == type_system.INT {
		object = this.arena.NewInt(value)
	} else if lvalue.TypeVariable().TypeVariableType() == type_system.LONG {
		object = this.arena.NewLong(value)
	} else {
		err := errors.New("lvalue and rvalue cannot be multiplied")
		panic(err)
	}

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(lvalue.TypeVariable().TypeVariableType(), 0)

	stack_item := new(stack.StackItem)
	stack_item.Init(type_variable, object.Address(), object.Size())

	byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

	if lvalue.Size() > byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			if byte_stream.Signbit() {
				byte_stream.Append(255)
			} else {
				byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < byte_stream.Size() {
		for lvalue.Size() != byte_stream.Size() {
			byte_stream.Remove(int(byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) AssignReturn() {
	lvalue := this.frame_chain.LastFrame().ReturnStack().Front(0)
	rvalue := this.frame_chain.LastFrame().Stack().Front(0)

	rvalue_byte_stream := this.arena.Pool().Memory().Read(rvalue.Address(), rvalue.Size())

	if lvalue.Size() > rvalue_byte_stream.Size() {
		for lvalue.Size() != rvalue_byte_stream.Size() {
			if rvalue_byte_stream.Signbit() {
				rvalue_byte_stream.Append(255)
			} else {
				rvalue_byte_stream.Append(0)
			}
		}
	} else if lvalue.Size() < rvalue_byte_stream.Size() {
		for lvalue.Size() != rvalue_byte_stream.Size() {
			rvalue_byte_stream.Remove(int(rvalue_byte_stream.Size() - 1))
		}
	}

	this.arena.Pool().Memory().Write(lvalue.Address(), lvalue.Size(), rvalue_byte_stream)

	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) Jump(label *abi.Label) {
	this.frame_chain.LastFrame().Pc().Jump(label)
}

func (this *VirtualMachine) JumpIfZero(label *abi.Label) {
	value := this.frame_chain.LastFrame().Stack().Front(0)

	value_value := this.arena.Pool().Memory().Read(value.Address(), value.Size()).SignedValue()

	if value_value == 0 {
		this.frame_chain.LastFrame().Pc().Jump(label)
	}

	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) JumpIfNonZero(label *abi.Label) {
	value := this.frame_chain.LastFrame().Stack().Front(0)

	value_value := this.arena.Pool().Memory().Read(value.Address(), value.Size()).SignedValue()

	if value_value != 0 {
		this.frame_chain.LastFrame().Pc().Jump(label)
	}

	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) Call(label *abi.Label) {
	this.frame_chain.NewFrame(label)
}

func (this *VirtualMachine) Return() {
	this.frame_chain.DeleteFrame()
}

func (this *VirtualMachine) Nop() {
}

func (this *VirtualMachine) DpuAlloc(dpu_id int64) {
	if int(dpu_id) > this.num_dpus {
		err := errors.New("DpuAlloc allocates more number of DPUs than available")
		panic(err)
	}
}

func (this *VirtualMachine) DpuLoad() {
	thread_pool := new(core.ThreadPool)
	thread_pool.Init()

	for _, dpu_ := range this.Dpus() {
		dpu_load_job := new(DpuLoadJob)
		dpu_load_job.Init(this.task, dpu_)

		thread_pool.Enque(dpu_load_job)
	}

	thread_pool.Start()
}

func (this *VirtualMachine) DpuPrepare() {
	dpu_id := this.frame_chain.LastFrame().Stack().Front(1)
	pointer := this.frame_chain.LastFrame().Stack().Front(0)

	dpu_id_value := this.arena.Pool().Memory().Read(dpu_id.Address(), dpu_id.Size()).SignedValue()
	pointer_value := this.arena.Pool().
		Memory().
		Read(pointer.Address(), pointer.Size()).
		SignedValue()

	dpu_ := this.Dpus()[dpu_id_value]

	this.prepare_xfer_buf[dpu_] = pointer_value

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) DpuTransfer() {
	direction := this.frame_chain.LastFrame().Stack().Front(4)
	base_ := this.frame_chain.LastFrame().Stack().Front(3)
	offset := this.frame_chain.LastFrame().Stack().Front(2)
	size := this.frame_chain.LastFrame().Stack().Front(1)

	direction_value := this.arena.Pool().
		Memory().
		Read(direction.Address(), direction.Size()).
		SignedValue()
	offset_value := this.arena.Pool().Memory().Read(offset.Address(), offset.Size()).SignedValue()
	size_value := this.arena.Pool().Memory().Read(size.Address(), size.Size()).SignedValue()

	for _, dpu_ := range this.Dpus() {
		if pointer_value, pointer_found := this.prepare_xfer_buf[dpu_]; pointer_found {
			delete(this.prepare_xfer_buf, dpu_)

			if direction_value == 0 {
				this.Checkpoint()

				if base_.TypeVariable().TypeVariableType() == type_system.STRING {
					base_string := this.DecodeString(
						this.arena.Pool().Memory().Read(base_.Address(), base_.Size()),
					)
					base_string = base_string[1 : len(base_string)-1]

					wram_address, wram_address_found := this.task.Addresses()[base_string]
					if !wram_address_found {
						err_msg := fmt.Sprintf("WRAM address (%s) is not found", base_string)
						err := errors.New(err_msg)
						panic(err)
					}

					byte_stream := this.PrepareByteStream(pointer_value, offset_value, size_value)
					dpu_.Dma().
						TransferToWram(wram_address+offset_value, byte_stream.Size(), byte_stream)
				} else {
					base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

					channel_id := dpu_.ChannelId()
					rank_id := dpu_.RankId()
					dpu_id := dpu_.DpuId()

					transfer_command := new(bank.TransferCommand)
					transfer_command.Init(bank.HOST_TO_DEVICE, pointer_value, channel_id, rank_id, dpu_id, base_value+offset_value, size_value)

					this.memory_controller.Push(transfer_command)
					this.push_xfer[transfer_command] = true
				}
			} else if direction_value == 1 {
				if base_.TypeVariable().TypeVariableType() == type_system.STRING {
					base_string := this.DecodeString(this.arena.Pool().Memory().Read(base_.Address(), base_.Size()))
					base_string = base_string[1 : len(base_string)-1]

					wram_address, wram_address_found := this.task.Addresses()[base_string]
					if !wram_address_found {
						err_msg := fmt.Sprintf("WRAM address (%s) is not found", base_string)
						err := errors.New(err_msg)
						panic(err)
					}

					byte_stream := dpu_.Dma().TransferFromWram(wram_address+offset_value, size_value)
					this.arena.Pool().Memory().Write(pointer_value, size_value, byte_stream)
				} else {
					base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

					channel_id := dpu_.ChannelId()
					rank_id := dpu_.RankId()
					dpu_id := dpu_.DpuId()

					transfer_command := new(bank.TransferCommand)
					transfer_command.Init(bank.DEVICE_TO_HOST, pointer_value, channel_id, rank_id, dpu_id, base_value+offset_value, size_value)

					this.memory_controller.Push(transfer_command)
					this.push_xfer[transfer_command] = true
				}
			} else {
				err := errors.New("direction value is not 0 nor 1")
				panic(err)
			}
		}
	}

	this.SimulateMemory()

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) DpuCopyTo() {
	dpu_ := this.frame_chain.LastFrame().Stack().Front(4)
	base_ := this.frame_chain.LastFrame().Stack().Front(3)
	offset := this.frame_chain.LastFrame().Stack().Front(2)
	pointer := this.frame_chain.LastFrame().Stack().Front(1)
	size := this.frame_chain.LastFrame().Stack().Front(0)

	dpu_value := this.arena.Pool().Memory().Read(dpu_.Address(), dpu_.Size()).SignedValue()
	offset_value := this.arena.Pool().Memory().Read(offset.Address(), offset.Size()).SignedValue()
	pointer_value := this.arena.Pool().
		Memory().
		Read(pointer.Address(), pointer.Size()).
		SignedValue()
	size_value := this.arena.Pool().Memory().Read(size.Address(), size.Size()).SignedValue()

	this.Checkpoint()

	if base_.TypeVariable().TypeVariableType() == type_system.STRING {
		base_string := this.DecodeString(
			this.arena.Pool().Memory().Read(base_.Address(), base_.Size()),
		)
		base_string = base_string[1 : len(base_string)-1]

		wram_address, wram_address_found := this.task.Addresses()[base_string]
		if !wram_address_found {
			err_msg := fmt.Sprintf("WRAM address (%s) is not found", base_string)
			err := errors.New(err_msg)
			panic(err)
		}

		byte_stream := this.PrepareByteStream(pointer_value, offset_value, size_value)
		this.Dpus()[dpu_value].Dma().
			TransferToWram(wram_address+offset_value, byte_stream.Size(), byte_stream)
	} else {
		base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

		channel_id := int(dpu_value) / (this.num_ranks_per_channel * this.num_dpus_per_rank)
		rank_id := (int(dpu_value) % (this.num_ranks_per_channel * this.num_dpus_per_rank)) / this.num_dpus_per_rank
		dpu_id := int(dpu_value) % this.num_dpus_per_rank

		transfer_command := new(bank.TransferCommand)
		transfer_command.Init(bank.HOST_TO_DEVICE, pointer_value, channel_id, rank_id, dpu_id, base_value+offset_value, size_value)

		this.memory_controller.Push(transfer_command)
		this.push_xfer[transfer_command] = true
	}

	this.SimulateMemory()

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) DpuCopyFrom() {
	dpu_ := this.frame_chain.LastFrame().Stack().Front(4)
	base_ := this.frame_chain.LastFrame().Stack().Front(3)
	offset := this.frame_chain.LastFrame().Stack().Front(2)
	pointer := this.frame_chain.LastFrame().Stack().Front(1)
	size := this.frame_chain.LastFrame().Stack().Front(0)

	dpu_value := this.arena.Pool().Memory().Read(dpu_.Address(), dpu_.Size()).SignedValue()
	offset_value := this.arena.Pool().Memory().Read(offset.Address(), offset.Size()).SignedValue()
	pointer_value := this.arena.Pool().
		Memory().
		Read(pointer.Address(), pointer.Size()).
		SignedValue()
	size_value := this.arena.Pool().Memory().Read(size.Address(), size.Size()).SignedValue()

	if base_.TypeVariable().TypeVariableType() == type_system.STRING {
		base_string := this.DecodeString(
			this.arena.Pool().Memory().Read(base_.Address(), base_.Size()),
		)
		base_string = base_string[1 : len(base_string)-1]

		wram_address, wram_address_found := this.task.Addresses()[base_string]
		if !wram_address_found {
			err_msg := fmt.Sprintf("WRAM address (%s) is not found", base_string)
			err := errors.New(err_msg)
			panic(err)
		}

		byte_stream := this.Dpus()[dpu_value].Dma().
			TransferFromWram(wram_address+offset_value, size_value)
		this.arena.Pool().Memory().Write(pointer_value, size_value, byte_stream)
	} else {
		base_value := this.arena.Pool().Memory().Read(base_.Address(), base_.Size()).SignedValue()

		channel_id := int(dpu_value) / (this.num_ranks_per_channel * this.num_dpus_per_rank)
		rank_id := (int(dpu_value) % (this.num_ranks_per_channel * this.num_dpus_per_rank)) / this.num_dpus_per_rank
		dpu_id := int(dpu_value) % this.num_dpus_per_rank

		transfer_command := new(bank.TransferCommand)
		transfer_command.Init(
			bank.DEVICE_TO_HOST,
			pointer_value,
			channel_id,
			rank_id,
			dpu_id,
			base_value+offset_value,
			size_value,
		)

		this.memory_controller.Push(transfer_command)
		this.push_xfer[transfer_command] = true
	}

	this.SimulateMemory()

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) DpuLaunch() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	for _, dpu_ := range this.Dpus() {
		threads := dpu_.Threads()

		for _, thread := range threads {
			bootstrap := config_loader.IramOffset()
			thread.RegFile().WritePcReg(bootstrap)
		}

		dpu_.Boot()
	}

	thread_pool := new(core.ThreadPool)
	thread_pool.Init()

	for _, dpu_ := range this.Dpus() {
		dpu_cycle_job := new(DpuComputeCycleJob)
		dpu_cycle_job.Init(this.task.SysEnd(), dpu_)

		thread_pool.Enque(dpu_cycle_job)
	}

	thread_pool.Start()

	for _, dpu_ := range this.Dpus() {
		dpu_.Unboot()
	}

	this.frame_chain.LastFrame().Stack().Pop()
	this.frame_chain.LastFrame().Stack().Pop()
}

func (this *VirtualMachine) DpuFree() {
}

func (this *VirtualMachine) Banks() []*bank.Bank {
	return this.memory_controller.Banks()
}

func (this *VirtualMachine) Dpus() []*dpu.Dpu {
	dpus := make([]*dpu.Dpu, 0)
	for _, channel_ := range this.channels {
		dpus = append(dpus, channel_.Dpus()...)
	}
	return dpus
}

func (this *VirtualMachine) TypeCast(
	loperand *stack.StackItem,
	roperand *stack.StackItem,
) type_system.TypeVariableType {
	if loperand.TypeVariable().NumStars() != 0 {
		err := errors.New("loperand's num stars != 0")
		panic(err)
	} else if roperand.TypeVariable().NumStars() != 0 {
		err := errors.New("roperand's num stars != 0")
		panic(err)
	}

	if loperand.TypeVariable().TypeVariableType() == type_system.VOID {
		err := errors.New("loperand's type variable type == void")
		panic(err)
	} else if loperand.TypeVariable().TypeVariableType() == type_system.CHAR {
		if roperand.TypeVariable().TypeVariableType() == type_system.VOID {
			err := errors.New("roperand's type variable type == void")
			panic(err)
		} else if roperand.TypeVariable().TypeVariableType() == type_system.CHAR {
			return type_system.CHAR
		} else if roperand.TypeVariable().TypeVariableType() == type_system.SHORT {
			return type_system.SHORT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.INT {
			return type_system.INT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.LONG {
			return type_system.LONG
		} else if roperand.TypeVariable().TypeVariableType() == type_system.STRUCT {
			err := errors.New("roperand's type variable type == struct")
			panic(err)
		} else {
			err := errors.New("roperand's type variable type is not valid")
			panic(err)
		}
	} else if loperand.TypeVariable().TypeVariableType() == type_system.SHORT {
		if roperand.TypeVariable().TypeVariableType() == type_system.VOID {
			err := errors.New("roperand's type variable type == void")
			panic(err)
		} else if roperand.TypeVariable().TypeVariableType() == type_system.CHAR {
			return type_system.SHORT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.SHORT {
			return type_system.SHORT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.INT {
			return type_system.INT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.LONG {
			return type_system.LONG
		} else if roperand.TypeVariable().TypeVariableType() == type_system.STRUCT {
			err := errors.New("roperand's type variable type == struct")
			panic(err)
		} else {
			err := errors.New("roperand's type variable type is not valid")
			panic(err)
		}
	} else if loperand.TypeVariable().TypeVariableType() == type_system.INT {
		if roperand.TypeVariable().TypeVariableType() == type_system.VOID {
			err := errors.New("roperand's type variable type == void")
			panic(err)
		} else if roperand.TypeVariable().TypeVariableType() == type_system.CHAR {
			return type_system.INT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.SHORT {
			return type_system.INT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.INT {
			return type_system.INT
		} else if roperand.TypeVariable().TypeVariableType() == type_system.LONG {
			return type_system.LONG
		} else if roperand.TypeVariable().TypeVariableType() == type_system.STRUCT {
			err := errors.New("roperand's type variable type == struct")
			panic(err)
		} else {
			err := errors.New("roperand's type variable type is not valid")
			panic(err)
		}
	} else if loperand.TypeVariable().TypeVariableType() == type_system.LONG {
		if roperand.TypeVariable().TypeVariableType() == type_system.VOID {
			err := errors.New("roperand's type variable type == void")
			panic(err)
		} else if roperand.TypeVariable().TypeVariableType() == type_system.CHAR {
			return type_system.LONG
		} else if roperand.TypeVariable().TypeVariableType() == type_system.SHORT {
			return type_system.LONG
		} else if roperand.TypeVariable().TypeVariableType() == type_system.INT {
			return type_system.LONG
		} else if roperand.TypeVariable().TypeVariableType() == type_system.LONG {
			return type_system.LONG
		} else if roperand.TypeVariable().TypeVariableType() == type_system.STRUCT {
			err := errors.New("roperand's type variable type == struct")
			panic(err)
		} else {
			err := errors.New("roperand's type variable type is not valid")
			panic(err)
		}
	} else if loperand.TypeVariable().TypeVariableType() == type_system.STRUCT {
		err := errors.New("loperand's type variable type == struct")
		panic(err)
	} else {
		err := errors.New("loperand's type variable type is not valid")
		panic(err)
	}
}

func (this *VirtualMachine) DecodeString(byte_stream *encoding.ByteStream) string {
	ascii_encoder := new(encoding.AsciiEncoder)
	ascii_encoder.Init()

	return ascii_encoder.Decode(byte_stream)
}

func (this *VirtualMachine) Checkpoint() {
	this.memory_controller.Flush()

	for _, object := range this.arena.Pool().Objects() {
		vm_address := object.Address()
		size := object.Size()

		byte_stream := this.arena.Pool().Memory().Read(object.Address(), object.Size())

		this.memory_controller.VmWrite(vm_address, size, byte_stream)
	}
}

func (this *VirtualMachine) PrepareByteStream(
	pointer int64,
	offset int64,
	size int64,
) *encoding.ByteStream {
	object := this.arena.Pool().Memory().Read(pointer+offset, size)

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < size; i++ {
		byte_stream.Append(object.Get(int(offset + i)))
	}

	return byte_stream
}

func (this *VirtualMachine) SimulateMemory() {
	for len(this.push_xfer) > 0 {
		thread_pool := new(core.ThreadPool)
		thread_pool.Init()

		for _, bank_ := range this.Banks() {
			bank_cycle_job := new(BankCycleJob)
			bank_cycle_job.Init(bank_)

			thread_pool.Enque(bank_cycle_job)
		}

		for _, dpu_ := range this.Dpus() {
			dpu_cycle_job := new(DpuCycleJob)
			dpu_cycle_job.Init(dpu_)

			thread_pool.Enque(dpu_cycle_job)
		}

		thread_pool.Start()

		for _, vm_channel_ := range this.memory_controller.VmChannels() {
			vm_channel_.Cycle()
		}

		for _, channel_ := range this.channels {
			channel_.Cycle()
		}

		this.memory_controller.Cycle()

		if this.memory_controller.CanPop() {
			transfer_command := this.memory_controller.Pop()

			if !transfer_command.IsVmReady() {
				err := errors.New("transfer command is not VM ready")
				panic(err)
			} else if !transfer_command.IsReady() {
				err := errors.New("transfer command is not ready")
				panic(err)
			}

			delete(this.push_xfer, transfer_command)

			if transfer_command.TransferCommandType() == bank.DEVICE_TO_HOST {
				this.arena.Pool().Memory().Write(
					transfer_command.VmAddress(),
					transfer_command.Size(),
					transfer_command.ByteStream(),
				)
			}
		}
	}
}

func (this *VirtualMachine) Dump() {
	file_dumper := new(misc.FileDumper)
	file_dumper.Init(filepath.Join(this.bin_dirpath, "log.txt"))

	lines := make([]string, 0)

	for _, dpu_ := range this.Dpus() {
		lines = append(lines, dpu_.StatFactory().ToLines()...)
		lines = append(lines, dpu_.ThreadScheduler().StatFactory().ToLines()...)
		lines = append(lines, dpu_.Logic().StatFactory().ToLines()...)
		lines = append(lines, dpu_.Logic().CycleRule().StatFactory().ToLines()...)
		lines = append(lines, dpu_.MemoryController().StatFactory().ToLines()...)
		lines = append(lines, dpu_.MemoryController().MemoryScheduler().StatFactory().ToLines()...)
		lines = append(lines, dpu_.MemoryController().RowBuffer().StatFactory().ToLines()...)
	}

	lines = append(lines, this.memory_controller.MemoryScheduler().StatFactory().ToLines()...)

	for _, vm_channel := range this.memory_controller.VmChannels() {
		for _, rank_ := range vm_channel.Ranks() {
			for _, bank_ := range rank_.Banks() {
				lines = append(lines, bank_.StatFactory().ToLines()...)
				lines = append(lines, bank_.RowBuffer().StatFactory().ToLines()...)
			}
		}
	}

	file_dumper.WriteLines(lines)
}

func (this *VirtualMachine) Stringify() string {
	ss := "\n=============== STACK ===============\n"

	for i := 0; i < this.frame_chain.LastFrame().Stack().Length(); i++ {
		stack_item := this.frame_chain.LastFrame().Stack().Front(i)

		byte_stream := this.arena.Pool().Memory().Read(stack_item.Address(), stack_item.Size())

		if byte_stream.Size() <= 8 {
			ss += fmt.Sprintf("%d: %d\n", i, byte_stream.SignedValue())
		} else {
			ss += fmt.Sprintf("%d: XXX\n", i)
		}
	}

	ss += "=====================================\n"

	return ss
}
