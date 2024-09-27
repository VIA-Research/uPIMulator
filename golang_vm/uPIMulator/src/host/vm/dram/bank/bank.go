package bank

import (
	"errors"
	"fmt"
	"uPIMulator/src/encoding"
	"uPIMulator/src/misc"
)

type Bank struct {
	channel_id int
	rank_id    int
	bank_id    int

	wordline_size int64

	array      *Array
	row_buffer *RowBuffer

	input_q *MemoryCommandQ
	ready_q *MemoryCommandQ

	stat_factory *misc.StatFactory
}

func (this *Bank) Init(
	channel_id int,
	rank_id int,
	bank_id int,
	command_line_parser *misc.CommandLineParser,
) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	} else if bank_id < 0 {
		err := errors.New("bank ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id
	this.bank_id = bank_id

	this.wordline_size = command_line_parser.IntParameter("wordline_size")

	this.array = new(Array)
	this.array.Init(command_line_parser)

	this.row_buffer = new(RowBuffer)
	this.row_buffer.Init(channel_id, rank_id, bank_id, command_line_parser)
	this.row_buffer.ConnectArray(this.array)

	this.input_q = new(MemoryCommandQ)
	this.input_q.Init(-1, 0)

	this.ready_q = new(MemoryCommandQ)
	this.ready_q.Init(-1, 0)

	name := fmt.Sprintf("VmBank[%d_%d_%d]", channel_id, rank_id, bank_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *Bank) Fini() {
	this.array.Fini()
	this.row_buffer.Fini()

	this.input_q.Fini()
	this.ready_q.Fini()
}

func (this *Bank) ChannelId() int {
	return this.channel_id
}

func (this *Bank) RankId() int {
	return this.rank_id
}

func (this *Bank) DpuId() int {
	return this.bank_id
}

func (this *Bank) RowBuffer() *RowBuffer {
	return this.row_buffer
}

func (this *Bank) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *Bank) IsEmpty() bool {
	return this.row_buffer.IsEmpty() && this.input_q.IsEmpty() && this.ready_q.IsEmpty()
}

func (this *Bank) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *Bank) Push(memory_command *MemoryCommand) {
	if !this.CanPush() {
		err := errors.New("bank cannot be pushed")
		panic(err)
	}

	this.input_q.Push(memory_command)
}

func (this *Bank) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *Bank) Pop() *MemoryCommand {
	if !this.CanPop() {
		err := errors.New("bank cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *Bank) Read(address int64, size int64) *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for cur_address := address; cur_address < address+size; {
		cur_wordline_address := this.WordlineAddress(cur_address)
		cur_size := this.Min(cur_wordline_address+this.wordline_size, address+size) - cur_address
		cur_offset := cur_address % this.wordline_size

		mram_byte_stream := this.array.Read(cur_wordline_address)

		for i := cur_offset; i < cur_offset+cur_size; i++ {
			byte_stream.Append(mram_byte_stream.Get(int(i)))
		}

		cur_address += cur_size
	}

	return byte_stream
}

func (this *Bank) Write(address int64, size int64, byte_stream *encoding.ByteStream) {
	if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	cur_byte_stream_offset := int64(0)
	for cur_address := address; cur_address < address+size; {
		cur_wordline_address := this.WordlineAddress(cur_address)
		cur_size := this.Min(cur_wordline_address+this.wordline_size, address+size) - cur_address
		cur_offset := cur_address % this.wordline_size

		mram_byte_stream := this.array.Read(cur_wordline_address)

		for i := int64(0); i < cur_size; i++ {
			mram_byte_stream.Set(int(i+cur_offset), byte_stream.Get(int(i+cur_byte_stream_offset)))
		}

		this.array.Write(cur_wordline_address, mram_byte_stream)

		cur_address += cur_size
		cur_byte_stream_offset += cur_size
	}
}

func (this *Bank) Flush() {
	this.row_buffer.Flush()
}

func (this *Bank) Cycle() {
	this.ServiceInputQ()
	this.ServiceRowBuffer()

	this.row_buffer.Cycle()

	this.input_q.Cycle()
	this.ready_q.Cycle()

	this.stat_factory.Increment("vm_memory_cycle", 1)
}

func (this *Bank) ServiceInputQ() {
	if this.input_q.CanPop(1) && this.row_buffer.CanPush() {
		memory_command := this.input_q.Pop()
		this.row_buffer.Push(memory_command)
	}
}

func (this *Bank) ServiceRowBuffer() {
	if this.row_buffer.CanPop() && this.ready_q.CanPush(1) {
		memory_command := this.row_buffer.Pop()

		memory_operation := memory_command.MemoryOperation()
		if memory_operation == ACTIVATION {
			return
		} else if memory_operation == PRECHARGE {
			return
		} else if memory_operation == READ {
			address := memory_command.BankAddress()
			size := memory_command.Size()
			byte_stream := memory_command.ByteStream()

			dma_command := memory_command.DmaCommand()
			dma_command.SetByteStream(address, size, byte_stream)
			dma_command.SetAck(address, size)
		} else if memory_operation == WRITE {
			address := memory_command.BankAddress()
			size := memory_command.Size()

			dma_command := memory_command.DmaCommand()
			dma_command.SetAck(address, size)
		} else {
			err := errors.New("memory operation is not valid")
			panic(err)
		}

		this.ready_q.Push(memory_command)
	}
}

func (this *Bank) WordlineAddress(address int64) int64 {
	return address / this.wordline_size * this.wordline_size
}

func (this *Bank) Min(x int64, y int64) int64 {
	if x <= y {
		return x
	} else {
		return y
	}
}
