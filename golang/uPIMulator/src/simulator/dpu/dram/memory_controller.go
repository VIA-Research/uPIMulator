package dram

import (
	"errors"
	"fmt"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/misc"
)

type MemoryController struct {
	channel_id int
	rank_id    int
	dpu_id     int

	wordline_size int64

	memory_scheduler *MemoryScheduler
	row_buffer       *RowBuffer
	mram             *Mram

	input_q          *DmaCommandQ
	wait_q           *DmaCommandQ
	memory_command_q *MemoryCommandQ
	ready_q          *DmaCommandQ

	stat_factory *misc.StatFactory
}

func (this *MemoryController) Init(
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

	this.wordline_size = command_line_parser.IntParameter("wordline_size")

	this.memory_scheduler = new(MemoryScheduler)
	this.memory_scheduler.Init(channel_id, rank_id, dpu_id, command_line_parser)

	this.row_buffer = new(RowBuffer)
	this.row_buffer.Init(channel_id, rank_id, dpu_id, command_line_parser)

	this.mram = nil

	this.input_q = new(DmaCommandQ)
	this.input_q.Init(-1, 0)

	this.wait_q = new(DmaCommandQ)
	this.wait_q.Init(-1, 0)

	this.memory_command_q = new(MemoryCommandQ)
	this.memory_command_q.Init(-1, 0)

	this.ready_q = new(DmaCommandQ)
	this.ready_q.Init(-1, 0)

	name := fmt.Sprintf("MemoryController[%d_%d_%d]", channel_id, rank_id, dpu_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *MemoryController) Fini() {
	this.memory_scheduler.Fini()
	this.row_buffer.Fini()

	this.input_q.Fini()
	this.wait_q.Fini()
	this.memory_command_q.Fini()
	this.ready_q.Fini()
}

func (this *MemoryController) ConnectMram(mram *Mram) {
	if this.mram != nil {
		err := errors.New("MRAM is already set")
		panic(err)
	}

	this.mram = mram
	this.row_buffer.ConnectMram(mram)
}

func (this *MemoryController) MemoryScheduler() *MemoryScheduler {
	return this.memory_scheduler
}

func (this *MemoryController) RowBuffer() *RowBuffer {
	return this.row_buffer
}

func (this *MemoryController) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *MemoryController) IsEmpty() bool {
	return this.memory_scheduler.IsEmpty() &&
		this.row_buffer.IsEmpty() &&
		this.input_q.IsEmpty() &&
		this.wait_q.IsEmpty() &&
		this.memory_command_q.IsEmpty() &&
		this.ready_q.IsEmpty()
}

func (this *MemoryController) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *MemoryController) Push(dma_command *DmaCommand) {
	if !this.CanPush() {
		err := errors.New("memory controller cannot be pushed")
		panic(err)
	}

	this.input_q.Push(dma_command)
}

func (this *MemoryController) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *MemoryController) Pop() *DmaCommand {
	if !this.CanPop() {
		err := errors.New("memory controller cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *MemoryController) Front() *DmaCommand {
	if !this.CanPop() {
		err := errors.New("memory controller cannot be popped")
		panic(err)
	}

	dma_command, _ := this.ready_q.Front(0)
	return dma_command
}

func (this *MemoryController) Read(address int64, size int64) *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for cur_address := address; cur_address < address+size; {
		cur_wordline_address := this.WordlineAddress(cur_address)
		cur_size := this.Min(cur_wordline_address+this.wordline_size, address+size) - cur_address
		cur_offset := cur_address % this.wordline_size

		mram_byte_stream := this.mram.Read(cur_wordline_address)

		for i := cur_offset; i < cur_offset+cur_size; i++ {
			byte_stream.Append(mram_byte_stream.Get(int(i)))
		}

		cur_address += cur_size
	}

	return byte_stream
}

func (this *MemoryController) Write(address int64, size int64, byte_stream *encoding.ByteStream) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	cur_byte_stream_offset := int64(0)
	for cur_address := address; cur_address < address+size; {
		cur_wordline_address := this.WordlineAddress(cur_address)
		cur_size := this.Min(cur_wordline_address+this.wordline_size, address+size) - cur_address
		cur_offset := cur_address % this.wordline_size

		mram_byte_stream := this.mram.Read(cur_wordline_address)

		for i := int64(0); i < cur_size; i++ {
			mram_byte_stream.Set(int(i+cur_offset), byte_stream.Get(int(i+cur_byte_stream_offset)))
		}

		this.mram.Write(cur_wordline_address, mram_byte_stream)

		cur_address += cur_size
		cur_byte_stream_offset += cur_size
	}
}

func (this *MemoryController) Flush() {
	this.memory_scheduler.Flush()
	this.row_buffer.Flush()
}

func (this *MemoryController) Cycle() {
	this.ServiceInputQ()
	this.ServiceScheduler()
	this.ServiceMemoryCommandQ()
	this.ServiceRowBuffer()
	this.ServiceWaitQ()

	this.memory_scheduler.Cycle()
	this.row_buffer.Cycle()

	this.input_q.Cycle()
	this.wait_q.Cycle()
	this.memory_command_q.Cycle()
	this.ready_q.Cycle()

	this.stat_factory.Increment("memory_cycle", 1)
}

func (this *MemoryController) ServiceInputQ() {
	if this.input_q.CanPop(1) && this.wait_q.CanPush(1) && this.memory_scheduler.CanPush() {
		dma_command := this.input_q.Pop()
		this.memory_scheduler.Push(dma_command)
		this.wait_q.Push(dma_command)
	}
}

func (this *MemoryController) ServiceScheduler() {
	if this.memory_scheduler.CanPop() && this.memory_command_q.CanPush(1) {
		memory_command := this.memory_scheduler.Pop()
		this.memory_command_q.Push(memory_command)
	}
}

func (this *MemoryController) ServiceMemoryCommandQ() {
	if this.memory_command_q.CanPop(1) && this.row_buffer.CanPush() {
		memory_command := this.memory_command_q.Pop()
		this.row_buffer.Push(memory_command)
	}
}

func (this *MemoryController) ServiceRowBuffer() {
	if this.row_buffer.CanPop() {
		memory_command := this.row_buffer.Pop()

		memory_operation := memory_command.MemoryOperation()
		if memory_operation == ACTIVATION {
			return
		} else if memory_operation == PRECHARGE {
			return
		} else if memory_operation == READ {
			address := memory_command.Address()
			size := memory_command.Size()
			byte_stream := memory_command.ByteStream()

			dma_command := memory_command.DmaCommand()
			dma_command.SetByteStream(address, size, byte_stream)
			dma_command.SetAck(address, size)
		} else if memory_operation == WRITE {
			address := memory_command.Address()
			size := memory_command.Size()

			dma_command := memory_command.DmaCommand()
			dma_command.SetAck(address, size)
		} else {
			err := errors.New("memory operation is not valid")
			panic(err)
		}
	}
}

func (this *MemoryController) ServiceWaitQ() {
	for i := 0; this.wait_q.CanPop(i + 1); i++ {
		dma_command, _ := this.wait_q.Front(i)

		if dma_command.IsReady() && this.ready_q.CanPush(1) {
			this.wait_q.Remove(i)
			this.ready_q.Push(dma_command)
		}
	}
}

func (this *MemoryController) WordlineAddress(address int64) int64 {
	return address / this.wordline_size * this.wordline_size
}

func (this *MemoryController) Min(x int64, y int64) int64 {
	if x <= y {
		return x
	} else {
		return y
	}
}
