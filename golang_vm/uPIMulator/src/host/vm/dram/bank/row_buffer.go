package bank

import (
	"errors"
	"fmt"
	"uPIMulator/src/encoding"
	"uPIMulator/src/misc"
)

type RowBuffer struct {
	channel_id int
	rank_id    int
	bank_id    int

	t_ras         int64
	t_rcd         int64
	t_cl          int64
	t_bl          int64
	t_rp          int64
	wordline_size int64

	array       *Array
	row_address *int64
	row_buffer  *encoding.ByteStream

	input_q *MemoryCommandQ
	ready_q *MemoryCommandQ

	activation_q *MemoryCommandQ
	io_q         *MemoryCommandQ
	bus_q        *MemoryCommandQ
	precharge_q  *MemoryCommandQ

	stat_factory *misc.StatFactory
}

func (this *RowBuffer) Init(
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

	this.t_ras = command_line_parser.IntParameter("t_ras")
	this.t_rcd = command_line_parser.IntParameter("t_rcd")
	this.t_cl = command_line_parser.IntParameter("t_cl")
	this.t_bl = command_line_parser.IntParameter("t_bl")
	this.t_rp = command_line_parser.IntParameter("t_rp")
	this.wordline_size = command_line_parser.IntParameter("wordline_size")

	this.array = nil
	this.row_address = nil
	this.row_buffer = nil

	this.input_q = new(MemoryCommandQ)
	this.input_q.Init(1, 0)

	this.ready_q = new(MemoryCommandQ)
	this.ready_q.Init(-1, 0)

	this.activation_q = new(MemoryCommandQ)
	this.activation_q.Init(1, this.t_ras)

	this.io_q = new(MemoryCommandQ)
	this.io_q.Init(1, this.t_cl)

	this.bus_q = new(MemoryCommandQ)
	this.bus_q.Init(1, this.t_bl)

	this.precharge_q = new(MemoryCommandQ)
	this.precharge_q.Init(1, this.t_rp)

	name := fmt.Sprintf("VmRowBuffer[%d_%d_%d]", channel_id, rank_id, bank_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *RowBuffer) Fini() {
	this.input_q.Fini()
	this.ready_q.Fini()
	this.activation_q.Fini()
	this.io_q.Fini()
	this.bus_q.Fini()
	this.precharge_q.Fini()
}

func (this *RowBuffer) ConnectArray(array *Array) {
	if this.array != nil {
		err := errors.New("array is already connected")
		panic(err)
	}

	this.array = array
}

func (this *RowBuffer) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *RowBuffer) IsEmpty() bool {
	return this.input_q.IsEmpty() && this.ready_q.IsEmpty() && this.activation_q.IsEmpty() &&
		this.io_q.IsEmpty() &&
		this.bus_q.IsEmpty() &&
		this.precharge_q.IsEmpty()
}

func (this *RowBuffer) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *RowBuffer) Push(memory_command *MemoryCommand) {
	if !this.CanPush() {
		err := errors.New("row buffer cannot be pushed")
		panic(err)
	}

	this.input_q.Push(memory_command)
}

