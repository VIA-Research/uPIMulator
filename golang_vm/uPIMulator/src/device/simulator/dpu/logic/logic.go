package logic

import (
	"errors"
	"fmt"
	"math"
	"uPIMulator/src/device/abi"
	"uPIMulator/src/device/linker/kernel/instruction"
	"uPIMulator/src/device/linker/kernel/instruction/cc"
	"uPIMulator/src/device/linker/kernel/instruction/reg_descriptor"
	"uPIMulator/src/device/simulator/dpu/sram"
	"uPIMulator/src/misc"
)

type Logic struct {
	channel_id int
	rank_id    int
	dpu_id     int

	num_channels          int
	num_ranks_per_channel int
	num_dpus_per_rank     int

	verbose int

	min_access_granularity int64

	thread_scheduler  *ThreadScheduler
	atomic            *sram.Atomic
	iram              *sram.Iram
	operand_collector *OperandCollector
	dma               *Dma

	scoreboard map[*instruction.Instruction]*Thread

	pipeline   *Pipeline
	cycle_rule *CycleRule

	alu    *Alu
	wait_q *InstructionQ

	stat_factory *misc.StatFactory
}

func (this *Logic) Init(
	channel_id int,
	rank_id int,
	dpu_id int,
	command_line_parser *misc.CommandLineParser,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	} else if dpu_id < 0 {
		err := errors.New("DPU ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.dpu_id = dpu_id

	this.num_channels = int(command_line_parser.IntParameter("num_channels"))
	this.num_ranks_per_channel = int(command_line_parser.IntParameter("num_ranks_per_channel"))
	this.num_dpus_per_rank = int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.min_access_granularity = command_line_parser.IntParameter("min_access_granularity")

	this.verbose = int(command_line_parser.IntParameter("verbose"))

	this.thread_scheduler = nil
	this.atomic = nil
	this.iram = nil
	this.operand_collector = nil
	this.dma = nil

	this.scoreboard = make(map[*instruction.Instruction]*Thread)

	this.pipeline = new(Pipeline)
	this.pipeline.Init(command_line_parser)

	this.cycle_rule = new(CycleRule)
	this.cycle_rule.Init(channel_id, rank_id, dpu_id, command_line_parser)

	this.alu = new(Alu)
	this.alu.Init()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.wait_q = new(InstructionQ)
	this.wait_q.Init(config_loader.MaxNumTasklets(), 0)

	name := fmt.Sprintf("Logic[%d_%d_%d]", channel_id, rank_id, dpu_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *Logic) Fini() {
	this.pipeline.Fini()
	this.cycle_rule.Fini()

	this.alu.Fini()
	this.wait_q.Fini()
}

func (this *Logic) ConnectThreadScheduler(thread_scheduler *ThreadScheduler) {
	if this.thread_scheduler != nil {
		err := errors.New("thread scheduler is already set")
		panic(err)
	}

	this.thread_scheduler = thread_scheduler
}

func (this *Logic) ConnectAtomic(atomic *sram.Atomic) {
	if this.atomic != nil {
		err := errors.New("atomic is already set")
		panic(err)
	}

	this.atomic = atomic
}

func (this *Logic) ConnectIram(iram *sram.Iram) {
	if this.iram != nil {
		err := errors.New("IRAM is already set")
		panic(err)
	}

	this.iram = iram
}

func (this *Logic) ConnectOperandCollector(operand_collector *OperandCollector) {
	if this.operand_collector != nil {
		err := errors.New("operand collector is already set")
		panic(err)
	}

	this.operand_collector = operand_collector
}

func (this *Logic) ConnectDma(dma *Dma) {
	if this.dma != nil {
		err := errors.New("DMA is already set")
		panic(err)
	}

	this.dma = dma
}

func (this *Logic) CycleRule() *CycleRule {
	return this.cycle_rule
}

func (this *Logic) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *Logic) IsEmpty() bool {
	return this.pipeline.IsEmpty() && this.cycle_rule.IsEmpty() && this.wait_q.IsEmpty()
}

func (this *Logic) Cycle() {
	this.ServiceThreadScheduler()
	this.ServicePipeline()
	this.ServiceCycleRule()
	this.ServiceLogic()
	this.ServiceDma()

	this.pipeline.Cycle()
	this.cycle_rule.Cycle()

	this.stat_factory.Increment("logic_cycle", 1)
}

func (this *Logic) ServiceThreadScheduler() {
	if this.pipeline.CanPush() && this.cycle_rule.CanPush() && this.wait_q.CanPush(1) {
		thread := this.thread_scheduler.Schedule()

		if thread != nil {
			pc := thread.RegFile().ReadPcReg()
			instruction_ := this.iram.Read(pc)

			this.scoreboard[instruction_] = thread

			this.pipeline.Push(instruction_)

			if instruction_.Suffix() != instruction.DMA_RRI {
				this.ExecuteInstruction(instruction_)
			} else {
				this.thread_scheduler.Block(thread.ThreadId())
				thread.RegFile().IncrementPcReg()
				this.wait_q.Push(instruction_)
			}

			this.stat_factory.Increment("num_instructions", 1)
		}

		active_tasklets := fmt.Sprintf(
			"active_tasklets_%d",
			this.thread_scheduler.NumIssuableThreads(),
		)
		this.stat_factory.Increment(active_tasklets, 1)
	} else {
		this.stat_factory.Increment("backpressure", 1)
		this.stat_factory.Increment("active_tasklets_0", 1)
	}
}

func (this *Logic) ServicePipeline() {
	if this.pipeline.CanPop() && this.cycle_rule.CanPush() {
		instruction_ := this.pipeline.Pop()
		thread := this.scoreboard[instruction_]

		if instruction_ != nil {
			this.cycle_rule.Push(instruction_, thread)
		}
	}
}

func (this *Logic) ServiceCycleRule() {
	if this.cycle_rule.CanPop() {
		instruction_ := this.cycle_rule.Pop()

		if instruction_.Suffix() != instruction.DMA_RRI {
			delete(this.scoreboard, instruction_)
		} else {
			this.ExecuteInstruction(instruction_)
		}
	}
}

func (this *Logic) ServiceLogic() {
}

func (this *Logic) ServiceDma() {
	if this.dma.CanPop() {
		dma_command := this.dma.Pop()

		has_waked_up := false
		for i := 0; this.wait_q.CanPop(i + 1); i++ {
			instruction_, _ := this.wait_q.Front(i)

			if dma_command.Instruction() == instruction_ {
				thread := this.scoreboard[instruction_]

				this.thread_scheduler.Awake(thread.ThreadId())

				this.wait_q.Remove(i)
				delete(this.scoreboard, instruction_)

				has_waked_up = true
				break
			}
		}

		if !has_waked_up {
			err := errors.New("DMA command has not waked up an instruction")
			panic(err)
		}
	}
}

func (this *Logic) ExecuteInstruction(instruction_ *instruction.Instruction) {
	thread := this.scoreboard[instruction_]

	unique_dpu_id := this.channel_id*this.num_ranks_per_channel*this.num_dpus_per_rank + this.rank_id*this.num_dpus_per_rank + this.dpu_id

	if this.verbose >= 1 {
		fmt.Printf(
			"{%d}[%d] %s\n",
			unique_dpu_id,
			thread.ThreadId(),
			instruction_.Stringify(),
		)
	}

	suffix := instruction_.Suffix()

	if suffix == instruction.RICI {
		this.ExecuteRici(instruction_)
	} else if suffix == instruction.RRI {
		this.ExecuteRri(instruction_)
	} else if suffix == instruction.RRIC {
		this.ExecuteRric(instruction_)
	} else if suffix == instruction.RRICI {
		this.ExecuteRrici(instruction_)
	} else if suffix == instruction.RRIF {
		this.ExecuteRrif(instruction_)
	} else if suffix == instruction.RRR {
		this.ExecuteRrr(instruction_)
	} else if suffix == instruction.RRRC {
		this.ExecuteRrrc(instruction_)
	} else if suffix == instruction.RRRCI {
		this.ExecuteRrrci(instruction_)
	} else if suffix == instruction.ZRI {
		this.ExecuteZri(instruction_)
	} else if suffix == instruction.ZRIC {
		this.ExecuteZric(instruction_)
	} else if suffix == instruction.ZRICI {
		this.ExecuteZrici(instruction_)
	} else if suffix == instruction.ZRIF {
		this.ExecuteZrif(instruction_)
	} else if suffix == instruction.ZRR {
		this.ExecuteZrr(instruction_)
	} else if suffix == instruction.ZRRC {
		this.ExecuteZrrc(instruction_)
	} else if suffix == instruction.ZRRCI {
		this.ExecuteZrrci(instruction_)
	} else if suffix == instruction.S_RRI || suffix == instruction.U_RRI {
		this.ExecuteSRri(instruction_)
	} else if suffix == instruction.S_RRIC || suffix == instruction.U_RRIC {
		this.ExecuteSRric(instruction_)
	} else if suffix == instruction.S_RRICI || suffix == instruction.U_RRICI {
		this.ExecuteSRrici(instruction_)
	} else if suffix == instruction.S_RRIF || suffix == instruction.U_RRIF {
		this.ExecuteSRrif(instruction_)
	} else if suffix == instruction.S_RRR || suffix == instruction.U_RRR {
		this.ExecuteSRrr(instruction_)
	} else if suffix == instruction.S_RRRC || suffix == instruction.U_RRRC {
		this.ExecuteSRrrc(instruction_)
	} else if suffix == instruction.S_RRRCI || suffix == instruction.U_RRRCI {
		this.ExecuteSRrrci(instruction_)
	} else if suffix == instruction.RR {
		this.ExecuteRr(instruction_)
	} else if suffix == instruction.RRC {
		this.ExecuteRrc(instruction_)
	} else if suffix == instruction.RRCI {
		this.ExecuteRrci(instruction_)
	} else if suffix == instruction.ZR {
		this.ExecuteZr(instruction_)
	} else if suffix == instruction.ZRC {
		this.ExecuteZrc(instruction_)
	} else if suffix == instruction.ZRCI {
		this.ExecuteZrci(instruction_)
	} else if suffix == instruction.S_RR || suffix == instruction.U_RR {
		this.ExecuteSRr(instruction_)
	} else if suffix == instruction.S_RRC || suffix == instruction.U_RRC {
		this.ExecuteSRrc(instruction_)
	} else if suffix == instruction.S_RRCI || suffix == instruction.U_RRCI {
		this.ExecuteSRrci(instruction_)
	} else if suffix == instruction.DRDICI {
		this.ExecuteDrdici(instruction_)
	} else if suffix == instruction.RRRI {
		this.ExecuteRrri(instruction_)
	} else if suffix == instruction.RRRICI {
		this.ExecuteRrrici(instruction_)
	} else if suffix == instruction.ZRRI {
		this.ExecuteZrri(instruction_)
	} else if suffix == instruction.ZRRICI {
		this.ExecuteZrrici(instruction_)
	} else if suffix == instruction.S_RRRI || suffix == instruction.U_RRRI {
		this.ExecuteSRrri(instruction_)
	} else if suffix == instruction.S_RRRICI || suffix == instruction.U_RRRICI {
		this.ExecuteSRrrici(instruction_)
	} else if suffix == instruction.RIR {
		this.ExecuteRir(instruction_)
	} else if suffix == instruction.RIRC {
		this.ExecuteRirc(instruction_)
	} else if suffix == instruction.RIRCI {
		this.ExecuteRirci(instruction_)
	} else if suffix == instruction.ZIR {
		this.ExecuteZir(instruction_)
	} else if suffix == instruction.ZIRC {
		this.ExecuteZirc(instruction_)
	} else if suffix == instruction.ZIRCI {
		this.ExecuteZirci(instruction_)
	} else if suffix == instruction.S_RIRC || suffix == instruction.U_RIRC {
		this.ExecuteSRirc(instruction_)
	} else if suffix == instruction.S_RIRCI || suffix == instruction.U_RIRCI {
		this.ExecuteSRirci(instruction_)
	} else if suffix == instruction.R {
		this.ExecuteR(instruction_)
	} else if suffix == instruction.RCI {
		this.ExecuteRci(instruction_)
	} else if suffix == instruction.Z {
		this.ExecuteZ(instruction_)
	} else if suffix == instruction.ZCI {
		this.ExecuteZci(instruction_)
	} else if suffix == instruction.S_R || suffix == instruction.U_R {
		this.ExecuteSR(instruction_)
	} else if suffix == instruction.S_RCI || suffix == instruction.U_RCI {
		this.ExecuteSRci(instruction_)
	} else if suffix == instruction.CI {
		this.ExecuteCi(instruction_)
	} else if suffix == instruction.I {
		this.ExecuteI(instruction_)
	} else if suffix == instruction.DDCI {
		this.ExecuteDdci(instruction_)
	} else if suffix == instruction.ERRI {
		this.ExecuteErri(instruction_)
	} else if suffix == instruction.S_ERRI || suffix == instruction.U_ERRI {
		this.ExecuteSErri(instruction_)
	} else if suffix == instruction.EDRI {
		this.ExecuteEdri(instruction_)
	} else if suffix == instruction.ERII {
		this.ExecuteErii(instruction_)
	} else if suffix == instruction.ERIR {
		this.ExecuteErir(instruction_)
	} else if suffix == instruction.ERID {
		this.ExecuteErid(instruction_)
	} else if suffix == instruction.DMA_RRI {
		this.ExecuteDmaRri(instruction_)
	} else {
		err := errors.New("suffix is not valid")
		panic(err)
	}

	if this.verbose >= 2 {
		fmt.Println(this.PrintRegFile(thread))
	}
}

func (this *Logic) ExecuteRici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RiciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RICI {
		err := errors.New("suffix is not RICI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_acquire_rici_op_code := instruction_.AcquireRiciOpCodes()[op_code]; is_acquire_rici_op_code {
		this.ExecuteAcquireRici(instruction_)
	} else if _, is_release_rici_op_code := instruction_.ReleaseRiciOpCodes()[op_code]; is_release_rici_op_code {
		this.ExecuteReleaseRici(instruction_)
	} else if _, is_boot_rici_op_code := instruction_.BootRiciOpCodes()[op_code]; is_boot_rici_op_code {
		this.ExecuteBootRici(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAcquireRici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AcquireRiciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid acquire RICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RICI {
		err := errors.New("suffix is not RICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.UNSIGNED)
	imm := instruction_.Imm().Value()

	atomic_address := this.alu.AtomicAddressHash(ra, imm)

	can_acquire := this.atomic.CanAcquire(atomic_address)
	if can_acquire {
		this.atomic.Acquire(atomic_address, thread.ThreadId())
	}

	thread.RegFile().ClearConditions()
	if can_acquire {
		this.SetAcquireCc(instruction_, 0)
	} else {
		this.SetAcquireCc(instruction_, 1)
	}

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	if can_acquire {
		this.SetFlags(instruction_, 0, false)
	} else {
		this.SetFlags(instruction_, 1, false)
	}
}

func (this *Logic) ExecuteReleaseRici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.ReleaseRiciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid release RICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RICI {
		err := errors.New("suffix is not RICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.UNSIGNED)
	imm := instruction_.Imm().Value()

	atomic_address := this.alu.AtomicAddressHash(ra, imm)

	can_release := this.atomic.CanRelease(atomic_address, thread.ThreadId())
	if can_release {
		this.atomic.Release(atomic_address, thread.ThreadId())
	}

	thread.RegFile().ClearConditions()
	if can_release {
		this.SetAcquireCc(instruction_, 0)
	} else {
		this.SetAcquireCc(instruction_, 1)
	}

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	if can_release {
		this.SetFlags(instruction_, 0, false)
	} else {
		this.SetFlags(instruction_, 1, false)
	}
}

func (this *Logic) ExecuteBootRici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.BootRiciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid boot RICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RICI {
		err := errors.New("suffix is not RICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.UNSIGNED)
	imm := instruction_.Imm().Value()

	thread_id := int(this.alu.AtomicAddressHash(ra, imm))

	thread.RegFile().ClearConditions()

	op_code := instruction_.OpCode()
	if op_code == instruction.BOOT {
		can_boot := this.thread_scheduler.Boot(thread_id)
		if can_boot {
			this.SetBootCc(instruction_, ra, 0)
			this.SetFlags(instruction_, 0, false)
		} else {
			this.SetBootCc(instruction_, ra, 1)
			this.SetFlags(instruction_, 1, false)
		}
	} else if op_code == instruction.RESUME {
		can_resume := this.thread_scheduler.Awake(thread_id)
		if can_resume {
			this.SetBootCc(instruction_, ra, 0)
			this.SetFlags(instruction_, 0, false)
		} else {
			this.SetBootCc(instruction_, ra, 1)
			this.SetFlags(instruction_, 1, false)
		}
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}
}

func (this *Logic) ExecuteRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRI {
		err := errors.New("suffix is not RRI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rri_op_code := instruction_.AddRriOpCodes()[op_code]; is_add_rri_op_code {
		this.ExecuteAddRri(instruction_)
	} else if _, is_asr_rri_op_code := instruction_.AsrRriOpCodes()[op_code]; is_asr_rri_op_code {
		this.ExecuteAsrRri(instruction_)
	} else if _, is_call_rri_op_code := instruction_.CallRriOpCodes()[op_code]; is_call_rri_op_code {
		this.ExecuteCallRri(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRI {
		err := errors.New("suffix is not RRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WriteGpReg(instruction_.Rc(), result)
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAsrRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRI {
		err := errors.New("suffix is not RRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WriteGpReg(instruction_.Rc(), result)
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteCallRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.CallRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid call RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRI {
		err := errors.New("suffix is not RRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	iram_data_size := int64(config_loader.IramDataWidth() / 8)

	if imm == 0 {
		result, carry, _ = this.alu.Add(ra, imm)
	} else {
		result, carry, _ = this.alu.Add(ra*iram_data_size, imm)
	}

	thread.RegFile().ClearConditions()

	pc := thread.RegFile().ReadPcReg()
	thread.RegFile().WriteGpReg(instruction_.Rc(), pc+iram_data_size)

	thread.RegFile().WritePcReg(result)

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRIC {
		err := errors.New("suffix is not RRIC")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rric_op_code := instruction_.AddRricOpCodes()[op_code]; is_add_rric_op_code {
		this.ExecuteAddRric(instruction_)
	} else if _, is_asrc_rric_op_code := instruction_.AsrRricOpCodes()[op_code]; is_asrc_rric_op_code {
		this.ExecuteAsrRric(instruction_)
	} else if _, is_sub_rric_op_code := instruction_.SubRricOpCodes()[op_code]; is_sub_rric_op_code {
		this.ExecuteSubRric(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddRric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRIC {
		err := errors.New("suffix is not RRIC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
		carry = false
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAsrRric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRIC {
		err := errors.New("suffix is not RRIC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteSubRric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SubRricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sub RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRIC {
		err := errors.New("suffix is not RRIC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result, carry, overflow = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetExtSubSetCc(instruction_, ra, imm, result, carry, overflow)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRICI {
		err := errors.New("suffix is not RRICI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rrici_op_code := instruction_.AddRriciOpCodes()[op_code]; is_add_rrici_op_code {
		this.ExecuteAddRrici(instruction_)
	} else if _, is_and_rrici_op_code := instruction_.AndRriciOpCodes()[op_code]; is_and_rrici_op_code {
		this.ExecuteAndRrici(instruction_)
	} else if _, is_asr_rrici_op_code := instruction_.AsrRriciOpCodes()[op_code]; is_asr_rrici_op_code {
		this.ExecuteAsrRrici(instruction_)
	} else if _, is_sub_rrici_op_code := instruction_.SubRriciOpCodes()[op_code]; is_sub_rrici_op_code {
		this.ExecuteSubRrici(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRICI {
		err := errors.New("suffix is not RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, overflow = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Addc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetAddNzCc(instruction_, ra, result, carry, overflow)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAndRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AndRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid and RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRICI {
		err := errors.New("suffix is not RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteAsrRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRICI {
		err := errors.New("suffix is not RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetImmShiftNzCc(instruction_, ra, result)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteSubRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SubRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sub RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRICI {
		err := errors.New("suffix is not RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubNzCc(instruction_, ra, imm, result, carry, overflow)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRrif(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrifOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRIF {
		err := errors.New("suffix is not RRIF")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
		carry = false
	} else if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(ra, imm, carry_flag)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WriteGpReg(instruction_.Rc(), result)
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRrr(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrrOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRR {
		err := errors.New("suffix is not RRR")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, rb)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, rb, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, rb)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, rb)
		carry = false
	} else if op_code == instruction.ASR {
		result = this.alu.Asr(ra, rb)
		carry = false
	} else if op_code == instruction.CMPB4 {
		result = this.alu.Cmpb4(ra, rb)
		carry = false
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, rb)
		carry = false
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, rb)
		carry = false
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, rb)
		carry = false
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SH {
		result = this.alu.MulShSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SL {
		result = this.alu.MulShSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UH {
		result = this.alu.MulShUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UL {
		result = this.alu.MulShUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SH {
		result = this.alu.MulSlSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SL {
		result = this.alu.MulSlSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UH {
		result = this.alu.MulSlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UL {
		result = this.alu.MulSlUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UH {
		result = this.alu.MulUhUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UL {
		result = this.alu.MulUhUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UH {
		result = this.alu.MulUlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UL {
		result = this.alu.MulUlUl(ra, rb)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, rb)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, rb)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, rb)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, rb)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, rb)
		carry = false
	} else if op_code == instruction.RSUB {
		result, carry, _ = this.alu.Sub(rb, ra)
	} else if op_code == instruction.RSUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(rb, ra, carry_flag)
	} else if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(ra, rb)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(ra, rb, carry_flag)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, rb)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, rb)
		carry = false
	} else if op_code == instruction.CALL {
		result, carry, _ = this.alu.Add(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()

	if op_code != instruction.CALL {
		thread.RegFile().WriteGpReg(instruction_.Rc(), result)
		thread.RegFile().IncrementPcReg()
	} else {
		config_loader := new(misc.ConfigLoader)
		config_loader.Init()

		pc := thread.RegFile().ReadPcReg()
		iram_data_size := int64(config_loader.IramDataWidth() / 8)

		thread.RegFile().WriteGpReg(instruction_.Rc(), pc+iram_data_size)
		thread.RegFile().WritePcReg(result)
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRC {
		err := errors.New("suffix is not RRRC")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rrrc_op_code := instruction_.AddRrrcOpCodes()[op_code]; is_add_rrrc_op_code {
		this.ExecuteAddRrrc(instruction_)
	} else if _, is_rsub_rrrc_op_code := instruction_.RsubRrrcOpCodes()[op_code]; is_rsub_rrrc_op_code {
		this.ExecuteRsubRrrc(instruction_)
	} else if _, is_sub_rrrc_op_code := instruction_.SubRrrcOpCodes()[op_code]; is_sub_rrrc_op_code {
		this.ExecuteSubRrrc(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddRrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRC {
		err := errors.New("suffix is not RRRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, rb)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, rb, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, rb)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, rb)
		carry = false
	} else if op_code == instruction.ASR {
		result = this.alu.Asr(ra, rb)
		carry = false
	} else if op_code == instruction.CMPB4 {
		result = this.alu.Cmpb4(ra, rb)
		carry = false
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, rb)
		carry = false
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, rb)
		carry = false
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, rb)
		carry = false
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SH {
		result = this.alu.MulShSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SL {
		result = this.alu.MulShSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UH {
		result = this.alu.MulShUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UL {
		result = this.alu.MulShUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SH {
		result = this.alu.MulSlSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SL {
		result = this.alu.MulSlSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UH {
		result = this.alu.MulSlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UL {
		result = this.alu.MulSlUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UH {
		result = this.alu.MulUhUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UL {
		result = this.alu.MulUhUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UH {
		result = this.alu.MulUlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UL {
		result = this.alu.MulUlUl(ra, rb)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, rb)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, rb)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, rb)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, rb)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, rb)
		carry = false
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, rb)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, rb)
		carry = false
	} else if op_code == instruction.CALL {
		result, carry, _ = this.alu.Add(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRsubRrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RsubRrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid rsub RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRC {
		err := errors.New("suffix is not RRRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.RSUB {
		result, carry, _ = this.alu.Sub(rb, ra)
	} else if op_code == instruction.RSUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(rb, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubSetCc(instruction_, ra, rb, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSubRrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SubRrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sub RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRC {
		err := errors.New("suffix is not RRRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(ra, rb)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, rb, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetExtSubSetCc(instruction_, ra, rb, result, carry, overflow)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rrrci_op_code := instruction_.AddRrrciOpCodes()[op_code]; is_add_rrrci_op_code {
		this.ExecuteAddRrrci(instruction_)
	} else if _, is_and_rrrci_op_code := instruction_.AndRrrciOpCodes()[op_code]; is_and_rrrci_op_code {
		this.ExecuteAndRrrci(instruction_)
	} else if _, is_asr_rrrci_op_code := instruction_.AsrRrrciOpCodes()[op_code]; is_asr_rrrci_op_code {
		this.ExecuteAsrRrrci(instruction_)
	} else if _, is_mul_rrrci_op_code := instruction_.MulRrrciOpCodes()[op_code]; is_mul_rrrci_op_code {
		this.ExecuteMulRrrci(instruction_)
	} else if _, is_rsub_rrrci_op_code := instruction_.RsubRrrciOpCodes()[op_code]; is_rsub_rrrci_op_code {
		this.ExecuteRsubRrrci(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddRrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, rb)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, rb, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetAddNzCc(instruction_, ra, result, carry, overflow)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAndRrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AndRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid and RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.AND {
		result = this.alu.And(ra, rb)
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, rb)
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, rb)
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, rb)
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, rb)
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, rb)
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, rb)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, rb)
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteAsrRrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, rb)
	} else if op_code == instruction.CMPB4 {
		result = this.alu.Cmpb4(ra, rb)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, rb)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, rb)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, rb)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, rb)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, rb)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, rb)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, rb)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, rb)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, rb)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteMulRrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.MulRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid mul RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.MUL_SH_SH {
		result = this.alu.MulShSh(ra, rb)
	} else if op_code == instruction.MUL_SH_SL {
		result = this.alu.MulShSl(ra, rb)
	} else if op_code == instruction.MUL_SH_UH {
		result = this.alu.MulShUh(ra, rb)
	} else if op_code == instruction.MUL_SH_UL {
		result = this.alu.MulShUl(ra, rb)
	} else if op_code == instruction.MUL_SL_SH {
		result = this.alu.MulSlSh(ra, rb)
	} else if op_code == instruction.MUL_SL_SL {
		result = this.alu.MulSlSl(ra, rb)
	} else if op_code == instruction.MUL_SL_UH {
		result = this.alu.MulSlUh(ra, rb)
	} else if op_code == instruction.MUL_SL_UL {
		result = this.alu.MulSlUl(ra, rb)
	} else if op_code == instruction.MUL_UH_UH {
		result = this.alu.MulUhUh(ra, rb)
	} else if op_code == instruction.MUL_UH_UL {
		result = this.alu.MulUhUl(ra, rb)
	} else if op_code == instruction.MUL_UL_UH {
		result = this.alu.MulUlUh(ra, rb)
	} else if op_code == instruction.MUL_UL_UL {
		result = this.alu.MulUlUl(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetMulNzCc(instruction_, ra, result)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteRsubRrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RsubRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid rsub RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.RSUB {
		result, carry, overflow = this.alu.Sub(rb, ra)
	} else if op_code == instruction.RSUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(rb, ra, carry_flag)
	} else if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(ra, rb)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, rb, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubNzCc(instruction_, ra, rb, result, carry, overflow)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRI {
		err := errors.New("suffix is not RRI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rri_op_code := instruction_.AddRriOpCodes()[op_code]; is_add_rri_op_code {
		this.ExecuteAddZri(instruction_)
	} else if _, is_asr_rri_op_code := instruction_.AsrRriOpCodes()[op_code]; is_asr_rri_op_code {
		this.ExecuteAsrZri(instruction_)
	} else if _, is_call_rri_op_code := instruction_.CallRriOpCodes()[op_code]; is_call_rri_op_code {
		this.ExecuteCallZri(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddZri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRI {
		err := errors.New("suffix is not ZRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAsrZri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRI {
		err := errors.New("suffix is not ZRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteCallZri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.CallRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid call RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRI {
		err := errors.New("suffix is not ZRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	iram_data_size := int64(config_loader.IramDataWidth() / 8)

	if imm == 0 {
		result, carry, _ = this.alu.Add(ra, imm)
	} else {
		result, carry, _ = this.alu.Add(ra*iram_data_size, imm)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WritePcReg(result)

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRIC {
		err := errors.New("suffix is not ZRIC")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rric_op_code := instruction_.AddRricOpCodes()[op_code]; is_add_rric_op_code {
		this.ExecuteAddZric(instruction_)
	} else if _, is_asrc_rric_op_code := instruction_.AsrRricOpCodes()[op_code]; is_asrc_rric_op_code {
		this.ExecuteAsrZric(instruction_)
	} else if _, is_sub_rric_op_code := instruction_.SubRricOpCodes()[op_code]; is_sub_rric_op_code {
		this.ExecuteSubZric(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddZric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRIC {
		err := errors.New("suffix is not ZRIC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
		carry = false
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAsrZric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRIC {
		err := errors.New("suffix is not ZRIC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteSubZric(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SubRricOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sub RRIC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRIC {
		err := errors.New("suffix is not ZRIC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result, carry, overflow = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetExtSubSetCc(instruction_, ra, imm, result, carry, overflow)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRICI {
		err := errors.New("suffix is not ZRICI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rrici_op_code := instruction_.AddRriciOpCodes()[op_code]; is_add_rrici_op_code {
		this.ExecuteAddZrici(instruction_)
	} else if _, is_and_rrici_op_code := instruction_.AndRriciOpCodes()[op_code]; is_and_rrici_op_code {
		this.ExecuteAndZrici(instruction_)
	} else if _, is_asr_rrici_op_code := instruction_.AsrRriciOpCodes()[op_code]; is_asr_rrici_op_code {
		this.ExecuteAsrZrici(instruction_)
	} else if _, is_sub_rrici_op_code := instruction_.SubRriciOpCodes()[op_code]; is_sub_rrici_op_code {
		this.ExecuteSubZrici(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddZrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRICI {
		err := errors.New("suffix is not ZRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, overflow = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Addc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetAddNzCc(instruction_, ra, result, carry, overflow)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAndZrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AndRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid and RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRICI {
		err := errors.New("suffix is not ZRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteAsrZrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRICI {
		err := errors.New("suffix is not ZRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetImmShiftNzCc(instruction_, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteSubZrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SubRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sub RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRICI {
		err := errors.New("suffix is not ZRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubNzCc(instruction_, ra, imm, result, carry, overflow)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZrif(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrifOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRIF {
		err := errors.New("suffix is not ZRIF")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
		carry = false
	} else if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(ra, imm, carry_flag)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZrr(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrrOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRR {
		err := errors.New("suffix is not ZRR")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, rb)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, rb, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, rb)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, rb)
		carry = false
	} else if op_code == instruction.ASR {
		result = this.alu.Asr(ra, rb)
		carry = false
	} else if op_code == instruction.CMPB4 {
		result = this.alu.Cmpb4(ra, rb)
		carry = false
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, rb)
		carry = false
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, rb)
		carry = false
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, rb)
		carry = false
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SH {
		result = this.alu.MulShSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SL {
		result = this.alu.MulShSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UH {
		result = this.alu.MulShUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UL {
		result = this.alu.MulShUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SH {
		result = this.alu.MulSlSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SL {
		result = this.alu.MulSlSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UH {
		result = this.alu.MulSlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UL {
		result = this.alu.MulSlUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UH {
		result = this.alu.MulUhUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UL {
		result = this.alu.MulUhUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UH {
		result = this.alu.MulUlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UL {
		result = this.alu.MulUlUl(ra, rb)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, rb)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, rb)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, rb)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, rb)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, rb)
		carry = false
	} else if op_code == instruction.RSUB {
		result, carry, _ = this.alu.Sub(rb, ra)
	} else if op_code == instruction.RSUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(rb, ra, carry_flag)
	} else if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(ra, rb)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(ra, rb, carry_flag)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, rb)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, rb)
		carry = false
	} else if op_code == instruction.CALL {
		result, carry, _ = this.alu.Add(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()

	if op_code != instruction.CALL {
		thread.RegFile().IncrementPcReg()
	} else {
		thread.RegFile().WritePcReg(result)
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRC {
		err := errors.New("suffix is not ZRRC")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rrrc_op_code := instruction_.AddRrrcOpCodes()[op_code]; is_add_rrrc_op_code {
		this.ExecuteAddZrrc(instruction_)
	} else if _, is_rsub_rrrc_op_code := instruction_.RsubRrrcOpCodes()[op_code]; is_rsub_rrrc_op_code {
		this.ExecuteRsubZrrc(instruction_)
	} else if _, is_sub_rrrc_op_code := instruction_.SubRrrcOpCodes()[op_code]; is_sub_rrrc_op_code {
		this.ExecuteSubZrrc(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddZrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRC {
		err := errors.New("suffix is not ZRRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, rb)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, rb, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, rb)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, rb)
		carry = false
	} else if op_code == instruction.ASR {
		result = this.alu.Asr(ra, rb)
		carry = false
	} else if op_code == instruction.CMPB4 {
		result = this.alu.Cmpb4(ra, rb)
		carry = false
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, rb)
		carry = false
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, rb)
		carry = false
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, rb)
		carry = false
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, rb)
		carry = false
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, rb)
		carry = false
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, rb)
		carry = false
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SH {
		result = this.alu.MulShSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_SL {
		result = this.alu.MulShSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UH {
		result = this.alu.MulShUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SH_UL {
		result = this.alu.MulShUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SH {
		result = this.alu.MulSlSh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_SL {
		result = this.alu.MulSlSl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UH {
		result = this.alu.MulSlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_SL_UL {
		result = this.alu.MulSlUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UH {
		result = this.alu.MulUhUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UH_UL {
		result = this.alu.MulUhUl(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UH {
		result = this.alu.MulUlUh(ra, rb)
		carry = false
	} else if op_code == instruction.MUL_UL_UL {
		result = this.alu.MulUlUl(ra, rb)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, rb)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, rb)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, rb)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, rb)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, rb)
		carry = false
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, rb)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, rb)
		carry = false
	} else if op_code == instruction.CALL {
		result, carry, _ = this.alu.Add(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRsubZrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RsubRrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid rsub RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRC {
		err := errors.New("suffix is not ZRRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.RSUB {
		result, carry, _ = this.alu.Sub(rb, ra)
	} else if op_code == instruction.RSUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(rb, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubSetCc(instruction_, ra, rb, result)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSubZrrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SubRrrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sub RRRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRC {
		err := errors.New("suffix is not ZRRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(ra, rb)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, rb, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetExtSubSetCc(instruction_, ra, rb, result, carry, overflow)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rrrci_op_code := instruction_.AddRrrciOpCodes()[op_code]; is_add_rrrci_op_code {
		this.ExecuteAddZrrci(instruction_)
	} else if _, is_and_rrrci_op_code := instruction_.AndRrrciOpCodes()[op_code]; is_and_rrrci_op_code {
		this.ExecuteAndZrrci(instruction_)
	} else if _, is_asr_rrrci_op_code := instruction_.AsrRrrciOpCodes()[op_code]; is_asr_rrrci_op_code {
		this.ExecuteAsrZrrci(instruction_)
	} else if _, is_mul_rrrci_op_code := instruction_.MulRrrciOpCodes()[op_code]; is_mul_rrrci_op_code {
		this.ExecuteMulZrrci(instruction_)
	} else if _, is_rsub_rrrci_op_code := instruction_.RsubRrrciOpCodes()[op_code]; is_rsub_rrrci_op_code {
		this.ExecuteRsubZrrci(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddZrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, rb)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, rb, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetAddNzCc(instruction_, ra, result, carry, overflow)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAndZrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AndRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid and RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.AND {
		result = this.alu.And(ra, rb)
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, rb)
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, rb)
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, rb)
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, rb)
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, rb)
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, rb)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, rb)
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteAsrZrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, rb)
	} else if op_code == instruction.CMPB4 {
		result = this.alu.Cmpb4(ra, rb)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, rb)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, rb)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, rb)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, rb)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, rb)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, rb)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, rb)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, rb)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, rb)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteMulZrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.MulRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid mul RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.MUL_SH_SH {
		result = this.alu.MulShSh(ra, rb)
	} else if op_code == instruction.MUL_SH_SL {
		result = this.alu.MulShSl(ra, rb)
	} else if op_code == instruction.MUL_SH_UH {
		result = this.alu.MulShUh(ra, rb)
	} else if op_code == instruction.MUL_SH_UL {
		result = this.alu.MulShUl(ra, rb)
	} else if op_code == instruction.MUL_SL_SH {
		result = this.alu.MulSlSh(ra, rb)
	} else if op_code == instruction.MUL_SL_SL {
		result = this.alu.MulSlSl(ra, rb)
	} else if op_code == instruction.MUL_SL_UH {
		result = this.alu.MulSlUh(ra, rb)
	} else if op_code == instruction.MUL_SL_UL {
		result = this.alu.MulSlUl(ra, rb)
	} else if op_code == instruction.MUL_UH_UH {
		result = this.alu.MulUhUh(ra, rb)
	} else if op_code == instruction.MUL_UH_UL {
		result = this.alu.MulUhUl(ra, rb)
	} else if op_code == instruction.MUL_UL_UH {
		result = this.alu.MulUlUh(ra, rb)
	} else if op_code == instruction.MUL_UL_UL {
		result = this.alu.MulUlUl(ra, rb)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetMulNzCc(instruction_, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteRsubZrrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RsubRrrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid rsub RRRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.RSUB {
		result, carry, overflow = this.alu.Sub(rb, ra)
	} else if op_code == instruction.RSUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(rb, ra, carry_flag)
	} else if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(ra, rb)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, rb, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubNzCc(instruction_, ra, rb, result, carry, overflow)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRI && instruction_.Suffix() != instruction.U_RRI {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rri_op_code := instruction_.AddRriOpCodes()[op_code]; is_add_rri_op_code {
		this.ExecuteAddSRri(instruction_)
	} else if _, is_asr_rri_op_code := instruction_.AsrRriOpCodes()[op_code]; is_asr_rri_op_code {
		this.ExecuteAsrSRri(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddSRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRI && instruction_.Suffix() != instruction.U_RRI {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()

	var even int64
	var odd int64
	if instruction_.Suffix() == instruction.S_RRI {
		even, odd = this.alu.SignedExtension(result)
	} else if instruction_.Suffix() == instruction.U_RRI {
		even, odd = this.alu.UnsignedExtension(result)
	} else {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAsrSRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRI && instruction_.Suffix() != instruction.U_RRI {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()

	var even int64
	var odd int64
	if instruction_.Suffix() == instruction.S_RRI {
		even, odd = this.alu.SignedExtension(result)
	} else if instruction_.Suffix() == instruction.U_RRI {
		even, odd = this.alu.UnsignedExtension(result)
	} else {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSRric(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRric is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRICI && instruction_.Suffix() != instruction.U_RRICI {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_add_rrici_op_code := instruction_.AddRriciOpCodes()[op_code]; is_add_rrici_op_code {
		this.ExecuteAddSRrici(instruction_)
	} else if _, is_and_rrici_op_code := instruction_.AndRriciOpCodes()[op_code]; is_and_rrici_op_code {
		this.ExecuteAndSRrici(instruction_)
	} else if _, is_asr_rrici_op_code := instruction_.AsrRriciOpCodes()[op_code]; is_asr_rrici_op_code {
		this.ExecuteAsrSRrici(instruction_)
	} else if _, is_sub_rrici_op_code := instruction_.SubRriciOpCodes()[op_code]; is_sub_rrici_op_code {
		this.ExecuteSubSRrici(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteAddSRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AddRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid add RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRICI && instruction_.Suffix() != instruction.U_RRICI {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, overflow = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Addc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetAddNzCc(instruction_, ra, result, carry, overflow)

	var even int64
	var odd int64
	if instruction_.Suffix() == instruction.S_RRICI {
		even, odd = this.alu.SignedExtension(result)
	} else if instruction_.Suffix() == instruction.U_RRICI {
		even, odd = this.alu.UnsignedExtension(result)
	} else {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteAndSRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AndRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid and RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRICI && instruction_.Suffix() != instruction.U_RRICI {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	var even int64
	var odd int64
	if instruction_.Suffix() == instruction.S_RRICI {
		even, odd = this.alu.SignedExtension(result)
	} else if instruction_.Suffix() == instruction.U_RRICI {
		even, odd = this.alu.UnsignedExtension(result)
	} else {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteAsrSRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.AsrRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid asr RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRICI && instruction_.Suffix() != instruction.U_RRICI {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result = this.alu.Asr(ra, imm)
	} else if op_code == instruction.LSL {
		result = this.alu.Lsl(ra, imm)
	} else if op_code == instruction.LSL1 {
		result = this.alu.Lsl1(ra, imm)
	} else if op_code == instruction.LSL1X {
		result = this.alu.Lsl1x(ra, imm)
	} else if op_code == instruction.LSLX {
		result = this.alu.Lslx(ra, imm)
	} else if op_code == instruction.LSR {
		result = this.alu.Lsr(ra, imm)
	} else if op_code == instruction.LSR1 {
		result = this.alu.Lsr1(ra, imm)
	} else if op_code == instruction.LSR1X {
		result = this.alu.Lsr1x(ra, imm)
	} else if op_code == instruction.LSRX {
		result = this.alu.Lsrx(ra, imm)
	} else if op_code == instruction.ROL {
		result = this.alu.Rol(ra, imm)
	} else if op_code == instruction.ROR {
		result = this.alu.Ror(ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetImmShiftNzCc(instruction_, ra, result)

	var even int64
	var odd int64
	if instruction_.Suffix() == instruction.S_RRICI {
		even, odd = this.alu.SignedExtension(result)
	} else if instruction_.Suffix() == instruction.U_RRICI {
		even, odd = this.alu.UnsignedExtension(result)
	} else {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteSubSRrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SubRriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sub RRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRICI && instruction_.Suffix() != instruction.U_RRICI {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ASR {
		result, carry, overflow = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(ra, imm, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubNzCc(instruction_, ra, imm, result, carry, overflow)

	var even int64
	var odd int64
	if instruction_.Suffix() == instruction.S_RRICI {
		even, odd = this.alu.SignedExtension(result)
	} else if instruction_.Suffix() == instruction.U_RRICI {
		even, odd = this.alu.UnsignedExtension(result)
	} else {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSRrif(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrifOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.S_RRIF && instruction_.Suffix() != instruction.U_RRIF {
		err := errors.New("suffix is not S_RRIF nor U_RRIF")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.ADD {
		result, carry, _ = this.alu.Add(ra, imm)
	} else if op_code == instruction.ADDC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Addc(ra, imm, carry_flag)
	} else if op_code == instruction.AND {
		result = this.alu.And(ra, imm)
		carry = false
	} else if op_code == instruction.ANDN {
		result = this.alu.Andn(ra, imm)
		carry = false
	} else if op_code == instruction.NAND {
		result = this.alu.Nand(ra, imm)
		carry = false
	} else if op_code == instruction.NOR {
		result = this.alu.Nor(ra, imm)
		carry = false
	} else if op_code == instruction.NXOR {
		result = this.alu.Nxor(ra, imm)
		carry = false
	} else if op_code == instruction.OR {
		result = this.alu.Or(ra, imm)
		carry = false
	} else if op_code == instruction.ORN {
		result = this.alu.Orn(ra, imm)
		carry = false
	} else if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(ra, imm)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(ra, imm, carry_flag)
	} else if op_code == instruction.XOR {
		result = this.alu.Xor(ra, imm)
		carry = false
	} else if op_code == instruction.HASH {
		result = this.alu.Hash(ra, imm)
		carry = false
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()

	var even int64
	var odd int64
	if instruction_.Suffix() == instruction.S_RRIF {
		even, odd = this.alu.SignedExtension(result)
	} else if instruction_.Suffix() == instruction.U_RRIF {
		even, odd = this.alu.UnsignedExtension(result)
	} else {
		err := errors.New("suffix is not S_RRIF nor U_RRIF")
		panic(err)
	}

	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSRrr(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRrr is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRrrc(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRrrc is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRrrci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRrrci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteRr(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RR {
		err := errors.New("suffix is not RR")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.CAO {
		result = this.alu.Cao(ra)
	} else if op_code == instruction.CLO {
		result = this.alu.Clo(ra)
	} else if op_code == instruction.CLS {
		result = this.alu.Cls(ra)
	} else if op_code == instruction.CLZ {
		result = this.alu.Clz(ra)
	} else if op_code == instruction.EXTSB {
		result = this.alu.Extsb(ra)
	} else if op_code == instruction.EXTSH {
		result = this.alu.Extsh(ra)
	} else if op_code == instruction.EXTUB {
		result = this.alu.Extub(ra)
	} else if op_code == instruction.EXTUH {
		result = this.alu.Extuh(ra)
	} else if op_code == instruction.SATS {
		result = this.alu.Sats(ra)
	} else if op_code == instruction.TIME_CFG {
		err := errors.New("TimeCfg is not yet implemented")
		panic(err)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WriteGpReg(instruction_.Rc(), result)
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteRrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRC {
		err := errors.New("suffix is not RRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.CAO {
		result = this.alu.Cao(ra)
	} else if op_code == instruction.CLO {
		result = this.alu.Clo(ra)
	} else if op_code == instruction.CLS {
		result = this.alu.Cls(ra)
	} else if op_code == instruction.CLZ {
		result = this.alu.Clz(ra)
	} else if op_code == instruction.EXTSB {
		result = this.alu.Extsb(ra)
	} else if op_code == instruction.EXTSH {
		result = this.alu.Extsh(ra)
	} else if op_code == instruction.EXTUB {
		result = this.alu.Extub(ra)
	} else if op_code == instruction.EXTUH {
		result = this.alu.Extuh(ra)
	} else if op_code == instruction.SATS {
		result = this.alu.Sats(ra)
	} else if op_code == instruction.TIME_CFG {
		err := errors.New("TimeCfg is not yet implemented")
		panic(err)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteRrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRCI {
		err := errors.New("suffix is not RRCI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_cao_rrci_op_code := instruction_.CaoRrciOpCodes()[op_code]; is_cao_rrci_op_code {
		this.ExecuteCaoRrci(instruction_)
	} else if _, is_extsb_rrci_op_code := instruction_.ExtsbRrciOpCodes()[op_code]; is_extsb_rrci_op_code {
		this.ExecuteExtsbRrci(instruction_)
	} else if _, is_time_cfg_rrci_op_code := instruction_.TimeCfgRrciOpCodes()[op_code]; is_time_cfg_rrci_op_code {
		this.ExecuteTimeCfgRrci(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteCaoRrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.CaoRrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid cao RRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRCI {
		err := errors.New("suffix is not RRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.CAO {
		result = this.alu.Cao(ra)
	} else if op_code == instruction.CLO {
		result = this.alu.Clo(ra)
	} else if op_code == instruction.CLS {
		result = this.alu.Cls(ra)
	} else if op_code == instruction.CLZ {
		result = this.alu.Clz(ra)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetCountNzCc(instruction_, ra, result)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteExtsbRrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.ExtsbRrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid extsb RRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRCI {
		err := errors.New("suffix is not RRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.EXTSB {
		result = this.alu.Extsb(ra)
	} else if op_code == instruction.EXTSH {
		result = this.alu.Extsh(ra)
	} else if op_code == instruction.EXTUB {
		result = this.alu.Extub(ra)
	} else if op_code == instruction.EXTUH {
		result = this.alu.Extuh(ra)
	} else if op_code == instruction.SATS {
		result = this.alu.Sats(ra)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteTimeCfgRrci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteTimeCfgRrci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteZr(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZR {
		err := errors.New("suffix is not ZR")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.CAO {
		result = this.alu.Cao(ra)
	} else if op_code == instruction.CLO {
		result = this.alu.Clo(ra)
	} else if op_code == instruction.CLS {
		result = this.alu.Cls(ra)
	} else if op_code == instruction.CLZ {
		result = this.alu.Clz(ra)
	} else if op_code == instruction.EXTSB {
		result = this.alu.Extsb(ra)
	} else if op_code == instruction.EXTSH {
		result = this.alu.Extsh(ra)
	} else if op_code == instruction.EXTUB {
		result = this.alu.Extub(ra)
	} else if op_code == instruction.EXTUH {
		result = this.alu.Extuh(ra)
	} else if op_code == instruction.SATS {
		result = this.alu.Sats(ra)
	} else if op_code == instruction.TIME_CFG {
		err := errors.New("TimeCfg is not yet implemented")
		panic(err)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteZrc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrcOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRC {
		err := errors.New("suffix is not RRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.CAO {
		result = this.alu.Cao(ra)
	} else if op_code == instruction.CLO {
		result = this.alu.Clo(ra)
	} else if op_code == instruction.CLS {
		result = this.alu.Cls(ra)
	} else if op_code == instruction.CLZ {
		result = this.alu.Clz(ra)
	} else if op_code == instruction.EXTSB {
		result = this.alu.Extsb(ra)
	} else if op_code == instruction.EXTSH {
		result = this.alu.Extsh(ra)
	} else if op_code == instruction.EXTUB {
		result = this.alu.Extub(ra)
	} else if op_code == instruction.EXTUH {
		result = this.alu.Extuh(ra)
	} else if op_code == instruction.SATS {
		result = this.alu.Sats(ra)
	} else if op_code == instruction.TIME_CFG {
		err := errors.New("TimeCfg is not yet implemented")
		panic(err)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogSetCc(instruction_, result)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteZrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRCI {
		err := errors.New("suffix is not ZRCI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_cao_rrci_op_code := instruction_.CaoRrciOpCodes()[op_code]; is_cao_rrci_op_code {
		this.ExecuteCaoZrci(instruction_)
	} else if _, is_extsb_rrci_op_code := instruction_.ExtsbRrciOpCodes()[op_code]; is_extsb_rrci_op_code {
		this.ExecuteExtsbZrci(instruction_)
	} else if _, is_time_cfg_rrci_op_code := instruction_.TimeCfgRrciOpCodes()[op_code]; is_time_cfg_rrci_op_code {
		this.ExecuteTimeCfgZrci(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteCaoZrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.CaoRrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid cao RRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRCI {
		err := errors.New("suffix is not ZRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.CAO {
		result = this.alu.Cao(ra)
	} else if op_code == instruction.CLO {
		result = this.alu.Clo(ra)
	} else if op_code == instruction.CLS {
		result = this.alu.Cls(ra)
	} else if op_code == instruction.CLZ {
		result = this.alu.Clz(ra)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetCountNzCc(instruction_, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteExtsbZrci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.ExtsbRrciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid extsb RRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRCI {
		err := errors.New("suffix is not ZRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.EXTSB {
		result = this.alu.Extsb(ra)
	} else if op_code == instruction.EXTSH {
		result = this.alu.Extsh(ra)
	} else if op_code == instruction.EXTUB {
		result = this.alu.Extub(ra)
	} else if op_code == instruction.EXTUH {
		result = this.alu.Extuh(ra)
	} else if op_code == instruction.SATS {
		result = this.alu.Sats(ra)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetLogNzCc(instruction_, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteTimeCfgZrci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteTimeCfgZrci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRr(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRr is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRrc(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRrc is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRrci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRrci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteDrdici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.DrdiciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid DRDICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DRDICI {
		err := errors.New("suffix is not DRDICI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_div_step_drdici_op_code := instruction_.DivStepDrdiciOpCodes()[op_code]; is_div_step_drdici_op_code {
		this.ExecuteDivStepDrdici(instruction_)
	} else if _, is_mul_step_drdici_op_code := instruction_.MulStepDrdiciOpCodes()[op_code]; is_mul_step_drdici_op_code {
		this.ExecuteMulStepDrdici(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteDivStepDrdici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.DivStepDrdiciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid div_step DRDICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DRDICI {
		err := errors.New("suffix is not DRDICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	dbe := thread.RegFile().ReadGpReg(instruction_.Db().EvenRegDescriptor(), abi.SIGNED)
	dbo := thread.RegFile().ReadGpReg(instruction_.Db().OddRegDescriptor(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	dbo_word := new(abi.Word)
	dbo_word.Init(mram_data_width)
	dbo_word.SetValue(dbo)

	ra_shift_value := this.alu.Lsl(ra, imm)
	ra_shift_word := new(abi.Word)
	ra_shift_word.Init(mram_data_width)
	ra_shift_word.SetValue(ra_shift_value)

	result, _, _ := this.alu.Sub(dbo, ra_shift_value)

	var dce int64
	var dco int64
	if dbo_word.Value(abi.UNSIGNED) >= ra_shift_word.Value(abi.UNSIGNED) {
		dce = this.alu.Lsl1(dbe, 1)
		dco = result
	} else {
		dce = this.alu.Lsl(dbe, 1)
		dco = thread.RegFile().ReadGpReg(instruction_.Dc().OddRegDescriptor(), abi.SIGNED)
	}

	thread.RegFile().ClearConditions()
	this.SetDivCc(instruction_, ra)

	thread.RegFile().WriteGpReg(instruction_.Dc().EvenRegDescriptor(), dce)
	thread.RegFile().WriteGpReg(instruction_.Dc().OddRegDescriptor(), dco)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, false)
}

func (this *Logic) ExecuteMulStepDrdici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.MulStepDrdiciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid mul_step DRDICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DRDICI {
		err := errors.New("suffix is not DRDICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	dbe := thread.RegFile().ReadGpReg(instruction_.Db().EvenRegDescriptor(), abi.SIGNED)
	dbo := thread.RegFile().ReadGpReg(instruction_.Db().OddRegDescriptor(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	result1 := this.alu.Lsr(dbe, 1)
	result2, _, _ := this.alu.Sub(this.alu.And(dbe, 1), 1)

	var dco int64
	if result2 == 0 {
		dco, _, _ = this.alu.Add(dbo, this.alu.Lsl(ra, imm))
	} else {
		dco = thread.RegFile().ReadGpReg(instruction_.Dc().OddRegDescriptor(), abi.SIGNED)
	}

	dce := this.alu.Lsr(dbe, 1)

	thread.RegFile().ClearConditions()
	this.SetBootCc(instruction_, ra, result1)

	thread.RegFile().WriteGpReg(instruction_.Dc().EvenRegDescriptor(), dce)
	thread.RegFile().WriteGpReg(instruction_.Dc().OddRegDescriptor(), dco)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result1, false)
}

func (this *Logic) ExecuteRrri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRI {
		err := errors.New("suffix is not RRRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.LSL_ADD {
		result, carry, _ = this.alu.LslAdd(ra, rb, imm)
	} else if op_code == instruction.LSL_SUB {
		result, carry, _ = this.alu.LslSub(ra, rb, imm)
	} else if op_code == instruction.LSR_ADD {
		result, carry, _ = this.alu.LsrAdd(ra, rb, imm)
	} else if op_code == instruction.ROL_ADD {
		result, carry, _ = this.alu.RolAdd(ra, rb, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WriteGpReg(instruction_.Rc(), result)
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRrrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RRRICI {
		err := errors.New("suffix is not RRRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.LSL_ADD {
		result, carry, _ = this.alu.LslAdd(ra, rb, imm)
	} else if op_code == instruction.LSL_SUB {
		result, carry, _ = this.alu.LslSub(ra, rb, imm)
	} else if op_code == instruction.LSR_ADD {
		result, carry, _ = this.alu.LsrAdd(ra, rb, imm)
	} else if op_code == instruction.ROL_ADD {
		result, carry, _ = this.alu.RolAdd(ra, rb, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetDivNzCc(instruction_, ra)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZrri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRI {
		err := errors.New("suffix is not ZRRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.LSL_ADD {
		result, carry, _ = this.alu.LslAdd(ra, rb, imm)
	} else if op_code == instruction.LSL_SUB {
		result, carry, _ = this.alu.LslSub(ra, rb, imm)
	} else if op_code == instruction.LSR_ADD {
		result, carry, _ = this.alu.LsrAdd(ra, rb, imm)
	} else if op_code == instruction.ROL_ADD {
		result, carry, _ = this.alu.RolAdd(ra, rb, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZrrici(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RrriciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZRRICI {
		err := errors.New("suffix is not ZRRICI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.LSL_ADD {
		result, carry, _ = this.alu.LslAdd(ra, rb, imm)
	} else if op_code == instruction.LSL_SUB {
		result, carry, _ = this.alu.LslSub(ra, rb, imm)
	} else if op_code == instruction.LSR_ADD {
		result, carry, _ = this.alu.LsrAdd(ra, rb, imm)
	} else if op_code == instruction.ROL_ADD {
		result, carry, _ = this.alu.RolAdd(ra, rb, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetDivNzCc(instruction_, ra)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSRrri(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRrri is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRrrici(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRrrici is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteRir(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RirOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RIR {
		err := errors.New("suffix is not RIR")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	imm := instruction_.Imm().Value()
	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(imm, ra)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(imm, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WriteGpReg(instruction_.Rc(), result)
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRirc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RircOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RIRC {
		err := errors.New("suffix is not RIRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	imm := instruction_.Imm().Value()
	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(imm, ra)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(imm, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubSetCc(instruction_, imm, ra, result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 1)
	} else {
		thread.RegFile().WriteGpReg(instruction_.Rc(), 0)
	}

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteRirci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RirciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.RIRCI {
		err := errors.New("suffix is not RIRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	imm := instruction_.Imm().Value()
	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(imm, ra)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(imm, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubNzCc(instruction_, imm, ra, result, carry, overflow)

	thread.RegFile().WriteGpReg(instruction_.Rc(), result)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZir(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RirOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZIR {
		err := errors.New("suffix is not ZIR")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	imm := instruction_.Imm().Value()
	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(imm, ra)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(imm, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZirc(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RircOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZIRC {
		err := errors.New("suffix is not ZIRC")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	imm := instruction_.Imm().Value()
	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64
	var carry bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, _ = this.alu.Sub(imm, ra)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, _ = this.alu.Subc(imm, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubSetCc(instruction_, imm, ra, result)

	thread.RegFile().IncrementPcReg()

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteZirci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.RirciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ZIRCI {
		err := errors.New("suffix is not ZIRCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	imm := instruction_.Imm().Value()
	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)

	var result int64
	var carry bool
	var overflow bool

	op_code := instruction_.OpCode()
	if op_code == instruction.SUB {
		result, carry, overflow = this.alu.Sub(imm, ra)
	} else if op_code == instruction.SUBC {
		carry_flag := thread.RegFile().ReadFlagReg(instruction.CARRY)
		result, carry, overflow = this.alu.Subc(imm, ra, carry_flag)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	this.SetSubNzCc(instruction_, imm, ra, result, carry, overflow)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}

	this.SetFlags(instruction_, result, carry)
}

func (this *Logic) ExecuteSRirc(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRirc is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRirci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRirci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteR(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteR is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteRci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteRci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteZ(instruction_ *instruction.Instruction) {
	if _, found := instruction_.ROpCodes()[instruction_.OpCode()]; !found &&
		instruction_.OpCode() != instruction.NOP {
		err := errors.New("op code is not a valid R op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.Z {
		err := errors.New("suffix is not Z")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	thread.RegFile().IncrementPcReg()
}

func (this *Logic) ExecuteZci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteZci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSR(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSR is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSRci(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSRci is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteCi(instruction_ *instruction.Instruction) {
	if _, found := instruction_.CiOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid CI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.CI {
		err := errors.New("suffix is not CI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	this.thread_scheduler.Sleep(thread.ThreadId())

	thread.RegFile().ClearConditions()

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}
}

func (this *Logic) ExecuteI(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteI is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteDdci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.DdciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid DDCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DDCI {
		err := errors.New("suffix is not DDCI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_movd_ddci_op_code := instruction_.MovdDdciOpCodes()[op_code]; is_movd_ddci_op_code {
		this.ExecuteMovdDdci(instruction_)
	} else if _, is_swapd_ddci_op_code := instruction_.SwapdDdciOpCodes()[op_code]; is_swapd_ddci_op_code {
		this.ExecuteSwapdDdciRri(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteMovdDdci(instruction_ *instruction.Instruction) {
	if _, found := instruction_.MovdDdciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid movd DDCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DDCI {
		err := errors.New("suffix is not DDCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	dbe := thread.RegFile().ReadGpReg(instruction_.Db().EvenRegDescriptor(), abi.SIGNED)
	dbo := thread.RegFile().ReadGpReg(instruction_.Db().OddRegDescriptor(), abi.SIGNED)

	thread.RegFile().ClearConditions()

	thread.RegFile().WriteGpReg(instruction_.Dc().EvenRegDescriptor(), dbe)
	thread.RegFile().WriteGpReg(instruction_.Dc().OddRegDescriptor(), dbo)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}
}

func (this *Logic) ExecuteSwapdDdciRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SwapdDdciOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid movd DDCI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DDCI {
		err := errors.New("suffix is not DDCI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	dbe := thread.RegFile().ReadGpReg(instruction_.Db().EvenRegDescriptor(), abi.SIGNED)
	dbo := thread.RegFile().ReadGpReg(instruction_.Db().OddRegDescriptor(), abi.SIGNED)

	thread.RegFile().ClearConditions()

	thread.RegFile().WriteGpReg(instruction_.Dc().EvenRegDescriptor(), dbo)
	thread.RegFile().WriteGpReg(instruction_.Dc().OddRegDescriptor(), dbe)

	if thread.RegFile().ReadConditionReg(instruction_.Condition()) {
		thread.RegFile().WritePcReg(instruction_.Pc().Value())
	} else {
		thread.RegFile().IncrementPcReg()
	}
}

func (this *Logic) ExecuteErri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.ErriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid ERRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ERRI {
		err := errors.New("suffix is not ERRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	off := instruction_.Off().Value()

	address, _, _ := this.alu.Add(ra, off)

	var result int64

	op_code := instruction_.OpCode()
	if op_code == instruction.LBS {
		result = this.operand_collector.Lbs(address)
	} else if op_code == instruction.LBU {
		result = this.operand_collector.Lbu(address)
	} else if op_code == instruction.LHS {
		result = this.operand_collector.Lhs(address)
	} else if op_code == instruction.LHU {
		result = this.operand_collector.Lhu(address)
	} else if op_code == instruction.LW {
		result = this.operand_collector.Lw(address)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WriteGpReg(instruction_.Rc(), result)
	thread.RegFile().IncrementPcReg()
}

func (this *Logic) ExecuteSErri(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteSErri is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteEdri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.EdriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid EDRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.EDRI {
		err := errors.New("suffix is not EDRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	off := instruction_.Off().Value()

	address, _, _ := this.alu.Add(ra, off)

	var even int64
	var odd int64

	op_code := instruction_.OpCode()
	if op_code == instruction.LD {
		even, odd = this.operand_collector.Ld(address)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().WritePairReg(instruction_.Dc(), even, odd)
	thread.RegFile().IncrementPcReg()
}

func (this *Logic) ExecuteErii(instruction_ *instruction.Instruction) {
	if _, found := instruction_.EriiOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid ERII op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ERII {
		err := errors.New("suffix is not ERII")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	off := instruction_.Off().Value()
	imm := instruction_.Imm().Value()

	address, _, _ := this.alu.Add(ra, off)

	op_code := instruction_.OpCode()
	if op_code == instruction.SB {
		this.operand_collector.Sb(address, imm)
	} else if op_code == instruction.SB_ID {
		this.operand_collector.Sb(address, this.alu.Or(int64(thread.ThreadId()), imm))
	} else if op_code == instruction.SH {
		this.operand_collector.Sh(address, imm)
	} else if op_code == instruction.SH_ID {
		this.operand_collector.Sh(address, this.alu.Or(int64(thread.ThreadId()), imm))
	} else if op_code == instruction.SW {
		this.operand_collector.Sw(address, imm)
	} else if op_code == instruction.SW_ID {
		this.operand_collector.Sw(address, this.alu.Or(int64(thread.ThreadId()), imm))
	} else if op_code == instruction.SD {
		even, odd := this.alu.UnsignedExtension(imm)
		this.operand_collector.Sd(address, even, odd)
	} else if op_code == instruction.SD_ID {
		even, odd := this.alu.UnsignedExtension(this.alu.Or(int64(thread.ThreadId()), imm))
		this.operand_collector.Sd(address, even, odd)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()
}

func (this *Logic) ExecuteErir(instruction_ *instruction.Instruction) {
	if _, found := instruction_.ErirOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid ERIR op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ERIR {
		err := errors.New("suffix is not ERIR")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	off := instruction_.Off().Value()
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)

	address, _, _ := this.alu.Add(ra, off)

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	rb_word := new(abi.Word)
	rb_word.Init(config_loader.MramDataWidth())
	rb_word.SetValue(rb)

	op_code := instruction_.OpCode()
	if op_code == instruction.SB {
		this.operand_collector.Sb(address, rb_word.BitSlice(abi.UNSIGNED, 0, 8))
	} else if op_code == instruction.SH {
		this.operand_collector.Sh(address, rb_word.BitSlice(abi.UNSIGNED, 0, 16))
	} else if op_code == instruction.SW {
		this.operand_collector.Sw(address, rb_word.BitSlice(abi.UNSIGNED, 0, 32))
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()
}

func (this *Logic) ExecuteErid(instruction_ *instruction.Instruction) {
	if _, found := instruction_.EridOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid ERID op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.ERID {
		err := errors.New("suffix is not ERID")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	off := instruction_.Off().Value()
	even, odd := thread.RegFile().ReadPairReg(instruction_.Db(), abi.SIGNED)

	address, _, _ := this.alu.Add(ra, off)

	op_code := instruction_.OpCode()
	if op_code == instruction.SD {
		this.operand_collector.Sd(address, even, odd)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	thread.RegFile().ClearConditions()
	thread.RegFile().IncrementPcReg()
}

func (this *Logic) ExecuteDmaRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.DmaRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid DMA_RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DMA_RRI {
		err := errors.New("suffix is not DMA_RRI")
		panic(err)
	}

	op_code := instruction_.OpCode()
	if _, is_ldma_dma_rri_op_code := instruction_.LdmaDmaRriOpCodes()[op_code]; is_ldma_dma_rri_op_code {
		this.ExecuteLdmaDmaRri(instruction_)
	} else if _, is_ldmai_dma_rri_op_code := instruction_.LdmaiDmaRriOpCodes()[op_code]; is_ldmai_dma_rri_op_code {
		this.ExecuteLdmaiDmaRri(instruction_)
	} else if _, is_sdma_dma_rri_op_code := instruction_.SdmaDmaRriOpCodes()[op_code]; is_sdma_dma_rri_op_code {
		this.ExecuteSdmaDmaRri(instruction_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Logic) ExecuteLdmaDmaRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.LdmaDmaRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid ldma DMA_RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DMA_RRI {
		err := errors.New("suffix is not DMA_RRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	wram_end_address := config_loader.WramOffset() + config_loader.WramSize()
	wram_end_address_width := int(math.Floor(math.Log2(float64(wram_end_address))) + 1)
	wram_mask := this.Pow2(wram_end_address_width) - 1
	wram_address := this.alu.And(ra, wram_mask)

	mram_end_address := config_loader.MramOffset() + config_loader.MramSize()
	mram_end_address_width := int(math.Floor(math.Log2(float64(mram_end_address))) + 1)
	mram_mask := this.Pow2(mram_end_address_width) - 1
	mram_address := this.alu.And(rb, mram_mask)

	size := (1 + this.alu.And(imm+this.alu.And(this.alu.Lsr(ra, 24), 255), 255)) * this.min_access_granularity

	this.dma.TransferFromMramToWram(wram_address, mram_address, size, instruction_)

	thread.RegFile().ClearConditions()
}

func (this *Logic) ExecuteLdmaiDmaRri(instruction_ *instruction.Instruction) {
	err := errors.New("ExecuteLdmaiDmaRri is not yet implemented")
	panic(err)
}

func (this *Logic) ExecuteSdmaDmaRri(instruction_ *instruction.Instruction) {
	if _, found := instruction_.SdmaDmaRriOpCodes()[instruction_.OpCode()]; !found {
		err := errors.New("op code is not a valid sdma DMA_RRI op code")
		panic(err)
	} else if instruction_.Suffix() != instruction.DMA_RRI {
		err := errors.New("suffix is not DMA_RRI")
		panic(err)
	}

	thread := this.scoreboard[instruction_]

	ra := thread.RegFile().ReadSrcReg(instruction_.Ra(), abi.SIGNED)
	rb := thread.RegFile().ReadSrcReg(instruction_.Rb(), abi.SIGNED)
	imm := instruction_.Imm().Value()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	wram_end_address := config_loader.WramOffset() + config_loader.WramSize()
	wram_end_address_width := int(math.Floor(math.Log2(float64(wram_end_address))) + 1)
	wram_mask := this.Pow2(wram_end_address_width) - 1
	wram_address := this.alu.And(ra, wram_mask)

	mram_end_address := config_loader.MramOffset() + config_loader.MramSize()
	mram_end_address_width := int(math.Floor(math.Log2(float64(mram_end_address))) + 1)
	mram_mask := this.Pow2(mram_end_address_width) - 1
	mram_address := this.alu.And(rb, mram_mask)

	size := (1 + this.alu.And(imm+this.alu.And(this.alu.Lsr(ra, 24), 255), 255)) * this.min_access_granularity

	this.dma.TransferFromWramToMram(wram_address, mram_address, size, instruction_)

	thread.RegFile().ClearConditions()
}

func (this *Logic) SetAcquireCc(instruction_ *instruction.Instruction, result int64) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}
}

func (this *Logic) SetAddNzCc(
	instruction_ *instruction.Instruction,
	operand1 int64,
	result int64,
	carry bool,
	overflow bool,
) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if carry {
		thread.RegFile().SetCondition(cc.C)
	} else {
		thread.RegFile().SetCondition(cc.NC)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if overflow {
		thread.RegFile().SetCondition(cc.OV)
	} else {
		thread.RegFile().SetCondition(cc.NOV)
	}

	if result >= 0 {
		thread.RegFile().SetCondition(cc.PL)
	} else {
		thread.RegFile().SetCondition(cc.MI)
	}

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	result_word := new(abi.Word)
	result_word.Init(config_loader.MramDataWidth())
	result_word.SetValue(result)

	if result_word.Bit(6) {
		thread.RegFile().SetCondition(cc.NC5)
	}

	if result_word.Bit(7) {
		thread.RegFile().SetCondition(cc.NC6)
	}

	if result_word.Bit(8) {
		thread.RegFile().SetCondition(cc.NC7)
	}

	if result_word.Bit(9) {
		thread.RegFile().SetCondition(cc.NC8)
	}

	if result_word.Bit(10) {
		thread.RegFile().SetCondition(cc.NC9)
	}

	if result_word.Bit(11) {
		thread.RegFile().SetCondition(cc.NC10)
	}

	if result_word.Bit(12) {
		thread.RegFile().SetCondition(cc.NC11)
	}

	if result_word.Bit(13) {
		thread.RegFile().SetCondition(cc.NC12)
	}

	if result_word.Bit(14) {
		thread.RegFile().SetCondition(cc.NC13)
	}

	if result_word.Bit(15) {
		thread.RegFile().SetCondition(cc.NC14)
	}
}

func (this *Logic) SetBootCc(instruction_ *instruction.Instruction, operand1 int64, result int64) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}
}

func (this *Logic) SetCountNzCc(
	instruction_ *instruction.Instruction,
	operand1 int64,
	result int64,
) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	if result == int64(config_loader.MramDataWidth()) {
		thread.RegFile().SetCondition(cc.MAX)
	} else {
		thread.RegFile().SetCondition(cc.NMAX)
	}
}

func (this *Logic) SetDivCc(instruction_ *instruction.Instruction, operand1 int64) {
	thread := this.scoreboard[instruction_]

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}
}

func (this *Logic) SetDivNzCc(instruction_ *instruction.Instruction, operand1 int64) {
	thread := this.scoreboard[instruction_]

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}
}

func (this *Logic) SetExtSubSetCc(
	instruction_ *instruction.Instruction,
	operand1 int64,
	operand2 int64,
	result int64,
	carry bool,
	overflow bool,
) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if carry {
		thread.RegFile().SetCondition(cc.C)
	} else {
		thread.RegFile().SetCondition(cc.NC)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if overflow {
		thread.RegFile().SetCondition(cc.OV)
	} else {
		thread.RegFile().SetCondition(cc.NOV)
	}

	if result >= 0 {
		thread.RegFile().SetCondition(cc.PL)
	} else {
		thread.RegFile().SetCondition(cc.MI)
	}

	if operand1 == operand2 {
		thread.RegFile().SetCondition(cc.EQ)
	} else {
		thread.RegFile().SetCondition(cc.NEQ)
	}

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	if word1.Value(abi.UNSIGNED) < word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.LTU)
	}

	if word1.Value(abi.UNSIGNED) <= word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.LEU)
	}

	if word1.Value(abi.UNSIGNED) > word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.GTU)
	}

	if word1.Value(abi.UNSIGNED) >= word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.GEU)
	}

	if word1.Value(abi.SIGNED) < word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.LTS)
	}

	if word1.Value(abi.SIGNED) <= word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.LES)
	}

	if word1.Value(abi.SIGNED) > word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.GTS)
	}

	if word1.Value(abi.SIGNED) >= word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.GES)
	}

	if carry || thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XLEU)
	}

	if carry || !thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XGTU)
	}

	if thread.RegFile().ReadFlagReg(instruction.ZERO) && (result < 0 || overflow) {
		thread.RegFile().SetCondition(cc.XLES)
	}

	if !thread.RegFile().ReadFlagReg(instruction.ZERO) && (result >= 0 || overflow) {
		thread.RegFile().SetCondition(cc.XGTS)
	}
}

func (this *Logic) SetImmShiftNzCc(
	instruction_ *instruction.Instruction,
	operand1 int64,
	result int64,
) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if result%2 == 0 {
		thread.RegFile().SetCondition(cc.E)
	} else {
		thread.RegFile().SetCondition(cc.O)
	}

	if result >= 0 {
		thread.RegFile().SetCondition(cc.PL)
	} else {
		thread.RegFile().SetCondition(cc.MI)
	}

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1%2 == 0 {
		thread.RegFile().SetCondition(cc.SE)
	} else {
		thread.RegFile().SetCondition(cc.SO)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}
}

func (this *Logic) SetLogNzCc(instruction_ *instruction.Instruction, operand1 int64, result int64) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if result >= 0 {
		thread.RegFile().SetCondition(cc.PL)
	} else {
		thread.RegFile().SetCondition(cc.MI)
	}

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}
}

func (this *Logic) SetLogSetCc(instruction_ *instruction.Instruction, result int64) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}
}

func (this *Logic) SetMulNzCc(instruction_ *instruction.Instruction, operand1 int64, result int64) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if operand1 == 0 {
		thread.RegFile().SetCondition(cc.SZ)
	} else {
		thread.RegFile().SetCondition(cc.SNZ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}

	if result < 256 {
		thread.RegFile().SetCondition(cc.SMALL)
	} else {
		thread.RegFile().SetCondition(cc.LARGE)
	}
}

func (this *Logic) SetSubNzCc(
	instruction_ *instruction.Instruction,
	operand1 int64,
	operand2 int64,
	result int64,
	carry bool,
	overflow bool,
) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if carry {
		thread.RegFile().SetCondition(cc.C)
	} else {
		thread.RegFile().SetCondition(cc.NC)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if overflow {
		thread.RegFile().SetCondition(cc.OV)
	} else {
		thread.RegFile().SetCondition(cc.NOV)
	}

	if result >= 0 {
		thread.RegFile().SetCondition(cc.PL)
	} else {
		thread.RegFile().SetCondition(cc.MI)
	}

	if operand1 == operand2 {
		thread.RegFile().SetCondition(cc.EQ)
	} else {
		thread.RegFile().SetCondition(cc.NEQ)
	}

	if operand1 >= 0 {
		thread.RegFile().SetCondition(cc.SPL)
	} else {
		thread.RegFile().SetCondition(cc.SMI)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	if word1.Value(abi.UNSIGNED) < word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.LTU)
	}

	if word1.Value(abi.UNSIGNED) <= word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.LEU)
	}

	if word1.Value(abi.UNSIGNED) > word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.GTU)
	}

	if word1.Value(abi.UNSIGNED) >= word2.Value(abi.UNSIGNED) {
		thread.RegFile().SetCondition(cc.GEU)
	}

	if word1.Value(abi.SIGNED) < word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.LTS)
	}

	if word1.Value(abi.SIGNED) <= word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.LES)
	}

	if word1.Value(abi.SIGNED) > word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.GTS)
	}

	if word1.Value(abi.SIGNED) >= word2.Value(abi.SIGNED) {
		thread.RegFile().SetCondition(cc.GES)
	}

	if carry || thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XLEU)
	}

	if carry || !thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XGTU)
	}

	if thread.RegFile().ReadFlagReg(instruction.ZERO) && (result < 0 || overflow) {
		thread.RegFile().SetCondition(cc.XLES)
	}

	if !thread.RegFile().ReadFlagReg(instruction.ZERO) && (result >= 0 || overflow) {
		thread.RegFile().SetCondition(cc.XGTS)
	}
}

func (this *Logic) SetSubSetCc(
	instruction_ *instruction.Instruction,
	operand1 int64,
	operand2 int64,
	result int64,
) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetCondition(cc.Z)
	} else {
		thread.RegFile().SetCondition(cc.NZ)
	}

	if result == 0 && thread.RegFile().ReadFlagReg(instruction.ZERO) {
		thread.RegFile().SetCondition(cc.XZ)
	} else {
		thread.RegFile().SetCondition(cc.XNZ)
	}

	if operand1 == operand2 {
		thread.RegFile().SetCondition(cc.EQ)
	} else {
		thread.RegFile().SetCondition(cc.NEQ)
	}
}

func (this *Logic) SetFlags(instruction_ *instruction.Instruction, result int64, carry bool) {
	thread := this.scoreboard[instruction_]

	if result == 0 {
		thread.RegFile().SetFlag(instruction.ZERO)
	} else {
		thread.RegFile().ClearFlag(instruction.ZERO)
	}

	if carry {
		thread.RegFile().SetFlag(instruction.CARRY)
	} else {
		thread.RegFile().ClearFlag(instruction.CARRY)
	}
}

func (this *Logic) Pow2(exponent int) int64 {
	if exponent < 0 {
		err := errors.New("exponent < 0")
		panic(err)
	}

	value := int64(1)
	for i := 0; i < exponent; i++ {
		value *= 2
	}
	return value
}

func (this *Logic) PrintRegFile(thread *Thread) string {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	lines := ""
	for i := 0; i < config_loader.NumGpRegisters(); i++ {
		gp_reg_descriptor := new(reg_descriptor.GpRegDescriptor)
		gp_reg_descriptor.Init(i)

		lines += fmt.Sprintf(
			"r%d: %d\n",
			i,
			thread.RegFile().ReadGpReg(gp_reg_descriptor, abi.SIGNED),
		)
	}
	return lines
}
