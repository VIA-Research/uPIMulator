package vm

import (
	"uPIMulator/src/device/simulator/dpu"
)

type DpuCycleJob struct {
	dpu *dpu.Dpu
}

func (this *DpuCycleJob) Init(dpu_ *dpu.Dpu) {
	this.dpu = dpu_
}

func (this *DpuCycleJob) Execute() {
	this.dpu.Cycle()
}
