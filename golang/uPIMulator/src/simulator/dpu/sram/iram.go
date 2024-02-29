package sram

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/linker/kernel/instruction"
	"uPIMulator/src/misc"
)

type Iram struct {
	address int64
	size    int64

	byte_stream *encoding.ByteStream
}

func (this *Iram) Init() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.address = config_loader.IramOffset()
	this.size = config_loader.IramSize()

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()
	for i := int64(0); i < this.size; i++ {
		this.byte_stream.Append(0)
	}
}

func (this *Iram) Fini() {
}

func (this *Iram) Address() int64 {
	return this.address
}

func (this *Iram) Size() int64 {
	return this.size
}

func (this *Iram) Read(address int64) *instruction.Instruction {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	iram_data_size := int64(config_loader.IramDataWidth() / 8)

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()
	for i := int64(0); i < iram_data_size; i++ {
		index := this.Index(address) + int(i)

		byte_stream.Append(this.byte_stream.Get(index))
	}

	instruction_ := new(instruction.Instruction)
	instruction_.Decode(byte_stream)
	return instruction_
}

func (this *Iram) Write(address int64, byte_stream *encoding.ByteStream) {
	for i := int64(0); i < byte_stream.Size(); i++ {
		index := this.Index(address) + int(i)

		this.byte_stream.Set(index, byte_stream.Get(int(i)))
	}
}

func (this *Iram) Index(address int64) int {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	iram_data_size := int64(config_loader.IramDataWidth() / 8)

	if address < this.address {
		err := errors.New("address < IRAM offset")
		panic(err)
	} else if address+iram_data_size > this.address+this.size {
		err := errors.New("address >= IRAM offset + IRAM size")
		panic(err)
	}

	if (address-this.address)%iram_data_size != 0 {
		err := errors.New("addresses are not aligned with IRAM data size")
		panic(err)
	}

	return int(address - this.address)
}
