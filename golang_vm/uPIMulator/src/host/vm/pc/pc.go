package pc

import (
	"uPIMulator/src/host/abi"
)

type Pc struct {
	label *abi.Label
	index int
}

func (this *Pc) Init() {
	this.label = nil
	this.index = 0
}

func (this *Pc) CanAdvance() bool {
	return this.index < this.label.Length()
}

func (this *Pc) Advance() *abi.Bytecode {
	bytecode := this.label.Get(this.index)

	this.index++

	return bytecode
}

func (this *Pc) Jump(label *abi.Label) {
	this.label = label
	this.index = 0
}
