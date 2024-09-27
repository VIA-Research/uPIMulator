package dpu

import (
	"errors"
)

type ControlInterface struct {
	boot bool
}

func (this *ControlInterface) Init() {
	this.boot = false
}

func (this *ControlInterface) Boot() bool {
	return this.boot
}

func (this *ControlInterface) SetBoot() {
	if this.boot {
		err := errors.New("control interface is already booted")
		panic(err)
	}

	this.boot = true
}

func (this *ControlInterface) UnsetBoot() {
	if !this.boot {
		err := errors.New("control interface is not booted")
		panic(err)
	}

	this.boot = false
}
