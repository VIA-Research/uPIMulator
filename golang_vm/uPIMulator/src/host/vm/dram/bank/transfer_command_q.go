package bank

import (
	"errors"
)

type TransferCommandQ struct {
	size  int
	timer int64

	transfer_commands []*TransferCommand
	cycles            []int64
}

func (this *TransferCommandQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.transfer_commands = make([]*TransferCommand, 0)
	this.cycles = make([]int64, 0)
}

func (this *TransferCommandQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("transfer command queue is not empty")
		panic(err)
	}
}

func (this *TransferCommandQ) Size() int {
	return this.size
}

func (this *TransferCommandQ) Length() int {
	return len(this.transfer_commands)
}

func (this *TransferCommandQ) Timer() int64 {
	return this.timer
}

func (this *TransferCommandQ) IsEmpty() bool {
	return len(this.transfer_commands) == 0
}

func (this *TransferCommandQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.transfer_commands) >= num_items
	} else {
		return true
	}
}

func (this *TransferCommandQ) Push(transfer_command *TransferCommand) {
	if !this.CanPush(1) {
		err := errors.New("transfer command queue cannot be pushed")
		panic(err)
	}

	this.transfer_commands = append(this.transfer_commands, transfer_command)
	this.cycles = append(this.cycles, this.timer)
}

func (this *TransferCommandQ) PushWithTimer(transfer_command *TransferCommand, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("transfer command queue cannot be pushed")
		panic(err)
	}

	this.transfer_commands = append(this.transfer_commands, transfer_command)
	this.cycles = append(this.cycles, timer)
}

func (this *TransferCommandQ) CanPop(num_items int) bool {
	if len(this.transfer_commands) < num_items {
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

func (this *TransferCommandQ) Pop() *TransferCommand {
	if !this.CanPop(1) {
		err := errors.New("transfer command queue cannot be popped")
		panic(err)
	}

	transfer_command := this.transfer_commands[0]

	this.transfer_commands = this.transfer_commands[1:]
	this.cycles = this.cycles[1:]

	return transfer_command
}

func (this *TransferCommandQ) Front(pos int) (*TransferCommand, int64) {
	return this.transfer_commands[pos], this.cycles[pos]
}

func (this *TransferCommandQ) Remove(pos int) {
	this.transfer_commands = append(this.transfer_commands[:pos], this.transfer_commands[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *TransferCommandQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
