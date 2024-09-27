package directive

import (
	"uPIMulator/src/device/abi"
	"uPIMulator/src/encoding"
)

type ByteDirective struct {
	immediate *abi.Immediate
}

func (this *ByteDirective) Init(value int64) {
	this.immediate = new(abi.Immediate)
	this.immediate.Init(abi.UNSIGNED, 8*int(this.Size()), value)
}

func (this *ByteDirective) Size() int64 {
	return 1
}

func (this *ByteDirective) Immediate() *abi.Immediate {
	return this.immediate
}

func (this *ByteDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
