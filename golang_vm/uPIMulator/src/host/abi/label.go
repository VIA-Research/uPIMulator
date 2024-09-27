package abi

import (
	"fmt"
)

type Label struct {
	name      string
	bytecodes []*Bytecode
}

func (this *Label) Init(name string) {
	this.name = name
}

func (this *Label) Name() string {
	return this.name
}

func (this *Label) Length() int {
	return len(this.bytecodes)
}

func (this *Label) Get(index int) *Bytecode {
	return this.bytecodes[index]
}

func (this *Label) Append(bytecode *Bytecode) {
	this.bytecodes = append(this.bytecodes, bytecode)
}

func (this *Label) Stringify() string {
	ss := fmt.Sprintf("%s:\n", this.name)

	for i := 0; i < len(this.bytecodes); i++ {
		bytecode := this.bytecodes[i]

		if i != len(this.bytecodes)-1 {
			ss += fmt.Sprintf("\t%s\n", bytecode.Stringify())
		} else {
			ss += fmt.Sprintf("\t%s", bytecode.Stringify())
		}
	}

	return ss
}
