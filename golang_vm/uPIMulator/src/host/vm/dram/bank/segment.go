package bank

import (
	"errors"
)

type Segment struct {
	vm_address int64

	channel_id   int
	rank_id      int
	bank_id      int
	bank_address int64

	size int64
}

func (this *Segment) Init(
	vm_address int64,
	channel_id int,
	rank_id int,
	bank_id int,
	bank_address int64,
	size int64,
) {
	if vm_address < 0 {
		err := errors.New("VM segment < 0")
		panic(err)
	} else if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	} else if bank_id < 0 {
		err := errors.New("DPU ID < 0")
		panic(err)
	} else if bank_address < 0 {
		err := errors.New("bank segment < 0")
		panic(err)
	} else if size < 0 {
		err := errors.New("size < 0")
		panic(err)
	}

	this.vm_address = vm_address

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.bank_id = bank_id
	this.bank_address = bank_address
	this.size = size
}

func (this *Segment) VmAddress() int64 {
	return this.vm_address
}

func (this *Segment) ChannelID() int {
	return this.channel_id
}

func (this *Segment) RankID() int {
	return this.rank_id
}

func (this *Segment) BankID() int {
	return this.bank_id
}

func (this *Segment) BankAddress() int64 {
	return this.bank_address
}

func (this *Segment) Size() int64 {
	return this.size
}
