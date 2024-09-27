package directive

import (
	"errors"
	"uPIMulator/src/device/abi"
	"uPIMulator/src/encoding"
)

type ZeroDirective struct {
	size      int64
	immediate *abi.Immediate
}

func (this *ZeroDirective) Init(size int64, value int64) {
	if size <= 0 {
		err := errors.New("size <= 0")
		panic(err)
	}

	this.size = size

	this.immediate = new(abi.Immediate)
	this.immediate.Init(abi.UNSIGNED, 8, value)
}

func (this *ZeroDirective) Size() int64 {
	return this.size
}

func (this *ZeroDirective) Immediate() *abi.Immediate {
	return this.immediate
}

func (this *ZeroDirective) Encode() *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < this.size; i++ {
		byte_stream.Merge(this.immediate.ToByteStream())
	}

	return byte_stream
}
