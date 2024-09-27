package arena

import (
	"errors"
	"uPIMulator/src/encoding"
)

type Memory struct {
	byte_stream *encoding.ByteStream
}

func (this *Memory) Init(size int64) {
	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()

	for i := int64(0); i < size; i++ {
		this.byte_stream.Append(0)
	}
}

func (this *Memory) Size() int64 {
	return this.byte_stream.Size()
}

func (this *Memory) Resize(size int64) {
	if size < this.byte_stream.Size() {
		err := errors.New("size < byte stream's size")
		panic(err)
	}

	for i := int64(0); i < size-this.byte_stream.Size(); i++ {
		this.byte_stream.Append(0)
	}
}

func (this *Memory) Read(address int64, size int64) *encoding.ByteStream {
	for address+size >= this.byte_stream.Size() {
		this.Resize(2 * this.byte_stream.Size())
	}

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < size; i++ {
		byte_stream.Append(this.byte_stream.Get(int(address + i)))
	}

	return byte_stream
}

func (this *Memory) Write(address int64, size int64, byte_stream *encoding.ByteStream) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	for address+size >= this.byte_stream.Size() {
		this.Resize(2 * this.byte_stream.Size())
	}

	for i := int64(0); i < size; i++ {
		this.byte_stream.Set(int(address+i), byte_stream.Get(int(i)))
	}
}
