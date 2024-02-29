package logic

import (
	"errors"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu/reg"
)

type ThreadState int

const (
	EMBRYO ThreadState = iota
	RUNNABLE
	SLEEP
	BLOCK
	ZOMBIE
)

type Thread struct {
	thread_id    int
	thread_state ThreadState
	reg_file     *reg.RegFile
	issue_cycle  int64
}

func (this *Thread) Init(thread_id int) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	if thread_id < 0 {
		err := errors.New("thread ID < 0")
		panic(err)
	} else if thread_id >= config_loader.MaxNumTasklets() {
		err := errors.New("thread ID >= max number of tasklets")
		panic(err)
	}

	this.thread_id = thread_id
	this.thread_state = EMBRYO

	this.reg_file = new(reg.RegFile)
	this.reg_file.Init(thread_id)

	this.issue_cycle = 0
}

func (this *Thread) Fini() {
	if this.thread_state != ZOMBIE {
		err := errors.New("thread state is not zombie")
		panic(err)
	}
}

func (this *Thread) ThreadId() int {
	return this.thread_id
}

func (this *Thread) ThreadState() ThreadState {
	return this.thread_state
}

func (this *Thread) SetThreadState(thread_state ThreadState) {
	this.thread_state = thread_state
}

func (this *Thread) RegFile() *reg.RegFile {
	return this.reg_file
}

func (this *Thread) IssueCycle() int64 {
	return this.issue_cycle
}

func (this *Thread) IncrementIssueCycle() {
	this.issue_cycle++
}

func (this *Thread) ResetIssueCycle() {
	this.issue_cycle = 0
}
