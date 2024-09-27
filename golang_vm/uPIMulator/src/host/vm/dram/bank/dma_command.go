package bank

import (
	"errors"
	"uPIMulator/src/encoding"
)

type DmaCommand struct {
	memory_operation MemoryOperation

	segment *Segment

	byte_stream *encoding.ByteStream
	acks        []bool

	transfer_command *TransferCommand
}

func (this *DmaCommand) InitRead(segment *Segment, transfer_command *TransferCommand) {
	this.memory_operation = READ

	this.segment = segment

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()

	this.acks = make([]bool, 0)
	for i := int64(0); i < segment.Size(); i++ {
		this.byte_stream.Append(0)
		this.acks = append(this.acks, false)
	}

	this.transfer_command = transfer_command
}

func (this *DmaCommand) InitWrite(segment *Segment, transfer_command *TransferCommand) {
	this.memory_operation = WRITE

	this.segment = segment

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()

	this.acks = make([]bool, 0)
	for i := int64(0); i < segment.Size(); i++ {
		this.byte_stream.Append(0)
		this.acks = append(this.acks, false)
	}

	this.transfer_command = transfer_command
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

func (this *DmaCommand) Segment() *Segment {
	return this.segment
}

func (this *DmaCommand) ByteStream(bank_address int64, size int64) *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < size; i++ {
		index := this.Index(bank_address) + int(i)

		byte_stream.Append(this.byte_stream.Get(index))
	}

	return byte_stream
}

func (this *DmaCommand) SetByteStream(
	bank_address int64,
	size int64,
	byte_stream *encoding.ByteStream,
) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	for i := int64(0); i < byte_stream.Size(); i++ {
		index := this.Index(bank_address) + int(i)

		this.byte_stream.Set(index, byte_stream.Get(int(i)))
	}
}

func (this *DmaCommand) SetAck(bank_address int64, size int64) {
	for i := int64(0); i < size; i++ {
		index := this.Index(bank_address) + int(i)

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

func (this *DmaCommand) TransferCommand() *TransferCommand {
	return this.transfer_command
}

func (this *DmaCommand) Index(bank_address int64) int {
	if bank_address < this.segment.BankAddress() {
		err := errors.New("bank address < DMA command's segment's bank address")
		panic(err)
	} else if bank_address >= this.segment.BankAddress()+this.segment.Size() {
		err := errors.New("bank address >= DMA command's segment's bank address + DMA command's size")
		panic(err)
	}

	return int(bank_address - this.segment.BankAddress())
}
