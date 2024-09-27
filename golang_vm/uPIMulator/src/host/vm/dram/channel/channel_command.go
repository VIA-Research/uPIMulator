package channel

import (
	"errors"
	"uPIMulator/src/host/vm/dram/bank"
	"uPIMulator/src/host/vm/dram/rank"
)

type ChannelCommand struct {
	channel_id int
	rank_id    int
	bank_id    int

	rank_command   *rank.RankCommand
	memory_command *bank.MemoryCommand
}

func (this *ChannelCommand) Init(
	channel_id int,
	rank_id int,
	bank_id int,
	memory_command *bank.MemoryCommand,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	} else if bank_id < 0 {
		err := errors.New("DPU ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.bank_id = bank_id

	this.rank_command = new(rank.RankCommand)
	this.rank_command.Init(channel_id, rank_id, bank_id, memory_command)

	this.memory_command = memory_command
}

func (this *ChannelCommand) ChannelId() int {
	return this.channel_id
}

func (this *ChannelCommand) RankId() int {
	return this.rank_id
}

func (this *ChannelCommand) BankId() int {
	return this.bank_id
}

func (this *ChannelCommand) RankCommand() *rank.RankCommand {
	return this.rank_command
}

func (this *ChannelCommand) MemoryCommand() *bank.MemoryCommand {
	return this.memory_command
}
