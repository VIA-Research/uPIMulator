package reg_descriptor

import (
	"errors"
	"uPIMulator/src/misc"
)

type GpRegDescriptor struct {
	index int
}

func (this *GpRegDescriptor) Init(index int) {
	if index < 0 {
		err := errors.New("index < 0")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	if index >= config_loader.NumGpRegisters() {
		err := errors.New("index >= num gp registers")
		panic(err)
	}

	this.index = index
}

func (this *GpRegDescriptor) Index() int {
	return this.index
}
