package dram

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/linker/kernel/instruction"
)

type DmaCommand struct {
	memory_operation MemoryOperation
	wram_address     *int64
	mram_address     *int64
	size             int64

	byte_stream *encoding.ByteStream
	acks        []bool

	instruction *instruction.Instruction
}

func (this *DmaCommand) InitReadFromMram(mram_address int64, size int64) {
	this.memory_operation = READ
	this.wram_address = nil

	this.mram_address = new(int64)
	*this.mram_address = mram_address

	this.size = size

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()
	for i := int64(0); i < size; i++ {
		this.byte_stream.Append(0)
	}

	this.acks = make([]bool, 0)
	for i := int64(0); i < size; i++ {
		this.acks = append(this.acks, false)
	}

	this.instruction = nil
}

func (this *DmaCommand) InitWriteToMram(
	mram_address int64,
	size int64,
	byte_stream *encoding.ByteStream,
) {
	this.memory_operation = WRITE
	this.wram_address = nil

	this.mram_address = new(int64)
	*this.mram_address = mram_address

	this.size = size
	this.byte_stream = byte_stream

	this.acks = make([]bool, 0)
	for i := int64(0); i < size; i++ {
		this.acks = append(this.acks, false)
	}

	this.instruction = nil
}

func (this *DmaCommand) InitReadFromMramToWram(
	wram_address int64,
	mram_address int64,
	size int64,
	instruction_ *instruction.Instruction,
) {
	if instruction_.OpCode() != instruction.LDMA {
		err := errors.New("instruction's op code != LDMA")
		panic(err)
	}

	this.memory_operation = READ

	this.wram_address = new(int64)
	*this.wram_address = wram_address

	this.mram_address = new(int64)
	*this.mram_address = mram_address

	this.size = size

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()
	for i := int64(0); i < size; i++ {
		this.byte_stream.Append(0)
	}

	this.acks = make([]bool, 0)
	for i := int64(0); i < size; i++ {
		this.acks = append(this.acks, false)
	}

	this.instruction = instruction_
}

func (this *DmaCommand) InitWriteToMramFromWram(
	wram_address int64,
	mram_address int64,
	size int64,
	byte_stream *encoding.ByteStream,
	instruction_ *instruction.Instruction,
) {
	if instruction_.OpCode() != instruction.SDMA {
		err := errors.New("instruction's op code != SDMA")
		panic(err)
	}

	this.memory_operation = WRITE

	this.wram_address = new(int64)
	*this.wram_address = wram_address

	this.mram_address = new(int64)
	*this.mram_address = mram_address

	this.size = size
	this.byte_stream = byte_stream

	this.acks = make([]bool, 0)
	for i := int64(0); i < size; i++ {
		this.acks = append(this.acks, false)
	}

	this.instruction = instruction_
}

func (this *DmaCommand) Fini() {
	if !this.IsReady() {
		err := errors.New("DMA command is not ready")
		panic(err)
	}
}

func (this *DmaCommand) MemoryOperation() MemoryOperation {
	return this.memory_operation
}

func (this *DmaCommand) WramAddress() int64 {
	if this.wram_address == nil {
		err := errors.New("DMA command does not have a WRAM address")
		panic(err)
	}

	return *this.wram_address
}

func (this *DmaCommand) MramAddress() int64 {
	if this.mram_address == nil {
		err := errors.New("DMA command does not have an MRAM address")
		panic(err)
	}

	return *this.mram_address
}

func (this *DmaCommand) Size() int64 {
	return this.size
}

func (this *DmaCommand) HasInstruction() bool {
	return this.instruction != nil
}

func (this *DmaCommand) Instruction() *instruction.Instruction {
	return this.instruction
}

func (this *DmaCommand) ByteStream(mram_address int64, size int64) *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < size; i++ {
		index := this.Index(mram_address) + int(i)

		byte_stream.Append(this.byte_stream.Get(index))
	}

	return byte_stream
}

func (this *DmaCommand) SetByteStream(
	mram_address int64,
	size int64,
	byte_stream *encoding.ByteStream,
) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	for i := int64(0); i < byte_stream.Size(); i++ {
		index := this.Index(mram_address) + int(i)

		this.byte_stream.Set(index, byte_stream.Get(int(i)))
	}
}

func (this *DmaCommand) SetAck(mram_address int64, size int64) {
	for i := int64(0); i < size; i++ {
		index := this.Index(mram_address) + int(i)

		if this.acks[index] {
			err := errors.New("ACK is already set")
			panic(err)
		}

		this.acks[index] = true
	}
}

func (this *DmaCommand) IsReady() bool {
	for _, ack := range this.acks {
		if !ack {
			return false
		}
	}
	return true
}

func (this *DmaCommand) Index(mram_address int64) int {
	if mram_address < this.MramAddress() {
		err := errors.New("MRAM address < DMA command's MRAM address")
		panic(err)
	} else if mram_address >= this.MramAddress()+this.Size() {
		err := errors.New("MRAM address >= DMA command's MRAM address + DMA command's size")
		panic(err)
	}

	return int(mram_address - this.MramAddress())
}
