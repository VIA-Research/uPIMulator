package logic

import (
	"errors"
	"uPIMulator/src/linker/kernel/instruction"
	"uPIMulator/src/linker/kernel/instruction/reg_descriptor"
)

type RegSet struct {
	thread_id int

	prev_write_gp_reg_set map[*reg_descriptor.GpRegDescriptor]bool
	cur_read_gp_reg_set   map[*reg_descriptor.GpRegDescriptor]bool
}

func (this *RegSet) Init(thread_id int) {
	if thread_id < 0 {
		err := errors.New("thread ID < 0")
		panic(err)
	}

	this.thread_id = thread_id

	this.prev_write_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
	this.cur_read_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
}

func (this *RegSet) ThreadId() int {
	return this.thread_id
}

func (this *RegSet) CollectReadGpRegs(instruction_ *instruction.Instruction) {
	suffix := instruction_.Suffix()
	if suffix == instruction.RICI ||
		suffix == instruction.RRI ||
		suffix == instruction.RRIC ||
		suffix == instruction.RRICI ||
		suffix == instruction.RRIF ||
		suffix == instruction.ZRI ||
		suffix == instruction.ZRIC ||
		suffix == instruction.ZRICI ||
		suffix == instruction.ZRIF ||
		suffix == instruction.S_RRI ||
		suffix == instruction.U_RRI ||
		suffix == instruction.S_RRIC ||
		suffix == instruction.U_RRIC ||
		suffix == instruction.S_RRICI ||
		suffix == instruction.U_RRICI ||
		suffix == instruction.S_RRIF ||
		suffix == instruction.U_RRIF ||
		suffix == instruction.RR ||
		suffix == instruction.RRC ||
		suffix == instruction.RRCI ||
		suffix == instruction.ZR ||
		suffix == instruction.ZRC ||
		suffix == instruction.ZRCI ||
		suffix == instruction.S_RR ||
		suffix == instruction.U_RR ||
		suffix == instruction.S_RRC ||
		suffix == instruction.U_RRC ||
		suffix == instruction.S_RRCI ||
		suffix == instruction.U_RRCI ||
		suffix == instruction.RIR ||
		suffix == instruction.RIRC ||
		suffix == instruction.RIRCI ||
		suffix == instruction.ZIR ||
		suffix == instruction.ZIRC ||
		suffix == instruction.ZIRCI ||
		suffix == instruction.S_RIRC ||
		suffix == instruction.U_RIRC ||
		suffix == instruction.S_RIRCI ||
		suffix == instruction.U_RIRCI ||
		suffix == instruction.ERRI ||
		suffix == instruction.S_ERRI || suffix == instruction.U_ERRI ||
		suffix == instruction.EDRI ||
		suffix == instruction.ERII {
		if instruction_.Ra().IsGpRegDescriptor() {
			this.cur_read_gp_reg_set[instruction_.Ra().GpRegDescriptor()] = true
		} else {
			this.cur_read_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
		}
	} else if suffix == instruction.RRR ||
		suffix == instruction.RRRC ||
		suffix == instruction.RRRCI ||
		suffix == instruction.ZRR ||
		suffix == instruction.ZRRC ||
		suffix == instruction.ZRRCI ||
		suffix == instruction.S_RRR ||
		suffix == instruction.U_RRR ||
		suffix == instruction.S_RRRC ||
		suffix == instruction.U_RRRC ||
		suffix == instruction.S_RRRCI ||
		suffix == instruction.U_RRRCI ||
		suffix == instruction.RRRI ||
		suffix == instruction.RRRICI ||
		suffix == instruction.ZRRI ||
		suffix == instruction.ZRRICI ||
		suffix == instruction.S_RRRI ||
		suffix == instruction.U_RRRI ||
		suffix == instruction.S_RRRICI ||
		suffix == instruction.U_RRRICI ||
		suffix == instruction.ERIR ||
		suffix == instruction.DMA_RRI {
		if instruction_.Ra().IsGpRegDescriptor() && instruction_.Rb().IsGpRegDescriptor() {
			this.cur_read_gp_reg_set[instruction_.Ra().GpRegDescriptor()] = true
			this.cur_read_gp_reg_set[instruction_.Rb().GpRegDescriptor()] = true
		} else if instruction_.Ra().IsGpRegDescriptor() {
			this.cur_read_gp_reg_set[instruction_.Ra().GpRegDescriptor()] = true
		} else if instruction_.Rb().IsGpRegDescriptor() {
			this.cur_read_gp_reg_set[instruction_.Rb().GpRegDescriptor()] = true
		} else {
			this.cur_read_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
		}
	} else if suffix == instruction.DRDICI || suffix == instruction.ERID {
		if instruction_.Ra().IsGpRegDescriptor() {
			this.cur_read_gp_reg_set[instruction_.Ra().GpRegDescriptor()] = true
			this.cur_read_gp_reg_set[instruction_.Db().EvenRegDescriptor()] = true
			this.cur_read_gp_reg_set[instruction_.Db().OddRegDescriptor()] = true
		} else {
			this.cur_read_gp_reg_set[instruction_.Db().EvenRegDescriptor()] = true
			this.cur_read_gp_reg_set[instruction_.Db().OddRegDescriptor()] = true
		}

	} else if suffix == instruction.R ||
		suffix == instruction.RCI ||
		suffix == instruction.Z ||
		suffix == instruction.ZCI ||
		suffix == instruction.S_R ||
		suffix == instruction.U_R ||
		suffix == instruction.S_RCI ||
		suffix == instruction.U_RCI ||
		suffix == instruction.CI ||
		suffix == instruction.I {
		this.cur_read_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
	} else if suffix == instruction.DDCI {
		this.cur_read_gp_reg_set[instruction_.Db().EvenRegDescriptor()] = true
		this.cur_read_gp_reg_set[instruction_.Db().OddRegDescriptor()] = true
	} else {
		err := errors.New("suffix is not valid")
		panic(err)
	}
}

