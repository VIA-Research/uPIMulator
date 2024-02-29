package channel

import (
	"errors"
	"uPIMulator/src/abi/encoding"
)

type ChannelOperation int

const (
	READ ChannelOperation = iota
	WRITE
)

type ChannelMessage struct {
	channel_operation ChannelOperation
	channel_id        int
	rank_id           int
	dpu_ids           []int
	address           int64
	size              int64
	byte_streams      []*encoding.ByteStream
}

func (this *ChannelMessage) InitRead(
	channel_id int,
	rank_id int,
	dpu_ids []int,
	address int64,
	size int64,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	}

	for _, dpu_id := range dpu_ids {
		if dpu_id < 0 {
			err := errors.New("DPU ID < 0")
			panic(err)
		} else if dpu_id%8 != dpu_ids[0]%8 {
			err := errors.New("DPU ID % 8 are different")
			panic(err)
		}
	}

	this.channel_operation = READ
	this.channel_id = channel_id
	this.rank_id = rank_id
	this.dpu_ids = dpu_ids
	this.address = address
	this.size = size
	this.byte_streams = make([]*encoding.ByteStream, 0)
}

func (this *ChannelMessage) InitWrite(
	channel_id int,
	rank_id int,
	dpu_ids []int,
	address int64,
	size int64,
	byte_streams []*encoding.ByteStream,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	}

	for _, dpu_id := range dpu_ids {
		if dpu_id < 0 {
			err := errors.New("DPU ID < 0")
			panic(err)
		} else if dpu_id%8 != dpu_ids[0]%8 {
			err := errors.New("DPU ID % 8 are different")
			panic(err)
		}
	}

	for _, byte_stream := range byte_streams {
		if byte_stream.Size() != size {
			err := errors.New("byte stream's size != size")
			panic(err)
		}
	}

	this.channel_operation = WRITE
	this.channel_id = channel_id
	this.rank_id = rank_id
	this.dpu_ids = dpu_ids
	this.address = address
	this.size = size
	this.byte_streams = byte_streams
}

func (this *ChannelMessage) SetByteStreams(byte_streams []*encoding.ByteStream) {
	for _, byte_stream := range byte_streams {
		if byte_stream.Size() != this.size {
			err := errors.New("byte stream's size != size")
			panic(err)
		}
	}

	this.byte_streams = byte_streams
}

func (this *ChannelMessage) ChannelOperation() ChannelOperation {
	return this.channel_operation
}

func (this *ChannelMessage) ChannelId() int {
	return this.channel_id
}

func (this *ChannelMessage) RankId() int {
	return this.rank_id
}

func (this *ChannelMessage) DpuIds() []int {
	return this.dpu_ids
}

func (this *ChannelMessage) Address() int64 {
	return this.address
}

func (this *ChannelMessage) Size() int64 {
	return this.size
}

func (this *ChannelMessage) ByteStreams() []*encoding.ByteStream {
	return this.byte_streams
}
