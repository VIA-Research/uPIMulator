package simulator

import (
	"uPIMulator/src/simulator/dpu"
)

type CycleJob struct {
	dpu *dpu.Dpu
}

func (this *CycleJob) Init(dpu_ *dpu.Dpu) {
	this.dpu = dpu_
}

func (this *CycleJob) Execute() {
	this.dpu.Cycle()
}
