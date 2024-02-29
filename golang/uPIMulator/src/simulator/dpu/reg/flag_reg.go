package reg

import (
	"uPIMulator/src/linker/kernel/instruction"
)

type FlagReg struct {
	flags map[instruction.Flag]bool
}

func (this *FlagReg) Init() {
	this.flags = make(map[instruction.Flag]bool, 0)

	this.ClearFlags()
}

func (this *FlagReg) Fini() {
}

func (this *FlagReg) Flag(flag instruction.Flag) bool {
	return this.flags[flag]
}

func (this *FlagReg) SetFlag(flag instruction.Flag) {
	this.flags[flag] = true
}

func (this *FlagReg) ClearFlag(flag instruction.Flag) {
	this.flags[flag] = false
}

func (this *FlagReg) ClearFlags() {
	for i := 0; i <= int(instruction.CARRY); i++ {
		flag := instruction.Flag(i)

		this.flags[flag] = false
	}
}
