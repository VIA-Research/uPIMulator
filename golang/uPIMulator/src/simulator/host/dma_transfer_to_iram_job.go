package host

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu"
)

type DmaTransferToIramJob struct {
	iram *encoding.ByteStream

	dpu *dpu.Dpu
}

func (this *DmaTransferToIramJob) Init(iram *encoding.ByteStream, dpu_ *dpu.Dpu) {
	this.iram = iram
	this.dpu = dpu_
}

func (this *DmaTransferToIramJob) Execute() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.dpu.Dma().TransferToIram(config_loader.IramOffset(), this.iram)
}
