package linker

import (
	"fmt"
	"uPIMulator/src/device/linker/kernel"
	"uPIMulator/src/device/linker/logic"
)

type AnalyzeLivenessJob struct {
	relocatable *kernel.Relocatable
}

func (this *AnalyzeLivenessJob) Init(relocatable *kernel.Relocatable) {
	this.relocatable = relocatable
}

func (this *AnalyzeLivenessJob) Execute() {
	fmt.Printf("Analyzing the liveness of %s...\n", this.relocatable.Path())

	liveness_analyzer := new(logic.LivenessAnalyzer)
	liveness_analyzer.Init()

	liveness := liveness_analyzer.Analyze(this.relocatable)
	this.relocatable.SetLiveness(liveness)
}
