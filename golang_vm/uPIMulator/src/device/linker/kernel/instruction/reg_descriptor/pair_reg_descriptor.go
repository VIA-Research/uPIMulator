package reg_descriptor

import (
	"errors"
)

type PairRegDescriptor struct {
	even_reg_descriptor *GpRegDescriptor
	odd_reg_descriptor  *GpRegDescriptor
}

func (this *PairRegDescriptor) Init(index int) {
	if index%2 != 0 {
		err := errors.New("index %2 != 0")
		panic(err)
	}

	this.even_reg_descriptor = new(GpRegDescriptor)
	this.even_reg_descriptor.Init(index)

	this.odd_reg_descriptor = new(GpRegDescriptor)
	this.odd_reg_descriptor.Init(index + 1)
}

func (this *PairRegDescriptor) Index() int {
	return this.even_reg_descriptor.Index()
}

func (this *PairRegDescriptor) EvenRegDescriptor() *GpRegDescriptor {
	return this.even_reg_descriptor
}

func (this *PairRegDescriptor) OddRegDescriptor() *GpRegDescriptor {
	return this.odd_reg_descriptor
}
