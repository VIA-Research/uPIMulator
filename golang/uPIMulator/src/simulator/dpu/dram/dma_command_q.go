package dram

import (
	"errors"
)

type DmaCommandQ struct {
	size  int
	timer int64

	dma_commands []*DmaCommand
	cycles       []int64
}

func (this *DmaCommandQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.dma_commands = make([]*DmaCommand, 0)
	this.cycles = make([]int64, 0)
}

func (this *DmaCommandQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("DMA command queue is not empty")
		panic(err)
	}
}

func (this *DmaCommandQ) Size() int {
	return this.size
}

func (this *DmaCommandQ) Timer() int64 {
	return this.timer
}

func (this *DmaCommandQ) IsEmpty() bool {
	return len(this.dma_commands) == 0
}

func (this *DmaCommandQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.dma_commands) >= num_items
	} else {
		return true
	}
}

func (this *DmaCommandQ) Push(dma_command *DmaCommand) {
	if !this.CanPush(1) {
		err := errors.New("DMA command queue cannot be pushed")
		panic(err)
	}

	this.dma_commands = append(this.dma_commands, dma_command)
	this.cycles = append(this.cycles, this.timer)
}

func (this *DmaCommandQ) PushWithTimer(dma_command *DmaCommand, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("DMA command queue cannot be pushed")
		panic(err)
	}

	this.dma_commands = append(this.dma_commands, dma_command)
	this.cycles = append(this.cycles, timer)
}

func (this *DmaCommandQ) CanPop(num_items int) bool {
	if len(this.dma_commands) < num_items {
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

func (this *DmaCommandQ) Pop() *DmaCommand {
	if !this.CanPop(1) {
		err := errors.New("DMA command queue cannot be popped")
		panic(err)
	}

	dma_command := this.dma_commands[0]

	this.dma_commands = this.dma_commands[1:]
	this.cycles = this.cycles[1:]

	return dma_command
}

func (this *DmaCommandQ) Front(pos int) (*DmaCommand, int64) {
	if this.IsEmpty() {
		err := errors.New("DMA command queue is empty")
		panic(err)
	}

	return this.dma_commands[pos], this.cycles[pos]
}

func (this *DmaCommandQ) Remove(pos int) {
	this.dma_commands = append(this.dma_commands[:pos], this.dma_commands[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *DmaCommandQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
