package channel

import (
	"errors"
)

type ChannelMessageQ struct {
	size  int
	timer int64

	channel_messages []*ChannelMessage
	cycles           []int64
}

func (this *ChannelMessageQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.channel_messages = make([]*ChannelMessage, 0)
	this.cycles = make([]int64, 0)
}

func (this *ChannelMessageQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("channel message queue is not empty")
		panic(err)
	}
}

func (this *ChannelMessageQ) Size() int {
	return this.size
}

func (this *ChannelMessageQ) Timer() int64 {
	return this.timer
}

func (this *ChannelMessageQ) IsEmpty() bool {
	return len(this.channel_messages) == 0
}

func (this *ChannelMessageQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.channel_messages) >= num_items
	} else {
		return true
	}
}

func (this *ChannelMessageQ) Push(channel_message *ChannelMessage) {
	if !this.CanPush(1) {
		err := errors.New("channel message queue cannot be pushed")
		panic(err)
	}

	this.channel_messages = append(this.channel_messages, channel_message)
	this.cycles = append(this.cycles, this.timer)
}

func (this *ChannelMessageQ) PushWithTimer(channel_message *ChannelMessage, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("channel message queue cannot be pushed")
		panic(err)
	}

	this.channel_messages = append(this.channel_messages, channel_message)
	this.cycles = append(this.cycles, timer)
}

func (this *ChannelMessageQ) CanPop(num_items int) bool {
	if len(this.channel_messages) < num_items {
		return false
	} else {
		for i := 0; i < num_items; i++ {
			cycle := this.cycles[i]

			if cycle > 0 {
				return false
			}
		}
		return true
	}
}

func (this *ChannelMessageQ) Pop() *ChannelMessage {
	if !this.CanPop(1) {
		err := errors.New("channel message queue cannot be popped")
		panic(err)
	}

	channel_message := this.channel_messages[0]

	this.channel_messages = this.channel_messages[1:]
	this.cycles = this.cycles[1:]

	return channel_message
}

func (this *ChannelMessageQ) Front(pos int) (*ChannelMessage, int64) {
	if this.IsEmpty() {
		err := errors.New("channel message queue is empty")
		panic(err)
	}

	return this.channel_messages[pos], this.cycles[pos]
}

func (this *ChannelMessageQ) Remove(pos int) {
	this.channel_messages = append(this.channel_messages[:pos], this.channel_messages[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *ChannelMessageQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
