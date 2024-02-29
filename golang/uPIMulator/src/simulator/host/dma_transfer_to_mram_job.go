package host

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu"
)

type DmaTransferToMramJob struct {
	mram *encoding.ByteStream

	dpu *dpu.Dpu
}

func (this *DmaTransferToMramJob) Init(mram *encoding.ByteStream, dpu_ *dpu.Dpu) {
	this.mram = mram
	this.dpu = dpu_
}

func (this *DmaTransferToMramJob) Execute() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.dpu.Dma().TransferToMram(config_loader.MramOffset(), this.mram)
}
