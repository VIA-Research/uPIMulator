package directive

import (
	"uPIMulator/src/device/abi"
	"uPIMulator/src/encoding"
)

type ShortDirective struct {
	immediate *abi.Immediate
}

func (this *ShortDirective) Init(value int64) {
	this.immediate = new(abi.Immediate)
	this.immediate.Init(abi.UNSIGNED, 8*int(this.Size()), value)
}

func (this *ShortDirective) Size() int64 {
	return 2
}

func (this *ShortDirective) Immediate() *abi.Immediate {
	return this.immediate
}

func (this *ShortDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
