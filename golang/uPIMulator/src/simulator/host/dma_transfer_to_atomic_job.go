package host

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu"
)

type DmaTransferToAtomicJob struct {
	atomic *encoding.ByteStream

	dpu *dpu.Dpu
}

func (this *DmaTransferToAtomicJob) Init(atomic *encoding.ByteStream, dpu_ *dpu.Dpu) {
	this.atomic = atomic
	this.dpu = dpu_
}

func (this *DmaTransferToAtomicJob) Execute() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.dpu.Dma().TransferToAtomic(config_loader.AtomicOffset(), this.atomic)
}
