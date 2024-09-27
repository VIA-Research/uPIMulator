package rank

import (
	"errors"
	"uPIMulator/src/device/simulator/dpu"
	"uPIMulator/src/device/simulator/dpu/dram"
	"uPIMulator/src/misc"
)

type Rank struct {
	channel_id int
	rank_id    int

	dpus []*dpu.Dpu

	input_q    *RankCommandQ
	ready_q    *RankCommandQ
	scoreboard map[*dram.DmaCommand]*RankCommand
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

	this.input_q = new(RankCommandQ)
	this.input_q.Init(-1, 0)

	this.ready_q = new(RankCommandQ)
	this.ready_q.Init(-1, 0)

	this.scoreboard = make(map[*dram.DmaCommand]*RankCommand)
}

func (this *Rank) Fini() {
	for _, dpu_ := range this.dpus {
		dpu_.Fini()
	}

	this.input_q.Fini()
	this.ready_q.Fini()
}

func (this *Rank) ChannelId() int {
	return this.channel_id
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

func (this *Rank) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *Rank) Push(rank_command *RankCommand) {
	if !this.CanPush() {
		err := errors.New("rank cannot be pushed")
		panic(err)
	}

	this.input_q.Push(rank_command)
}

func (this *Rank) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *Rank) Pop() *RankCommand {
	if !this.CanPop() {
		err := errors.New("rank cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *Rank) Cycle() {
	this.ServiceInputQ()
	this.ServiceReadyQ()

	this.input_q.Cycle()
	this.ready_q.Cycle()
}

func (this *Rank) ServiceInputQ() {
	if this.input_q.CanPop(1) {
		rank_command, _ := this.input_q.Front(0)

		dpu_id := rank_command.DpuId()
		if this.dpus[dpu_id].Dma().CanPush() {
			this.input_q.Pop()

			dma_command := rank_command.DmaCommand()

			this.dpus[dpu_id].Dma().Push(dma_command)
			this.scoreboard[dma_command] = rank_command
		}
	}
}

func (this *Rank) ServiceReadyQ() {
	for _, dpu_ := range this.dpus {
		if dpu_.Dma().CanPop() && this.ready_q.CanPush(1) {
			dma_command := dpu_.Dma().Pop()
			rank_command := this.scoreboard[dma_command]

			this.ready_q.Push(rank_command)
			delete(this.scoreboard, dma_command)
		}
	}
}
