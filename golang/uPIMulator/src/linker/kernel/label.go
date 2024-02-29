package kernel

import (
	"errors"
	"uPIMulator/src/abi/encoding"
)

type Label struct {
	name        string
	address     *int64
	size        int64
	byte_stream *encoding.ByteStream
}

func (this *Label) Init(name string) {
	this.name = name
	this.address = nil
	this.size = 0

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()
}

func (this *Label) Name() string {
	return this.name
}

func (this *Label) Address() int64 {
	if this.address == nil {
		err := errors.New("address is not yet set")
		panic(err)
	}

	return *this.address
}

func (this *Label) BeginAddress() int64 {
	return this.Address()
}

func (this *Label) EndAddress() int64 {
	return this.Address() + this.Size()
}

func (this *Label) SetAddress(address int64) {
	if this.address != nil {
		err := errors.New("address is already set")
		panic(err)
	}

	this.address = new(int64)
	*this.address = address
}

func (this *Label) Size() int64 {
	return this.size
}

func (this *Label) SetSize(size int64) {
	this.size = size
}

func (this *Label) ToByteStream() *encoding.ByteStream {
	if this.size != this.byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	return this.byte_stream
}

func (this *Label) Append(encodable Encodable) {
	this.byte_stream.Merge(encodable.Encode())

	if this.byte_stream.Size() > this.size {
		err := errors.New("byte stream's size > size")
		panic(err)
	}
}
