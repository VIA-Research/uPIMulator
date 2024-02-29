package logic

import (
	"errors"
	"fmt"
	"uPIMulator/src/linker/kernel/instruction"
	"uPIMulator/src/misc"
)

type CycleRule struct {
	channel_id int
	rank_id    int
	dpu_id     int

	input_q *InstructionQ
	wait_q  *InstructionQ
	ready_q *InstructionQ

	scoreboard map[*instruction.Instruction]*Thread
	reg_sets   []*RegSet

	stat_factory *misc.StatFactory
}

func (this *CycleRule) Init(
	channel_id int,
	rank_id int,
	dpu_id int,
	command_line_parser *misc.CommandLineParser,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	} else if dpu_id < 0 {
		err := errors.New("DPU ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.dpu_id = dpu_id

	num_tasklets := int(command_line_parser.IntParameter("num_tasklets"))

	this.input_q = new(InstructionQ)
	this.input_q.Init(1, 0)

	this.wait_q = new(InstructionQ)
	this.wait_q.Init(1, 0)

	this.ready_q = new(InstructionQ)
	this.ready_q.Init(1, 0)

	this.scoreboard = make(map[*instruction.Instruction]*Thread, 0)

	for i := 0; i < num_tasklets; i++ {
		reg_set := new(RegSet)
		reg_set.Init(i)

		this.reg_sets = append(this.reg_sets, reg_set)
	}

	name := fmt.Sprintf("CycleRule[%d_%d_%d]", channel_id, rank_id, dpu_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *CycleRule) Fini() {
	this.input_q.Fini()
	this.wait_q.Fini()
	this.ready_q.Fini()
}

func (this *CycleRule) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *CycleRule) IsEmpty() bool {
	return this.input_q.IsEmpty() && this.wait_q.IsEmpty() && this.ready_q.IsEmpty()
}

func (this *CycleRule) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *CycleRule) Push(instruction_ *instruction.Instruction, thread *Thread) {
	if !this.CanPush() {
		err := errors.New("cycle rule cannot be pushed")
		panic(err)
	}

	this.input_q.Push(instruction_)
	this.scoreboard[instruction_] = thread
}

func (this *CycleRule) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *CycleRule) Pop() *instruction.Instruction {
	if !this.CanPop() {
		err := errors.New("cycle rule cannot be popped")
		panic(err)
	}

	instruction_ := this.ready_q.Pop()

	delete(this.scoreboard, instruction_)
	return instruction_
}

func (this *CycleRule) Cycle() {
	this.ServiceInputQ()
	this.ServiceWaitQ()

	this.input_q.Cycle()
	this.wait_q.Cycle()
	this.ready_q.Cycle()
}

func (this *CycleRule) ServiceInputQ() {
	if this.input_q.CanPop(1) && this.wait_q.CanPush(1) {
		instruction_ := this.input_q.Pop()

		thread_id := this.scoreboard[instruction_].ThreadId()
		this.reg_sets[thread_id].CollectReadGpRegs(instruction_)

		extra_cycles := this.CalculateExtraCycles(instruction_)

		this.wait_q.PushWithTimer(instruction_, extra_cycles)

		this.stat_factory.Increment("cycle_rule", extra_cycles)
	}
}

func (this *CycleRule) ServiceWaitQ() {
	if this.wait_q.CanPop(1) && this.ready_q.CanPush(1) {
		instruction_ := this.wait_q.Pop()
		this.ready_q.Push(instruction_)

		thread_id := this.scoreboard[instruction_].ThreadId()
		this.reg_sets[thread_id].Clear()
		this.reg_sets[thread_id].CollectWriteGpRegs(instruction_)
	}
}

func (this *CycleRule) CalculateExtraCycles(instruction_ *instruction.Instruction) int64 {
	thread_id := this.scoreboard[instruction_].ThreadId()

	reg_set := this.reg_sets[thread_id]

	if reg_set.ThreadId() != thread_id {
		err := errors.New("reg set's thread ID != thread ID")
		panic(err)
	}

	reg_indicies := reg_set.RegIndices()

	even_counter := 0
	odd_counter := 0
	for reg_index, _ := range reg_indicies {
		if reg_index%2 == 0 {
			even_counter++
		} else {
			odd_counter++
		}
	}

	return int64(even_counter/2 + odd_counter/2)
}