func (this *RowBuffer) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *RowBuffer) Pop() *MemoryCommand {
	if !this.CanPop() {
		err := errors.New("row buffer cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *RowBuffer) Flush() {
	if this.row_address != nil {
		this.WriteToBank()

		this.row_address = nil
		this.row_buffer = nil
	}
}

func (this *RowBuffer) Cycle() {
	this.ServiceInputQ()
	this.ServiceActivationQ()
	this.ServiceIoQ()
	this.ServiceBusQ()
	this.ServicePrechargeQ()

	this.input_q.Cycle()
	this.ready_q.Cycle()

	this.activation_q.Cycle()
	this.io_q.Cycle()
	this.bus_q.Cycle()
	this.precharge_q.Cycle()
}

func (this *RowBuffer) ServiceInputQ() {
	if this.input_q.CanPop(1) {
		memory_command, _ := this.input_q.Front(0)

		memory_operation := memory_command.MemoryOperation()
		if memory_operation == ACTIVATION {
			if this.activation_q.IsEmpty() && this.row_address == nil {
				this.activation_q.Push(memory_command)
				this.input_q.Pop()
			}
		} else if memory_operation == READ {
			if this.io_q.CanPush(1) && this.row_address != nil {
				this.io_q.Push(memory_command)
				this.input_q.Pop()
			}
		} else if memory_operation == WRITE {
			if this.io_q.CanPush(1) && this.row_address != nil {
				this.io_q.Push(memory_command)
				this.input_q.Pop()
			}
		} else if memory_operation == PRECHARGE {
			if this.activation_q.IsEmpty() && this.io_q.IsEmpty() && this.bus_q.IsEmpty() && this.precharge_q.IsEmpty() {
				this.precharge_q.Push(memory_command)
				this.input_q.Pop()
			}
		} else {
			err := errors.New("memory operation is not valid")
			panic(err)
		}
	}
}

func (this *RowBuffer) ServiceActivationQ() {
	if !this.activation_q.IsEmpty() {
		memory_command, cycle := this.activation_q.Front(0)

		if cycle == this.t_ras-this.t_rcd {
			if this.row_address != nil {
				err := errors.New("row buffer is not precharged")
				panic(err)
			} else if memory_command.BankAddress()%this.wordline_size != 0 {
				err := errors.New("memory command is not aligned with wordline size")
				panic(err)
			}

			this.row_address = new(int64)
			*this.row_address = memory_command.BankAddress()

			this.row_buffer = this.ReadFromBank()
		}
	}

	if this.activation_q.CanPop(1) && this.ready_q.CanPush(1) {
		memory_command := this.activation_q.Pop()
		this.ready_q.Push(memory_command)

		this.stat_factory.Increment("num_activations", 1)
	}
}

func (this *RowBuffer) ServiceIoQ() {
	if this.io_q.CanPop(1) && this.bus_q.CanPush(1) {
		memory_command := this.io_q.Pop()
		this.bus_q.Push(memory_command)
	}
}

func (this *RowBuffer) ServiceBusQ() {
	if this.bus_q.CanPop(1) && this.ready_q.CanPush(1) {
		memory_command := this.bus_q.Pop()
		this.ready_q.Push(memory_command)

		memory_operation := memory_command.MemoryOperation()
		if memory_operation == READ {
			byte_stream := this.ReadFromRowBuffer(
				memory_command.BankAddress(),
				memory_command.Size(),
			)
			memory_command.SetByteStream(byte_stream)

			this.stat_factory.Increment("num_reads", 1)
			this.stat_factory.Increment("read_bytes", memory_command.Size())
		} else if memory_operation == WRITE {
			this.WriteToRowBuffer(memory_command.BankAddress(), memory_command.Size(), memory_command.ByteStream())

			this.stat_factory.Increment("num_writes", 1)
			this.stat_factory.Increment("write_bytes", memory_command.Size())
		} else {
			err := errors.New("memory operation is not valid")
			panic(err)
		}
	}
}

func (this *RowBuffer) ServicePrechargeQ() {
	if this.precharge_q.CanPop(1) && this.ready_q.CanPush(1) {
		memory_command := this.precharge_q.Pop()

		address := memory_command.BankAddress()
		if address%this.wordline_size != 0 {
			err := errors.New("address is not aligned with wordline size")
			panic(err)
		} else if address != *this.row_address {
			err := errors.New("address != row address")
			panic(err)
		}

		this.WriteToBank()
		this.row_address = nil
		this.ready_q.Push(memory_command)

		this.stat_factory.Increment("num_precharges", 1)
	}
}

func (this *RowBuffer) ReadFromBank() *encoding.ByteStream {
	if this.row_address == nil {
		err := errors.New("row address is not set")
		panic(err)
	}

	return this.array.Read(*this.row_address)
}

func (this *RowBuffer) ReadFromRowBuffer(address int64, size int64) *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < size; i++ {
		index := this.Index(address) + int(i)

		byte_stream.Append(this.row_buffer.Get(index))
	}

	return byte_stream
}

func (this *RowBuffer) WriteToBank() {
	if this.row_address == nil {
		err := errors.New("row address is not set")
		panic(err)
	}

	this.array.Write(*this.row_address, this.row_buffer)
}

func (this *RowBuffer) WriteToRowBuffer(
	address int64,
	size int64,
	byte_stream *encoding.ByteStream,
) {
	if this.row_address == nil {
		err := errors.New("row address is not set")
		panic(err)
	} else if size != byte_stream.Size() {
		err := errors.New("size != byte stream's size")
		panic(err)
	}

	for i := int64(0); i < byte_stream.Size(); i++ {
		index := this.Index(address) + int(i)

		this.row_buffer.Set(index, byte_stream.Get(int(i)))
	}
}

func (this *RowBuffer) Index(address int64) int {
	if this.row_address == nil {
		err := errors.New("row address is not set")
		panic(err)
	} else if address < *this.row_address {
		err := errors.New("address < row address")
		panic(err)
	} else if address >= *this.row_address+this.wordline_size {
		err := errors.New("address >= row address + wordline size")
		panic(err)
	}

	return int(address - *this.row_address)
}
