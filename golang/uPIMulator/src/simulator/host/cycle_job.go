package host

import (
	"uPIMulator/src/simulator/dpu"
	"uPIMulator/src/simulator/dpu/logic"
)

type CycleJob struct {
	sys_end int64

	dpu *dpu.Dpu
}

func (this *CycleJob) Init(sys_end int64, dpu_ *dpu.Dpu) {
	this.sys_end = sys_end
	this.dpu = dpu_
}

func (this *CycleJob) Execute() {
	threads := this.dpu.Threads()

	for _, thread := range threads {
		if thread.RegFile().ReadPcReg() == this.sys_end && thread.ThreadState() == logic.SLEEP {
			this.dpu.ThreadScheduler().Shutdown(thread.ThreadId())
		}
	}
}
