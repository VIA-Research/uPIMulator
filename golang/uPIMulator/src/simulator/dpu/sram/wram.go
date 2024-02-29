package sram

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
)

type Wram struct {
	address int64
	size    int64

	byte_stream *encoding.ByteStream
}

func (this *Wram) Init() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.address = config_loader.WramOffset()
	this.size = config_loader.WramSize()

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()
	for i := int64(0); i < this.size; i++ {
		this.byte_stream.Append(0)
	}
}

func (this *Wram) Fini() {
}

func (this *Wram) Address() int64 {
	return this.address
}

func (this *Wram) Size() int64 {
	return this.size
}

func (this *Wram) Read(address int64, size int64) *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < size; i++ {
		index := this.Index(address) + int(i)

		byte_stream.Append(this.byte_stream.Get(index))
	}

	return byte_stream
}

func (this *Wram) Write(address int64, size int64, byte_stream *encoding.ByteStream) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	for i := int64(0); i < size; i++ {
		index := this.Index(address) + int(i)

		this.byte_stream.Set(index, byte_stream.Get(int(i)))
	}
}

func (this *Wram) Index(address int64) int {
	if address < this.address {
		err := errors.New("address < WRAM offset")
		panic(err)
	} else if address >= this.address+this.size {
		err := errors.New("address >= WRAM offset + WRAM size")
		panic(err)
	}

	return int(address - this.address)
}
