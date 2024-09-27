package reg

import (
	"uPIMulator/src/device/abi"
	"uPIMulator/src/device/linker/kernel/instruction/reg_descriptor"
	"uPIMulator/src/misc"
)

type GpReg struct {
	gp_reg_descriptor *reg_descriptor.GpRegDescriptor
	word              *abi.Word
}

func (this *GpReg) Init(index int) {
	this.gp_reg_descriptor = new(reg_descriptor.GpRegDescriptor)
	this.gp_reg_descriptor.Init(index)

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.word = new(abi.Word)
	this.word.Init(config_loader.MramDataWidth())
}

func (this *GpReg) Fini() {
}

func (this *GpReg) Index() int {
	return this.gp_reg_descriptor.Index()
}

func (this *GpReg) Read(representation abi.Representation) int64 {
	return this.word.Value(representation)
}

func (this *GpReg) Write(value int64) {
	this.word.SetValue(value)
}
