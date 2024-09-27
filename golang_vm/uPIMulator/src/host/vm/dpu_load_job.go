package vm

import (
	"uPIMulator/src/device/simulator/dpu"
	"uPIMulator/src/misc"
	"uPIMulator/src/program"
)

type DpuLoadJob struct {
	task *program.Task

	dpu *dpu.Dpu
}

func (this *DpuLoadJob) Init(task *program.Task, dpu_ *dpu.Dpu) {
	this.task = task
	this.dpu = dpu_
}

func (this *DpuLoadJob) Execute() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.dpu.Dma().
		TransferToAtomic(config_loader.AtomicOffset(), this.task.Atomic().Size(), this.task.Atomic())
	this.dpu.Dma().
		TransferToIram(config_loader.IramOffset(), this.task.Iram().Size(), this.task.Iram())
	this.dpu.Dma().
		TransferToWram(config_loader.WramOffset(), this.task.Wram().Size(), this.task.Wram())
	this.dpu.Dma().
		TransferToMram(config_loader.MramOffset(), this.task.Mram().Size(), this.task.Mram())
}
