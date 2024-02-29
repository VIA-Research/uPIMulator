package logic

import (
	"errors"
)

type ThreadQ struct {
	size  int
	timer int64

	threads []*Thread
	cycles  []int64
}

func (this *ThreadQ) Init(size int, timer int64) {
	if size == 0 {
		err := errors.New("size == 0")
		panic(err)
	} else if timer < 0 {
		err := errors.New("timer < 0")
		panic(err)
	}

	this.size = size
	this.timer = timer

	this.threads = make([]*Thread, 0)
	this.cycles = make([]int64, 0)
}

func (this *ThreadQ) Fini() {
	if !this.IsEmpty() {
		err := errors.New("thread queue is not empty")
		panic(err)
	}
}

func (this *ThreadQ) Size() int {
	return this.size
}

func (this *ThreadQ) Timer() int64 {
	return this.timer
}

func (this *ThreadQ) IsEmpty() bool {
	return len(this.threads) == 0
}

func (this *ThreadQ) CanPush(num_items int) bool {
	if this.size >= 0 {
		return this.size-len(this.threads) >= num_items
	} else {
		return true
	}
}

func (this *ThreadQ) Push(thread *Thread) {
	if !this.CanPush(1) {
		err := errors.New("thread queue cannot be pushed")
		panic(err)
	}

	this.threads = append(this.threads, thread)
	this.cycles = append(this.cycles, this.timer)
}

func (this *ThreadQ) PushWithTimer(thread *Thread, timer int64) {
	if !this.CanPush(1) {
		err := errors.New("thread queue cannot be pushed")
		panic(err)
	}

	this.threads = append(this.threads, thread)
	this.cycles = append(this.cycles, timer)
}

func (this *ThreadQ) CanPop(num_items int) bool {
	if len(this.threads) < num_items {
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

func (this *ThreadQ) Pop() *Thread {
	if !this.CanPop(1) {
		err := errors.New("thread queue cannot be popped")
		panic(err)
	}

	thread := this.threads[0]

	this.threads = this.threads[1:]
	this.cycles = this.cycles[1:]

	return thread
}

func (this *ThreadQ) Front(pos int) (*Thread, int64) {
	if this.IsEmpty() {
		err := errors.New("thread queue is empty")
		panic(err)
	}

	return this.threads[pos], this.cycles[pos]
}

func (this *ThreadQ) Remove(pos int) {
	this.threads = append(this.threads[:pos], this.threads[pos+1:]...)
	this.cycles = append(this.cycles[:pos], this.cycles[pos+1:]...)
}

func (this *ThreadQ) Cycle() {
	if !this.IsEmpty() {
		this.cycles[0] -= 1
	}
}
