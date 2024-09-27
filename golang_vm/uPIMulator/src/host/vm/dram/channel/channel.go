package channel

import (
	"errors"
	"uPIMulator/src/host/vm/dram/bank"
	"uPIMulator/src/host/vm/dram/rank"
	"uPIMulator/src/misc"
)

type Channel struct {
	channel_id int
	ranks      []*rank.Rank

	input_q    *ChannelCommandQ
	ready_q    *ChannelCommandQ
	scoreboard map[*rank.RankCommand]*ChannelCommand
}

func (this *Channel) Init(channel_id int, command_line_parser *misc.CommandLineParser) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	}

	this.channel_id = channel_id

	this.ranks = make([]*rank.Rank, 0)
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	for i := 0; i < num_ranks_per_channel; i++ {
		rank_ := new(rank.Rank)
		rank_.Init(channel_id, i, command_line_parser)
		this.ranks = append(this.ranks, rank_)
	}

	this.input_q = new(ChannelCommandQ)
	this.input_q.Init(-1, 0)

	this.ready_q = new(ChannelCommandQ)
	this.ready_q.Init(-1, 0)

	this.scoreboard = make(map[*rank.RankCommand]*ChannelCommand)
}

func (this *Channel) Fini() {
	for _, rank_ := range this.ranks {
		rank_.Fini()
	}

	this.input_q.Fini()
	this.ready_q.Fini()
}

func (this *Channel) ChannelId() int {
	return this.channel_id
}

func (this *Channel) NumRanks() int {
	return len(this.ranks)
}

func (this *Channel) Ranks() []*rank.Rank {
	return this.ranks
}

func (this *Channel) Banks() []*bank.Bank {
	banks := make([]*bank.Bank, 0)

	for _, rank_ := range this.ranks {
		banks = append(banks, rank_.Banks()...)
	}

	return banks
}

func (this *Channel) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *Channel) Push(channel_command *ChannelCommand) {
	if !this.CanPush() {
		err := errors.New("channel cannot be pushed")
		panic(err)
	}

	this.input_q.Push(channel_command)
}

func (this *Channel) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *Channel) Pop() *ChannelCommand {
	if !this.CanPop() {
		err := errors.New("channel cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *Channel) Read(channel_command *ChannelCommand) {
	rank_id := channel_command.RankId()

	rank_command := channel_command.RankCommand()

	this.ranks[rank_id].Read(rank_command)
}

func (this *Channel) Write(channel_command *ChannelCommand) {
	rank_id := channel_command.RankId()

	rank_command := channel_command.RankCommand()

	this.ranks[rank_id].Write(rank_command)
}

func (this *Channel) Flush() {
	for _, rank_ := range this.ranks {
		rank_.Flush()
	}
}

func (this *Channel) Cycle() {
	this.ServiceInputQ()
	this.ServiceReadyQ()

	for _, rank_ := range this.ranks {
		rank_.Cycle()
	}

	this.input_q.Cycle()
	this.ready_q.Cycle()
}

func (this *Channel) ServiceInputQ() {
	if this.input_q.CanPop(1) {
		channel_command, _ := this.input_q.Front(0)

		rank_id := channel_command.RankId()
		if this.ranks[rank_id].CanPush() {
			this.input_q.Pop()

			rank_command := channel_command.RankCommand()

			this.ranks[rank_id].Push(rank_command)
			this.scoreboard[rank_command] = channel_command
		}
	}
}

func (this *Channel) ServiceReadyQ() {
	for _, rank_ := range this.ranks {
		if rank_.CanPop() && this.ready_q.CanPush(1) {
			rank_command := rank_.Pop()
			channel_command := this.scoreboard[rank_command]

			this.ready_q.Push(channel_command)
			delete(this.scoreboard, rank_command)
		}
	}
}
