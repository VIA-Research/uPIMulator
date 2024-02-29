package logic

import (
	"errors"
	"uPIMulator/src/linker/kernel/instruction"
	"uPIMulator/src/misc"
)

type Pipeline struct {
	num_pipeline_stages int

	input_q *InstructionQ
	wait_q  *InstructionQ
	ready_q *InstructionQ
}

func (this *Pipeline) Init(command_line_parser *misc.CommandLineParser) {
	this.num_pipeline_stages = int(command_line_parser.IntParameter("num_pipeline_stages"))

	this.input_q = new(InstructionQ)
	this.input_q.Init(1, 0)

	this.wait_q = new(InstructionQ)
	this.wait_q.Init(this.num_pipeline_stages-1, 0)
	for this.wait_q.CanPush(1) {
		this.wait_q.Push(nil)
	}

	this.ready_q = new(InstructionQ)
	this.ready_q.Init(1, 0)
	for this.ready_q.CanPush(1) {
		this.ready_q.Push(nil)
	}
}

func (this *Pipeline) Fini() {
	this.input_q.Fini()

	for this.wait_q.CanPop(1) {
		if this.wait_q.Pop() != nil {
			err := errors.New("wait queue is not empty")
			panic(err)
		}
	}
	this.wait_q.Fini()

	for this.ready_q.CanPop(1) {
		if this.ready_q.Pop() != nil {
			err := errors.New("ready queue is not empty")
			panic(err)
		}
	}
	this.ready_q.Fini()
}

func (this *Pipeline) IsEmpty() bool {
	return this.IsInputQEmpty() && this.IsWaitQEmpty() && this.IsReadyQEmpty()
}

func (this *Pipeline) IsInputQEmpty() bool {
	return this.input_q.IsEmpty()
}

func (this *Pipeline) IsWaitQEmpty() bool {
	if this.wait_q.IsEmpty() {
		return true
	} else {
		for i := 0; this.wait_q.CanPop(i + 1); i++ {
			instruction_, _ := this.wait_q.Front(i)

			if instruction_ != nil {
				return false
			}
		}

		return true
	}
}

func (this *Pipeline) IsReadyQEmpty() bool {
	if this.ready_q.IsEmpty() {
		return true
	} else {
		for i := 0; this.ready_q.CanPop(i + 1); i++ {
			instruction_, _ := this.ready_q.Front(i)

			if instruction_ != nil {
				return false
			}
		}
		return true
	}
}

func (this *Pipeline) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *Pipeline) Push(instruction_ *instruction.Instruction) {
	if !this.CanPush() {
		err := errors.New("pipeline cannot be pushed")
		panic(err)
	} else if instruction_ == nil {
		err := errors.New("instruction == nil")
		panic(err)
	}

	this.input_q.Push(instruction_)
}

func (this *Pipeline) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *Pipeline) Pop() *instruction.Instruction {
	if !this.CanPop() {
		err := errors.New("pipeline cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *Pipeline) Cycle() {
	this.ServiceInputQ()
	this.ServiceWaitQ()

	this.input_q.Cycle()
	this.wait_q.Cycle()
	this.ready_q.Cycle()
}

func (this *Pipeline) ServiceInputQ() {
	if this.input_q.CanPop(1) && this.wait_q.CanPush(1) {
		instruction_ := this.input_q.Pop()
		this.wait_q.Push(instruction_)
	} else if this.wait_q.CanPush(1) {
		this.wait_q.Push(nil)
	}
}

func (this *Pipeline) ServiceWaitQ() {
	if this.wait_q.CanPop(1) && this.ready_q.CanPush(1) {
		instruction_ := this.wait_q.Pop()
		this.ready_q.Push(instruction_)
	}
}
