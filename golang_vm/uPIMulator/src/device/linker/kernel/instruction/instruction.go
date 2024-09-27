package instruction

import (
	"errors"
	"math"
	"strconv"
	"uPIMulator/src/device/abi"
	"uPIMulator/src/device/linker/kernel/instruction/cc"
	"uPIMulator/src/device/linker/kernel/instruction/reg_descriptor"
	"uPIMulator/src/encoding"
	"uPIMulator/src/misc"
)

type Instruction struct {
	op_code OpCode
	suffix  Suffix

	rc *reg_descriptor.GpRegDescriptor
	ra *reg_descriptor.SrcRegDescriptor
	rb *reg_descriptor.SrcRegDescriptor

	dc *reg_descriptor.PairRegDescriptor
	db *reg_descriptor.PairRegDescriptor

	condition *cc.Condition

	imm *abi.Immediate
	off *abi.Immediate
	pc  *abi.Immediate

	endian *Endian
}

func (this *Instruction) InitRici(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RiciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RICI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RICI
	this.ra = ra

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 16, imm)

	this.condition = new(cc.Condition)
	*this.condition = condition

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitRri(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
) {
	if _, found := this.RriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRI
	this.rc = rc
	this.ra = ra

	if _, is_add_rri_op_code := this.AddRriOpCodes()[op_code]; is_add_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 32, imm)
	} else if _, is_asr_rri_op_code := this.AsrRriOpCodes()[op_code]; is_asr_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_call_rri_op_code := this.CallRriOpCodes()[op_code]; is_call_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}
}

func (this *Instruction) InitRric(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
) {
	if _, found := this.RricOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRIC
	this.rc = rc
	this.ra = ra

	if _, is_add_rric_op_code := this.AddRricOpCodes()[op_code]; is_add_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)

		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_asr_rric_op_code := this.AsrRricOpCodes()[op_code]; is_asr_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)

		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_sub_rric_op_code := this.SubRricOpCodes()[op_code]; is_sub_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)

		ext_sub_setcc := new(cc.ExtSubSetCc)
		ext_sub_setcc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = ext_sub_setcc.Condition()
	} else {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}
}

func (this *Instruction) InitRrici(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RriciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRICI
	this.rc = rc
	this.ra = ra

	if _, is_add_rrici_op_code := this.AddRriciOpCodes()[op_code]; is_add_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)

		add_nz_cc := new(cc.AddNzCc)
		add_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = add_nz_cc.Condition()
	} else if _, is_and_rrici_op_code := this.AndRriciOpCodes()[op_code]; is_and_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)

		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_asr_rrici_op_code := this.AsrRriciOpCodes()[op_code]; is_asr_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)

		imm_shift_nz_cc := new(cc.ImmShiftNzCc)
		imm_shift_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = imm_shift_nz_cc.Condition()
	} else if _, is_sub_rrici_op_code := this.SubRriciOpCodes()[op_code]; is_sub_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)

		sub_nz_cc := new(cc.SubNzCc)
		sub_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_nz_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitRrif(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
) {
	if _, found := this.RrifOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRIF
	this.rc = rc
	this.ra = ra

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	false_cc := new(cc.FalseCc)
	false_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = false_cc.Condition()
}

func (this *Instruction) InitRrr(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
) {
	if _, found := this.RrrOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRR
	this.rc = rc
	this.ra = ra
	this.rb = rb
}

func (this *Instruction) InitRrrc(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RrrcOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRRC
	this.rc = rc
	this.ra = ra
	this.rb = rb

	if _, is_add_rrrc_op_code := this.AddRrrcOpCodes()[op_code]; is_add_rrrc_op_code {
		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_rsub_rrrc_op_code := this.RsubRrrcOpCodes()[op_code]; is_rsub_rrrc_op_code {
		sub_set_cc := new(cc.SubSetCc)
		sub_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_set_cc.Condition()
	} else if _, is_sub_rrrc_op_code := this.SubRrrcOpCodes()[op_code]; is_sub_rrrc_op_code {
		ext_sub_set_cc := new(cc.ExtSubSetCc)
		ext_sub_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = ext_sub_set_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}
}

func (this *Instruction) InitRrrci(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrrciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRRCI
	this.rc = rc
	this.ra = ra
	this.rb = rb

	if _, is_add_rrrci_op_code := this.AddRrrciOpCodes()[op_code]; is_add_rrrci_op_code {
		add_nz_cc := new(cc.AddNzCc)
		add_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = add_nz_cc.Condition()
	} else if _, is_and_rrrci_op_code := this.AndRrrciOpCodes()[op_code]; is_and_rrrci_op_code {
		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_asr_rrrci_op_code := this.AsrRrrciOpCodes()[op_code]; is_asr_rrrci_op_code {
		shift_nz_cc := new(cc.ShiftNzCc)
		shift_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = shift_nz_cc.Condition()
	} else if _, is_mul_rrrci_op_code := this.MulRrrciOpCodes()[op_code]; is_mul_rrrci_op_code {
		mul_nz_cc := new(cc.MulNzCc)
		mul_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = mul_nz_cc.Condition()
	} else if _, is_rsub_rrrci_op_code := this.RsubRrrciOpCodes()[op_code]; is_rsub_rrrci_op_code {
		sub_nz_cc := new(cc.SubNzCc)
		sub_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_nz_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitZri(op_code OpCode, ra *reg_descriptor.SrcRegDescriptor, imm int64) {
	if _, found := this.RriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRI
	this.ra = ra

	if _, is_add_rri_op_code := this.AddRriOpCodes()[op_code]; is_add_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 32, imm)
	} else if _, is_asr_rri_op_code := this.AsrRriOpCodes()[op_code]; is_asr_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_call_rri_op_code := this.CallRriOpCodes()[op_code]; is_call_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 28, imm)
	} else {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}
}

func (this *Instruction) InitZric(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
) {
	if _, found := this.RricOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRIC
	this.ra = ra

	if _, is_add_rric_op_code := this.AddRricOpCodes()[op_code]; is_add_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 27, imm)

		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_asr_rric_op_code := this.AsrRricOpCodes()[op_code]; is_asr_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)

		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_sub_rric_op_code := this.SubRricOpCodes()[op_code]; is_sub_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 27, imm)

		ext_sub_setcc := new(cc.ExtSubSetCc)
		ext_sub_setcc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = ext_sub_setcc.Condition()
	} else {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}
}

func (this *Instruction) InitZrici(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RriciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRICI
	this.ra = ra

	if _, is_add_rrici_op_code := this.AddRriciOpCodes()[op_code]; is_add_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 11, imm)

		add_nz_cc := new(cc.AddNzCc)
		add_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = add_nz_cc.Condition()
	} else if _, is_and_rrici_op_code := this.AndRriciOpCodes()[op_code]; is_and_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 11, imm)

		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_sub_rrici_op_code := this.SubRriciOpCodes()[op_code]; is_sub_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 11, imm)

		sub_nz_cc := new(cc.SubNzCc)
		sub_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_nz_cc.Condition()
	} else if _, is_asr_rrici_op_code := this.AsrRriciOpCodes()[op_code]; is_asr_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)

		imm_shift_nz_cc := new(cc.ImmShiftNzCc)
		imm_shift_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = imm_shift_nz_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitZrif(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
) {
	if _, found := this.RrifOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRIF
	this.ra = ra

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 27, imm)

	false_cc := new(cc.LogSetCc)
	false_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = false_cc.Condition()
}

func (this *Instruction) InitZrr(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
) {
	if _, found := this.RrrOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRR
	this.ra = ra
	this.rb = rb
}

func (this *Instruction) InitZrrc(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RrrcOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRRC
	this.ra = ra
	this.rb = rb

	if _, is_add_rrrc_op_code := this.AddRrrcOpCodes()[op_code]; is_add_rrrc_op_code {
		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_rsub_rrrc_op_code := this.RsubRrrcOpCodes()[op_code]; is_rsub_rrrc_op_code {
		sub_set_cc := new(cc.SubSetCc)
		sub_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_set_cc.Condition()
	} else if _, is_sub_rrrc_op_code := this.SubRrrcOpCodes()[op_code]; is_sub_rrrc_op_code {
		ext_sub_set_cc := new(cc.ExtSubSetCc)
		ext_sub_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = ext_sub_set_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}
}

func (this *Instruction) InitZrrci(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrrciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRRCI
	this.ra = ra
	this.rb = rb

	if _, is_add_rrrci_op_code := this.AddRrrciOpCodes()[op_code]; is_add_rrrci_op_code {
		add_nz_cc := new(cc.AddNzCc)
		add_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = add_nz_cc.Condition()
	} else if _, is_and_rrrci_op_code := this.AndRrrciOpCodes()[op_code]; is_and_rrrci_op_code {
		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_asr_rrrci_op_code := this.AsrRrrciOpCodes()[op_code]; is_asr_rrrci_op_code {
		shift_nz_cc := new(cc.ShiftNzCc)
		shift_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = shift_nz_cc.Condition()
	} else if _, is_mul_rrrci_op_code := this.MulRrrciOpCodes()[op_code]; is_mul_rrrci_op_code {
		mul_nz_cc := new(cc.MulNzCc)
		mul_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = mul_nz_cc.Condition()
	} else if _, is_rsub_rrrci_op_code := this.RsubRrrciOpCodes()[op_code]; is_rsub_rrrci_op_code {
		sub_nz_cc := new(cc.SubNzCc)
		sub_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_nz_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitSRri(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
) {
	if _, found := this.RriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	if suffix != S_RRI && suffix != U_RRI {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra

	if _, is_add_rri_op_code := this.AddRriOpCodes()[op_code]; is_add_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 32, imm)
	} else if _, is_asr_rri_op_code := this.AsrRriOpCodes()[op_code]; is_asr_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_call_rri_op_code := this.CallRriOpCodes()[op_code]; is_call_rri_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else {
		err := errors.New("op code is not a valid S_RRI nor U_RRI op code")
		panic(err)
	}
}

func (this *Instruction) InitSRric(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
) {
	if _, found := this.RricOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	if suffix != S_RRIC && suffix != U_RRIC {
		err := errors.New("suffix is not S_RRIC nor U_RRIC")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra

	if _, is_add_rric_op_code := this.AddRricOpCodes()[op_code]; is_add_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)

		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_asr_rric_op_code := this.AsrRricOpCodes()[op_code]; is_asr_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)

		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_sub_rric_op_code := this.SubRricOpCodes()[op_code]; is_sub_rric_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)

		ext_sub_set_cc := new(cc.ExtSubSetCc)
		ext_sub_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = ext_sub_set_cc.Condition()
	} else {
		err := errors.New("op code is not a valid S_RRI nor U_RRI op code")
		panic(err)
	}
}

