package bank

import (
	"errors"
	"uPIMulator/src/device/simulator/dpu/dram"
	"uPIMulator/src/encoding"
)

type TransferCommandType int

const (
	HOST_TO_DEVICE TransferCommandType = iota
	DEVICE_TO_HOST
)

type TransferCommandState int

const (
	BEGIN TransferCommandState = iota
	MIDDLE
	END
)

type TransferCommand struct {
	transfer_command_type  TransferCommandType
	transfer_command_state TransferCommandState

	vm_address int64

	channel_id int
	rank_id    int
	dpu_id     int

	mram_address int64

	size int64

	byte_stream     *encoding.ByteStream
	vm_dma_commands map[*DmaCommand]bool
	dma_commands    map[*dram.DmaCommand]bool
}

func (this *TransferCommand) Init(
	transfer_command_type TransferCommandType,
	vm_address int64,
	channel_id int,
	rank_id int,
	dpu_id int,
	mram_address int64,
	size int64,
) {
	this.transfer_command_type = transfer_command_type
	this.transfer_command_state = BEGIN

	this.vm_address = vm_address

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.dpu_id = dpu_id
	this.mram_address = mram_address

	this.size = size

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()

	for i := int64(0); i < size; i++ {
		this.byte_stream.Append(0)
	}

	this.vm_dma_commands = make(map[*DmaCommand]bool)
	this.dma_commands = make(map[*dram.DmaCommand]bool)
}

func (this *TransferCommand) InitFast(vm_address int64, size int64) {
	this.vm_address = vm_address
	this.size = size

	this.byte_stream = new(encoding.ByteStream)
	this.byte_stream.Init()

	for i := int64(0); i < size; i++ {
		this.byte_stream.Append(0)
	}
}

func (this *TransferCommand) TransferCommandType() TransferCommandType {
	return this.transfer_command_type
}

func (this *TransferCommand) VmAddress() int64 {
	return this.vm_address
}

func (this *TransferCommand) ChannelId() int {
	return this.channel_id
}

func (this *TransferCommand) RankId() int {
	return this.rank_id
}

func (this *TransferCommand) DpuId() int {
	return this.dpu_id
}

func (this *TransferCommand) MramAddress() int64 {
	return this.mram_address
}

func (this *TransferCommand) Size() int64 {
	return this.size
}

func (this *TransferCommand) ByteStream() *encoding.ByteStream {
	return this.byte_stream
}

func (this *TransferCommand) SetByteStream(
	vm_address int64,
	size int64,
	byte_stream *encoding.ByteStream,
) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	for i := int64(0); i < byte_stream.Size(); i++ {
		index := this.Index(vm_address) + int(i)

		this.byte_stream.Set(index, byte_stream.Get(int(i)))
	}
}

func (this *TransferCommand) AppendVmDmaCommand(vm_dma_command *DmaCommand) {
	this.vm_dma_commands[vm_dma_command] = true
}

func (this *TransferCommand) AppendDmaCommand(dma_command *dram.DmaCommand) {
	this.dma_commands[dma_command] = true
}

func (this *TransferCommand) AckVmDmaCommand(vm_dma_command *DmaCommand) {
	if !vm_dma_command.IsReady() {
		err := errors.New("VM DMA command is not ready")
		panic(err)
	}

	if _, found := this.vm_dma_commands[vm_dma_command]; found {
		delete(this.vm_dma_commands, vm_dma_command)
	} else {
		err := errors.New("VM DMA command is not found")
		panic(err)
	}
}

func (this *TransferCommand) AckDmaCommand(dma_command *dram.DmaCommand) {
	if !dma_command.IsReady() {
		err := errors.New("DMA command is not ready")
		panic(err)
	}

	if _, found := this.dma_commands[dma_command]; found {
		delete(this.dma_commands, dma_command)
	} else {
		err := errors.New("DMA command is not found")
		panic(err)
	}
}

func (this *TransferCommand) IsVmReady() bool {
	return len(this.vm_dma_commands) == 0
}

func (this *TransferCommand) IsReady() bool {
	return len(this.dma_commands) == 0
}

func (this *TransferCommand) TransferCommandState() TransferCommandState {
	return this.transfer_command_state
}

func (this *TransferCommand) SetTransferCommandState(transfer_command_state TransferCommandState) {
	if transfer_command_state == BEGIN {
		err := errors.New("transfer command state is begin")
		panic(err)
	} else if transfer_command_state == MIDDLE {
		if this.transfer_command_state != BEGIN {
			err := errors.New("transfer command state is middle, but current transfer command state is not begin")
			panic(err)
		}

		this.transfer_command_state = transfer_command_state
	} else if transfer_command_state == END {
		if this.transfer_command_state != MIDDLE {
			err := errors.New("transfer command state is end, but current transfer command state is not middle")
			panic(err)
		}

		this.transfer_command_state = transfer_command_state
	} else {
		err := errors.New("transfer command state is not valid")
		panic(err)
	}
}

func (this *TransferCommand) Index(vm_address int64) int {
	if vm_address < this.vm_address {
		err := errors.New("VM address < transfer command's VM address")
		panic(err)
	} else if vm_address >= this.vm_address+this.size {
		err := errors.New("VM address >= transfer command's VM address + transfer command's size")
		panic(err)
	}

	return int(vm_address - this.vm_address)
}
