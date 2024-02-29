package host

import (
	"errors"
	"uPIMulator/src/simulator/channel"
)

type ChannelTransferWriteJob struct {
	channel_message *channel.ChannelMessage

	channel *channel.Channel
}

func (this *ChannelTransferWriteJob) Init(
	channel_message *channel.ChannelMessage,
	channel_ *channel.Channel,
) {
	if channel_message.ChannelOperation() != channel.WRITE {
		err := errors.New("channel operation is not write")
		panic(err)
	}

	this.channel_message = channel_message
	this.channel = channel_
}

func (this *ChannelTransferWriteJob) Execute() {
	this.channel.Lock()

	this.channel.Push(this.channel_message)

	for !this.channel.CanPopChannelMessage(this.channel_message) {
		this.channel.Cycle()
	}

	this.channel.PopChannelMessage(this.channel_message)

	this.channel.Unlock()
}