func (this *RegSet) CollectWriteGpRegs(instruction_ *instruction.Instruction) {
	suffix := instruction_.Suffix()
	if suffix == instruction.RICI ||
		suffix == instruction.ZRI ||
		suffix == instruction.ZRIC ||
		suffix == instruction.ZRICI ||
		suffix == instruction.ZRIF ||
		suffix == instruction.ZRR ||
		suffix == instruction.ZRRC ||
		suffix == instruction.ZRRCI ||
		suffix == instruction.ZR ||
		suffix == instruction.ZRC ||
		suffix == instruction.ZRCI ||
		suffix == instruction.ZRRI ||
		suffix == instruction.ZRRICI ||
		suffix == instruction.ZIR ||
		suffix == instruction.ZIRC ||
		suffix == instruction.ZIRCI ||
		suffix == instruction.Z ||
		suffix == instruction.ZCI ||
		suffix == instruction.CI ||
		suffix == instruction.I ||
		suffix == instruction.ERII ||
		suffix == instruction.ERIR ||
		suffix == instruction.ERID ||
		suffix == instruction.DMA_RRI {
		this.prev_write_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
	} else if suffix == instruction.RRI ||
		suffix == instruction.RRIC ||
		suffix == instruction.RRICI ||
		suffix == instruction.RRIF ||
		suffix == instruction.RRR ||
		suffix == instruction.RRRC ||
		suffix == instruction.RRRCI ||
		suffix == instruction.RR ||
		suffix == instruction.RRC ||
		suffix == instruction.RRCI ||
		suffix == instruction.RRRI ||
		suffix == instruction.RRRICI ||
		suffix == instruction.RIR ||
		suffix == instruction.RIRC ||
		suffix == instruction.RIRCI ||
		suffix == instruction.R ||
		suffix == instruction.RCI ||
		suffix == instruction.ERRI {
		this.prev_write_gp_reg_set[instruction_.Rc()] = true
	} else if suffix == instruction.S_RRI ||
		suffix == instruction.U_RRI ||
		suffix == instruction.S_RRIC ||
		suffix == instruction.U_RRIC ||
		suffix == instruction.S_RRICI ||
		suffix == instruction.U_RRICI ||
		suffix == instruction.S_RRIF ||
		suffix == instruction.U_RRIF ||
		suffix == instruction.S_RRR ||
		suffix == instruction.U_RRR ||
		suffix == instruction.S_RRRC ||
		suffix == instruction.U_RRRC ||
		suffix == instruction.S_RRRCI ||
		suffix == instruction.U_RRRCI ||
		suffix == instruction.S_RR ||
		suffix == instruction.U_RR ||
		suffix == instruction.S_RRC ||
		suffix == instruction.U_RRC ||
		suffix == instruction.S_RRCI ||
		suffix == instruction.U_RRCI ||
		suffix == instruction.DRDICI ||
		suffix == instruction.S_RRRI ||
		suffix == instruction.U_RRRI ||
		suffix == instruction.S_RRRICI ||
		suffix == instruction.U_RRRICI ||
		suffix == instruction.S_RIRC ||
		suffix == instruction.U_RIRC ||
		suffix == instruction.S_RIRCI ||
		suffix == instruction.U_RIRCI ||
		suffix == instruction.S_R ||
		suffix == instruction.U_R ||
		suffix == instruction.DDCI ||
		suffix == instruction.S_ERRI ||
		suffix == instruction.U_ERRI ||
		suffix == instruction.EDRI {
		this.prev_write_gp_reg_set[instruction_.Dc().EvenRegDescriptor()] = true
		this.prev_write_gp_reg_set[instruction_.Dc().OddRegDescriptor()] = true
	} else {
		err := errors.New("suffix is not valid")
		panic(err)
	}
}

func (this *RegSet) Clear() {
	this.prev_write_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
	this.cur_read_gp_reg_set = make(map[*reg_descriptor.GpRegDescriptor]bool, 0)
}

func (this *RegSet) RegIndices() map[int]bool {
	reg_indices := make(map[int]bool, 0)

	for gp_reg_descriptor, _ := range this.prev_write_gp_reg_set {
		reg_indices[gp_reg_descriptor.Index()] = true
	}

	for gp_reg_descriptor, _ := range this.cur_read_gp_reg_set {
		reg_indices[gp_reg_descriptor.Index()] = true
	}

	return reg_indices
}
