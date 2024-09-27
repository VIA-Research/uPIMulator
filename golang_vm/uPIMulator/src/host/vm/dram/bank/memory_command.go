package bank

import (
	"errors"
	"uPIMulator/src/encoding"
)

type MemoryOperation int

const (
	ACTIVATION MemoryOperation = iota
	READ
	WRITE
	PRECHARGE
)

type MemoryCommand struct {
	memory_operation MemoryOperation
	bank_address     int64
	size             int64
	byte_stream      *encoding.ByteStream
	dma_command      *DmaCommand
}

func (this *MemoryCommand) InitActivation(memory_operation MemoryOperation, bank_address int64) {
	this.memory_operation = memory_operation
	this.bank_address = bank_address
	this.size = 0
	this.byte_stream = nil
	this.dma_command = nil
}

func (this *MemoryCommand) InitRead(
	memory_operation MemoryOperation,
	bank_address int64,
	size int64,
	dma_command *DmaCommand,
) {
	if memory_operation != READ {
		err := errors.New("memory operation != READ")
		panic(err)
	} else if dma_command.MemoryOperation() != READ {
		err := errors.New("DMA operation != READ")
		panic(err)
	}

	this.memory_operation = memory_operation
	this.bank_address = bank_address
	this.size = size

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()
	for i := int64(0); i < size; i++ {
		this.byte_stream.Append(0)
	}

	this.dma_command = dma_command
}

func (this *MemoryCommand) InitWrite(
	memory_operation MemoryOperation,
	bank_address int64,
	size int64,
	byte_stream *encoding.ByteStream,
	dma_command *DmaCommand,
) {
	if memory_operation != WRITE {
		err := errors.New("memory operation != WRITE")
		panic(err)
	} else if dma_command.MemoryOperation() != WRITE {
		err := errors.New("DMA operation != WRITE")
		panic(err)
	} else if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	this.memory_operation = memory_operation
	this.bank_address = bank_address
	this.size = size
	this.byte_stream = byte_stream
	this.dma_command = dma_command
}

func (this *MemoryCommand) MemoryOperation() MemoryOperation {
	return this.memory_operation
}

func (this *MemoryCommand) BankAddress() int64 {
	return this.bank_address
}

func (this *MemoryCommand) Size() int64 {
	return this.size
}

func (this *MemoryCommand) ByteStream() *encoding.ByteStream {
	if this.byte_stream == nil {
		err := errors.New("byte stream == nil")
		panic(err)
	}

	return this.byte_stream
}

func (this *MemoryCommand) SetByteStream(byte_stream *encoding.ByteStream) {
	if this.memory_operation != READ {
		err := errors.New("memory operation != READ")
		panic(err)
	}

	for i := int64(0); i < byte_stream.Size(); i++ {
		this.byte_stream.Set(int(i), byte_stream.Get(int(i)))
	}
}

func (this *MemoryCommand) DmaCommand() *DmaCommand {
	if this.dma_command == nil {
		err := errors.New("DMA command == nil")
		panic(err)
	}

	return this.dma_command
}
