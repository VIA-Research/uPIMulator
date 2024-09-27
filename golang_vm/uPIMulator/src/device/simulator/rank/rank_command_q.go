package rank

import (
	"errors"
)

type RankCommandQ struct {
	size  int
	timer int64

	rank_commands []*RankCommand
	cycles        []int64
}

func (this *RankCommandQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.rank_commands = make([]*RankCommand, 0)
	this.cycles = make([]int64, 0)
}

func (this *RankCommandQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("rank command queue is not empty")
		panic(err)
	}
}

func (this *RankCommandQ) Size() int {
	return this.size
}

func (this *RankCommandQ) Length() int {
	return len(this.rank_commands)
}

func (this *RankCommandQ) Timer() int64 {
	return this.timer
}

func (this *RankCommandQ) IsEmpty() bool {
	return len(this.rank_commands) == 0
}

func (this *RankCommandQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.rank_commands) >= num_items
	} else {
		return true
	}
}

func (this *RankCommandQ) Push(rank_command *RankCommand) {
	if !this.CanPush(1) {
		err := errors.New("rank command queue cannot be pushed")
		panic(err)
	}

	this.rank_commands = append(this.rank_commands, rank_command)
	this.cycles = append(this.cycles, this.timer)
}

func (this *RankCommandQ) PushWithTimer(rank_command *RankCommand, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("rank command queue cannot be pushed")
		panic(err)
	}

	this.rank_commands = append(this.rank_commands, rank_command)
	this.cycles = append(this.cycles, timer)
}

func (this *RankCommandQ) CanPop(num_items int) bool {
	if len(this.rank_commands) < num_items {
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

func (this *RankCommandQ) Pop() *RankCommand {
	if !this.CanPop(1) {
		err := errors.New("rank command queue cannot be popped")
		panic(err)
	}

	rank_command := this.rank_commands[0]

	this.rank_commands = this.rank_commands[1:]
	this.cycles = this.cycles[1:]

	return rank_command
}

func (this *RankCommandQ) Front(pos int) (*RankCommand, int64) {
	return this.rank_commands[pos], this.cycles[pos]
}

func (this *RankCommandQ) Remove(pos int) {
	this.rank_commands = append(this.rank_commands[:pos], this.rank_commands[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *RankCommandQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
