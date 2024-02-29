package directive

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
)

type ByteDirective struct {
	immediate *word.Immediate
}

func (this *ByteDirective) Init(value int64) {
	this.immediate = new(word.Immediate)
	this.immediate.Init(word.UNSIGNED, 8*int(this.Size()), value)
}

func (this *ByteDirective) Size() int64 {
	return 1
}

func (this *ByteDirective) Immediate() *word.Immediate {
	return this.immediate
}

func (this *ByteDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
