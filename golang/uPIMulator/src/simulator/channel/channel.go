package channel

import (
	"errors"
	"sync"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu"
	"uPIMulator/src/simulator/rank"
)

type Channel struct {
	mutex sync.Mutex

	channel_id int
	ranks      []*rank.Rank

	read_bandwidth  int64
	write_bandwidth int64

	input_q         *ChannelMessageQ
	communication_q *ChannelMessageQ
	ready_q         *ChannelMessageQ
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

	this.read_bandwidth = command_line_parser.IntParameter("read_bandwidth")
	this.write_bandwidth = command_line_parser.IntParameter("write_bandwidth")

	this.input_q = new(ChannelMessageQ)
	this.input_q.Init(-1, 0)

	this.communication_q = new(ChannelMessageQ)
	this.communication_q.Init(-1, 0)

	this.ready_q = new(ChannelMessageQ)
	this.ready_q.Init(-1, 0)
}

func (this *Channel) Fini() {
	for _, rank_ := range this.ranks {
		rank_.Fini()
	}

	this.input_q.Fini()
	this.communication_q.Fini()
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

func (this *Channel) Dpus() []*dpu.Dpu {
	dpus := make([]*dpu.Dpu, 0)

	for _, rank_ := range this.ranks {
		dpus = append(dpus, rank_.Dpus()...)
	}

	return dpus
}

func (this *Channel) Lock() {
	this.mutex.Lock()
}

func (this *Channel) Unlock() {
	this.mutex.Unlock()
}

func (this *Channel) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *Channel) Push(channel_message *ChannelMessage) {
	if !this.CanPush() {
		err := errors.New("channel cannot be pushed")
		panic(err)
	}

	this.input_q.Push(channel_message)
}

func (this *Channel) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *Channel) Pop() *ChannelMessage {
	if !this.CanPop() {
		err := errors.New("channel cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *Channel) CanPopChannelMessage(channel_message *ChannelMessage) bool {
	for i := 0; this.ready_q.CanPop(i + 1); i++ {
		cm, _ := this.ready_q.Front(i)

		if channel_message == cm {
			return true
		}
	}

	return false
}

func (this *Channel) PopChannelMessage(channel_message *ChannelMessage) {
	if !this.CanPopChannelMessage(channel_message) {
		err := errors.New("channel cannot be popped with the channel message")
		panic(err)
	}

	for i := 0; this.ready_q.CanPop(i + 1); i++ {
		cm, _ := this.ready_q.Front(i)

		if channel_message == cm {
			this.ready_q.Remove(i)
		}
	}
}

func (this *Channel) Cycle() {
	this.ServiceInputQ()
	this.ServiceCommunicationQ()

	this.input_q.Cycle()
	this.communication_q.Cycle()
	this.ready_q.Cycle()
}

func (this *Channel) ServiceInputQ() {
	if this.input_q.CanPop(1) && this.communication_q.CanPush(1) {
		channel_messaage := this.input_q.Pop()

		var latency int64

		channel_operation := channel_messaage.ChannelOperation()
		if channel_operation == READ {
			latency = channel_messaage.Size() / this.read_bandwidth
		} else if channel_operation == WRITE {
			latency = channel_messaage.Size() / this.write_bandwidth
		} else {
			err := errors.New("channel operation is not valid")
			panic(err)
		}
		this.communication_q.PushWithTimer(channel_messaage, latency)
	}
}

func (this *Channel) ServiceCommunicationQ() {
	if this.communication_q.CanPop(1) && this.ready_q.CanPush(1) {
		channel_messaage := this.communication_q.Pop()

		rank_id := channel_messaage.RankId()
		rank_ := this.ranks[rank_id]

		if rank_.RankId() != rank_id {
			err := errors.New("rank's rank ID != rank ID")
			panic(err)
		}

		dpu_ids := channel_messaage.DpuIds()

		address := channel_messaage.Address()
		size := channel_messaage.Size()

		channel_operation := channel_messaage.ChannelOperation()
		if channel_operation == READ {
			byte_streams := make([]*encoding.ByteStream, 0)

			for _, dpu_id := range dpu_ids {
				byte_stream := rank_.Read(dpu_id, address, size)
				byte_streams = append(byte_streams, byte_stream)
			}

			channel_messaage.SetByteStreams(byte_streams)
		} else if channel_operation == WRITE {
			for i, _ := range dpu_ids {
				dpu_id := dpu_ids[i]
				byte_stream := channel_messaage.ByteStreams()[i]

				rank_.Write(dpu_id, address, byte_stream)
			}
		} else {
			err := errors.New("channel operation is not valid")
			panic(err)
		}

		this.ready_q.Push(channel_messaage)
	}
}