func (this *Instruction) InitSRrici(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RriciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	if suffix != S_RRICI && suffix != U_RRICI {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra

	if _, is_add_rrici_op_code := this.AddRriciOpCodes()[op_code]; is_add_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)

		add_nz_cc := new(cc.AddNzCc)
		add_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = add_nz_cc.Condition()
	} else if _, is_and_rrici_op_code := this.AndRriciOpCodes()[op_code]; is_and_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)

		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_asr_rrici_op_code := this.AsrRriciOpCodes()[op_code]; is_asr_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)

		imm_shift_nz_cc := new(cc.ImmShiftNzCc)
		imm_shift_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = imm_shift_nz_cc.Condition()
	} else if _, is_sub_rrici_op_code := this.SubRriciOpCodes()[op_code]; is_sub_rrici_op_code {
		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)

		sub_nz_cc := new(cc.SubNzCc)
		sub_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_nz_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitSRrif(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
) {
	if _, found := this.RrifOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	if suffix != S_RRIF && suffix != U_RRIF {
		err := errors.New("suffix is not S_RRIF nor U_RRIF")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	false_cc := new(cc.FalseCc)
	false_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = false_cc.Condition()
}

func (this *Instruction) InitSRrr(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
) {
	if _, found := this.RrrOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	if suffix != S_RRR && suffix != U_RRR {
		err := errors.New("suffix is not S_RRR nor U_RRR")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb
}

func (this *Instruction) InitSRrrc(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RrrcOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	if suffix != S_RRRC && suffix != U_RRRC {
		err := errors.New("suffix is not S_RRRC nor U_RRRC")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb

	if _, is_add_rrrc_op_code := this.AddRrrcOpCodes()[op_code]; is_add_rrrc_op_code {
		log_set_cc := new(cc.LogSetCc)
		log_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_set_cc.Condition()
	} else if _, is_rsub_rrrc_op_code := this.RsubRrrcOpCodes()[op_code]; is_rsub_rrrc_op_code {
		sub_set_cc := new(cc.SubSetCc)
		sub_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_set_cc.Condition()
	} else if _, is_sub_rrrc_op_code := this.SubRrrcOpCodes()[op_code]; is_sub_rrrc_op_code {
		ext_sub_set_cc := new(cc.ExtSubSetCc)
		ext_sub_set_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = ext_sub_set_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}
}

func (this *Instruction) InitSRrrci(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrrciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	if suffix != S_RRRCI && suffix != U_RRRCI {
		err := errors.New("suffix is not S_RRRCI nor U_RRRCI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb

	if _, is_add_rrrci_op_code := this.AddRrrciOpCodes()[op_code]; is_add_rrrci_op_code {
		add_nz_cc := new(cc.AddNzCc)
		add_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = add_nz_cc.Condition()
	} else if _, is_and_rrrci_op_code := this.AndRrrciOpCodes()[op_code]; is_and_rrrci_op_code {
		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_asr_rrrci_op_code := this.AsrRrrciOpCodes()[op_code]; is_asr_rrrci_op_code {
		shift_nz_cc := new(cc.ShiftNzCc)
		shift_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = shift_nz_cc.Condition()
	} else if _, is_mul_rrrci_op_code := this.MulRrrciOpCodes()[op_code]; is_mul_rrrci_op_code {
		mul_nz_cc := new(cc.MulNzCc)
		mul_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = mul_nz_cc.Condition()
	} else if _, is_rsub_rrrci_op_code := this.RsubRrrciOpCodes()[op_code]; is_rsub_rrrci_op_code {
		sub_nz_cc := new(cc.SubNzCc)
		sub_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = sub_nz_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitRr(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
) {
	if _, found := this.RrOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RR
	this.rc = rc
	this.ra = ra
}

func (this *Instruction) InitRrc(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RrcOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRC
	this.rc = rc
	this.ra = ra

	log_set_cc := new(cc.LogSetCc)
	log_set_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = log_set_cc.Condition()
}

func (this *Instruction) InitRrci(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRCI
	this.rc = rc
	this.ra = ra

	if _, is_cao_rrci_op_code := this.CaoRrciOpCodes()[op_code]; is_cao_rrci_op_code {
		count_nz_cc := new(cc.CountNzCc)
		count_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = count_nz_cc.Condition()
	} else if _, is_extsb_rrci_op_code := this.ExtsbRrciOpCodes()[op_code]; is_extsb_rrci_op_code {
		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_time_cfg_rrci_op_code := this.TimeCfgRrciOpCodes()[op_code]; is_time_cfg_rrci_op_code {
		true_cc := new(cc.TrueCc)
		true_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = true_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitZr(op_code OpCode, ra *reg_descriptor.SrcRegDescriptor) {
	if _, found := this.RrOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZR
	this.ra = ra
}

func (this *Instruction) InitZrc(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RrcOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRC
	this.ra = ra

	log_set_cc := new(cc.LogSetCc)
	log_set_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = log_set_cc.Condition()
}

func (this *Instruction) InitZrci(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRCI
	this.ra = ra

	if _, is_cao_rrci_op_code := this.CaoRrciOpCodes()[op_code]; is_cao_rrci_op_code {
		count_nz_cc := new(cc.CountNzCc)
		count_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = count_nz_cc.Condition()
	} else if _, is_extsb_rrci_op_code := this.ExtsbRrciOpCodes()[op_code]; is_extsb_rrci_op_code {
		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_time_cfg_rrci_op_code := this.TimeCfgRrciOpCodes()[op_code]; is_time_cfg_rrci_op_code {
		true_cc := new(cc.TrueCc)
		true_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = true_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitSRr(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
) {
	if _, found := this.RrOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	if suffix != S_RR && suffix != U_RR {
		err := errors.New("suffix is not S_RR nor U_RR")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
}

func (this *Instruction) InitSRrc(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RrcOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	if suffix != S_RRC && suffix != U_RRC {
		err := errors.New("suffix is not S_RRC nor U_RRC")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra

	log_set_cc := new(cc.LogSetCc)
	log_set_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = log_set_cc.Condition()
}

func (this *Instruction) InitSRrci(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	if suffix != S_RRCI && suffix != U_RRCI {
		err := errors.New("suffix is not S_RRCI nor U_RRCI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra

	if _, is_cao_rrci_op_code := this.CaoRrciOpCodes()[op_code]; is_cao_rrci_op_code {
		count_nz_cc := new(cc.CountNzCc)
		count_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = count_nz_cc.Condition()
	} else if _, is_extsb_rrci_op_code := this.ExtsbRrciOpCodes()[op_code]; is_extsb_rrci_op_code {
		log_nz_cc := new(cc.LogNzCc)
		log_nz_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = log_nz_cc.Condition()
	} else if _, is_time_cfg_rrci_op_code := this.TimeCfgRrciOpCodes()[op_code]; is_time_cfg_rrci_op_code {
		true_cc := new(cc.TrueCc)
		true_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = true_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitDrdici(
	op_code OpCode,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	db *reg_descriptor.PairRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.DrdiciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid DRDICI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = DRDICI
	this.dc = dc
	this.ra = ra
	this.db = db

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	if _, is_div_step_drdici_op_code := this.DivStepDrdiciOpCodes()[op_code]; is_div_step_drdici_op_code {
		div_cc := new(cc.DivCc)
		div_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = div_cc.Condition()
	} else if _, is_mul_step_drdici_op_code := this.MulStepDrdiciOpCodes()[op_code]; is_mul_step_drdici_op_code {
		boot_cc := new(cc.BootCc)
		boot_cc.Init(condition)

		this.condition = new(cc.Condition)
		*this.condition = boot_cc.Condition()
	} else {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitRrri(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	imm int64,
) {
	if _, found := this.RrriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRRI
	this.rc = rc
	this.ra = ra
	this.rb = rb

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)
}

func (this *Instruction) InitRrrici(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrriciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RRRICI
	this.rc = rc
	this.ra = ra
	this.rb = rb

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	div_nz_cc := new(cc.DivNzCc)
	div_nz_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = div_nz_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitZrri(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	imm int64,
) {
	if _, found := this.RrriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRRI
	this.ra = ra
	this.rb = rb

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)
}

func (this *Instruction) InitZrrici(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrriciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZRRICI
	this.ra = ra
	this.rb = rb

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	div_nz_cc := new(cc.DivNzCc)
	div_nz_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = div_nz_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitSRrri(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	imm int64,
) {
	if _, found := this.RrriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	if suffix != S_RRRI && suffix != U_RRRI {
		err := errors.New("suffix is not S_RRRI nor U_RRRI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)
}

func (this *Instruction) InitSRrrici(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	imm int64,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RrriciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	if suffix != S_RRRICI && suffix != U_RRRICI {
		err := errors.New("suffix is not S_RRRICI nor U_RRRICI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	div_nz_cc := new(cc.DivNzCc)
	div_nz_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = div_nz_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitRir(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	imm int64,
	ra *reg_descriptor.SrcRegDescriptor,
) {
	if _, found := this.RirOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RIR
	this.rc = rc

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 32, imm)

	this.ra = ra
}

func (this *Instruction) InitRirc(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	imm int64,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RircOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RIRC
	this.rc = rc

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	this.ra = ra

	sub_set_cc := new(cc.SubSetCc)
	sub_set_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = sub_set_cc.Condition()
}

func (this *Instruction) InitRirci(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	imm int64,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RirciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RIRCI
	this.rc = rc

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 8, imm)

	this.ra = ra

	sub_nz_cc := new(cc.SubNzCc)
	sub_nz_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = sub_nz_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitZir(op_code OpCode, imm int64, ra *reg_descriptor.SrcRegDescriptor) {
	if _, found := this.RirOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZIR

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 32, imm)

	this.ra = ra
}

func (this *Instruction) InitZirc(
	op_code OpCode,
	imm int64,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RircOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZIRC

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 27, imm)

	this.ra = ra

	sub_set_cc := new(cc.SubSetCc)
	sub_set_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = sub_set_cc.Condition()
}

func (this *Instruction) InitZirci(
	op_code OpCode,
	imm int64,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RirciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZIRCI

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 11, imm)

	this.ra = ra

	sub_nz_cc := new(cc.SubNzCc)
	sub_nz_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = sub_nz_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitSRirc(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	imm int64,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
) {
	if _, found := this.RircOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	if suffix != S_RIRC && suffix != U_RIRC {
		err := errors.New("suffix is not S_RIRC nor U_RIRC")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	this.ra = ra

	sub_set_cc := new(cc.SubSetCc)
	sub_set_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = sub_set_cc.Condition()
}

func (this *Instruction) InitSRirci(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	imm int64,
	ra *reg_descriptor.SrcRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RirciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	if suffix != S_RIRCI && suffix != U_RIRCI {
		err := errors.New("suffix is not S_RIRCI nor U_RIRCI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 8, imm)

	this.ra = ra

	sub_nz_cc := new(cc.SubNzCc)
	sub_nz_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = sub_nz_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitR(op_code OpCode, rc *reg_descriptor.GpRegDescriptor) {
	if _, found := this.ROpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid R op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = R
	this.rc = rc
}

func (this *Instruction) InitRci(
	op_code OpCode,
	rc *reg_descriptor.GpRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = RCI
	this.rc = rc

	true_cc := new(cc.TrueCc)
	true_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = true_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitZ(op_code OpCode) {
	if _, found := this.ROpCodes()[op_code]; !found && op_code != NOP {
		err := errors.New("op code is not a valid R op code nor NOP")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = Z
}

func (this *Instruction) InitZci(
	op_code OpCode,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ZCI

	true_cc := new(cc.TrueCc)
	true_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = true_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitSR(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
) {
	if _, found := this.ROpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid R op code")
		panic(err)
	}

	if suffix != S_R && suffix != U_R {
		err := errors.New("suffix is not S_R nor U_R")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
}

func (this *Instruction) InitSRci(
	op_code OpCode,
	suffix Suffix,
	dc *reg_descriptor.PairRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.RciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	if suffix != S_RCI && suffix != U_RCI {
		err := errors.New("suffix is not S_RCI nor U_RCI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc

	true_cc := new(cc.TrueCc)
	true_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = true_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitCi(op_code OpCode, condition cc.Condition, pc int64) {
	if _, found := this.CiOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid CI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = CI

	boot_cc := new(cc.BootCc)
	boot_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = boot_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitI(op_code OpCode, imm int64) {
	if _, found := this.IOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid I op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = I

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)
}

func (this *Instruction) InitDdci(
	op_code OpCode,
	dc *reg_descriptor.PairRegDescriptor,
	db *reg_descriptor.PairRegDescriptor,
	condition cc.Condition,
	pc int64,
) {
	if _, found := this.DdciOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid DDCI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = DDCI

	this.dc = dc
	this.db = db

	true_false_cc := new(cc.TrueFalseCc)
	true_false_cc.Init(condition)

	this.condition = new(cc.Condition)
	*this.condition = true_false_cc.Condition()

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, config_loader.AddressWidth(), pc)
}

func (this *Instruction) InitErri(
	op_code OpCode,
	endian Endian,
	rc *reg_descriptor.GpRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	off int64,
) {
	if _, found := this.ErriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid ERRI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ERRI

	this.endian = new(Endian)
	*this.endian = endian

	this.rc = rc
	this.ra = ra

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)
}

func (this *Instruction) InitSErri(
	op_code OpCode,
	suffix Suffix,
	endian Endian,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	off int64,
) {
	if _, found := this.ErriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid ERRI op code")
		panic(err)
	}

	if suffix != S_ERRI && suffix != U_ERRI {
		err := errors.New("suffix is not S_ERRI nor U_ERRI")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix

	this.endian = new(Endian)
	*this.endian = endian

	this.dc = dc
	this.ra = ra

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)
}

func (this *Instruction) InitEdri(
	op_code OpCode,
	endian Endian,
	dc *reg_descriptor.PairRegDescriptor,
	ra *reg_descriptor.SrcRegDescriptor,
	off int64,
) {
	if _, found := this.EdriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid EDRI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = EDRI

	this.endian = new(Endian)
	*this.endian = endian

	this.dc = dc
	this.ra = ra

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)
}

func (this *Instruction) InitErii(
	op_code OpCode,
	endian Endian,
	ra *reg_descriptor.SrcRegDescriptor,
	off int64,
	imm int64,
) {
	if _, found := this.EriiOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid ERII op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ERII

	this.endian = new(Endian)
	*this.endian = endian

	this.ra = ra

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 16, imm)
}

func (this *Instruction) InitErir(
	op_code OpCode,
	endian Endian,
	ra *reg_descriptor.SrcRegDescriptor,
	off int64,
	rb *reg_descriptor.SrcRegDescriptor,
) {
	if _, found := this.ErirOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid ERIR op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ERIR

	this.endian = new(Endian)
	*this.endian = endian

	this.ra = ra

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)

	this.rb = rb
}

func (this *Instruction) InitErid(
	op_code OpCode,
	endian Endian,
	ra *reg_descriptor.SrcRegDescriptor,
	off int64,
	db *reg_descriptor.PairRegDescriptor,
) {
	if _, found := this.EridOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid ERID op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = ERID

	this.endian = new(Endian)
	*this.endian = endian

	this.ra = ra

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)

	this.db = db
}

func (this *Instruction) InitDmaRri(
	op_code OpCode,
	ra *reg_descriptor.SrcRegDescriptor,
	rb *reg_descriptor.SrcRegDescriptor,
	imm int64,
) {
	if _, found := this.DmaRriOpCodes()[op_code]; !found {
		err := errors.New("op code is not a valid DMA RRI op code")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = DMA_RRI

	this.ra = ra
	this.rb = rb

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 8, imm)
}

func (this *Instruction) OpCode() OpCode {
	return this.op_code
}

func (this *Instruction) Suffix() Suffix {
	return this.suffix
}

func (this *Instruction) Rc() *reg_descriptor.GpRegDescriptor {
	return this.rc
}

func (this *Instruction) Ra() *reg_descriptor.SrcRegDescriptor {
	return this.ra
}

func (this *Instruction) Rb() *reg_descriptor.SrcRegDescriptor {
	return this.rb
}

func (this *Instruction) Dc() *reg_descriptor.PairRegDescriptor {
	return this.dc
}

func (this *Instruction) Db() *reg_descriptor.PairRegDescriptor {
	return this.db
}

func (this *Instruction) Condition() cc.Condition {
	if this.condition == nil {
		err := errors.New("condition == nil")
		panic(err)
	}

	return *this.condition
}

func (this *Instruction) Imm() *abi.Immediate {
	return this.imm
}

func (this *Instruction) Off() *abi.Immediate {
	return this.off
}

func (this *Instruction) Pc() *abi.Immediate {
	return this.pc
}

func (this *Instruction) Endian() Endian {
	if this.endian == nil {
		err := errors.New("endian == nil")
		panic(err)
	}

	return *this.endian
}

func (this *Instruction) Encode() *encoding.ByteStream {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	word_ := new(abi.Word)
	word_.Init(config_loader.IramDataWidth())

	this.EncodeOpCode(word_)
	this.EncodeSuffix(word_)

	if this.suffix == RICI {
		this.EncodeRici(word_)
	} else if this.suffix == RRI {
		this.EncodeRri(word_)
	} else if this.suffix == RRIC {
		this.EncodeRric(word_)
	} else if this.suffix == RRICI {
		this.EncodeRrici(word_)
	} else if this.suffix == RRIF {
		this.EncodeRrif(word_)
	} else if this.suffix == RRR {
		this.EncodeRrr(word_)
	} else if this.suffix == RRRC {
		this.EncodeRrrc(word_)
	} else if this.suffix == RRRCI {
		this.EncodeRrrci(word_)
	} else if this.suffix == ZRI {
		this.EncodeZri(word_)
	} else if this.suffix == ZRIC {
		this.EncodeZric(word_)
	} else if this.suffix == ZRICI {
		this.EncodeZrici(word_)
	} else if this.suffix == ZRIF {
		this.EncodeZrif(word_)
	} else if this.suffix == ZRR {
		this.EncodeZrr(word_)
	} else if this.suffix == ZRRC {
		this.EncodeZrrc(word_)
	} else if this.suffix == ZRRCI {
		this.EncodeZrrci(word_)
	} else if this.suffix == S_RRI || this.suffix == U_RRI {
		this.EncodeSRri(word_)
	} else if this.suffix == S_RRIC || this.suffix == U_RRIC {
		this.EncodeSRric(word_)
	} else if this.suffix == S_RRICI || this.suffix == U_RRICI {
		this.EncodeSRrici(word_)
	} else if this.suffix == S_RRIF || this.suffix == U_RRIF {
		this.EncodeSRrif(word_)
	} else if this.suffix == S_RRR || this.suffix == U_RRR {
		this.EncodeSRrr(word_)
	} else if this.suffix == S_RRRC || this.suffix == U_RRRC {
		this.EncodeSRrrc(word_)
	} else if this.suffix == S_RRRCI || this.suffix == U_RRRCI {
		this.EncodeSRrrci(word_)
	} else if this.suffix == RR {
		this.EncodeRr(word_)
	} else if this.suffix == RRC {
		this.EncodeRrc(word_)
	} else if this.suffix == RRCI {
		this.EncodeRrci(word_)
	} else if this.suffix == ZR {
		this.EncodeZr(word_)
	} else if this.suffix == ZRC {
		this.EncodeZrc(word_)
	} else if this.suffix == ZRCI {
		this.EncodeZrci(word_)
	} else if this.suffix == S_RR || this.suffix == U_RR {
		this.EncodeSRr(word_)
	} else if this.suffix == S_RRC || this.suffix == U_RRC {
		this.EncodeSRrc(word_)
	} else if this.suffix == S_RRCI || this.suffix == U_RRCI {
		this.EncodeSRrci(word_)
	} else if this.suffix == DRDICI {
		this.EncodeDrdici(word_)
	} else if this.suffix == RRRI {
		this.EncodeRrri(word_)
	} else if this.suffix == RRRICI {
		this.EncodeRrrici(word_)
	} else if this.suffix == ZRRI {
		this.EncodeZrri(word_)
	} else if this.suffix == ZRRICI {
		this.EncodeZrrici(word_)
	} else if this.suffix == S_RRRI || this.suffix == U_RRRI {
		this.EncodeSRrri(word_)
	} else if this.suffix == S_RRRICI || this.suffix == U_RRRICI {
		this.EncodeSRrrici(word_)
	} else if this.suffix == RIR {
		this.EncodeRir(word_)
	} else if this.suffix == RIRC {
		this.EncodeRirc(word_)
	} else if this.suffix == RIRCI {
		this.EncodeRirci(word_)
	} else if this.suffix == ZIR {
		this.EncodeZir(word_)
	} else if this.suffix == ZIRC {
		this.EncodeZirc(word_)
	} else if this.suffix == ZIRCI {
		this.EncodeZirci(word_)
	} else if this.suffix == S_RIRC || this.suffix == U_RIRC {
		this.EncodeSRirc(word_)
	} else if this.suffix == S_RIRCI || this.suffix == U_RIRCI {
		this.EncodeSRirci(word_)
	} else if this.suffix == R {
		this.EncodeR(word_)
	} else if this.suffix == RCI {
		this.EncodeRci(word_)
	} else if this.suffix == Z {
		this.EncodeZ(word_)
	} else if this.suffix == ZCI {
		this.EncodeZci(word_)
	} else if this.suffix == S_R || this.suffix == U_R {
		this.EncodeSR(word_)
	} else if this.suffix == S_RCI || this.suffix == U_RCI {
		this.EncodeSRci(word_)
	} else if this.suffix == CI {
		this.EncodeCi(word_)
	} else if this.suffix == I {
		this.EncodeI(word_)
	} else if this.suffix == DDCI {
		this.EncodeDdci(word_)
	} else if this.suffix == ERRI {
		this.EncodeErri(word_)
	} else if this.suffix == S_ERRI || this.suffix == U_ERRI {
		this.EncodeSErri(word_)
	} else if this.suffix == EDRI {
		this.EncodeEdri(word_)
	} else if this.suffix == ERII {
		this.EncodeErii(word_)
	} else if this.suffix == ERIR {
		this.EncodeErir(word_)
	} else if this.suffix == ERID {
		this.EncodeErid(word_)
	} else if this.suffix == DMA_RRI {
		this.EncodeDmaRri(word_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}

	return word_.ToByteStream()
}

func (this *Instruction) EncodeRici(word_ *abi.Word) {
	if _, found := this.RiciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RICI op code")
		panic(err)
	}

	if this.suffix != RICI {
		err := errors.New("suffix is not RICI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.pc.Width()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeRri(word_ *abi.Word) {
	if _, found := this.RriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	if this.suffix != RRI {
		err := errors.New("suffix is not RRI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeRric(word_ *abi.Word) {
	if _, found := this.RricOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	if this.suffix != RRIC {
		err := errors.New("suffix is not RRIC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeRrici(word_ *abi.Word) {
	if _, found := this.RriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	if this.suffix != RRICI {
		err := errors.New("suffix is not RRICI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeRrif(word_ *abi.Word) {
	if _, found := this.RrifOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	if this.suffix != RRIF {
		err := errors.New("suffix is not RRIF")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeRrr(word_ *abi.Word) {
	if _, found := this.RrrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	if this.suffix != RRR {
		err := errors.New("suffix is not RRR")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)
}

func (this *Instruction) EncodeRrrc(word_ *abi.Word) {
	if _, found := this.RrrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	if this.suffix != RRRC {
		err := errors.New("suffix is not RRRC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeRrrci(word_ *abi.Word) {
	if _, found := this.RrrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	if this.suffix != RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeZri(word_ *abi.Word) {
	if _, found := this.RriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	if this.suffix != ZRI {
		err := errors.New("suffix is not ZRI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeZric(word_ *abi.Word) {
	if _, found := this.RricOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	if this.suffix != ZRIC {
		err := errors.New("suffix is not ZRIC")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeZrici(word_ *abi.Word) {
	if _, found := this.RriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	if this.suffix != ZRICI {
		err := errors.New("suffix is not ZRICI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeZrif(word_ *abi.Word) {
	if _, found := this.RrifOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	if this.suffix != ZRIF {
		err := errors.New("suffix is not RRIF")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeZrr(word_ *abi.Word) {
	if _, found := this.RrrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	if this.suffix != ZRR {
		err := errors.New("suffix is not RRR")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)
}

func (this *Instruction) EncodeZrrc(word_ *abi.Word) {
	if _, found := this.RrrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	if this.suffix != ZRRC {
		err := errors.New("suffix is not ZRRC")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeZrrci(word_ *abi.Word) {
	if _, found := this.RrrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	if this.suffix != ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeSRri(word_ *abi.Word) {
	if _, found := this.RriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	if this.suffix != S_RRI && this.suffix != U_RRI {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeSRric(word_ *abi.Word) {
	if _, found := this.RricOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	if this.suffix != S_RRIC && this.suffix != U_RRIC {
		err := errors.New("suffix is not S_RRIC nor U_RRIC")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeSRrici(word_ *abi.Word) {
	if _, found := this.RriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	if this.suffix != S_RRICI && this.suffix != U_RRICI {
		err := errors.New("suffix is not RRICI nor U_RRICI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeSRrif(word_ *abi.Word) {
	if _, found := this.RrifOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	if this.suffix != S_RRIF && this.suffix != U_RRIF {
		err := errors.New("suffix is not S_RRIF nor U_RRIF")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	imm_begin := ra_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeSRrr(word_ *abi.Word) {
	if _, found := this.RrrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	if this.suffix != RRR {
		err := errors.New("suffix is not S_RRR nor U_RRR")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)
}

func (this *Instruction) EncodeSRrrc(word_ *abi.Word) {
	if _, found := this.RrrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	if this.suffix != S_RRRC && this.suffix != U_RRRC {
		err := errors.New("suffix is not S_RRRC nor U_RRRC")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeSRrrci(word_ *abi.Word) {
	if _, found := this.RrrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	if this.suffix != S_RRRCI && this.suffix != U_RRRCI {
		err := errors.New("suffix is not S_RRRCI nor U_RRRCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeRr(word_ *abi.Word) {
	if _, found := this.RrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	if this.suffix != RR {
		err := errors.New("suffix is not RR")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)
}

func (this *Instruction) EncodeRrc(word_ *abi.Word) {
	if _, found := this.RrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	if this.suffix != RRC {
		err := errors.New("suffix is not RRC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeRrci(word_ *abi.Word) {
	if _, found := this.RrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	if this.suffix != RRCI {
		err := errors.New("suffix is not RRCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeZr(word_ *abi.Word) {
	if _, found := this.RrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	if this.suffix != ZR {
		err := errors.New("suffix is not ZR")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)
}

func (this *Instruction) EncodeZrc(word_ *abi.Word) {
	if _, found := this.RrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	if this.suffix != ZRC {
		err := errors.New("suffix is not RRC")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeZrci(word_ *abi.Word) {
	if _, found := this.RrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	if this.suffix != ZRCI {
		err := errors.New("suffix is not ZRCI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeSRr(word_ *abi.Word) {
	if _, found := this.RrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	if this.suffix != S_RR && this.suffix != U_RR {
		err := errors.New("suffix is not S_RR nor U_RR")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)
}

func (this *Instruction) EncodeSRrc(word_ *abi.Word) {
	if _, found := this.RrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	if this.suffix != S_RRC && this.suffix != U_RRC {
		err := errors.New("suffix is not S_RRC nor U_RRC")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeSRrci(word_ *abi.Word) {
	if _, found := this.RrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	if this.suffix != S_RRCI && this.suffix != U_RRCI {
		err := errors.New("suffix is not S_RRCI nor U_RRCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeDrdici(word_ *abi.Word) {
	if _, found := this.DrdiciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid DRDICI op code")
		panic(err)
	}

	if this.suffix != DRDICI {
		err := errors.New("suffix is not DRDICI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	db_begin := ra_end
	db_end := db_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, db_begin, db_end, this.db)

	imm_begin := db_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeRrri(word_ *abi.Word) {
	if _, found := this.RrriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	if this.suffix != RRRI {
		err := errors.New("suffix is not RRRI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	imm_begin := rb_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeRrrici(word_ *abi.Word) {
	if _, found := this.RrriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	if this.suffix != RRRICI {
		err := errors.New("suffix is not RRRICI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	imm_begin := rb_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeZrri(word_ *abi.Word) {
	if _, found := this.RrriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	if this.suffix != ZRRI {
		err := errors.New("suffix is not ZRRI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	imm_begin := rb_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeZrrici(word_ *abi.Word) {
	if _, found := this.RrriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	if this.suffix != ZRRICI {
		err := errors.New("suffix is not ZRRICI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	imm_begin := rb_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeSRrri(word_ *abi.Word) {
	if _, found := this.RrriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	if this.suffix != S_RRRI && this.suffix != U_RRRI {
		err := errors.New("suffix is not S_RRRI nor U_RRRI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	imm_begin := rb_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeSRrrici(word_ *abi.Word) {
	if _, found := this.RrriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	if this.suffix != S_RRRICI && this.suffix != U_RRRICI {
		err := errors.New("suffix is not S_RRRICI nor U_RRRICI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	imm_begin := rb_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeRir(word_ *abi.Word) {
	if _, found := this.RirOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	}

	if this.suffix != RIR {
		err := errors.New("suffix is not RIR")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	imm_begin := rc_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)
}

func (this *Instruction) EncodeRirc(word_ *abi.Word) {
	if _, found := this.RircOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	if this.suffix != RIRC {
		err := errors.New("suffix is not RIRC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	imm_begin := rc_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeRirci(word_ *abi.Word) {
	if _, found := this.RirciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	if this.suffix != RIRCI {
		err := errors.New("suffix is not RIRCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	imm_begin := rc_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeZir(word_ *abi.Word) {
	if _, found := this.RirOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	}

	if this.suffix != ZIR {
		err := errors.New("suffix is not ZIR")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)
}

func (this *Instruction) EncodeZirc(word_ *abi.Word) {
	if _, found := this.RircOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	if this.suffix != ZIRC {
		err := errors.New("suffix is not ZIRC")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeZirci(word_ *abi.Word) {
	if _, found := this.RirciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	if this.suffix != ZIRCI {
		err := errors.New("suffix is not ZIRCI")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeSRirc(word_ *abi.Word) {
	if _, found := this.RircOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	if this.suffix != S_RIRC && this.suffix != U_RIRC {
		err := errors.New("suffix is not S_RIRC nor U_RIRC")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	imm_begin := dc_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)
}

func (this *Instruction) EncodeSRirci(word_ *abi.Word) {
	if _, found := this.RirciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	if this.suffix != S_RIRCI && this.suffix != U_RIRCI {
		err := errors.New("suffix is not S_RIRCI nor U_RIRCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	imm_begin := dc_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeR(word_ *abi.Word) {
	if _, found := this.ROpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid R op code")
		panic(err)
	}

	if this.suffix != R {
		err := errors.New("suffix is not R")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)
}

func (this *Instruction) EncodeRci(word_ *abi.Word) {
	if _, found := this.RciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	if this.suffix != RCI {
		err := errors.New("suffix is not RCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	condition_begin := rc_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeZ(word_ *abi.Word) {
	if _, found := this.ROpCodes()[this.op_code]; !found && this.op_code != NOP {
		err := errors.New("op code is not a valid R op code nor NOP")
		panic(err)
	}

	if this.suffix != Z {
		err := errors.New("suffix is not Z")
		panic(err)
	}
}

func (this *Instruction) EncodeZci(word_ *abi.Word) {
	if _, found := this.RciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	if this.suffix != ZCI {
		err := errors.New("suffix is not ZCI")
		panic(err)
	}

	condition_begin := this.SuffixEnd()
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeSR(word_ *abi.Word) {
	if _, found := this.ROpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid R op code")
		panic(err)
	}

	if this.suffix != S_R && this.suffix != U_R {
		err := errors.New("suffix is not S_R nor U_R")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)
}

func (this *Instruction) EncodeSRci(word_ *abi.Word) {
	if _, found := this.RciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	if this.suffix != S_RCI && this.suffix != U_RCI {
		err := errors.New("suffix is not S_RCI nor U_RCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	condition_begin := dc_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeCi(word_ *abi.Word) {
	if _, found := this.CiOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid CI op code")
		panic(err)
	}

	if this.suffix != CI {
		err := errors.New("suffix is not CI")
		panic(err)
	}

	condition_begin := this.SuffixEnd()
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeI(word_ *abi.Word) {
	if _, found := this.IOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid I op code")
		panic(err)
	}

	if this.suffix != I {
		err := errors.New("suffix is not I")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeDdci(word_ *abi.Word) {
	if _, found := this.DdciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid DDCI op code")
		panic(err)
	}

	if this.suffix != DDCI {
		err := errors.New("suffix is not DDCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	db_begin := dc_end
	db_end := db_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, db_begin, db_end, this.db)

	condition_begin := db_end
	condition_end := condition_begin + this.ConditionWidth()
	this.EncodeCondition(word_, condition_begin, condition_end, *this.condition)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	this.EncodePc(word_, pc_begin, pc_end, this.pc.Value())
}

func (this *Instruction) EncodeErri(word_ *abi.Word) {
	if _, found := this.ErriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERRI op code")
		panic(err)
	}

	if this.suffix != ERRI {
		err := errors.New("suffix is not ERRI")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()
	this.EncodeEndian(word_, endian_begin, endian_end, *this.endian)

	rc_begin := endian_end
	rc_end := rc_begin + this.RegisterWidth()
	this.EncodeGpRegDescriptor(word_, rc_begin, rc_end, this.rc)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	off_begin := ra_end
	off_end := off_begin + this.off.Width()
	this.EncodeOff(word_, off_begin, off_end, this.off.Value())
}

func (this *Instruction) EncodeSErri(word_ *abi.Word) {
	if _, found := this.ErriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERRI op code")
		panic(err)
	}

	if this.suffix != S_ERRI && this.suffix != U_ERRI {
		err := errors.New("suffix is not S_ERRI nor U_ERRI")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()
	this.EncodeEndian(word_, endian_begin, endian_end, *this.endian)

	dc_begin := endian_end
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	off_begin := ra_end
	off_end := off_begin + this.off.Width()
	this.EncodeOff(word_, off_begin, off_end, this.off.Value())
}

func (this *Instruction) EncodeEdri(word_ *abi.Word) {
	if _, found := this.EdriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid EDRI op code")
		panic(err)
	}

	if this.suffix != EDRI {
		err := errors.New("suffix is not EDRI")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()
	this.EncodeEndian(word_, endian_begin, endian_end, *this.endian)

	dc_begin := endian_end
	dc_end := dc_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, dc_begin, dc_end, this.dc)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	off_begin := ra_end
	off_end := off_begin + this.off.Width()
	this.EncodeOff(word_, off_begin, off_end, this.off.Value())
}

func (this *Instruction) EncodeErii(word_ *abi.Word) {
	if _, found := this.EriiOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERII op code")
		panic(err)
	}

	if this.suffix != ERII {
		err := errors.New("suffix is not ERII")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()
	this.EncodeEndian(word_, endian_begin, endian_end, *this.endian)

	ra_begin := endian_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	off_begin := ra_end
	off_end := off_begin + this.off.Width()
	this.EncodeOff(word_, off_begin, off_end, this.off.Value())

	imm_begin := off_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeErir(word_ *abi.Word) {
	if _, found := this.ErirOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERIR op code")
		panic(err)
	}

	if this.suffix != ERIR {
		err := errors.New("suffix is not ERIR")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()
	this.EncodeEndian(word_, endian_begin, endian_end, *this.endian)

	ra_begin := endian_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	off_begin := ra_end
	off_end := off_begin + this.off.Width()
	this.EncodeOff(word_, off_begin, off_end, this.off.Value())

	rb_begin := off_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)
}

func (this *Instruction) EncodeErid(word_ *abi.Word) {
	if _, found := this.EridOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERID op code")
		panic(err)
	}

	if this.suffix != ERID {
		err := errors.New("suffix is not ERID")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()
	this.EncodeEndian(word_, endian_begin, endian_end, *this.endian)

	ra_begin := endian_end
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	off_begin := ra_end
	off_end := off_begin + this.off.Width()
	this.EncodeOff(word_, off_begin, off_end, this.off.Value())

	db_begin := off_end
	db_end := db_begin + this.RegisterWidth()
	this.EncodePairRegDescriptor(word_, db_begin, db_end, this.db)
}

func (this *Instruction) EncodeDmaRri(word_ *abi.Word) {
	if _, found := this.DmaRriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid DMA_RRI op code")
		panic(err)
	}

	if this.suffix != DMA_RRI {
		err := errors.New("suffix is not DMA_RRI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, ra_begin, ra_end, this.ra)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.EncodeSrcRegDescriptor(word_, rb_begin, rb_end, this.rb)

	imm_begin := rb_end
	imm_end := imm_begin + this.imm.Width()
	this.EncodeImm(word_, imm_begin, imm_end, this.imm.Value())
}

func (this *Instruction) EncodeOpCode(word_ *abi.Word) {
	word_.SetBitSlice(this.OpCodeBegin(), this.OpCodeEnd(), int64(this.op_code))
}

func (this *Instruction) EncodeSuffix(word_ *abi.Word) {
	word_.SetBitSlice(this.SuffixBegin(), this.SuffixEnd(), int64(this.suffix))
}

func (this *Instruction) EncodeGpRegDescriptor(
	word_ *abi.Word,
	begin int,
	end int,
	gp_reg_descriptor *reg_descriptor.GpRegDescriptor,
) {
	word_.SetBitSlice(begin, end, int64(gp_reg_descriptor.Index()))
}

func (this *Instruction) EncodeSrcRegDescriptor(
	word_ *abi.Word,
	begin int,
	end int,
	src_reg_descriptor *reg_descriptor.SrcRegDescriptor,
) {
	if src_reg_descriptor.IsGpRegDescriptor() {
		word_.SetBitSlice(begin, end, int64(src_reg_descriptor.GpRegDescriptor().Index()))
	} else if src_reg_descriptor.IsSpRegDescriptor() {
		config_loader := new(misc.ConfigLoader)
		config_loader.Init()

		index := config_loader.NumGpRegisters() + int(*src_reg_descriptor.SpRegDescriptor())
		word_.SetBitSlice(begin, end, int64(index))
	} else {
		err := errors.New("sp reg descriptor is corrupted")
		panic(err)
	}
}

func (this *Instruction) EncodePairRegDescriptor(
	word_ *abi.Word,
	begin int,
	end int,
	pair_reg_descriptor *reg_descriptor.PairRegDescriptor,
) {
	word_.SetBitSlice(begin, end, int64(pair_reg_descriptor.Index()))
}

func (this *Instruction) EncodeImm(word_ *abi.Word, begin int, end int, value int64) {
	word_.SetBitSlice(begin, end, value)
}

func (this *Instruction) EncodeCondition(
	word_ *abi.Word,
	begin int,
	end int,
	condition cc.Condition,
) {
	word_.SetBitSlice(begin, end, int64(condition))
}

func (this *Instruction) EncodePc(word_ *abi.Word, begin int, end int, pc int64) {
	this.EncodeImm(word_, begin, end, pc)
}

func (this *Instruction) EncodeEndian(word_ *abi.Word, begin int, end int, endian Endian) {
	word_.SetBitSlice(begin, end, int64(endian))
}

func (this *Instruction) EncodeOff(word_ *abi.Word, begin int, end int, value int64) {
	this.EncodeImm(word_, begin, end, value)
}

func (this *Instruction) Decode(byte_stream *encoding.ByteStream) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	word_ := new(abi.Word)
	word_.Init(config_loader.IramDataWidth())
	word_.FromByteStream(byte_stream)

	this.op_code = this.DecodeOpCode(word_)
	this.suffix = this.DecodeSuffix(word_)

	if this.suffix == RICI {
		this.DecodeRici(word_)
	} else if this.suffix == RRI {
		this.DecodeRri(word_)
	} else if this.suffix == RRIC {
		this.DecodeRric(word_)
	} else if this.suffix == RRICI {
		this.DecodeRrici(word_)
	} else if this.suffix == RRIF {
		this.DecodeRrif(word_)
	} else if this.suffix == RRR {
		this.DecodeRrr(word_)
	} else if this.suffix == RRRC {
		this.DecodeRrrc(word_)
	} else if this.suffix == RRRCI {
		this.DecodeRrrci(word_)
	} else if this.suffix == ZRI {
		this.DecodeZri(word_)
	} else if this.suffix == ZRIC {
		this.DecodeZric(word_)
	} else if this.suffix == ZRICI {
		this.DecodeZrici(word_)
	} else if this.suffix == ZRIF {
		this.DecodeZrif(word_)
	} else if this.suffix == ZRR {
		this.DecodeZrr(word_)
	} else if this.suffix == ZRRC {
		this.DecodeZrrc(word_)
	} else if this.suffix == ZRRCI {
		this.DecodeZrrci(word_)
	} else if this.suffix == S_RRI || this.suffix == U_RRI {
		this.DecodeSRri(word_)
	} else if this.suffix == S_RRIC || this.suffix == U_RRIC {
		this.DecodeSRric(word_)
	} else if this.suffix == S_RRICI || this.suffix == U_RRICI {
		this.DecodeSRrici(word_)
	} else if this.suffix == S_RRIF || this.suffix == U_RRIF {
		this.DecodeSRrif(word_)
	} else if this.suffix == S_RRR || this.suffix == U_RRR {
		this.DecodeSRrr(word_)
	} else if this.suffix == S_RRRC || this.suffix == U_RRRC {
		this.DecodeSRrrc(word_)
	} else if this.suffix == S_RRRCI || this.suffix == U_RRRCI {
		this.DecodeSRrrci(word_)
	} else if this.suffix == RR {
		this.DecodeRr(word_)
	} else if this.suffix == RRC {
		this.DecodeRrc(word_)
	} else if this.suffix == RRCI {
		this.DecodeRrci(word_)
	} else if this.suffix == ZR {
		this.DecodeZr(word_)
	} else if this.suffix == ZRC {
		this.DecodeZrc(word_)
	} else if this.suffix == ZRCI {
		this.DecodeZrci(word_)
	} else if this.suffix == S_RR || this.suffix == U_RR {
		this.DecodeSRr(word_)
	} else if this.suffix == S_RRC || this.suffix == U_RRC {
		this.DecodeSRrc(word_)
	} else if this.suffix == S_RRCI || this.suffix == U_RRCI {
		this.DecodeSRrci(word_)
	} else if this.suffix == DRDICI {
		this.DecodeDrdici(word_)
	} else if this.suffix == RRRI {
		this.DecodeRrri(word_)
	} else if this.suffix == RRRICI {
		this.DecodeRrrici(word_)
	} else if this.suffix == ZRRI {
		this.DecodeZrri(word_)
	} else if this.suffix == ZRRICI {
		this.DecodeZrrici(word_)
	} else if this.suffix == S_RRRI || this.suffix == U_RRRI {
		this.DecodeSRrri(word_)
	} else if this.suffix == S_RRRICI || this.suffix == U_RRRICI {
		this.DecodeSRrrici(word_)
	} else if this.suffix == RIR {
		this.DecodeRir(word_)
	} else if this.suffix == RIRC {
		this.DecodeRirc(word_)
	} else if this.suffix == RIRCI {
		this.DecodeRirci(word_)
	} else if this.suffix == ZIR {
		this.DecodeZir(word_)
	} else if this.suffix == ZIRC {
		this.DecodeZirc(word_)
	} else if this.suffix == ZIRCI {
		this.DecodeZirci(word_)
	} else if this.suffix == S_RIRC || this.suffix == U_RIRC {
		this.DecodeSRirc(word_)
	} else if this.suffix == S_RIRCI || this.suffix == U_RIRCI {
		this.DecodeSRirci(word_)
	} else if this.suffix == R {
		this.DecodeR(word_)
	} else if this.suffix == RCI {
		this.DecodeRci(word_)
	} else if this.suffix == Z {
		this.DecodeZ(word_)
	} else if this.suffix == ZCI {
		this.DecodeZci(word_)
	} else if this.suffix == S_R || this.suffix == U_R {
		this.DecodeSR(word_)
	} else if this.suffix == S_RCI || this.suffix == U_RCI {
		this.DecodeSRci(word_)
	} else if this.suffix == CI {
		this.DecodeCi(word_)
	} else if this.suffix == I {
		this.DecodeI(word_)
	} else if this.suffix == DDCI {
		this.DecodeDdci(word_)
	} else if this.suffix == ERRI {
		this.DecodeErri(word_)
	} else if this.suffix == S_ERRI || this.suffix == U_ERRI {
		this.DecodeSErri(word_)
	} else if this.suffix == EDRI {
		this.DecodeEdri(word_)
	} else if this.suffix == ERII {
		this.DecodeErii(word_)
	} else if this.suffix == ERIR {
		this.DecodeErir(word_)
	} else if this.suffix == ERID {
		this.DecodeErid(word_)
	} else if this.suffix == DMA_RRI {
		this.DecodeDmaRri(word_)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Instruction) DecodeRici(word_ *abi.Word) {
	if _, found := this.RiciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RICI op code")
		panic(err)
	}

	if this.suffix != RICI {
		err := errors.New("suffix is not RICI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	imm_end := imm_begin + 16
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 16, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()
	condition := this.DecodeCondition(word_, condition_begin, condition_end)

	this.condition = new(cc.Condition)
	*this.condition = condition

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.SIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeRri(word_ *abi.Word) {
	if _, found := this.RriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	if this.suffix != RRI {
		err := errors.New("suffix is not RRI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	if _, is_add_rri_op_code := this.AddRriOpCodes()[this.op_code]; is_add_rri_op_code {
		imm_end := imm_begin + 32
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 32, imm)
	} else if _, is_asr_rri_op_code := this.AsrRriOpCodes()[this.op_code]; is_asr_rri_op_code {
		imm_end := imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_call_rri_op_code := this.CallRriOpCodes()[this.op_code]; is_call_rri_op_code {
		imm_end := imm_begin + 24
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}
}

func (this *Instruction) DecodeRric(word_ *abi.Word) {
	if _, found := this.RricOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	if this.suffix != RRIC {
		err := errors.New("suffix is not RRIC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	var imm_end int
	if _, is_add_rric_op_code := this.AddRricOpCodes()[this.op_code]; is_add_rric_op_code {
		imm_end = imm_begin + 24
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else if _, is_asr_rric_op_code := this.AsrRricOpCodes()[this.op_code]; is_asr_rric_op_code {
		imm_end = imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_sub_rric_op_code := this.SubRricOpCodes()[this.op_code]; is_sub_rric_op_code {
		imm_end = imm_begin + 24
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	condition_begin := imm_end
	condition_end := imm_end + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeRrici(word_ *abi.Word) {
	if _, found := this.RriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	if this.suffix != RRICI {
		err := errors.New("suffix is not RRICI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	var imm_end int
	if _, is_add_rrici_op_code := this.AddRriciOpCodes()[this.op_code]; is_add_rrici_op_code {
		imm_end = imm_begin + 8
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)
	} else if _, is_asr_rrici_op_code := this.AsrRriciOpCodes()[this.op_code]; is_asr_rrici_op_code {
		imm_end = imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_and_rrici_op_code := this.AndRriciOpCodes()[this.op_code]; is_and_rrici_op_code {
		imm_end = imm_begin + 8
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)
	} else if _, is_sub_rrici_op_code := this.SubRriciOpCodes()[this.op_code]; is_sub_rrici_op_code {
		imm_end = imm_begin + 8
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)
	} else {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	condition_begin := imm_end
	condition_end := imm_end + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeRrif(word_ *abi.Word) {
	if _, found := this.RrifOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	if this.suffix != RRIF {
		err := errors.New("suffix is not RRIF")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	imm_end := imm_begin + 24
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeRrr(word_ *abi.Word) {
	if _, found := this.RrrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	if this.suffix != RRR {
		err := errors.New("suffix is not RRR")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)
}

func (this *Instruction) DecodeRrrc(word_ *abi.Word) {
	if _, found := this.RrrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	if this.suffix != RRRC {
		err := errors.New("suffix is not RRRC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeRrrci(word_ *abi.Word) {
	if _, found := this.RrrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	if this.suffix != RRRCI {
		err := errors.New("suffix is not RRRCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeZri(word_ *abi.Word) {
	if _, found := this.RriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	if this.suffix != ZRI {
		err := errors.New("suffix is not ZRI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	if _, is_add_rri_op_code := this.AddRriOpCodes()[this.op_code]; is_add_rri_op_code {
		imm_end := imm_begin + 32
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 32, imm)
	} else if _, is_asr_rri_op_code := this.AsrRriOpCodes()[this.op_code]; is_asr_rri_op_code {
		imm_end := imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_call_rri_op_code := this.CallRriOpCodes()[this.op_code]; is_call_rri_op_code {
		imm_end := imm_begin + 28
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 28, imm)
	} else {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}
}

func (this *Instruction) DecodeZric(word_ *abi.Word) {
	if _, found := this.RricOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	if this.suffix != ZRIC {
		err := errors.New("suffix is not RRIC")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	var imm_end int
	if _, is_add_rric_op_code := this.AddRricOpCodes()[this.op_code]; is_add_rric_op_code {
		imm_end = imm_begin + 27
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 27, imm)
	} else if _, is_asr_rric_op_code := this.AsrRricOpCodes()[this.op_code]; is_asr_rric_op_code {
		imm_end = imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_sub_rric_op_code := this.SubRricOpCodes()[this.op_code]; is_sub_rric_op_code {
		imm_end = imm_begin + 27
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 27, imm)
	} else {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	condition_begin := imm_end
	condition_end := imm_end + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeZrici(word_ *abi.Word) {
	if _, found := this.RriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	if this.suffix != ZRICI {
		err := errors.New("suffix is not ZRICI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	var imm_end int
	if _, is_add_rrici_op_code := this.AddRriciOpCodes()[this.op_code]; is_add_rrici_op_code {
		imm_end = imm_begin + 11
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 11, imm)
	} else if _, is_asr_rrici_op_code := this.AsrRriciOpCodes()[this.op_code]; is_asr_rrici_op_code {
		imm_end = imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_and_rrici_op_code := this.AndRriciOpCodes()[this.op_code]; is_and_rrici_op_code {
		imm_end = imm_begin + 11
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 11, imm)
	} else if _, is_sub_rrici_op_code := this.SubRriciOpCodes()[this.op_code]; is_sub_rrici_op_code {
		imm_end = imm_begin + 11
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 11, imm)
	} else {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	condition_begin := imm_end
	condition_end := imm_end + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeZrif(word_ *abi.Word) {
	if _, found := this.RrifOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	if this.suffix != ZRIF {
		err := errors.New("suffix is not ZRIF")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	imm_end := imm_begin + 27
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 27, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeZrr(word_ *abi.Word) {
	if _, found := this.RrrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	if this.suffix != ZRR {
		err := errors.New("suffix is not ZRR")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)
}

func (this *Instruction) DecodeZrrc(word_ *abi.Word) {
	if _, found := this.RrrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	if this.suffix != ZRRC {
		err := errors.New("suffix is not ZRRC")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeZrrci(word_ *abi.Word) {
	if _, found := this.RrrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	if this.suffix != ZRRCI {
		err := errors.New("suffix is not ZRRCI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeSRri(word_ *abi.Word) {
	if _, found := this.RriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}

	if this.suffix != S_RRI && this.suffix != U_RRI {
		err := errors.New("suffix is not S_RRI nor U_RRI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	if _, is_add_rri_op_code := this.AddRriOpCodes()[this.op_code]; is_add_rri_op_code {
		imm_end := imm_begin + 32
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 32, imm)
	} else if _, is_asr_rri_op_code := this.AsrRriOpCodes()[this.op_code]; is_asr_rri_op_code {
		imm_end := imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_call_rri_op_code := this.CallRriOpCodes()[this.op_code]; is_call_rri_op_code {
		imm_end := imm_begin + 24
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else {
		err := errors.New("op code is not a valid RRI op code")
		panic(err)
	}
}

func (this *Instruction) DecodeSRric(word_ *abi.Word) {
	if _, found := this.RricOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	if this.suffix != S_RRIC && this.suffix != U_RRIC {
		err := errors.New("suffix is not S_RRIC nor U_RRIC")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	var imm_end int
	if _, is_add_rric_op_code := this.AddRricOpCodes()[this.op_code]; is_add_rric_op_code {
		imm_end = imm_begin + 24
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else if _, is_asr_rric_op_code := this.AsrRricOpCodes()[this.op_code]; is_asr_rric_op_code {
		imm_end = imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_sub_rric_op_code := this.SubRricOpCodes()[this.op_code]; is_sub_rric_op_code {
		imm_end = imm_begin + 24
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 24, imm)
	} else {
		err := errors.New("op code is not a valid RRIC op code")
		panic(err)
	}

	condition_begin := imm_end
	condition_end := imm_end + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeSRrici(word_ *abi.Word) {
	if _, found := this.RriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	if this.suffix != S_RRICI && this.suffix != U_RRICI {
		err := errors.New("suffix is not S_RRICI nor U_RRICI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	var imm_end int
	if _, is_add_rrici_op_code := this.AddRriciOpCodes()[this.op_code]; is_add_rrici_op_code {
		imm_end = imm_begin + 8
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)
	} else if _, is_asr_rrici_op_code := this.AsrRriciOpCodes()[this.op_code]; is_asr_rrici_op_code {
		imm_end = imm_begin + 5
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.UNSIGNED, 5, imm)
	} else if _, is_and_rrici_op_code := this.AndRriciOpCodes()[this.op_code]; is_and_rrici_op_code {
		imm_end = imm_begin + 8
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)
	} else if _, is_sub_rrici_op_code := this.SubRriciOpCodes()[this.op_code]; is_sub_rrici_op_code {
		imm_end = imm_begin + 8
		imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

		this.imm = new(abi.Immediate)
		this.imm.Init(abi.SIGNED, 8, imm)
	} else {
		err := errors.New("op code is not a valid RRICI op code")
		panic(err)
	}

	condition_begin := imm_end
	condition_end := imm_end + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeSRrif(word_ *abi.Word) {
	if _, found := this.RrifOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRIF op code")
		panic(err)
	}

	if this.suffix != S_RRIF && this.suffix != U_RRIF {
		err := errors.New("suffix is not S_RRIF nor U_RRIF")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	imm_begin := ra_end
	imm_end := imm_begin + 24
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeSRrr(word_ *abi.Word) {
	if _, found := this.RrrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRR op code")
		panic(err)
	}

	if this.suffix != S_RRR {
		err := errors.New("suffix is not S_RRR nor U_RRR")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)
}

func (this *Instruction) DecodeSRrrc(word_ *abi.Word) {
	if _, found := this.RrrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRC op code")
		panic(err)
	}

	if this.suffix != S_RRRC && this.suffix != U_RRRC {
		err := errors.New("suffix is not S_RRRC nor U_RRRC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeSRrrci(word_ *abi.Word) {
	if _, found := this.RrrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRCI op code")
		panic(err)
	}

	if this.suffix != S_RRRCI && this.suffix != U_RRRCI {
		err := errors.New("suffix is not S_RRRCI nor U_RRRCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	condition_begin := rb_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeRr(word_ *abi.Word) {
	if _, found := this.RrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	if this.suffix != RR {
		err := errors.New("suffix is not RR")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)
}

func (this *Instruction) DecodeRrc(word_ *abi.Word) {
	if _, found := this.RrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	if this.suffix != RRC {
		err := errors.New("suffix is not RRC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeRrci(word_ *abi.Word) {
	if _, found := this.RrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	if this.suffix != RRCI {
		err := errors.New("suffix is not RRCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeZr(word_ *abi.Word) {
	if _, found := this.RrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	if this.suffix != ZR {
		err := errors.New("suffix is not ZR")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)
}

func (this *Instruction) DecodeZrc(word_ *abi.Word) {
	if _, found := this.RrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	if this.suffix != ZRC {
		err := errors.New("suffix is not ZRC")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeZrci(word_ *abi.Word) {
	if _, found := this.RrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	if this.suffix != ZRCI {
		err := errors.New("suffix is not ZRCI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeSRr(word_ *abi.Word) {
	if _, found := this.RrOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RR op code")
		panic(err)
	}

	if this.suffix != S_RR && this.suffix != U_RR {
		err := errors.New("suffix is not S_RR nor U_RR")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)
}

func (this *Instruction) DecodeSRrc(word_ *abi.Word) {
	if _, found := this.RrcOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRC op code")
		panic(err)
	}

	if this.suffix != S_RRC && this.suffix != U_RRC {
		err := errors.New("suffix is not S_RRC nor U_RRC")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeSRrci(word_ *abi.Word) {
	if _, found := this.RrciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRCI op code")
		panic(err)
	}

	if this.suffix != S_RRCI && this.suffix != U_RRCI {
		err := errors.New("suffix is not S_RRCI nor U_RRCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodePc(word_, pc_begin, pc_end)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeDrdici(word_ *abi.Word) {
	if _, found := this.DrdiciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid DRDICI op code")
		panic(err)
	}

	if this.suffix != DRDICI {
		err := errors.New("suffix is not DRDICI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	db_begin := ra_end
	db_end := db_begin + this.RegisterWidth()
	this.db = this.DecodePairRegDescriptor(word_, db_begin, db_end)

	imm_begin := db_end
	imm_end := imm_begin + 5
	imm := this.DecodePc(word_, imm_begin, imm_end)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeRrri(word_ *abi.Word) {
	if _, found := this.RrriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	if this.suffix != RRRI {
		err := errors.New("suffix is not RRRI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	imm_begin := rb_end
	imm_end := imm_begin + 5
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)
}

func (this *Instruction) DecodeRrrici(word_ *abi.Word) {
	if _, found := this.RrriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	if this.suffix != RRRICI {
		err := errors.New("suffix is not RRRICI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	imm_begin := rb_end
	imm_end := imm_begin + 5
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeZrri(word_ *abi.Word) {
	if _, found := this.RrriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	if this.suffix != ZRRI {
		err := errors.New("suffix is not ZRRI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	imm_begin := rb_end
	imm_end := imm_begin + 5
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)
}

func (this *Instruction) DecodeZrrici(word_ *abi.Word) {
	if _, found := this.RrriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	if this.suffix != ZRRICI {
		err := errors.New("suffix is not ZRRICI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	imm_begin := rb_end
	imm_end := imm_begin + 5
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeSRrri(word_ *abi.Word) {
	if _, found := this.RrriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRI op code")
		panic(err)
	}

	if this.suffix != S_RRRI && this.suffix != U_RRRI {
		err := errors.New("suffix is not S_RRRI nor U_RRRI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	imm_begin := rb_end
	imm_end := imm_begin + 5
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)
}

func (this *Instruction) DecodeSRrrici(word_ *abi.Word) {
	if _, found := this.RrriciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RRRICI op code")
		panic(err)
	}

	if this.suffix != S_RRRICI && this.suffix != U_RRRICI {
		err := errors.New("suffix is not S_RRRICI nor U_RRRICI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	imm_begin := rb_end
	imm_end := imm_begin + 5
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 5, imm)

	condition_begin := imm_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeRir(word_ *abi.Word) {
	if _, found := this.RirOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	}

	if this.suffix != RIR {
		err := errors.New("suffix is not RIR")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	imm_begin := rc_end
	imm_end := imm_begin + 32
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 32, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)
}

func (this *Instruction) DecodeRirc(word_ *abi.Word) {
	if _, found := this.RircOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	if this.suffix != RIRC {
		err := errors.New("suffix is not RIRC")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	imm_begin := rc_end
	imm_end := imm_begin + 24
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeRirci(word_ *abi.Word) {
	if _, found := this.RirciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	if this.suffix != RIRCI {
		err := errors.New("suffix is not RIRCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	imm_begin := rc_end
	imm_end := imm_begin + 8
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 8, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeZir(word_ *abi.Word) {
	if _, found := this.RirOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIR op code")
		panic(err)
	}

	if this.suffix != ZIR {
		err := errors.New("suffix is not ZIR")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + 32
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 32, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)
}

func (this *Instruction) DecodeZirc(word_ *abi.Word) {
	if _, found := this.RircOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	if this.suffix != ZIRC {
		err := errors.New("suffix is not ZIRC")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + 24
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeZirci(word_ *abi.Word) {
	if _, found := this.RirciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	if this.suffix != ZIRCI {
		err := errors.New("suffix is not ZIRCI")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + 8
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 8, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeSRirc(word_ *abi.Word) {
	if _, found := this.RircOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRC op code")
		panic(err)
	}

	if this.suffix != S_RIRC && this.suffix != U_RIRC {
		err := errors.New("suffix is not S_RIRC nor U_RIRC")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	imm_begin := dc_end
	imm_end := imm_begin + 24
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)
}

func (this *Instruction) DecodeSRirci(word_ *abi.Word) {
	if _, found := this.RirciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RIRCI op code")
		panic(err)
	}

	if this.suffix != S_RIRCI && this.suffix != U_RIRCI {
		err := errors.New("suffix is not S_RIRCI nor U_RIRCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	imm_begin := dc_end
	imm_end := imm_begin + 8
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 8, imm)

	ra_begin := imm_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	condition_begin := ra_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeR(word_ *abi.Word) {
	if _, found := this.ROpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid R op code")
		panic(err)
	}

	if this.suffix != R {
		err := errors.New("suffix is not R")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)
}

func (this *Instruction) DecodeRci(word_ *abi.Word) {
	if _, found := this.RciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	if this.suffix != RCI {
		err := errors.New("suffix is not RCI")
		panic(err)
	}

	rc_begin := this.SuffixEnd()
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	condition_begin := rc_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeZ(word_ *abi.Word) {
	if _, found := this.ROpCodes()[this.op_code]; !found && this.op_code != NOP {
		err := errors.New("op code is not a valid R op code nor NOP")
		panic(err)
	}

	if this.suffix != Z {
		err := errors.New("suffix is not R")
		panic(err)
	}
}

func (this *Instruction) DecodeZci(word_ *abi.Word) {
	if _, found := this.RciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	if this.suffix != ZCI {
		err := errors.New("suffix is not ZCI")
		panic(err)
	}

	condition_begin := this.SuffixEnd()
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeSR(word_ *abi.Word) {
	if _, found := this.ROpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid R op code")
		panic(err)
	}

	if this.suffix != S_R && this.suffix != U_R {
		err := errors.New("suffix is not S_R nor U_R")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)
}

func (this *Instruction) DecodeSRci(word_ *abi.Word) {
	if _, found := this.RciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid RCI op code")
		panic(err)
	}

	if this.suffix != S_RCI && this.suffix != U_RCI {
		err := errors.New("suffix is not S_RCI nor U_RCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	condition_begin := dc_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeCi(word_ *abi.Word) {
	if _, found := this.CiOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid CI op code")
		panic(err)
	}

	if this.suffix != CI {
		err := errors.New("suffix is not CI")
		panic(err)
	}

	condition_begin := this.SuffixEnd()
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeI(word_ *abi.Word) {
	if _, found := this.IOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid I op code")
		panic(err)
	}

	if this.suffix != I {
		err := errors.New("suffix is not I")
		panic(err)
	}

	imm_begin := this.SuffixEnd()
	imm_end := imm_begin + 24
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 24, imm)
}

func (this *Instruction) DecodeDdci(word_ *abi.Word) {
	if _, found := this.DdciOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid DDCI op code")
		panic(err)
	}

	if this.suffix != DDCI {
		err := errors.New("suffix is not DDCI")
		panic(err)
	}

	dc_begin := this.SuffixEnd()
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	db_begin := dc_end
	db_end := db_begin + this.RegisterWidth()
	this.db = this.DecodePairRegDescriptor(word_, db_begin, db_end)

	condition_begin := db_end
	condition_end := condition_begin + this.ConditionWidth()

	this.condition = new(cc.Condition)
	*this.condition = this.DecodeCondition(word_, condition_begin, condition_end)

	pc_begin := condition_end
	pc_end := pc_begin + this.PcWidth()
	pc := this.DecodeImm(word_, pc_begin, pc_end, abi.UNSIGNED)

	this.pc = new(abi.Immediate)
	this.pc.Init(abi.UNSIGNED, this.PcWidth(), pc)
}

func (this *Instruction) DecodeErri(word_ *abi.Word) {
	if _, found := this.ErriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERRI op code")
		panic(err)
	}

	if this.suffix != ERRI {
		err := errors.New("suffix is not ERRI")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()

	this.endian = new(Endian)
	*this.endian = this.DecodeEndian(word_, endian_begin, endian_end)

	rc_begin := endian_end
	rc_end := rc_begin + this.RegisterWidth()
	this.rc = this.DecodeGpRegDescriptor(word_, rc_begin, rc_end)

	ra_begin := rc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	off_begin := ra_end
	off_end := off_begin + 24
	off := this.DecodeOff(word_, off_begin, off_end, abi.SIGNED)

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)
}

func (this *Instruction) DecodeSErri(word_ *abi.Word) {
	if _, found := this.ErriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERRI op code")
		panic(err)
	}

	if this.suffix != S_ERRI && this.suffix != U_ERRI {
		err := errors.New("suffix is not S_ERRI nor U_ERRI")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()

	this.endian = new(Endian)
	*this.endian = this.DecodeEndian(word_, endian_begin, endian_end)

	dc_begin := endian_end
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	off_begin := ra_end
	off_end := off_begin + 24
	off := this.DecodeOff(word_, off_begin, off_end, abi.SIGNED)

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)
}

func (this *Instruction) DecodeEdri(word_ *abi.Word) {
	if _, found := this.EdriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid EDRI op code")
		panic(err)
	}

	if this.suffix != EDRI {
		err := errors.New("suffix is not EDRI")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()

	this.endian = new(Endian)
	*this.endian = this.DecodeEndian(word_, endian_begin, endian_end)

	dc_begin := endian_end
	dc_end := dc_begin + this.RegisterWidth()
	this.dc = this.DecodePairRegDescriptor(word_, dc_begin, dc_end)

	ra_begin := dc_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	off_begin := ra_end
	off_end := off_begin + 24
	off := this.DecodeOff(word_, off_begin, off_end, abi.SIGNED)

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)
}

func (this *Instruction) DecodeErii(word_ *abi.Word) {
	if _, found := this.EriiOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERII op code")
		panic(err)
	}

	if this.suffix != ERII {
		err := errors.New("suffix is not ERII")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()

	this.endian = new(Endian)
	*this.endian = this.DecodeEndian(word_, endian_begin, endian_end)

	ra_begin := endian_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	off_begin := ra_end
	off_end := off_begin + 24
	off := this.DecodeOff(word_, off_begin, off_end, abi.SIGNED)

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)

	imm_begin := off_end
	imm_end := imm_begin + 16
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.SIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.SIGNED, 16, imm)
}

func (this *Instruction) DecodeErir(word_ *abi.Word) {
	if _, found := this.ErirOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERIR op code")
		panic(err)
	}

	if this.suffix != ERIR {
		err := errors.New("suffix is not ERIR")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()

	this.endian = new(Endian)
	*this.endian = this.DecodeEndian(word_, endian_begin, endian_end)

	ra_begin := endian_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	off_begin := ra_end
	off_end := off_begin + 24
	off := this.DecodeOff(word_, off_begin, off_end, abi.SIGNED)

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)

	rb_begin := off_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)
}

func (this *Instruction) DecodeErid(word_ *abi.Word) {
	if _, found := this.EridOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid ERID op code")
		panic(err)
	}

	if this.suffix != ERID {
		err := errors.New("suffix is not ERID")
		panic(err)
	}

	endian_begin := this.SuffixEnd()
	endian_end := endian_begin + this.EndianWidth()

	this.endian = new(Endian)
	*this.endian = this.DecodeEndian(word_, endian_begin, endian_end)

	ra_begin := endian_end
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	off_begin := ra_end
	off_end := off_begin + 24
	off := this.DecodeOff(word_, off_begin, off_end, abi.SIGNED)

	this.off = new(abi.Immediate)
	this.off.Init(abi.SIGNED, 24, off)

	db_begin := off_end
	db_end := db_begin + this.RegisterWidth()
	this.db = this.DecodePairRegDescriptor(word_, db_begin, db_end)
}

func (this *Instruction) DecodeDmaRri(word_ *abi.Word) {
	if _, found := this.DmaRriOpCodes()[this.op_code]; !found {
		err := errors.New("op code is not a valid DMA_RRI op code")
		panic(err)
	}

	if this.suffix != DMA_RRI {
		err := errors.New("suffix is not DMA_RRI")
		panic(err)
	}

	ra_begin := this.SuffixEnd()
	ra_end := ra_begin + this.RegisterWidth()
	this.ra = this.DecodeSrcRegDescriptor(word_, ra_begin, ra_end)

	rb_begin := ra_end
	rb_end := rb_begin + this.RegisterWidth()
	this.rb = this.DecodeSrcRegDescriptor(word_, rb_begin, rb_end)

	imm_begin := rb_end
	imm_end := imm_begin + 8
	imm := this.DecodeImm(word_, imm_begin, imm_end, abi.UNSIGNED)

	this.imm = new(abi.Immediate)
	this.imm.Init(abi.UNSIGNED, 8, imm)
}

func (this *Instruction) DecodeOpCode(word_ *abi.Word) OpCode {
	return OpCode(word_.BitSlice(abi.UNSIGNED, this.OpCodeBegin(), this.OpCodeEnd()))
}

func (this *Instruction) DecodeSuffix(word_ *abi.Word) Suffix {
	return Suffix(word_.BitSlice(abi.UNSIGNED, this.SuffixBegin(), this.SuffixEnd()))
}

func (this *Instruction) DecodeGpRegDescriptor(
	word_ *abi.Word,
	begin int,
	end int,
) *reg_descriptor.GpRegDescriptor {
	index := int(word_.BitSlice(abi.UNSIGNED, begin, end))

	gp_reg_descriptor := new(reg_descriptor.GpRegDescriptor)
	gp_reg_descriptor.Init(index)

	return gp_reg_descriptor
}

func (this *Instruction) DecodeSrcRegDescriptor(
	word_ *abi.Word,
	begin int,
	end int,
) *reg_descriptor.SrcRegDescriptor {
	index := int(word_.BitSlice(abi.UNSIGNED, begin, end))

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	if index < config_loader.NumGpRegisters() {
		gp_reg_descriptor := new(reg_descriptor.GpRegDescriptor)
		gp_reg_descriptor.Init(index)

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitGpRegDescriptor(gp_reg_descriptor)

		return src_reg_descriptor
	} else {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.SpRegDescriptor(index - config_loader.NumGpRegisters())

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)

		return src_reg_descriptor
	}
}

func (this *Instruction) DecodePairRegDescriptor(
	word_ *abi.Word,
	begin int,
	end int,
) *reg_descriptor.PairRegDescriptor {
	index := int(word_.BitSlice(abi.UNSIGNED, begin, end))

	pair_reg_descriptor := new(reg_descriptor.PairRegDescriptor)
	pair_reg_descriptor.Init(index)

	return pair_reg_descriptor
}

func (this *Instruction) DecodeImm(
	word_ *abi.Word,
	begin int,
	end int,
	representation abi.Representation,
) int64 {
	return word_.BitSlice(representation, begin, end)
}

func (this *Instruction) DecodeCondition(word_ *abi.Word, begin int, end int) cc.Condition {
	return cc.Condition(word_.BitSlice(abi.UNSIGNED, begin, end))
}

func (this *Instruction) DecodePc(word_ *abi.Word, begin int, end int) int64 {
	return this.DecodeImm(word_, begin, end, abi.UNSIGNED)
}

func (this *Instruction) DecodeEndian(word_ *abi.Word, begin int, end int) Endian {
	return Endian(word_.BitSlice(abi.UNSIGNED, begin, end))
}

func (this *Instruction) DecodeOff(
	word_ *abi.Word,
	begin int,
	end int,
	representation abi.Representation,
) int64 {
	return this.DecodeImm(word_, begin, end, representation)
}

func (this *Instruction) OpCodeBegin() int {
	return 0
}

func (this *Instruction) OpCodeEnd() int {
	return this.OpCodeBegin() + this.OpCodeWidth()
}

func (this *Instruction) OpCodeWidth() int {
	return int(math.Ceil(math.Log2(1.0 + float64(SDMA))))
}

func (this *Instruction) SuffixBegin() int {
	return this.OpCodeEnd()
}

func (this *Instruction) SuffixEnd() int {
	return this.SuffixBegin() + this.SuffixWidth()
}

func (this *Instruction) SuffixWidth() int {
	return int(math.Ceil(math.Log2(1.0 + float64(DMA_RRI))))
}

func (this *Instruction) RegisterWidth() int {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	return int(
		math.Ceil(
			math.Log2(float64(config_loader.NumGpRegisters()) + 1.0 + float64(reg_descriptor.ID8)),
		),
	)
}

func (this *Instruction) ConditionWidth() int {
	return int(math.Ceil(math.Log2(1.0 + float64(cc.LARGE))))
}

func (this *Instruction) PcWidth() int {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	return config_loader.AddressWidth()
}

func (this *Instruction) EndianWidth() int {
	return int(math.Ceil(math.Log2(1.0 + float64(BIG))))
}

func (this *Instruction) AcquireRiciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ACQUIRE: true,
	}
}

func (this *Instruction) ReleaseRiciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		RELEASE: true,
	}
}

func (this *Instruction) BootRiciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		BOOT:   true,
		RESUME: true,
	}
}

func (this *Instruction) RiciOpCodes() map[OpCode]bool {
	acquire_rici_op_codes := this.AcquireRiciOpCodes()
	release_rici_op_codes := this.ReleaseRiciOpCodes()
	boot_rici_op_codes := this.BootRiciOpCodes()

	rici_op_codes := make(map[OpCode]bool)

	for k, v := range acquire_rici_op_codes {
		rici_op_codes[k] = v
	}

	for k, v := range release_rici_op_codes {
		rici_op_codes[k] = v
	}

	for k, v := range boot_rici_op_codes {
		rici_op_codes[k] = v
	}

	return rici_op_codes
}

func (this *Instruction) AddRriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ADD:  true,
		ADDC: true,
		AND:  true,
		OR:   true,
		XOR:  true,
	}
}

func (this *Instruction) AsrRriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ASR:   true,
		LSL:   true,
		LSL1:  true,
		LSL1X: true,
		LSLX:  true,
		LSR:   true,
		LSR1:  true,
		LSR1X: true,
		LSRX:  true,
		ROL:   true,
		ROR:   true,
	}
}

func (this *Instruction) CallRriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		CALL: true,
	}
}

func (this *Instruction) RriOpCodes() map[OpCode]bool {
	add_rri_op_codes := this.AddRriOpCodes()
	asr_rri_op_codes := this.AsrRriOpCodes()
	call_rri_op_codes := this.CallRriOpCodes()

	rri_op_codes := make(map[OpCode]bool)

	for k, v := range add_rri_op_codes {
		rri_op_codes[k] = v
	}

	for k, v := range asr_rri_op_codes {
		rri_op_codes[k] = v
	}

	for k, v := range call_rri_op_codes {
		rri_op_codes[k] = v
	}

	return rri_op_codes
}

func (this *Instruction) AddRricOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ADD:  true,
		ADDC: true,
		AND:  true,
		ANDN: true,
		NAND: true,
		NOR:  true,
		NXOR: true,
		OR:   true,
		ORN:  true,
		XOR:  true,
		HASH: true,
	}
}

func (this *Instruction) AsrRricOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ASR:   true,
		LSL:   true,
		LSL1:  true,
		LSL1X: true,
		LSLX:  true,
		LSR:   true,
		LSR1:  true,
		LSR1X: true,
		LSRX:  true,
		ROL:   true,
		ROR:   true,
	}
}

func (this *Instruction) SubRricOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SUB:  true,
		SUBC: true,
	}
}

func (this *Instruction) RricOpCodes() map[OpCode]bool {
	add_rric_op_codes := this.AddRricOpCodes()
	asr_rric_op_codes := this.AsrRricOpCodes()
	sub_rric_op_codes := this.SubRricOpCodes()

	rric_op_codes := make(map[OpCode]bool)

	for k, v := range add_rric_op_codes {
		rric_op_codes[k] = v
	}

	for k, v := range asr_rric_op_codes {
		rric_op_codes[k] = v
	}

	for k, v := range sub_rric_op_codes {
		rric_op_codes[k] = v
	}

	return rric_op_codes
}

func (this *Instruction) AddRriciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ADD:  true,
		ADDC: true,
	}
}

func (this *Instruction) AndRriciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		AND:  true,
		ANDN: true,
		NAND: true,
		NOR:  true,
		NXOR: true,
		OR:   true,
		ORN:  true,
		XOR:  true,
		HASH: true,
	}
}

func (this *Instruction) AsrRriciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ASR:   true,
		LSL:   true,
		LSL1:  true,
		LSL1X: true,
		LSLX:  true,
		LSR:   true,
		LSR1:  true,
		LSR1X: true,
		LSRX:  true,
		ROL:   true,
		ROR:   true,
	}
}

func (this *Instruction) SubRriciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SUB:  true,
		SUBC: true,
	}
}

func (this *Instruction) RriciOpCodes() map[OpCode]bool {
	add_rrici_op_codes := this.AddRriciOpCodes()
	and_rrici_op_codes := this.AndRriciOpCodes()
	asr_rrici_op_codes := this.AsrRriciOpCodes()
	sub_rrici_op_codes := this.SubRriciOpCodes()

	rrici_op_codes := make(map[OpCode]bool)

	for k, v := range add_rrici_op_codes {
		rrici_op_codes[k] = v
	}

	for k, v := range and_rrici_op_codes {
		rrici_op_codes[k] = v
	}

	for k, v := range asr_rrici_op_codes {
		rrici_op_codes[k] = v
	}

	for k, v := range sub_rrici_op_codes {
		rrici_op_codes[k] = v
	}

	return rrici_op_codes
}

func (this *Instruction) RrifOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ADD:  true,
		ADDC: true,
		AND:  true,
		ANDN: true,
		NAND: true,
		NOR:  true,
		NXOR: true,
		OR:   true,
		ORN:  true,
		SUB:  true,
		SUBC: true,
		XOR:  true,
		HASH: true,
	}
}

func (this *Instruction) RrrOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ADD:       true,
		ADDC:      true,
		AND:       true,
		ANDN:      true,
		ASR:       true,
		CMPB4:     true,
		LSL:       true,
		LSL1:      true,
		LSL1X:     true,
		LSLX:      true,
		LSR:       true,
		LSR1:      true,
		LSR1X:     true,
		LSRX:      true,
		MUL_SH_SH: true,
		MUL_SH_SL: true,
		MUL_SH_UH: true,
		MUL_SH_UL: true,
		MUL_SL_SH: true,
		MUL_SL_SL: true,
		MUL_SL_UH: true,
		MUL_SL_UL: true,
		MUL_UH_UH: true,
		MUL_UH_UL: true,
		MUL_UL_UH: true,
		MUL_UL_UL: true,
		NAND:      true,
		NOR:       true,
		NXOR:      true,
		OR:        true,
		ORN:       true,
		ROL:       true,
		ROR:       true,
		RSUB:      true,
		RSUBC:     true,
		SUB:       true,
		SUBC:      true,
		XOR:       true,
		HASH:      true,
		CALL:      true,
	}
}

func (this *Instruction) AddRrrcOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ADD:       true,
		ADDC:      true,
		AND:       true,
		ANDN:      true,
		ASR:       true,
		CMPB4:     true,
		LSL:       true,
		LSL1:      true,
		LSL1X:     true,
		LSLX:      true,
		LSR:       true,
		LSR1:      true,
		LSR1X:     true,
		LSRX:      true,
		MUL_SH_SH: true,
		MUL_SH_SL: true,
		MUL_SH_UH: true,
		MUL_SH_UL: true,
		MUL_SL_SH: true,
		MUL_SL_SL: true,
		MUL_SL_UH: true,
		MUL_SL_UL: true,
		MUL_UH_UH: true,
		MUL_UH_UL: true,
		MUL_UL_UH: true,
		MUL_UL_UL: true,
		NAND:      true,
		NOR:       true,
		NXOR:      true,
		ROL:       true,
		ROR:       true,
		OR:        true,
		ORN:       true,
		XOR:       true,
		HASH:      true,
		CALL:      true,
	}
}

func (this *Instruction) RsubRrrcOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		RSUB:  true,
		RSUBC: true,
	}
}

func (this *Instruction) SubRrrcOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SUB:  true,
		SUBC: true,
	}
}

func (this *Instruction) RrrcOpCodes() map[OpCode]bool {
	add_rrrc_op_codes := this.AddRrrcOpCodes()
	rsub_rrrc_op_codes := this.RsubRrrcOpCodes()
	sub_rrrc_op_codes := this.SubRrrcOpCodes()

	rrrc_op_codes := make(map[OpCode]bool)

	for k, v := range add_rrrc_op_codes {
		rrrc_op_codes[k] = v
	}

	for k, v := range rsub_rrrc_op_codes {
		rrrc_op_codes[k] = v
	}

	for k, v := range sub_rrrc_op_codes {
		rrrc_op_codes[k] = v
	}

	return rrrc_op_codes
}

func (this *Instruction) AddRrrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ADD:  true,
		ADDC: true,
	}
}

func (this *Instruction) AndRrrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		AND:  true,
		ANDN: true,
		NAND: true,
		NOR:  true,
		NXOR: true,
		OR:   true,
		ORN:  true,
		XOR:  true,
		HASH: true,
	}
}

func (this *Instruction) AsrRrrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		ASR:   true,
		CMPB4: true,
		LSL:   true,
		LSL1:  true,
		LSL1X: true,
		LSLX:  true,
		LSR:   true,
		LSR1:  true,
		LSR1X: true,
		LSRX:  true,
		ROL:   true,
		ROR:   true,
	}
}

func (this *Instruction) MulRrrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		MUL_SH_SH: true,
		MUL_SH_SL: true,
		MUL_SH_UH: true,
		MUL_SH_UL: true,
		MUL_SL_SH: true,
		MUL_SL_SL: true,
		MUL_SL_UH: true,
		MUL_SL_UL: true,
		MUL_UH_UH: true,
		MUL_UH_UL: true,
		MUL_UL_UH: true,
		MUL_UL_UL: true,
	}
}

func (this *Instruction) RsubRrrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		RSUB:  true,
		RSUBC: true,
		SUB:   true,
		SUBC:  true,
	}
}

func (this *Instruction) RrrciOpCodes() map[OpCode]bool {
	add_rrrci_op_codes := this.AddRrrciOpCodes()
	and_rrrci_op_codes := this.AndRrrciOpCodes()
	asr_rrrci_op_codes := this.AsrRrrciOpCodes()
	mul_rrrci_op_codes := this.MulRrrciOpCodes()
	rsub_rrrci_op_codes := this.RsubRrrciOpCodes()

	rrrci_op_codes := make(map[OpCode]bool)

	for k, v := range add_rrrci_op_codes {
		rrrci_op_codes[k] = v
	}

	for k, v := range and_rrrci_op_codes {
		rrrci_op_codes[k] = v
	}

	for k, v := range asr_rrrci_op_codes {
		rrrci_op_codes[k] = v
	}

	for k, v := range mul_rrrci_op_codes {
		rrrci_op_codes[k] = v
	}

	for k, v := range rsub_rrrci_op_codes {
		rrrci_op_codes[k] = v
	}

	return rrrci_op_codes
}

func (this *Instruction) RrOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		CAO:      true,
		CLO:      true,
		CLS:      true,
		CLZ:      true,
		EXTSB:    true,
		EXTSH:    true,
		EXTUB:    true,
		EXTUH:    true,
		SATS:     true,
		TIME_CFG: true,
	}
}

func (this *Instruction) RrcOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		CAO:   true,
		CLO:   true,
		CLS:   true,
		CLZ:   true,
		EXTSB: true,
		EXTSH: true,
		EXTUB: true,
		EXTUH: true,
		SATS:  true,
	}
}

func (this *Instruction) CaoRrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		CAO: true,
		CLO: true,
		CLS: true,
		CLZ: true,
	}
}

func (this *Instruction) ExtsbRrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		EXTSB: true,
		EXTSH: true,
		EXTUB: true,
		EXTUH: true,
		SATS:  true,
	}
}

func (this *Instruction) TimeCfgRrciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		TIME_CFG: true,
	}
}

func (this *Instruction) RrciOpCodes() map[OpCode]bool {
	cao_rrci_op_codes := this.CaoRrciOpCodes()
	extsb_rrci_op_codes := this.ExtsbRrciOpCodes()
	time_cfg_rrci_op_codes := this.TimeCfgRrciOpCodes()

	rrci_op_codes := make(map[OpCode]bool)

	for k, v := range cao_rrci_op_codes {
		rrci_op_codes[k] = v
	}

	for k, v := range extsb_rrci_op_codes {
		rrci_op_codes[k] = v
	}

	for k, v := range time_cfg_rrci_op_codes {
		rrci_op_codes[k] = v
	}

	return rrci_op_codes
}

func (this *Instruction) DivStepDrdiciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		DIV_STEP: true,
	}
}

func (this *Instruction) MulStepDrdiciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		MUL_STEP: true,
	}
}

func (this *Instruction) DrdiciOpCodes() map[OpCode]bool {
	div_step_drdici_op_codes := this.DivStepDrdiciOpCodes()
	mul_step_drdici_op_codes := this.MulStepDrdiciOpCodes()

	drdici_op_codes := make(map[OpCode]bool)

	for k, v := range div_step_drdici_op_codes {
		drdici_op_codes[k] = v
	}

	for k, v := range mul_step_drdici_op_codes {
		drdici_op_codes[k] = v
	}

	return drdici_op_codes
}

func (this *Instruction) RrriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		LSL_ADD: true,
		LSL_SUB: true,
		LSR_ADD: true,
		ROL_ADD: true,
	}
}

func (this *Instruction) RrriciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		LSL_ADD: true,
		LSL_SUB: true,
		LSR_ADD: true,
		ROL_ADD: true,
	}
}

func (this *Instruction) RirOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SUB:  true,
		SUBC: true,
	}
}

func (this *Instruction) RircOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SUB:  true,
		SUBC: true,
	}
}

func (this *Instruction) RirciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SUB:  true,
		SUBC: true,
	}
}

func (this *Instruction) ROpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		TIME: true,
	}
}

func (this *Instruction) RciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		TIME: true,
	}
}

func (this *Instruction) CiOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		STOP: true,
	}
}

func (this *Instruction) IOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		FAULT: true,
	}
}

func (this *Instruction) MovdDdciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		MOVD: true,
	}
}

func (this *Instruction) SwapdDdciOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SWAPD: true,
	}
}

func (this *Instruction) DdciOpCodes() map[OpCode]bool {
	movd_ddci_op_codes := this.MovdDdciOpCodes()
	swapd_ddci_op_codes := this.SwapdDdciOpCodes()

	ddci_op_codes := make(map[OpCode]bool)

	for k, v := range movd_ddci_op_codes {
		ddci_op_codes[k] = v
	}

	for k, v := range swapd_ddci_op_codes {
		ddci_op_codes[k] = v
	}

	return ddci_op_codes
}

func (this *Instruction) ErriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		LBS: true,
		LBU: true,
		LHS: true,
		LHU: true,
		LW:  true,
	}
}

func (this *Instruction) EdriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		LD: true,
	}
}

func (this *Instruction) EriiOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SB:    true,
		SB_ID: true,
		SD:    true,
		SD_ID: true,
		SH:    true,
		SH_ID: true,
		SW:    true,
		SW_ID: true,
	}
}

func (this *Instruction) ErirOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SB: true,
		SH: true,
		SW: true,
	}
}

func (this *Instruction) EridOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SD: true,
	}
}

func (this *Instruction) LdmaDmaRriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		LDMA: true,
	}
}

func (this *Instruction) LdmaiDmaRriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		LDMAI: true,
	}
}

func (this *Instruction) SdmaDmaRriOpCodes() map[OpCode]bool {
	return map[OpCode]bool{
		SDMA: true,
	}
}

func (this *Instruction) DmaRriOpCodes() map[OpCode]bool {
	ldma_dma_rri_op_codes := this.LdmaDmaRriOpCodes()
	ldmai_dma_rri_op_codes := this.LdmaiDmaRriOpCodes()
	sdma_dma_rri_op_codes := this.SdmaDmaRriOpCodes()

	dma_rri_op_codes := make(map[OpCode]bool)

	for k, v := range ldma_dma_rri_op_codes {
		dma_rri_op_codes[k] = v
	}

	for k, v := range ldmai_dma_rri_op_codes {
		dma_rri_op_codes[k] = v
	}

	for k, v := range sdma_dma_rri_op_codes {
		dma_rri_op_codes[k] = v
	}

	return dma_rri_op_codes
}

func (this *Instruction) Stringify() string {
	str := this.StringifyOpCode() + ", "
	str += this.StringifySuffix() + ", "

	if this.suffix == RICI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == RRI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == RRIC {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == RRICI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == RRIF {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == RRR {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb)
	} else if this.suffix == RRRC {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == RRRCI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == ZRI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == ZRIC {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == ZRICI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == ZRIF {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == ZRR {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb)
	} else if this.suffix == ZRRC {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == ZRRCI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == S_RRI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == S_RRIC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == S_RRICI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == S_RRIF {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == S_RRR {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb)
	} else if this.suffix == S_RRRC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == S_RRRCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == U_RRI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == U_RRIC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == U_RRICI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == U_RRIF {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == U_RRR {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb)
	} else if this.suffix == U_RRRC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == U_RRRCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == RR {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra)
	} else if this.suffix == RRC {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == RRCI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == ZR {
		str += this.StringifySrcRegDescriptor(this.ra)
	} else if this.suffix == ZRC {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == ZRCI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == S_RR {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra)
	} else if this.suffix == S_RRC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == S_RRCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == U_RR {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra)
	} else if this.suffix == U_RRC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == U_RRCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == DRDICI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyPairRegDescriptor(this.db) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == RRRI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == RRRICI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == ZRRI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == ZRRICI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == S_RRRI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == S_RRRICI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == U_RRRI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == U_RRRICI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == RIR {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra)
	} else if this.suffix == RIRC {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == RIRCI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == ZIR {
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra)
	} else if this.suffix == ZIRC {
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == ZIRCI {
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == S_RIRC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == S_RIRCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == U_RIRC {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition)
	} else if this.suffix == U_RIRCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifyImm(this.imm) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == R {
		str += this.StringifyGpRegDescriptor(this.rc)
	} else if this.suffix == RCI {
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == Z {
		str = str[:len(str)-2]
	} else if this.suffix == ZCI {
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == S_R {
		str += this.StringifyPairRegDescriptor(this.dc)
	} else if this.suffix == S_RCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == U_R {
		str += this.StringifyPairRegDescriptor(this.dc)
	} else if this.suffix == U_RCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == CI {
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == I {
		str += this.StringifyImm(this.imm)
	} else if this.suffix == DDCI {
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifyPairRegDescriptor(this.db) + ", "
		str += this.StringifyCondition(*this.condition) + ", "
		str += this.StringifyPc(this.pc)
	} else if this.suffix == ERRI {
		str += this.StringifyEndian(*this.endian) + ", "
		str += this.StringifyGpRegDescriptor(this.rc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyOff(this.off)
	} else if this.suffix == S_ERRI {
		str += this.StringifyEndian(*this.endian) + ", "
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyOff(this.off)
	} else if this.suffix == U_ERRI {
		str += this.StringifyEndian(*this.endian) + ", "
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyOff(this.off)
	} else if this.suffix == EDRI {
		str += this.StringifyEndian(*this.endian) + ", "
		str += this.StringifyPairRegDescriptor(this.dc) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyOff(this.off)
	} else if this.suffix == ERII {
		str += this.StringifyEndian(*this.endian) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyOff(this.off) + ", "
		str += this.StringifyImm(this.imm)
	} else if this.suffix == ERIR {
		str += this.StringifyEndian(*this.endian) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyOff(this.off) + ", "
		str += this.StringifySrcRegDescriptor(this.rb)
	} else if this.suffix == ERID {
		str += this.StringifyEndian(*this.endian) + ", "
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifyOff(this.off) + ", "
		str += this.StringifyPairRegDescriptor(this.db)
	} else if this.suffix == DMA_RRI {
		str += this.StringifySrcRegDescriptor(this.ra) + ", "
		str += this.StringifySrcRegDescriptor(this.rb) + ", "
		str += this.StringifyImm(this.imm)
	} else {
		err := errors.New("suffix is not valid")
		panic(err)
	}

	return str
}

func (this *Instruction) StringifyOpCode() string {
	if this.op_code == ACQUIRE {
		return "acquire"
	} else if this.op_code == RELEASE {
		return "release"
	} else if this.op_code == RELEASE {
		return "release"
	} else if this.op_code == ADD {
		return "add"
	} else if this.op_code == ADDC {
		return "addc"
	} else if this.op_code == AND {
		return "and"
	} else if this.op_code == ANDN {
		return "andn"
	} else if this.op_code == ASR {
		return "asr"
	} else if this.op_code == CAO {
		return "cao"
	} else if this.op_code == CLO {
		return "clo"
	} else if this.op_code == CLS {
		return "cls"
	} else if this.op_code == CLZ {
		return "clz"
	} else if this.op_code == CMPB4 {
		return "cmpb4"
	} else if this.op_code == DIV_STEP {
		return "div_step"
	} else if this.op_code == EXTSB {
		return "extsb"
	} else if this.op_code == EXTSH {
		return "extsh"
	} else if this.op_code == EXTUB {
		return "extub"
	} else if this.op_code == EXTUH {
		return "extuh"
	} else if this.op_code == LSL {
		return "lsl"
	} else if this.op_code == LSL_ADD {
		return "lsl_add"
	} else if this.op_code == LSL_SUB {
		return "lsl_sub"
	} else if this.op_code == LSL1 {
		return "lsl1"
	} else if this.op_code == LSL1X {
		return "lsl1x"
	} else if this.op_code == LSLX {
		return "lslx"
	} else if this.op_code == LSR {
		return "lsr"
	} else if this.op_code == LSR_ADD {
		return "lsr_add"
	} else if this.op_code == LSR1 {
		return "lsr1"
	} else if this.op_code == LSR1X {
		return "lsr1x"
	} else if this.op_code == LSRX {
		return "lsrx"
	} else if this.op_code == MUL_SH_SH {
		return "mul_sh_sh"
	} else if this.op_code == MUL_SH_SL {
		return "mul_sh_sl"
	} else if this.op_code == MUL_SH_UH {
		return "mul_sh_uh"
	} else if this.op_code == MUL_SH_UL {
		return "mul_sh_ul"
	} else if this.op_code == MUL_SL_SH {
		return "mul_sl_sh"
	} else if this.op_code == MUL_SL_SL {
		return "mul_sl_sl"
	} else if this.op_code == MUL_SL_UH {
		return "mul_sl_uh"
	} else if this.op_code == MUL_SL_UL {
		return "mul_sl_ul"
	} else if this.op_code == MUL_STEP {
		return "mul_step"
	} else if this.op_code == MUL_UH_UH {
		return "mul_uh_uh"
	} else if this.op_code == MUL_UH_UL {
		return "mul_uh_ul"
	} else if this.op_code == MUL_UL_UH {
		return "mul_ul_uh"
	} else if this.op_code == MUL_UL_UL {
		return "mul_ul_ul"
	} else if this.op_code == NAND {
		return "nand"
	} else if this.op_code == NOR {
		return "nor"
	} else if this.op_code == NXOR {
		return "nxor"
	} else if this.op_code == OR {
		return "or"
	} else if this.op_code == ORN {
		return "orn"
	} else if this.op_code == ROL {
		return "rol"
	} else if this.op_code == ROL_ADD {
		return "rol_add"
	} else if this.op_code == ROR {
		return "ror"
	} else if this.op_code == RSUB {
		return "rsub"
	} else if this.op_code == RSUBC {
		return "rsubc"
	} else if this.op_code == SUB {
		return "sub"
	} else if this.op_code == SUBC {
		return "subc"
	} else if this.op_code == XOR {
		return "xor"
	} else if this.op_code == BOOT {
		return "boot"
	} else if this.op_code == RESUME {
		return "resume"
	} else if this.op_code == STOP {
		return "stop"
	} else if this.op_code == CALL {
		return "call"
	} else if this.op_code == FAULT {
		return "fault"
	} else if this.op_code == NOP {
		return "nop"
	} else if this.op_code == SATS {
		return "sats"
	} else if this.op_code == MOVD {
		return "movd"
	} else if this.op_code == SWAPD {
		return "swapd"
	} else if this.op_code == HASH {
		return "hash"
	} else if this.op_code == TIME {
		return "time"
	} else if this.op_code == TIME_CFG {
		return "time_cfg"
	} else if this.op_code == LBS {
		return "lbs"
	} else if this.op_code == LBU {
		return "lbu"
	} else if this.op_code == LD {
		return "ld"
	} else if this.op_code == LHS {
		return "lhs"
	} else if this.op_code == LHU {
		return "lhu"
	} else if this.op_code == LW {
		return "lw"
	} else if this.op_code == SB {
		return "sb"
	} else if this.op_code == SB_ID {
		return "sb_id"
	} else if this.op_code == SD {
		return "sd"
	} else if this.op_code == SD_ID {
		return "sd_id"
	} else if this.op_code == SH {
		return "sh"
	} else if this.op_code == SH_ID {
		return "sh_id"
	} else if this.op_code == SW {
		return "sw"
	} else if this.op_code == SW_ID {
		return "sw_id"
	} else if this.op_code == LDMA {
		return "ldma"
	} else if this.op_code == LDMAI {
		return "ldmai"
	} else if this.op_code == SDMA {
		return "sdma"
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}

func (this *Instruction) StringifySuffix() string {
	if this.suffix == RICI {
		return "rici"
	} else if this.suffix == RRI {
		return "rri"
	} else if this.suffix == RRIC {
		return "rric"
	} else if this.suffix == RRICI {
		return "rrici"
	} else if this.suffix == RRIF {
		return "rrif"
	} else if this.suffix == RRR {
		return "rrr"
	} else if this.suffix == RRRC {
		return "rrrc"
	} else if this.suffix == RRRCI {
		return "rrrci"
	} else if this.suffix == ZRI {
		return "zri"
	} else if this.suffix == ZRIC {
		return "zric"
	} else if this.suffix == ZRICI {
		return "zrici"
	} else if this.suffix == ZRIF {
		return "zrif"
	} else if this.suffix == ZRR {
		return "zrr"
	} else if this.suffix == ZRRC {
		return "zrrc"
	} else if this.suffix == ZRRCI {
		return "zrrci"
	} else if this.suffix == S_RRI {
		return "s_rri"
	} else if this.suffix == S_RRIC {
		return "s_rric"
	} else if this.suffix == S_RRICI {
		return "s_rrici"
	} else if this.suffix == S_RRIF {
		return "s_rrif"
	} else if this.suffix == S_RRR {
		return "s_rrr"
	} else if this.suffix == S_RRRC {
		return "s_rrrc"
	} else if this.suffix == S_RRRCI {
		return "s_rrrci"
	} else if this.suffix == U_RRI {
		return "u_rri"
	} else if this.suffix == U_RRIC {
		return "u_rric"
	} else if this.suffix == U_RRICI {
		return "u_rrici"
	} else if this.suffix == U_RRIF {
		return "u_rrif"
	} else if this.suffix == U_RRR {
		return "u_rrr"
	} else if this.suffix == U_RRRC {
		return "u_rrrc"
	} else if this.suffix == U_RRRCI {
		return "u_rrrci"
	} else if this.suffix == RR {
		return "rr"
	} else if this.suffix == RRC {
		return "rrc"
	} else if this.suffix == RRCI {
		return "rrci"
	} else if this.suffix == ZR {
		return "zr"
	} else if this.suffix == ZRC {
		return "zrc"
	} else if this.suffix == ZRCI {
		return "zrci"
	} else if this.suffix == S_RR {
		return "s_rr"
	} else if this.suffix == S_RRC {
		return "s_rrc"
	} else if this.suffix == S_RRCI {
		return "s_rrci"
	} else if this.suffix == U_RR {
		return "u_rr"
	} else if this.suffix == U_RRC {
		return "u_rrc"
	} else if this.suffix == U_RRCI {
		return "u_rrci"
	} else if this.suffix == DRDICI {
		return "drdici"
	} else if this.suffix == RRRI {
		return "rrri"
	} else if this.suffix == RRRICI {
		return "rrrici"
	} else if this.suffix == ZRRI {
		return "zrri"
	} else if this.suffix == ZRRICI {
		return "zrrici"
	} else if this.suffix == S_RRRI {
		return "s_rrri"
	} else if this.suffix == S_RRRICI {
		return "s_rrrici"
	} else if this.suffix == U_RRRI {
		return "u_rrri"
	} else if this.suffix == U_RRRICI {
		return "u_rrrici"
	} else if this.suffix == RIR {
		return "rir"
	} else if this.suffix == RIRC {
		return "rirc"
	} else if this.suffix == RIRCI {
		return "rirci"
	} else if this.suffix == ZIR {
		return "zir"
	} else if this.suffix == ZIRC {
		return "zirc"
	} else if this.suffix == ZIRCI {
		return "zirci"
	} else if this.suffix == S_RIRC {
		return "s_zirc"
	} else if this.suffix == S_RIRCI {
		return "s_zirci"
	} else if this.suffix == U_RIRC {
		return "u_zirc"
	} else if this.suffix == U_RIRCI {
		return "u_zirci"
	} else if this.suffix == R {
		return "r"
	} else if this.suffix == RCI {
		return "rci"
	} else if this.suffix == Z {
		return "z"
	} else if this.suffix == ZCI {
		return "zci"
	} else if this.suffix == S_R {
		return "s_r"
	} else if this.suffix == S_RCI {
		return "s_rci"
	} else if this.suffix == U_R {
		return "u_r"
	} else if this.suffix == U_RCI {
		return "u_rci"
	} else if this.suffix == CI {
		return "ci"
	} else if this.suffix == I {
		return "i"
	} else if this.suffix == DDCI {
		return "ddci"
	} else if this.suffix == ERRI {
		return "erri"
	} else if this.suffix == S_ERRI {
		return "s_erri"
	} else if this.suffix == U_ERRI {
		return "u_erri"
	} else if this.suffix == EDRI {
		return "edri"
	} else if this.suffix == ERII {
		return "erii"
	} else if this.suffix == ERIR {
		return "erir"
	} else if this.suffix == ERID {
		return "erid"
	} else if this.suffix == DMA_RRI {
		return "dma_rri"
	} else {
		err := errors.New("suffix is not valid")
		panic(err)
	}
}

func (this *Instruction) StringifyGpRegDescriptor(
	gp_reg_descriptor *reg_descriptor.GpRegDescriptor,
) string {
	return "r" + strconv.Itoa(gp_reg_descriptor.Index())
}

func (this *Instruction) StringifySrcRegDescriptor(
	src_reg_descriptor *reg_descriptor.SrcRegDescriptor,
) string {
	if src_reg_descriptor.IsGpRegDescriptor() {
		return this.StringifyGpRegDescriptor(src_reg_descriptor.GpRegDescriptor())
	} else {
		sp_reg_descriptor := src_reg_descriptor.SpRegDescriptor()

		if *sp_reg_descriptor == reg_descriptor.ZERO {
			return "zero"
		} else if *sp_reg_descriptor == reg_descriptor.ONE {
			return "one"
		} else if *sp_reg_descriptor == reg_descriptor.LNEG {
			return "lneg"
		} else if *sp_reg_descriptor == reg_descriptor.MNEG {
			return "mneg"
		} else if *sp_reg_descriptor == reg_descriptor.ID {
			return "id"
		} else if *sp_reg_descriptor == reg_descriptor.ID2 {
			return "id2"
		} else if *sp_reg_descriptor == reg_descriptor.ID4 {
			return "id4"
		} else if *sp_reg_descriptor == reg_descriptor.ID8 {
			return "id8"
		} else {
			err := errors.New("sp reg descriptor is not valid")
			panic(err)
		}
	}
}

func (this *Instruction) StringifyPairRegDescriptor(
	pair_reg_descriptor *reg_descriptor.PairRegDescriptor,
) string {
	return "d" + strconv.Itoa(pair_reg_descriptor.Index())
}

func (this *Instruction) StringifyImm(imm *abi.Immediate) string {
	return strconv.FormatInt(imm.Value(), 10)
}

func (this *Instruction) StringifyCondition(condition cc.Condition) string {
	if condition == cc.TRUE {
		return "true"
	} else if condition == cc.FALSE {
		return "false"
	} else if condition == cc.Z {
		return "z"
	} else if condition == cc.NZ {
		return "nz"
	} else if condition == cc.E {
		return "e"
	} else if condition == cc.O {
		return "o"
	} else if condition == cc.PL {
		return "pl"
	} else if condition == cc.MI {
		return "mi"
	} else if condition == cc.OV {
		return "ov"
	} else if condition == cc.NOV {
		return "nov"
	} else if condition == cc.C {
		return "c"
	} else if condition == cc.NC {
		return "nc"
	} else if condition == cc.SZ {
		return "sz"
	} else if condition == cc.SNZ {
		return "snz"
	} else if condition == cc.SPL {
		return "spl"
	} else if condition == cc.SMI {
		return "smi"
	} else if condition == cc.SO {
		return "so"
	} else if condition == cc.SE {
		return "se"
	} else if condition == cc.NC5 {
		return "nc5"
	} else if condition == cc.NC6 {
		return "nc6"
	} else if condition == cc.NC7 {
		return "nc7"
	} else if condition == cc.NC8 {
		return "nc8"
	} else if condition == cc.NC9 {
		return "nc9"
	} else if condition == cc.NC10 {
		return "nc10"
	} else if condition == cc.NC11 {
		return "nc11"
	} else if condition == cc.NC12 {
		return "nc12"
	} else if condition == cc.NC13 {
		return "nc13"
	} else if condition == cc.NC14 {
		return "nc14"
	} else if condition == cc.MAX {
		return "max"
	} else if condition == cc.NMAX {
		return "nmax"
	} else if condition == cc.SH32 {
		return "sh32"
	} else if condition == cc.NSH32 {
		return "nsh32"
	} else if condition == cc.EQ {
		return "eq"
	} else if condition == cc.NEQ {
		return "neq"
	} else if condition == cc.LTU {
		return "ltu"
	} else if condition == cc.LEU {
		return "leu"
	} else if condition == cc.GTU {
		return "gtu"
	} else if condition == cc.GEU {
		return "geu"
	} else if condition == cc.LTS {
		return "lts"
	} else if condition == cc.LES {
		return "les"
	} else if condition == cc.GTS {
		return "gts"
	} else if condition == cc.GES {
		return "ges"
	} else if condition == cc.XZ {
		return "xz"
	} else if condition == cc.XNZ {
		return "xnz"
	} else if condition == cc.XLEU {
		return "xleu"
	} else if condition == cc.XGTU {
		return "xgtu"
	} else if condition == cc.XLES {
		return "xles"
	} else if condition == cc.XGTS {
		return "xgts"
	} else if condition == cc.SMALL {
		return "small"
	} else if condition == cc.LARGE {
		return "large"
	} else {
		err := errors.New("condition is not valid")
		panic(err)
	}
}

func (this *Instruction) StringifyEndian(endian Endian) string {
	if endian == LITTLE {
		return "!little"
	} else if endian == BIG {
		return "!big"
	} else {
		err := errors.New("endian is not valid")
		panic(err)
	}
}

func (this *Instruction) StringifyOff(off *abi.Immediate) string {
	return strconv.FormatInt(off.Value(), 10)
}

func (this *Instruction) StringifyPc(pc *abi.Immediate) string {
	return strconv.FormatInt(pc.Value(), 10)
}
