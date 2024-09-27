package directive

import (
	"uPIMulator/src/device/abi"
	"uPIMulator/src/encoding"
)

type LongDirective struct {
	immediate *abi.Immediate
}

func (this *LongDirective) Init(value int64) {
	this.immediate = new(abi.Immediate)
	this.immediate.Init(abi.UNSIGNED, 8*int(this.Size()), value)
}

func (this *LongDirective) Size() int64 {
	return 4
}

func (this *LongDirective) Immediate() *abi.Immediate {
	return this.immediate
}

func (this *LongDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
