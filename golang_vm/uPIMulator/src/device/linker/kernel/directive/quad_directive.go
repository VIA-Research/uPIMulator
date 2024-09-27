package directive

import (
	"uPIMulator/src/device/abi"
	"uPIMulator/src/encoding"
)

type QuadDirective struct {
	immediate *abi.Immediate
}

func (this *QuadDirective) Init(value int64) {
	this.immediate = new(abi.Immediate)
	this.immediate.Init(abi.UNSIGNED, 8*int(this.Size()), value)
}

func (this *QuadDirective) Size() int64 {
	return 8
}

func (this *QuadDirective) Immediate() *abi.Immediate {
	return this.immediate
}

func (this *QuadDirective) Encode() *encoding.ByteStream {
	return this.immediate.ToByteStream()
}
