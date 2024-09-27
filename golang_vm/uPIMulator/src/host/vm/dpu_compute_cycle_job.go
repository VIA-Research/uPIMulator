package vm

import (
	"uPIMulator/src/device/simulator/dpu"
	"uPIMulator/src/device/simulator/dpu/logic"
)

type DpuComputeCycleJob struct {
	sys_end int64

	dpu *dpu.Dpu
}

func (this *DpuComputeCycleJob) Init(sys_end int64, dpu_ *dpu.Dpu) {
	this.sys_end = sys_end

	this.dpu = dpu_
}

func (this *DpuComputeCycleJob) Execute() {
	for !this.dpu.IsZombie() {
		for _, thread := range this.dpu.Threads() {
			if thread.RegFile().ReadPcReg() == this.sys_end && thread.ThreadState() == logic.SLEEP {
				this.dpu.ThreadScheduler().Shutdown(thread.ThreadId())
			}
		}

		this.dpu.Cycle()
	}
}
