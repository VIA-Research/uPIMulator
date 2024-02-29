package dram

import (
	"errors"
)

type MemoryCommandQ struct {
	size  int
	timer int64

	memory_commands []*MemoryCommand
	cycles          []int64
}

func (this *MemoryCommandQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.memory_commands = make([]*MemoryCommand, 0)
	this.cycles = make([]int64, 0)
}

func (this *MemoryCommandQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("memory command queue is not empty")
		panic(err)
	}
}

func (this *MemoryCommandQ) Size() int {
	return this.size
}

func (this *MemoryCommandQ) Timer() int64 {
	return this.timer
}

func (this *MemoryCommandQ) IsEmpty() bool {
	return len(this.memory_commands) == 0
}

func (this *MemoryCommandQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.memory_commands) >= num_items
	} else {
		return true
	}
}

func (this *MemoryCommandQ) Push(memory_command *MemoryCommand) {
	if !this.CanPush(1) {
		err := errors.New("memory command queue cannot be pushed")
		panic(err)
	}

	this.memory_commands = append(this.memory_commands, memory_command)
	this.cycles = append(this.cycles, this.timer)
}

func (this *MemoryCommandQ) PushWithTimer(memory_command *MemoryCommand, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("memory command queue cannot be pushed")
		panic(err)
	}

	this.memory_commands = append(this.memory_commands, memory_command)
	this.cycles = append(this.cycles, timer)
}

func (this *MemoryCommandQ) CanPop(num_items int) bool {
	if len(this.memory_commands) < num_items {
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

func (this *MemoryCommandQ) Pop() *MemoryCommand {
	if !this.CanPop(1) {
		err := errors.New("memory command queue cannot be popped")
		panic(err)
	}

	memory_command := this.memory_commands[0]

	this.memory_commands = this.memory_commands[1:]
	this.cycles = this.cycles[1:]

	return memory_command
}

func (this *MemoryCommandQ) Front(pos int) (*MemoryCommand, int64) {
	if this.IsEmpty() {
		err := errors.New("memory command queue is empty")
		panic(err)
	}

	return this.memory_commands[pos], this.cycles[pos]
}

func (this *MemoryCommandQ) Remove(pos int) {
	this.memory_commands = append(this.memory_commands[:pos], this.memory_commands[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *MemoryCommandQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
