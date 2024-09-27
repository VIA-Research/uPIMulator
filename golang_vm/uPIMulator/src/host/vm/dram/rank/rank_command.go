package rank

import (
	"errors"
	"uPIMulator/src/host/vm/dram/bank"
)

type RankCommand struct {
	channel_id int
	rank_id    int
	bank_id    int

	memory_command *bank.MemoryCommand
}

func (this *RankCommand) Init(
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
		err := errors.New("bank ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.bank_id = bank_id

	this.memory_command = memory_command
}

func (this *RankCommand) ChannelId() int {
	return this.channel_id
}

func (this *RankCommand) RankId() int {
	return this.rank_id
}

func (this *RankCommand) BankId() int {
	return this.bank_id
}

func (this *RankCommand) MemoryCommand() *bank.MemoryCommand {
	return this.memory_command
}
