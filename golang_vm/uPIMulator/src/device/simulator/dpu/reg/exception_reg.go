package reg

import (
	"uPIMulator/src/device/linker/kernel/instruction"
)

type ExceptionReg struct {
	exceptions map[instruction.Exception]bool
}

func (this *ExceptionReg) Init() {
	this.exceptions = make(map[instruction.Exception]bool)

	this.ClearExceptions()
}

func (this *ExceptionReg) Fini() {
}

func (this *ExceptionReg) Exception(exception instruction.Exception) bool {
	return this.exceptions[exception]
}

func (this *ExceptionReg) SetException(exception instruction.Exception) {
	this.exceptions[exception] = true
}

func (this *ExceptionReg) ClearException(exception instruction.Exception) {
	this.exceptions[exception] = false
}

func (this *ExceptionReg) ClearExceptions() {
	for i := 0; i <= int(instruction.NOT_PROFILING); i++ {
		exception := instruction.Exception(i)

		this.exceptions[exception] = false
	}
}
