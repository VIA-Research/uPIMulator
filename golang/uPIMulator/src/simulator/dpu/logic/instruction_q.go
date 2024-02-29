package logic

import (
	"errors"
	"uPIMulator/src/linker/kernel/instruction"
)

type InstructionQ struct {
	size  int
	timer int64

	instructions []*instruction.Instruction
	cycles       []int64
}

func (this *InstructionQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.instructions = make([]*instruction.Instruction, 0)
	this.cycles = make([]int64, 0)
}

func (this *InstructionQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("instruction queue is not empty")
		panic(err)
	}
}

func (this *InstructionQ) Size() int {
	return this.size
}

func (this *InstructionQ) Timer() int64 {
	return this.timer
}

func (this *InstructionQ) IsEmpty() bool {
	return len(this.instructions) == 0
}

func (this *InstructionQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.instructions) >= num_items
	} else {
		return true
	}
}

func (this *InstructionQ) Push(instruction_ *instruction.Instruction) {
	if !this.CanPush(1) {
		err := errors.New("instruction queue cannot be pushed")
		panic(err)
	}

	this.instructions = append(this.instructions, instruction_)
	this.cycles = append(this.cycles, this.timer)
}

func (this *InstructionQ) PushWithTimer(instruction_ *instruction.Instruction, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("instruction queue cannot be pushed")
		panic(err)
	}

	this.instructions = append(this.instructions, instruction_)
	this.cycles = append(this.cycles, timer)
}

func (this *InstructionQ) CanPop(num_items int) bool {
	if len(this.instructions) < num_items {
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

func (this *InstructionQ) Pop() *instruction.Instruction {
	if !this.CanPop(1) {
		err := errors.New("instruction queue cannot be popped")
		panic(err)
	}

	instruction_ := this.instructions[0]

	this.instructions = this.instructions[1:]
	this.cycles = this.cycles[1:]

	return instruction_
}

func (this *InstructionQ) Front(pos int) (*instruction.Instruction, int64) {
	if this.IsEmpty() {
		err := errors.New("instruction queue is empty")
		panic(err)
	}

	return this.instructions[pos], this.cycles[pos]
}

func (this *InstructionQ) Remove(pos int) {
	this.instructions = append(this.instructions[:pos], this.instructions[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *InstructionQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
