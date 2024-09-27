package channel

import (
	"errors"
)

type ChannelCommandQ struct {
	size  int
	timer int64

	channel_commands []*ChannelCommand
	cycles           []int64
}

func (this *ChannelCommandQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.channel_commands = make([]*ChannelCommand, 0)
	this.cycles = make([]int64, 0)
}

func (this *ChannelCommandQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("channel command queue is not empty")
		panic(err)
	}
}

func (this *ChannelCommandQ) Size() int {
	return this.size
}

func (this *ChannelCommandQ) Length() int {
	return len(this.channel_commands)
}

func (this *ChannelCommandQ) Timer() int64 {
	return this.timer
}

func (this *ChannelCommandQ) IsEmpty() bool {
	return len(this.channel_commands) == 0
}

func (this *ChannelCommandQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.channel_commands) >= num_items
	} else {
		return true
	}
}

func (this *ChannelCommandQ) Push(channel_command *ChannelCommand) {
	if !this.CanPush(1) {
		err := errors.New("channel command queue cannot be pushed")
		panic(err)
	}

	this.channel_commands = append(this.channel_commands, channel_command)
	this.cycles = append(this.cycles, this.timer)
}

func (this *ChannelCommandQ) PushWithTimer(channel_command *ChannelCommand, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("channel command queue cannot be pushed")
		panic(err)
	}

	this.channel_commands = append(this.channel_commands, channel_command)
	this.cycles = append(this.cycles, timer)
}

func (this *ChannelCommandQ) CanPop(num_items int) bool {
	if len(this.channel_commands) < num_items {
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

func (this *ChannelCommandQ) Pop() *ChannelCommand {
	if !this.CanPop(1) {
		err := errors.New("channel command queue cannot be popped")
		panic(err)
	}

	channel_command := this.channel_commands[0]

	this.channel_commands = this.channel_commands[1:]
	this.cycles = this.cycles[1:]

	return channel_command
}

func (this *ChannelCommandQ) Front(pos int) (*ChannelCommand, int64) {
	return this.channel_commands[pos], this.cycles[pos]
}

func (this *ChannelCommandQ) Remove(pos int) {
	this.channel_commands = append(this.channel_commands[:pos], this.channel_commands[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *ChannelCommandQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
