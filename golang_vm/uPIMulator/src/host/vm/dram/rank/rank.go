package rank

import (
	"errors"
	"uPIMulator/src/host/vm/dram/bank"
	"uPIMulator/src/misc"
)

type Rank struct {
	channel_id int
	rank_id    int

	banks []*bank.Bank

	input_q    *RankCommandQ
	ready_q    *RankCommandQ
	scoreboard map[*bank.MemoryCommand]*RankCommand
}

func (this *Rank) Init(channel_id int, rank_id int, command_line_parser *misc.CommandLineParser) {
	if channel_id < 0 {
		err := errors.New("channel ID < 0")
		panic(err)
	} else if rank_id < 0 {
		err := errors.New("rank ID < 0")
		panic(err)
	}

	this.channel_id = channel_id
	this.rank_id = rank_id

	this.banks = make([]*bank.Bank, 0)
	num_vm_banks_per_rank := int(command_line_parser.IntParameter("num_vm_banks_per_rank"))
	for i := 0; i < num_vm_banks_per_rank; i++ {
		bank_ := new(bank.Bank)
		bank_.Init(channel_id, rank_id, i, command_line_parser)

		this.banks = append(this.banks, bank_)
	}

	this.input_q = new(RankCommandQ)
	this.input_q.Init(-1, 0)

	this.ready_q = new(RankCommandQ)
	this.ready_q.Init(-1, 0)

	this.scoreboard = make(map[*bank.MemoryCommand]*RankCommand)
}

func (this *Rank) Fini() {
	for _, bank_ := range this.banks {
		bank_.Fini()
	}

	this.input_q.Fini()
	this.ready_q.Fini()
}

func (this *Rank) ChannelID() int {
	return this.channel_id
}

func (this *Rank) RankId() int {
	return this.rank_id
}

func (this *Rank) NumBanks() int {
	return len(this.banks)
}

func (this *Rank) Banks() []*bank.Bank {
	return this.banks
}

func (this *Rank) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *Rank) Push(rank_command *RankCommand) {
	if !this.CanPush() {
		err := errors.New("rank cannot be pushed")
		panic(err)
	}

	this.input_q.Push(rank_command)
}

func (this *Rank) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *Rank) Pop() *RankCommand {
	if !this.CanPop() {
		err := errors.New("rank cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *Rank) Read(rank_command *RankCommand) {
	bank_id := rank_command.BankId()

	memory_command := rank_command.MemoryCommand()

	bank_address := memory_command.BankAddress()
	size := memory_command.Size()

	byte_stream := this.banks[bank_id].Read(bank_address, size)
	memory_command.SetByteStream(byte_stream)
}

func (this *Rank) Write(rank_command *RankCommand) {
	bank_id := rank_command.BankId()

	memory_command := rank_command.MemoryCommand()

	bank_address := memory_command.BankAddress()
	size := memory_command.Size()
	byte_stream := memory_command.ByteStream()

	this.banks[bank_id].Write(bank_address, size, byte_stream)
}

func (this *Rank) Flush() {
	for _, bank_ := range this.banks {
		bank_.Flush()
	}
}

func (this *Rank) Cycle() {
	this.ServiceInputQ()
	this.ServiceReadyQ()

	this.input_q.Cycle()
	this.ready_q.Cycle()
}

func (this *Rank) ServiceInputQ() {
	if this.input_q.CanPop(1) {
		rank_command, _ := this.input_q.Front(0)

		bank_id := rank_command.BankId()
		if this.banks[bank_id].CanPush() {
			this.input_q.Pop()

			memory_command := rank_command.MemoryCommand()

			this.banks[bank_id].Push(memory_command)
			this.scoreboard[memory_command] = rank_command
		}
	}
}

func (this *Rank) ServiceReadyQ() {
	for _, bank_ := range this.banks {
		if bank_.CanPop() && this.ready_q.CanPush(1) {
			memory_command := bank_.Pop()
			rank_command := this.scoreboard[memory_command]

			this.ready_q.Push(rank_command)
			delete(this.scoreboard, memory_command)
		}
	}
}
