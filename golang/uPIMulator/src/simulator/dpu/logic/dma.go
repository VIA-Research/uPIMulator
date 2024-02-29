package logic

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/linker/kernel/instruction"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/dpu/dram"
	"uPIMulator/src/simulator/dpu/sram"
)

type Dma struct {
	atomic            *sram.Atomic
	iram              *sram.Iram
	operand_collector *OperandCollector
	memory_controller *dram.MemoryController

	input_q *dram.DmaCommandQ
	ready_q *dram.DmaCommandQ
}

func (this *Dma) Init() {
	this.atomic = nil
	this.iram = nil
	this.operand_collector = nil
	this.memory_controller = nil

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	max_num_tasklets := config_loader.MaxNumTasklets()

	this.input_q = new(dram.DmaCommandQ)
	this.input_q.Init(max_num_tasklets, 0)

	this.ready_q = new(dram.DmaCommandQ)
	this.ready_q.Init(max_num_tasklets, 0)
}

func (this *Dma) Fini() {
	this.input_q.Fini()
	this.ready_q.Fini()
}

func (this *Dma) ConnectAtomic(atomic *sram.Atomic) {
	if this.atomic != nil {
		err := errors.New("atomic is already set")
		panic(err)
	}

	this.atomic = atomic
}

func (this *Dma) ConnectIram(iram *sram.Iram) {
	if this.iram != nil {
		err := errors.New("IRAM is already set")
		panic(err)
	}

	this.iram = iram
}

func (this *Dma) ConnectOperandCollector(operand_collector *OperandCollector) {
	if this.operand_collector != nil {
		err := errors.New("operand collector is already set")
		panic(err)
	}

	this.operand_collector = operand_collector
}

func (this *Dma) ConnectMemoryController(memory_controller *dram.MemoryController) {
	if this.memory_controller != nil {
		err := errors.New("memory controller is already set")
		panic(err)
	}

	this.memory_controller = memory_controller
}

func (this *Dma) IsEmpty() bool {
	return this.input_q.IsEmpty() && this.ready_q.IsEmpty()
}

func (this *Dma) TransferToAtomic(address int64, byte_stream *encoding.ByteStream) {
	for i := int64(0); i < byte_stream.Size(); i++ {
		if byte_stream.Get(int(i)) != 0 {
			err := errors.New("atomic byte is not set to 0")
			panic(err)
		}
	}
}

func (this *Dma) TransferToIram(address int64, byte_stream *encoding.ByteStream) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	iram_offset := config_loader.IramOffset()
	iram_data_size := int64(config_loader.IramDataWidth() / 8)

	if address != this.iram.Address() {
		err := errors.New("address != IRAM's address")
		panic(err)
	} else if byte_stream.Size()%iram_data_size != 0 {
		err := errors.New("byte stream's size is not aligned with IRAM data size")
		panic(err)
	}

	this.iram.Write(iram_offset, byte_stream)
}

func (this *Dma) TransferFromWram(address int64, size int64) *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < size; i++ {
		value := this.operand_collector.Lbu(address + i)
		byte_stream.Append(uint8(value))
	}

	return byte_stream
}

func (this *Dma) TransferToWram(address int64, byte_stream *encoding.ByteStream) {
	for i := int64(0); i < byte_stream.Size(); i++ {
		value := byte_stream.Get(int(i))
		this.operand_collector.Sb(address+i, int64(value))
	}
}

func (this *Dma) TransferFromMram(address int64, size int64) *encoding.ByteStream {
	this.memory_controller.Flush()
	return this.memory_controller.Read(address, size)
}

func (this *Dma) TransferToMram(address int64, byte_stream *encoding.ByteStream) {
	this.memory_controller.Write(address, byte_stream.Size(), byte_stream)
}

func (this *Dma) TransferFromWramToMram(
	wram_address int64,
	mram_address int64,
	size int64,
	instruction_ *instruction.Instruction,
) {
	if !this.CanPush() {
		err := errors.New("DMA cannot be pushed")
		panic(err)
	}

	byte_stream := this.TransferFromWram(wram_address, size)

	dma_command := new(dram.DmaCommand)
	dma_command.InitWriteToMramFromWram(wram_address, mram_address, size, byte_stream, instruction_)

	this.Push(dma_command)
}

func (this *Dma) TransferFromMramToWram(
	wram_address int64,
	mram_address int64,
	size int64,
	instruction_ *instruction.Instruction,
) {
	if !this.CanPush() {
		err := errors.New("DMA cannot be pushed")
		panic(err)
	}

	dma_command := new(dram.DmaCommand)
	dma_command.InitReadFromMramToWram(wram_address, mram_address, size, instruction_)

	this.Push(dma_command)
}

func (this *Dma) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *Dma) Push(dma_command *dram.DmaCommand) {
	if !this.CanPush() {
		err := errors.New("DMA cannot be pushed")
		panic(err)
	}

	this.input_q.Push(dma_command)
}

func (this *Dma) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *Dma) Pop() *dram.DmaCommand {
	if !this.CanPop() {
		err := errors.New("DMA cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *Dma) Cycle() {
	this.ServiceInputQ()
	this.ServiceReadyQ()
}

func (this *Dma) ServiceInputQ() {
	if this.input_q.CanPop(1) && this.memory_controller.CanPush() {
		dma_command := this.input_q.Pop()
		this.memory_controller.Push(dma_command)
	}
}

func (this *Dma) ServiceReadyQ() {
	if this.memory_controller.CanPop() && this.ready_q.CanPush(1) {
		dma_command := this.memory_controller.Pop()
		this.ready_q.Push(dma_command)

		if dma_command.MemoryOperation() == dram.READ {
			wram_address := dma_command.WramAddress()
			mram_address := dma_command.MramAddress()
			size := dma_command.Size()
			byte_stream := dma_command.ByteStream(mram_address, size)

			this.TransferToWram(wram_address, byte_stream)
		}
	}
}
