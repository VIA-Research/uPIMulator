package word

import (
	"uPIMulator/src/abi/encoding"
)

type Immediate struct {
	representation Representation
	value          int64
	word           *Word
}

func (this *Immediate) Init(representation Representation, width int, value int64) {
	this.representation = representation
	this.value = value

	this.word = new(Word)
	this.word.Init(width)
	this.word.SetValue(value)
}

func (this *Immediate) Representation() Representation {
	return this.representation
}

func (this *Immediate) Width() int {
	return this.word.Width()
}

func (this *Immediate) Bit(pos int) bool {
	return this.word.Bit(pos)
}

func (this *Immediate) BitSlice(begin int, end int) int64 {
	return this.word.BitSlice(this.representation, begin, end)
}

func (this *Immediate) Value() int64 {
	return this.value
}

func (this *Immediate) ToByteStream() *encoding.ByteStream {
	return this.word.ToByteStream()
}
