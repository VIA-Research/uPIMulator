package channel

import (
	"errors"
	"uPIMulator/src/device/simulator/dpu/dram"
	"uPIMulator/src/device/simulator/rank"
)

type ChannelCommand struct {
	channel_id int
	rank_id    int
	dpu_id     int

	rank_command *rank.RankCommand
	dma_command  *dram.DmaCommand
}

func (this *ChannelCommand) Init(
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

	this.rank_command = new(rank.RankCommand)
	this.rank_command.Init(channel_id, rank_id, dpu_id, dma_command)

	this.dma_command = dma_command
}

func (this *ChannelCommand) ChannelId() int {
	return this.channel_id
}

func (this *ChannelCommand) RankId() int {
	return this.rank_id
}

func (this *ChannelCommand) DpuId() int {
	return this.dpu_id
}

func (this *ChannelCommand) RankCommand() *rank.RankCommand {
	return this.rank_command
}

func (this *ChannelCommand) DmaCommand() *dram.DmaCommand {
	return this.dma_command
}
