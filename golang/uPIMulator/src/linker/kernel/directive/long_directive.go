package directive

import (
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
)

type LongDirective struct {
	immediate *word.Immediate
}

func (this *LongDirective) Init(value int64) {
	this.immediate = new(word.Immediate)
	this.immediate.Init(word.UNSIGNED, 8*int(this.Size()), value)
}

func (this *LongDirective) Size() int64 {
	return 4
}

func (this *LongDirective) Immediate() *word.Immediate {
	return this.immediate
}

func (this *LongDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
