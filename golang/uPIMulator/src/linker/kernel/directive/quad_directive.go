package directive

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
)

type QuadDirective struct {
	immediate *word.Immediate
}

func (this *QuadDirective) Init(value int64) {
	this.immediate = new(word.Immediate)
	this.immediate.Init(word.UNSIGNED, 8*int(this.Size()), value)
}

func (this *QuadDirective) Size() int64 {
	return 8
}

func (this *QuadDirective) Immediate() *word.Immediate {
	return this.immediate
}

func (this *QuadDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
