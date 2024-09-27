package rank

import (
	"errors"
	"uPIMulator/src/device/simulator/dpu/dram"
)

type RankCommand struct {
	channel_id int
	rank_id    int
	dpu_id     int

	dma_command *dram.DmaCommand
}

func (this *RankCommand) Init(
	channel_id int,
	rank_id int,
	dpu_id int,
	dma_command *dram.DmaCommand,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	} else if dpu_id < 0 {
		err := errors.New("DPU ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.dpu_id = dpu_id

	this.dma_command = dma_command
}

func (this *RankCommand) RankId() int {
	return this.rank_id
}

func (this *RankCommand) DpuId() int {
	return this.dpu_id
}

func (this *RankCommand) DmaCommand() *dram.DmaCommand {
	return this.dma_command
}
