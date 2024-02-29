package logic

import (
	"errors"
	"fmt"
	"uPIMulator/src/misc"
)

type ThreadScheduler struct {
	channel_id int
	rank_id    int
	dpu_id     int

	num_revolver_scheduling_cycles int64

	threads  []*Thread
	thread_q *ThreadQ

	stat_factory *misc.StatFactory
}

func (this *ThreadScheduler) Init(
	channel_id int,
	rank_id int,
	dpu_id int,
	threads []*Thread,
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

	this.num_revolver_scheduling_cycles = command_line_parser.IntParameter(
		"num_revolver_scheduling_cycles",
	)

	this.threads = threads

	this.thread_q = new(ThreadQ)
	this.thread_q.Init(len(this.threads), 0)
	for _, thread := range threads {
		this.thread_q.Push(thread)
	}

	name := fmt.Sprintf("ThreadScheduler[%d_%d_%d]", channel_id, rank_id, dpu_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *ThreadScheduler) Fini() {
	for this.thread_q.CanPop(1) {
		this.thread_q.Pop()
	}

	this.thread_q.Fini()
}

func (this *ThreadScheduler) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *ThreadScheduler) NumIssuableThreads() int {
	num_issuable_threads := 0

	for _, thread := range this.threads {
		if thread.ThreadState() == RUNNABLE {
			num_issuable_threads++
		}
	}

	return num_issuable_threads
}

func (this *ThreadScheduler) Schedule() *Thread {
	var is_blocked bool
	is_blocked = false
	for i := 0; i < this.thread_q.Size(); i++ {
		thread := this.thread_q.Pop()
		this.thread_q.Push(thread)

		if thread.IssueCycle() >= this.num_revolver_scheduling_cycles {
			if thread.ThreadState() == RUNNABLE {
				thread.ResetIssueCycle()

				this.stat_factory.Increment("breakdown_run", 1)

				return thread
			} else if thread.ThreadState() == BLOCK {
				is_blocked = true
			}
		}
	}

	if is_blocked {
		this.stat_factory.Increment("breakdown_dma", 1)
	} else {
		this.stat_factory.Increment("breakdown_etc", 1)
	}

	return nil
}

func (this *ThreadScheduler) Boot(thread_id int) bool {
	thread := this.threads[thread_id]

	if thread.ThreadId() != thread_id {
		err := errors.New("thread's thread ID != thread ID")
		panic(err)
	}

	thread_state := thread.ThreadState()
	if thread_state == EMBRYO {
		thread.SetThreadState(RUNNABLE)
		return true
	} else if thread_state == ZOMBIE {
		thread.SetThreadState(RUNNABLE)
		return true
	} else {
		err := errors.New("thread is not bootable")
		panic(err)
	}
}

func (this *ThreadScheduler) Sleep(thread_id int) bool {
	thread := this.threads[thread_id]

	if thread.ThreadId() != thread_id {
		err := errors.New("thread's thread ID != thread ID")
		panic(err)
	}

	thread_state := thread.ThreadState()
	if thread_state == RUNNABLE {
		thread.SetThreadState(SLEEP)
		return true
	} else {
		err := errors.New("thread is not sleepable")
		panic(err)
	}
}

func (this *ThreadScheduler) Block(thread_id int) bool {
	thread := this.threads[thread_id]

	if thread.ThreadId() != thread_id {
		err := errors.New("thread's thread ID != thread ID")
		panic(err)
	}

	thread_state := thread.ThreadState()
	if thread_state == RUNNABLE {
		thread.SetThreadState(BLOCK)
		return true
	} else {
		err := errors.New("thread is not blockable")
		panic(err)
	}
}

func (this *ThreadScheduler) Awake(thread_id int) bool {
	thread := this.threads[thread_id]

	if thread.ThreadId() != thread_id {
		err := errors.New("thread's thread ID != thread ID")
		panic(err)
	}

	thread_state := thread.ThreadState()
	if thread_state == EMBRYO {
		thread.SetThreadState(RUNNABLE)
		return true
	} else if thread_state == SLEEP {
		thread.SetThreadState(RUNNABLE)
		return true
	} else if thread_state == BLOCK {
		thread.SetThreadState(RUNNABLE)
		return true
	} else {
		err := errors.New("thread is not awakable")
		panic(err)
	}
}

func (this *ThreadScheduler) Shutdown(thread_id int) bool {
	thread := this.threads[thread_id]

	if thread.ThreadId() != thread_id {
		err := errors.New("thread's thread ID != thread ID")
		panic(err)
	}

	thread_state := thread.ThreadState()
	if thread_state == SLEEP {
		thread.SetThreadState(ZOMBIE)
		return true
	} else {
		err := errors.New("thread is not shotdownable")
		panic(err)
	}
}

func (this *ThreadScheduler) Cycle() {
}
