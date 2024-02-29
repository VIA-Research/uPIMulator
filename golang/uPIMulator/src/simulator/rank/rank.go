package rank

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu"
)

type Rank struct {
	channel_id int
	rank_id    int

	dpus []*dpu.Dpu
}

func (this *Rank) Init(channel_id int, rank_id int, command_line_parser *misc.CommandLineParser) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id

	this.dpus = make([]*dpu.Dpu, 0)
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))
	for i := 0; i < num_dpus_per_rank; i++ {
		dpu_ := new(dpu.Dpu)
		dpu_.Init(channel_id, rank_id, i, command_line_parser)

		this.dpus = append(this.dpus, dpu_)
	}
}

func (this *Rank) Fini() {
	for _, dpu_ := range this.dpus {
		dpu_.Fini()
	}
}

func (this *Rank) RankId() int {
	return this.rank_id
}

func (this *Rank) NumDpus() int {
	return len(this.dpus)
}

func (this *Rank) Dpus() []*dpu.Dpu {
	return this.dpus
}

func (this *Rank) Read(dpu_id int, address int64, size int64) *encoding.ByteStream {
	dpu_ := this.dpus[dpu_id]

	if dpu_.DpuId() != dpu_id {
		err := errors.New("DPU's DPU ID != DPU ID")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	if config_loader.WramOffset() <= address &&
		address+size <= config_loader.WramOffset()+config_loader.WramSize() {
		return dpu_.Dma().TransferFromWram(address, size)
	} else if config_loader.MramOffset() <= address && address+size <= config_loader.MramOffset()+config_loader.MramSize() {
		return dpu_.Dma().TransferFromMram(address, size)
	} else {
		err := errors.New("address does not fall under WRAM nor MRAM region")
		panic(err)
	}
}

func (this *Rank) Write(dpu_id int, address int64, byte_stream *encoding.ByteStream) {
	dpu_ := this.dpus[dpu_id]

	if dpu_.DpuId() != dpu_id {
		err := errors.New("DPU's DPU ID != DPU ID")
		panic(err)
	}

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	if config_loader.WramOffset() <= address &&
		address+byte_stream.Size() <= config_loader.WramOffset()+config_loader.WramSize() {
		dpu_.Dma().TransferToWram(address, byte_stream)
	} else if config_loader.MramOffset() <= address && address+byte_stream.Size() <= config_loader.MramOffset()+config_loader.MramSize() {
		dpu_.Dma().TransferToMram(address, byte_stream)
	} else {
		err := errors.New("address does not fall under WRAM nor MRAM region")
		panic(err)
	}
}
