package directive

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
)

type ShortDirective struct {
	immediate *word.Immediate
}

func (this *ShortDirective) Init(value int64) {
	this.immediate = new(word.Immediate)
	this.immediate.Init(word.UNSIGNED, 8*int(this.Size()), value)
}

func (this *ShortDirective) Size() int64 {
	return 2
}

func (this *ShortDirective) Immediate() *word.Immediate {
	return this.immediate
}

func (this *ShortDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
