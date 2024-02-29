package dram

import (
	"errors"
	"fmt"
	"uPIMulator/src/misc"
)

type MemoryScheduler struct {
	channel_id int
	rank_id    int
	dpu_id     int

	input_q        *DmaCommandQ
	reorder_buffer *MemoryCommandQ
	ready_q        *MemoryCommandQ

	row_address            *int64
	wordline_size          int64
	min_access_granularity int64

	stat_factory *misc.StatFactory
}

func (this *MemoryScheduler) Init(
	channel_id int,
	rank_id int,
	dpu_id int,
	command_line_parser *misc.CommandLineParser,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	} else if dpu_id < 0 {
		err := errors.New("DPU ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.dpu_id = dpu_id

	this.input_q = new(DmaCommandQ)
	this.input_q.Init(-1, 0)

	this.reorder_buffer = new(MemoryCommandQ)
	this.reorder_buffer.Init(-1, 0)

	this.ready_q = new(MemoryCommandQ)
	this.ready_q.Init(-1, 0)

	this.row_address = nil
	this.wordline_size = command_line_parser.IntParameter("wordline_size")
	this.min_access_granularity = command_line_parser.IntParameter("min_access_granularity")

	name := fmt.Sprintf("MemoryScheduler[%d_%d_%d]", channel_id, rank_id, dpu_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *MemoryScheduler) Fini() {
	this.input_q.Fini()
	this.reorder_buffer.Fini()
	this.ready_q.Fini()
}

func (this *MemoryScheduler) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *MemoryScheduler) IsEmpty() bool {
	return this.input_q.IsEmpty() && this.reorder_buffer.IsEmpty() && this.ready_q.IsEmpty()
}

func (this *MemoryScheduler) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *MemoryScheduler) Push(dma_command *DmaCommand) {
	if !this.CanPush() {
		err := errors.New("memory scheduler cannot be pushed")
		panic(err)
	}

	this.input_q.Push(dma_command)
}

func (this *MemoryScheduler) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *MemoryScheduler) Pop() *MemoryCommand {
	if !this.CanPop() {
		err := errors.New("memory scheduler cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *MemoryScheduler) Flush() {
	if !this.IsEmpty() {
		err := errors.New("memory scheduler cannot be flushed")
		panic(err)
	}

	this.row_address = nil
}

func (this *MemoryScheduler) Cycle() {
	this.ServiceInputQ()

	if !this.ReorderFr() {
		this.ReorderFcFs()
	}

	this.input_q.Cycle()
	this.reorder_buffer.Cycle()
	this.ready_q.Cycle()
}

func (this *MemoryScheduler) ServiceInputQ() {
	if this.input_q.CanPop(1) {
		dma_command := this.input_q.Pop()

		this.PopulateMemoryCommands(dma_command)
	}
}

func (this *MemoryScheduler) PopulateMemoryCommands(dma_command *DmaCommand) {
	begin_address := dma_command.MramAddress()
	end_address := dma_command.MramAddress() + dma_command.Size()

	for address := begin_address; address < end_address; {
		wordline_address := this.WordlineAddress(address)

		size := this.Min(
			this.Min(
				address+this.min_access_granularity,
				wordline_address+this.wordline_size,
			),
			end_address,
		) - address

		memory_operation := dma_command.MemoryOperation()
		memory_command := new(MemoryCommand)
		if memory_operation == READ {
			memory_command.InitRead(READ, address, size, dma_command)
		} else if memory_operation == WRITE {
			byte_stream := dma_command.ByteStream(address, size)
			memory_command.InitWrite(WRITE, address, size, byte_stream, dma_command)
		} else {
			err := errors.New("memory operation is not valid")
			panic(err)
		}

		this.reorder_buffer.Push(memory_command)

		address += size
	}
}

func (this *MemoryScheduler) ReorderFr() bool {
	if this.row_address != nil && this.ready_q.CanPush(1) {
		for i := 0; this.reorder_buffer.CanPop(i + 1); i++ {
			memory_command, _ := this.reorder_buffer.Front(i)

			if this.WordlineAddress(memory_command.Address()) == *this.row_address &&
				this.ready_q.CanPush(1) {
				if i != 0 {
					this.stat_factory.Increment("num_fr", 1)
				} else {
					this.stat_factory.Increment("num_fcfs", 1)
				}

				this.reorder_buffer.Remove(i)
				this.ready_q.Push(memory_command)

				return true
			}
		}
	}

	return false
}

func (this *MemoryScheduler) ReorderFcFs() bool {
	if this.reorder_buffer.CanPop(1) && this.row_address != nil && this.ready_q.CanPush(3) {
		memory_command := this.reorder_buffer.Pop()
		wordline_address := this.WordlineAddress(memory_command.Address())

		if wordline_address == *this.row_address {
			err := errors.New("FR has not worked correctly")
			panic(err)
		}

		precharge := new(MemoryCommand)
		precharge.InitActivation(PRECHARGE, *this.row_address)

		activation := new(MemoryCommand)
		activation.InitActivation(ACTIVATION, wordline_address)

		this.ready_q.Push(precharge)
		this.ready_q.Push(activation)
		this.ready_q.Push(memory_command)

		*this.row_address = wordline_address

		return true
	} else if this.reorder_buffer.CanPop(1) && this.row_address == nil && this.ready_q.CanPush(2) {
		memory_command := this.reorder_buffer.Pop()
		wordline_address := this.WordlineAddress(memory_command.Address())

		activation := new(MemoryCommand)
		activation.InitActivation(ACTIVATION, wordline_address)

		this.ready_q.Push(activation)
		this.ready_q.Push(memory_command)

		this.row_address = new(int64)
		*this.row_address = wordline_address

		return true
	} else {
		return false
	}
}

func (this *MemoryScheduler) WordlineAddress(address int64) int64 {
	return address / this.wordline_size * this.wordline_size
}

func (this *MemoryScheduler) Min(x int64, y int64) int64 {
	if x <= y {
		return x
	} else {
		return y
	}
}
