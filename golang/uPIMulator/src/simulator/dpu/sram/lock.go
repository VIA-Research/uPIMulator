package sram

import (
	"errors"
)

type Lock struct {
	thread_id *int
}

func (this *Lock) Init() {
	this.thread_id = nil
}

func (this *Lock) Fini() {
	if this.thread_id != nil {
		err := errors.New("thread ID != nil")
		panic(err)
	}
}

func (this *Lock) CanAcquire() bool {
	return this.thread_id == nil
}

func (this *Lock) Acquire(thread_id int) {
	if !this.CanAcquire() {
		err := errors.New("lock cannot be acquired")
		panic(err)
	}

	this.thread_id = new(int)
	*this.thread_id = thread_id
}

func (this *Lock) CanRelease(thread_id int) bool {
	return this.thread_id == nil || *this.thread_id == thread_id
}

func (this *Lock) Release(thread_id int) {
	if !this.CanRelease(thread_id) {
		err := errors.New("lock cannot be released")
		panic(err)
	}

	this.thread_id = nil
}
