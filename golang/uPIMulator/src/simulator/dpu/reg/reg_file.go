package reg

import (
	"uPIMulator/src/abi/word"
	"uPIMulator/src/linker/kernel/instruction"
	"uPIMulator/src/linker/kernel/instruction/cc"
	"uPIMulator/src/linker/kernel/instruction/reg_descriptor"
	"uPIMulator/src/misc"
)

type RegFile struct {
	gp_regs       []*GpReg
	sp_reg        *SpReg
	pc_reg        *PcReg
	condition_reg *ConditionReg
	flag_reg      *FlagReg
	exception_reg *ExceptionReg
}

func (this *RegFile) Init(thread_id int) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.gp_regs = make([]*GpReg, 0)
	for i := 0; i < config_loader.NumGpRegisters(); i++ {
		gp_reg := new(GpReg)
		gp_reg.Init(i)

		this.gp_regs = append(this.gp_regs, gp_reg)
	}

	this.sp_reg = new(SpReg)
	this.sp_reg.Init(thread_id)

	this.pc_reg = new(PcReg)
	this.pc_reg.Init()

	this.condition_reg = new(ConditionReg)
	this.condition_reg.Init()

	this.flag_reg = new(FlagReg)
	this.flag_reg.Init()

	this.exception_reg = new(ExceptionReg)
	this.exception_reg.Init()
}

func (this *RegFile) Fini() {
	for _, gp_reg := range this.gp_regs {
		gp_reg.Fini()
	}

	this.sp_reg.Fini()
	this.pc_reg.Fini()
	this.condition_reg.Fini()
	this.flag_reg.Fini()
	this.exception_reg.Fini()
}

func (this *RegFile) ReadGpReg(
	gp_reg_descriptor *reg_descriptor.GpRegDescriptor,
	representation word.Representation,
) int64 {
	return this.gp_regs[gp_reg_descriptor.Index()].Read(representation)
}

func (this *RegFile) ReadSpReg(
	sp_reg_descriptor *reg_descriptor.SpRegDescriptor,
	representation word.Representation,
) int64 {
	return this.sp_reg.Read(sp_reg_descriptor, representation)
}

func (this *RegFile) ReadPairReg(
	pair_reg_descriptor *reg_descriptor.PairRegDescriptor,
	representation word.Representation,
) (int64, int64) {
	even := this.ReadGpReg(pair_reg_descriptor.EvenRegDescriptor(), representation)
	odd := this.ReadGpReg(pair_reg_descriptor.OddRegDescriptor(), word.UNSIGNED)

	return even, odd
}

func (this *RegFile) ReadSrcReg(
	src_reg_descriptor *reg_descriptor.SrcRegDescriptor,
	representation word.Representation,
) int64 {
	if src_reg_descriptor.IsGpRegDescriptor() {
		return this.ReadGpReg(src_reg_descriptor.GpRegDescriptor(), representation)
	} else {
		return this.ReadSpReg(src_reg_descriptor.SpRegDescriptor(), representation)
	}
}

func (this *RegFile) ReadPcReg() int64 {
	return this.pc_reg.Read()
}

func (this *RegFile) ReadConditionReg(condition cc.Condition) bool {
	return this.condition_reg.Condition(condition)
}

func (this *RegFile) ReadFlagReg(flag instruction.Flag) bool {
	return this.flag_reg.Flag(flag)
}

func (this *RegFile) ReadExceptionReg(exception instruction.Exception) bool {
	return this.exception_reg.Exception(exception)
}

func (this *RegFile) WriteGpReg(gp_reg_descriptor *reg_descriptor.GpRegDescriptor, value int64) {
	this.gp_regs[gp_reg_descriptor.Index()].Write(value)
}

func (this *RegFile) WritePairReg(
	pair_reg_descriptor *reg_descriptor.PairRegDescriptor,
	even int64,
	odd int64,
) {
	this.WriteGpReg(pair_reg_descriptor.EvenRegDescriptor(), even)
	this.WriteGpReg(pair_reg_descriptor.OddRegDescriptor(), odd)
}

func (this *RegFile) WritePcReg(value int64) {
	this.pc_reg.Write(value)
}

func (this *RegFile) IncrementPcReg() {
	this.pc_reg.Increment()
}

func (this *RegFile) SetCondition(condition cc.Condition) {
	this.condition_reg.SetCondition(condition)
}

func (this *RegFile) ClearCondition(condition cc.Condition) {
	this.condition_reg.ClearCondition(condition)
}

func (this *RegFile) ClearConditions() {
	this.condition_reg.ClearConditions()
}

func (this *RegFile) SetFlag(flag instruction.Flag) {
	this.flag_reg.SetFlag(flag)
}

func (this *RegFile) ClearFlag(flag instruction.Flag) {
	this.flag_reg.ClearFlag(flag)
}

func (this *RegFile) ClearFlags() {
	this.flag_reg.ClearFlags()
}

func (this *RegFile) SetException(exception instruction.Exception) {
	this.exception_reg.SetException(exception)
}

func (this *RegFile) ClearException(exception instruction.Exception) {
	this.exception_reg.ClearException(exception)
}

func (this *RegFile) ClearExceptions() {
	this.exception_reg.ClearExceptions()
}
