package sram

import (
	"errors"
	"uPIMulator/src/misc"
)

type Atomic struct {
	address int64
	size    int64

	locks []*Lock
}

func (this *Atomic) Init() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.address = config_loader.AtomicOffset()
	this.size = config_loader.AtomicSize()

	this.locks = make([]*Lock, 0)
	for i := int64(0); i < this.size; i++ {
		lock := new(Lock)
		lock.Init()

		this.locks = append(this.locks, lock)
	}
}

func (this *Atomic) Fini() {
	for _, lock := range this.locks {
		lock.Fini()
	}
}

func (this *Atomic) Address() int64 {
	return this.address
}

func (this *Atomic) Size() int64 {
	return this.size
}

func (this *Atomic) CanAcquire(address int64) bool {
	return this.locks[this.Index(address)].CanAcquire()
}

func (this *Atomic) Acquire(address int64, thread_id int) {
	this.locks[this.Index(address)].Acquire(thread_id)
}

func (this *Atomic) CanRelease(address int64, thread_id int) bool {
	return this.locks[this.Index(address)].CanRelease(thread_id)
}

func (this *Atomic) Release(address int64, thread_id int) {
	this.locks[this.Index(address)].Release(thread_id)
}

func (this *Atomic) Index(address int64) int {
	if address < this.address {
		err := errors.New("address < atomic offset")
		panic(err)
	} else if address >= this.address+this.size {
		err := errors.New("address >= atomic offset + atomic size")
		panic(err)
	}

	return int(address - this.address)
}
