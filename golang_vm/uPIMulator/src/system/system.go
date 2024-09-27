package system

import (
	"uPIMulator/src/host/vm"
	"uPIMulator/src/misc"
	"uPIMulator/src/program"
)

type System struct {
	vm *vm.VirtualMachine
}

func (this *System) Init(command_line_parser *misc.CommandLineParser) {
	this.vm = new(vm.VirtualMachine)
	this.vm.Init(command_line_parser)
}

func (this *System) Fini() {
	this.vm.Fini()
}

func (this *System) Simulate(app *program.App, task *program.Task) {
	this.vm.Load(app, task)

	for this.vm.CanAdvance() {
		this.vm.Advance()
	}
}

func (this *System) Dump() {
	this.vm.Dump()
}
