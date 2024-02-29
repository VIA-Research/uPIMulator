package kernel

import (
	"errors"
	"uPIMulator/src/abi/encoding"
)

type Kernel struct {
	addresses map[string]int64

	atomic *encoding.ByteStream
	iram   *encoding.ByteStream
	wram   *encoding.ByteStream
	mram   *encoding.ByteStream
}

func (this *Kernel) Init() {
	this.addresses = make(map[string]int64, 0)
}

func (this *Kernel) Address(label_name string) int64 {
	if address, found := this.addresses[label_name]; found {
		return address
	} else {
		err := errors.New("address is not found")
		panic(err)
	}
}

func (this *Kernel) Atomic() *encoding.ByteStream {
	return this.atomic
}

func (this *Kernel) SetAtomic(atomic *encoding.ByteStream) {
	this.atomic = atomic
}

func (this *Kernel) Iram() *encoding.ByteStream {
	return this.iram
}

func (this *Kernel) SetIram(iram *encoding.ByteStream) {
	this.iram = iram
}

func (this *Kernel) Wram() *encoding.ByteStream {
	return this.wram
}

func (this *Kernel) SetWram(wram *encoding.ByteStream) {
	this.wram = wram
}

func (this *Kernel) Mram() *encoding.ByteStream {
	return this.mram
}

func (this *Kernel) SetMram(mram *encoding.ByteStream) {
	this.mram = mram
}
