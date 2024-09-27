package dram

import (
	"errors"
	"uPIMulator/src/host/vm/dram/bank"
	"uPIMulator/src/host/vm/dram/channel"
	"uPIMulator/src/misc"
)

type MemoryScheduler struct {
	num_vm_channels          int
	num_vm_ranks_per_channel int
	num_vm_banks_per_rank    int

	input_q        *bank.DmaCommandQ
	reorder_buffer *channel.ChannelCommandQ
	ready_q        *channel.ChannelCommandQ

	row_addresses          map[int]int64
	wordline_size          int64
	min_access_granularity int64
	reorder_window_size    int

	stat_factory *misc.StatFactory
}

func (this *MemoryScheduler) Init(command_line_parser *misc.CommandLineParser) {
	this.num_vm_channels = int(command_line_parser.IntParameter("num_vm_channels"))
	this.num_vm_ranks_per_channel = int(
		command_line_parser.IntParameter("num_vm_ranks_per_channel"),
	)
	this.num_vm_banks_per_rank = int(command_line_parser.IntParameter("num_vm_banks_per_rank"))

	this.input_q = new(bank.DmaCommandQ)
	this.input_q.Init(-1, 0)

	this.reorder_buffer = new(channel.ChannelCommandQ)
	this.reorder_buffer.Init(-1, 0)

	this.ready_q = new(channel.ChannelCommandQ)
	this.ready_q.Init(-1, 0)

	this.row_addresses = make(map[int]int64)
	this.wordline_size = command_line_parser.IntParameter("wordline_size")
	this.min_access_granularity = command_line_parser.IntParameter("min_access_granularity")
	this.reorder_window_size = int(command_line_parser.IntParameter("reorder_window_size"))

	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init("MemoryScheduler")
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

func (this *MemoryScheduler) Push(dma_command *bank.DmaCommand) {
	if !this.CanPush() {
		err := errors.New("memory scheduler cannot be pushed")
		panic(err)
	}

	this.input_q.Push(dma_command)
}

func (this *MemoryScheduler) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *MemoryScheduler) Pop() *channel.ChannelCommand {
	if !this.CanPop() {
		err := errors.New("memory scheduler cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *MemoryScheduler) Generate(dma_command *bank.DmaCommand) []*channel.ChannelCommand {
	begin_address := dma_command.Segment().BankAddress()
	end_address := dma_command.Segment().BankAddress() + dma_command.Segment().Size()

	channel_commands := make([]*channel.ChannelCommand, 0)
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
		memory_command := new(bank.MemoryCommand)
		if memory_operation == bank.READ {
			memory_command.InitRead(bank.READ, address, size, dma_command)
		} else if memory_operation == bank.WRITE {
			byte_stream := dma_command.ByteStream(address, size)
			memory_command.InitWrite(bank.WRITE, address, size, byte_stream, dma_command)
		} else {
			err := errors.New("memory operation is not valid")
			panic(err)
		}

		channel_command := new(channel.ChannelCommand)
		channel_command.Init(
			dma_command.Segment().ChannelID(),
			dma_command.Segment().RankID(),
			dma_command.Segment().BankID(),
			memory_command,
		)

		channel_commands = append(channel_commands, channel_command)

		address += size
	}
	return channel_commands
}

func (this *MemoryScheduler) Flush() {
	if !this.IsEmpty() {
		err := errors.New("memory scheduler cannot be flushed")
		panic(err)
	}

	this.row_addresses = make(map[int]int64)
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

func (this *MemoryScheduler) PopulateMemoryCommands(dma_command *bank.DmaCommand) {
	begin_address := dma_command.Segment().BankAddress()
	end_address := dma_command.Segment().BankAddress() + dma_command.Segment().Size()

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
		memory_command := new(bank.MemoryCommand)
		if memory_operation == bank.READ {
			memory_command.InitRead(bank.READ, address, size, dma_command)
		} else if memory_operation == bank.WRITE {
			byte_stream := dma_command.ByteStream(address, size)
			memory_command.InitWrite(bank.WRITE, address, size, byte_stream, dma_command)
		} else {
			err := errors.New("memory operation is not valid")
			panic(err)
		}

		channel_command := new(channel.ChannelCommand)
		channel_command.Init(
			dma_command.Segment().ChannelID(),
			dma_command.Segment().RankID(),
			dma_command.Segment().BankID(),
			memory_command,
		)

		this.reorder_buffer.Push(channel_command)

		address += size
	}
}

func (this *MemoryScheduler) ReorderFr() bool {
	if this.ready_q.CanPush(1) {
		for i := 0; this.reorder_buffer.CanPop(i+1) && i < this.reorder_window_size; i++ {
			channel_command, _ := this.reorder_buffer.Front(i)
			memory_command := channel_command.MemoryCommand()

			wordline_address := this.WordlineAddress(memory_command.BankAddress())

			if this.IsOpened(channel_command) &&
				this.RowAddress(channel_command) == wordline_address {
				if i != 0 {
					this.stat_factory.Increment("num_fr", 1)
				} else {
					this.stat_factory.Increment("num_fcfs", 1)
				}

				this.reorder_buffer.Remove(i)
				this.ready_q.Push(channel_command)

				return true
			}
		}
	}

	return false
}

func (this *MemoryScheduler) ReorderFcFs() bool {
	if this.reorder_buffer.CanPop(1) && this.ready_q.CanPush(3) {
		channel_command, _ := this.reorder_buffer.Front(0)

		if this.IsOpened(channel_command) {
			this.reorder_buffer.Remove(0)

			memory_command := channel_command.MemoryCommand()

			wordline_address := this.WordlineAddress(memory_command.BankAddress())

			precharge := new(bank.MemoryCommand)
			precharge.InitActivation(bank.PRECHARGE, this.RowAddress(channel_command))

			activation := new(bank.MemoryCommand)
			activation.InitActivation(bank.ACTIVATION, wordline_address)

			channel_id := channel_command.ChannelId()
			rank_id := channel_command.RankId()
			bank_id := channel_command.BankId()

			precharge_channel_command := new(channel.ChannelCommand)
			precharge_channel_command.Init(channel_id, rank_id, bank_id, precharge)

			activation_channel_command := new(channel.ChannelCommand)
			activation_channel_command.Init(channel_id, rank_id, bank_id, activation)

			this.ready_q.Push(precharge_channel_command)
			this.ready_q.Push(activation_channel_command)
			this.ready_q.Push(channel_command)

			this.row_addresses[this.UniqueBankId(channel_id, rank_id, bank_id)] = wordline_address

			return true
		} else {
			this.reorder_buffer.Remove(0)

			memory_command := channel_command.MemoryCommand()

			wordline_address := this.WordlineAddress(memory_command.BankAddress())

			activation := new(bank.MemoryCommand)
			activation.InitActivation(bank.ACTIVATION, wordline_address)

			channel_id := channel_command.ChannelId()
			rank_id := channel_command.RankId()
			bank_id := channel_command.BankId()

			activation_channel_command := new(channel.ChannelCommand)
			activation_channel_command.Init(channel_id, rank_id, bank_id, activation)

			this.ready_q.Push(activation_channel_command)
			this.ready_q.Push(channel_command)

			this.row_addresses[this.UniqueBankId(channel_id, rank_id, bank_id)] = wordline_address

			return true
		}
	} else {
		return false
	}
}

func (this *MemoryScheduler) IsOpened(channel_command *channel.ChannelCommand) bool {
	unique_bank_id := this.UniqueBankId(
		channel_command.ChannelId(),
		channel_command.RankId(),
		channel_command.BankId(),
	)

	if _, found := this.row_addresses[unique_bank_id]; found {
		return true
	} else {
		return false
	}
}

func (this *MemoryScheduler) RowAddress(channel_command *channel.ChannelCommand) int64 {
	unique_bank_id := this.UniqueBankId(
		channel_command.ChannelId(),
		channel_command.RankId(),
		channel_command.BankId(),
	)

	if row_address, found := this.row_addresses[unique_bank_id]; found {
		return row_address
	} else {
		err := errors.New("unique bank is not opened")
		panic(err)
	}
}

func (this *MemoryScheduler) UniqueBankId(channel_id int, rank_id int, bank_id int) int {
	return channel_id*this.num_vm_ranks_per_channel*this.num_vm_banks_per_rank + rank_id*this.num_vm_banks_per_rank + bank_id
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
