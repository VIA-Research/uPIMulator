package host

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/simulator/channel"
)

type ChannelTransferReadJob struct {
	channel_message *channel.ChannelMessage

	byte_streams []*encoding.ByteStream

	channel *channel.Channel
}

func (this *ChannelTransferReadJob) Init(
	channel_message *channel.ChannelMessage,
	byte_streams []*encoding.ByteStream,
	channel_ *channel.Channel,
) {
	if channel_message.ChannelOperation() != channel.READ {
		err := errors.New("channel operation is not read")
		panic(err)
	}

	this.channel_message = channel_message
	this.byte_streams = byte_streams
	this.channel = channel_
}

func (this *ChannelTransferReadJob) Execute() {
	this.channel.Lock()

	this.channel.Push(this.channel_message)

	for !this.channel.CanPopChannelMessage(this.channel_message) {
		this.channel.Cycle()
	}

	this.channel.PopChannelMessage(this.channel_message)

	this.CompareByteStreams(this.byte_streams, this.channel_message.ByteStreams())

	this.channel.Unlock()
}

func (this *ChannelTransferReadJob) CompareByteStreams(
	byte_streams_1 []*encoding.ByteStream,
	byte_streams_2 []*encoding.ByteStream,
) {
	if len(byte_streams_1) != len(byte_streams_2) {
		err := errors.New("byte streams 1's length != byte streams 2's length")
		panic(err)
	}

	for i := 0; i < len(byte_streams_1); i++ {
		byte_stream_1 := byte_streams_1[i]
		byte_stream_2 := byte_streams_2[i]

		if byte_stream_1.Size() != byte_stream_2.Size() {
			err := errors.New("byte stream 1's size != byte stream 2's size")
			panic(err)
		}

		for j := int64(0); j < byte_stream_1.Size(); j++ {
			if byte_stream_1.Get(int(j)) != byte_stream_2.Get(int(j)) {
				err := errors.New("bytes are different")
				panic(err)
			}
		}
	}
}
