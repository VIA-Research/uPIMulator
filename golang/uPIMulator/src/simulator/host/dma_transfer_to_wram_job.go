package host

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu"
)

type DmaTransferToWramJob struct {
	wram *encoding.ByteStream

	dpu *dpu.Dpu
}

func (this *DmaTransferToWramJob) Init(wram *encoding.ByteStream, dpu_ *dpu.Dpu) {
	this.wram = wram
	this.dpu = dpu_
}

func (this *DmaTransferToWramJob) Execute() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.dpu.Dma().TransferToWram(config_loader.WramOffset(), this.wram)
}
