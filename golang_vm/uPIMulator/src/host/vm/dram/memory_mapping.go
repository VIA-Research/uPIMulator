package dram

import (
	"uPIMulator/src/host/vm/dram/bank"
	"uPIMulator/src/misc"
)

type MemoryMapping struct {
	num_vm_channels          int
	num_vm_ranks_per_channel int
	num_vm_banks_per_rank    int

	vm_bg0  int
	vm_bg1  int
	vm_bank int

	vm_bank_offset int64
	vm_bank_size   int64
}

func (this *MemoryMapping) Init(command_line_parser *misc.CommandLineParser) {
	this.num_vm_channels = int(command_line_parser.IntParameter("num_vm_channels"))
	this.num_vm_ranks_per_channel = int(
		command_line_parser.IntParameter("num_vm_ranks_per_channel"),
	)
	this.num_vm_banks_per_rank = int(command_line_parser.IntParameter("num_vm_banks_per_rank"))

	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.vm_bg0 = config_loader.VmBg0()
	this.vm_bg1 = config_loader.VmBg1()
	this.vm_bank = config_loader.VmBank()

	this.vm_bank_offset = config_loader.VmBankOffset()
	this.vm_bank_size = config_loader.VmBankSize()
}

func (this *MemoryMapping) Fini() {
}

func (this *MemoryMapping) Map(vm_address int64, vm_size int64) []*bank.Segment {
	begin_address := vm_address
	end_address := vm_address + vm_size

	segments := make([]*bank.Segment, 0)

	for address := begin_address; address < end_address; {
		offset_address := this.OffsetAddress(address)

		size := this.Min(
			this.Min(
				address+this.Offset(),
				offset_address+this.Offset(),
			),
			end_address,
		) - address

		channel_id := this.ChannelId(address)
		rank_id := this.RankId(address)
		bank_id := this.BankId(address)
		bank_address := this.BankAddress(address)

		segment := new(bank.Segment)
		segment.Init(address, channel_id, rank_id, bank_id, bank_address, size)

		segments = append(segments, segment)

		address += size
	}

	return segments
}

func (this *MemoryMapping) ChannelId(address int64) int {
	return int(address / this.ChannelSize())
}

func (this *MemoryMapping) RankId(address int64) int {
	return int((address % this.ChannelSize()) / this.RankSize())
}

func (this *MemoryMapping) BankId(address int64) int {
	return this.Bank(address)*4 + this.Bg1(address)*2 + this.Bg0(address)
}

func (this *MemoryMapping) BankAddress(address int64) int64 {
	return address % this.vm_bank_size
}

func (this *MemoryMapping) ChannelSize() int64 {
	return int64(this.num_vm_ranks_per_channel) * this.RankSize()
}

func (this *MemoryMapping) RankSize() int64 {
	return int64(this.num_vm_banks_per_rank) * this.vm_bank_size
}

func (this *MemoryMapping) Bg0(address int64) int {
	return int(address & (1 << this.vm_bg0) >> this.vm_bg0)
}

func (this *MemoryMapping) Bg1(address int64) int {
	return int(address & (1 << this.vm_bg1) >> this.vm_bg1)
}

func (this *MemoryMapping) Bank(address int64) int {
	return int(address & (3 << this.vm_bank) >> this.vm_bank)
}

func (this *MemoryMapping) Offset() int64 {
	return this.Pow2(this.vm_bg0)
}

func (this *MemoryMapping) OffsetAddress(address int64) int64 {
	return (address / this.Offset()) * this.Offset()
}

func (this *MemoryMapping) Pow2(exponent int) int64 {
	value := int64(1)
	for i := 0; i < exponent; i++ {
		value *= 2
	}
	return value
}

func (this *MemoryMapping) Min(x int64, y int64) int64 {
	if x <= y {
		return x
	} else {
		return y
	}
}
