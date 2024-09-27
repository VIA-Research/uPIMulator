package dram

import (
	"errors"
	"uPIMulator/src/device/simulator/channel"
	"uPIMulator/src/device/simulator/dpu/dram"
	"uPIMulator/src/encoding"
	"uPIMulator/src/host/vm/dram/bank"
	vm_channel "uPIMulator/src/host/vm/dram/channel"
	"uPIMulator/src/misc"
)

type MemoryController struct {
	vm_channels []*vm_channel.Channel
	channels    []*channel.Channel

	memory_scheduler *MemoryScheduler
	memory_mapping   *MemoryMapping

	input_q *bank.TransferCommandQ

	vm_dma_command_q     *bank.DmaCommandQ
	vm_channel_command_q *vm_channel.ChannelCommandQ

	channel_command_q *channel.ChannelCommandQ

	vm_wait_q *bank.TransferCommandQ

	scoreboard map[*dram.DmaCommand]*bank.TransferCommand
	wait_q     *bank.TransferCommandQ

	ready_q *bank.TransferCommandQ
}

func (this *MemoryController) Init(command_line_parser *misc.CommandLineParser) {
	num_vm_channels := int(command_line_parser.IntParameter("num_vm_channels"))

	this.vm_channels = make([]*vm_channel.Channel, 0)
	for i := 0; i < num_vm_channels; i++ {
		vm_channel_ := new(vm_channel.Channel)
		vm_channel_.Init(i, command_line_parser)

		this.vm_channels = append(this.vm_channels, vm_channel_)
	}

	this.memory_scheduler = new(MemoryScheduler)
	this.memory_scheduler.Init(command_line_parser)

	this.memory_mapping = new(MemoryMapping)
	this.memory_mapping.Init(command_line_parser)

	this.input_q = new(bank.TransferCommandQ)
	this.input_q.Init(-1, 0)

	this.vm_dma_command_q = new(bank.DmaCommandQ)
	this.vm_dma_command_q.Init(-1, 0)

	this.vm_channel_command_q = new(vm_channel.ChannelCommandQ)
	this.vm_channel_command_q.Init(-1, 0)

	this.channel_command_q = new(channel.ChannelCommandQ)
	this.channel_command_q.Init(-1, 0)

	this.vm_wait_q = new(bank.TransferCommandQ)
	this.vm_wait_q.Init(-1, 0)

	this.scoreboard = make(map[*dram.DmaCommand]*bank.TransferCommand)

	this.wait_q = new(bank.TransferCommandQ)
	this.wait_q.Init(-1, 0)

	this.ready_q = new(bank.TransferCommandQ)
	this.ready_q.Init(-1, 0)
}

func (this *MemoryController) Fini() {
	for _, vm_channel_ := range this.vm_channels {
		vm_channel_.Fini()
	}

	this.memory_scheduler.Fini()
	this.memory_mapping.Fini()

	this.input_q.Fini()
	this.vm_dma_command_q.Fini()
	this.vm_channel_command_q.Fini()
	this.channel_command_q.Fini()
	this.vm_wait_q.Fini()

	if len(this.scoreboard) != 0 {
		err := errors.New("scoreboard is not empty")
		panic(err)
	}

	this.wait_q.Fini()

	this.ready_q.Fini()
}

func (this *MemoryController) ConnectChannels(channels []*channel.Channel) {
	this.channels = channels
}

func (this *MemoryController) VmChannels() []*vm_channel.Channel {
	return this.vm_channels
}

func (this *MemoryController) MemoryScheduler() *MemoryScheduler {
	return this.memory_scheduler
}

func (this *MemoryController) Banks() []*bank.Bank {
	banks := make([]*bank.Bank, 0)
	for _, vm_channel_ := range this.vm_channels {
		banks = append(banks, vm_channel_.Banks()...)
	}
	return banks
}

func (this *MemoryController) CanPush() bool {
	return this.input_q.CanPush(1)
}

func (this *MemoryController) Push(transfer_command *bank.TransferCommand) {
	if !this.CanPush() {
		err := errors.New("memory controller cannot be pushed")
		panic(err)
	}

	this.input_q.Push(transfer_command)
}

func (this *MemoryController) CanPop() bool {
	return this.ready_q.CanPop(1)
}

func (this *MemoryController) Pop() *bank.TransferCommand {
	if !this.CanPop() {
		err := errors.New("memory controller cannot be popped")
		panic(err)
	}

	return this.ready_q.Pop()
}

func (this *MemoryController) VmRead(vm_address int64, size int64) *encoding.ByteStream {
	segments := this.memory_mapping.Map(vm_address, size)

	transfer_command := new(bank.TransferCommand)
	transfer_command.InitFast(vm_address, size)

	for _, segment := range segments {
		dma_command := new(bank.DmaCommand)
		dma_command.InitRead(segment, transfer_command)

		channel_commands := this.memory_scheduler.Generate(dma_command)

		for _, channel_command := range channel_commands {
			channel_id := channel_command.ChannelId()

			this.vm_channels[channel_id].Read(channel_command)
		}

		byte_stream := dma_command.ByteStream(segment.BankAddress(), segment.Size())

		transfer_command.SetByteStream(
			dma_command.Segment().VmAddress(),
			dma_command.Segment().Size(),
			byte_stream,
		)
	}

	return transfer_command.ByteStream()
}

func (this *MemoryController) VmWrite(
	vm_address int64,
	size int64,
	byte_stream *encoding.ByteStream,
) {
	segments := this.memory_mapping.Map(vm_address, size)

	transfer_command := new(bank.TransferCommand)
	transfer_command.InitFast(vm_address, size)
	transfer_command.SetByteStream(vm_address, size, byte_stream)

	for _, segment := range segments {
		dma_command := new(bank.DmaCommand)
		dma_command.InitWrite(segment, transfer_command)

		dma_command_byte_stream := new(encoding.ByteStream)
		dma_command_byte_stream.Init()
		for i := int64(0); i < dma_command.Segment().Size(); i++ {
			byte_ := transfer_command.ByteStream().Get(int(segment.VmAddress() + i - vm_address))
			dma_command_byte_stream.Append(byte_)
		}

		dma_command.SetByteStream(segment.BankAddress(), segment.Size(), dma_command_byte_stream)

		channel_commands := this.memory_scheduler.Generate(dma_command)

		for _, channel_command := range channel_commands {
			channel_id := channel_command.ChannelId()

			this.vm_channels[channel_id].Write(channel_command)
		}
	}
}

func (this *MemoryController) Flush() {
	this.memory_scheduler.Flush()
	for _, vm_channel_ := range this.vm_channels {
		vm_channel_.Flush()
	}
}

func (this *MemoryController) Cycle() {
	this.ServiceInputQ()
	this.ServiceVmDmaCommandQ()
	this.ServiceMemoryScheduler()
	this.ServiceVmChannelCommandQ()
	this.ServiceVmChannels()
	this.ServiceChannelCommandQ()
	this.ServiceChannels()
	this.ServiceVmWaitQ()
	this.ServiceWaitQ()

	this.memory_scheduler.Cycle()

	this.input_q.Cycle()
	this.vm_dma_command_q.Cycle()
	this.vm_channel_command_q.Cycle()
	this.channel_command_q.Cycle()
	this.vm_wait_q.Cycle()
	this.wait_q.Cycle()
	this.ready_q.Cycle()

	for _, channel_ := range this.channels {
		channel_.Cycle()
	}
}

func (this *MemoryController) ServiceInputQ() {
	// TODO(bongjoon.hyun@gmail.com): need to check if VM DMA queue or channel command queue can be pushed
	if this.input_q.CanPop(1) {
		transfer_command := this.input_q.Pop()

		if transfer_command.TransferCommandType() == bank.HOST_TO_DEVICE {
			this.vm_wait_q.Push(transfer_command)

			segments := this.memory_mapping.Map(
				transfer_command.VmAddress(),
				transfer_command.Size(),
			)

			for _, segment := range segments {
				dma_command := new(bank.DmaCommand)
				dma_command.InitRead(segment, transfer_command)

				transfer_command.AppendVmDmaCommand(dma_command)

				this.vm_dma_command_q.Push(dma_command)
			}
		} else if transfer_command.TransferCommandType() == bank.DEVICE_TO_HOST {
			this.wait_q.Push(transfer_command)

			channel_id := transfer_command.ChannelId()
			rank_id := transfer_command.RankId()
			dpu_id := transfer_command.DpuId()
			mram_address := transfer_command.MramAddress()
			size := transfer_command.Size()

			dma_command := new(dram.DmaCommand)
			dma_command.InitReadFromMram(mram_address, size)

			this.scoreboard[dma_command] = transfer_command

			channel_command := new(channel.ChannelCommand)
			channel_command.Init(channel_id, rank_id, dpu_id, dma_command)

			transfer_command.AppendDmaCommand(dma_command)

			this.channel_command_q.Push(channel_command)
		} else {
			err := errors.New("transfer command type is not valid")
			panic(err)
		}
	}
}

func (this *MemoryController) ServiceVmDmaCommandQ() {
	if this.vm_dma_command_q.CanPop(1) && this.memory_scheduler.CanPush() {
		dma_command := this.vm_dma_command_q.Pop()

		this.memory_scheduler.Push(dma_command)
	}
}

func (this *MemoryController) ServiceMemoryScheduler() {
	if this.memory_scheduler.CanPop() && this.vm_channel_command_q.CanPush(1) {
		vm_channel_command := this.memory_scheduler.Pop()

		this.vm_channel_command_q.Push(vm_channel_command)
	}
}

func (this *MemoryController) ServiceVmChannelCommandQ() {
	if this.vm_channel_command_q.CanPop(1) {
		vm_channel_command, _ := this.vm_channel_command_q.Front(0)

		channel_id := vm_channel_command.ChannelId()

		if this.vm_channels[channel_id].CanPush() {
			this.vm_channel_command_q.Pop()

			this.vm_channels[channel_id].Push(vm_channel_command)
		}
	}
}

func (this *MemoryController) ServiceVmChannels() {
	for _, vm_channel_ := range this.vm_channels {
		if vm_channel_.CanPop() {
			vm_channel_command := vm_channel_.Pop()

			dma_command := vm_channel_command.MemoryCommand().DmaCommand()
			transfer_command := dma_command.TransferCommand()

			if transfer_command.TransferCommandType() == bank.HOST_TO_DEVICE {
				segment := dma_command.Segment()
				vm_address := segment.VmAddress()
				size := segment.Size()
				byte_stream := dma_command.ByteStream(segment.BankAddress(), segment.Size())

				transfer_command.SetByteStream(vm_address, size, byte_stream)

				if dma_command.IsReady() {
					transfer_command.AckVmDmaCommand(dma_command)
				}
			} else if transfer_command.TransferCommandType() == bank.DEVICE_TO_HOST {
				if dma_command.IsReady() {
					transfer_command.AckVmDmaCommand(dma_command)
				}
			} else {
				err := errors.New("transfer command type is not valid")
				panic(err)
			}
		}
	}
}

func (this *MemoryController) ServiceChannelCommandQ() {
	if this.channel_command_q.CanPop(1) {
		channel_command, _ := this.channel_command_q.Front(0)

		channel_id := channel_command.ChannelId()

		if this.channels[channel_id].CanPush() {
			this.channel_command_q.Pop()

			this.channels[channel_id].Push(channel_command)
		}
	}
}

func (this *MemoryController) ServiceChannels() {
	for _, channel_ := range this.channels {
		if channel_.CanPop() {
			channel_command := channel_.Pop()

			dma_command := channel_command.DmaCommand()

			transfer_command := this.scoreboard[dma_command]
			delete(this.scoreboard, dma_command)

			if transfer_command.TransferCommandType() == bank.HOST_TO_DEVICE {
				if dma_command.IsReady() {
					transfer_command.AckDmaCommand(dma_command)
				}
			} else if transfer_command.TransferCommandType() == bank.DEVICE_TO_HOST {
				vm_address := transfer_command.VmAddress()
				size := transfer_command.Size()
				byte_stream := dma_command.ByteStream(dma_command.MramAddress(), dma_command.Size())

				transfer_command.SetByteStream(vm_address, size, byte_stream)

				if dma_command.IsReady() {
					transfer_command.AckDmaCommand(dma_command)
				}
			} else {
				err := errors.New("transfer command type is not valid")
				panic(err)
			}
		}
	}
}

func (this *MemoryController) ServiceVmWaitQ() {
	for i := 0; i < this.vm_wait_q.Length(); i++ {
		transfer_command, _ := this.vm_wait_q.Front(i)

		if transfer_command.IsVmReady() && transfer_command.IsReady() &&
			transfer_command.TransferCommandState() == bank.MIDDLE &&
			this.ready_q.CanPush(1) {
			this.vm_wait_q.Remove(i)

			transfer_command.SetTransferCommandState(bank.END)

			this.ready_q.Push(transfer_command)
		} else if transfer_command.IsVmReady() {
			this.vm_wait_q.Remove(i)

			transfer_command.SetTransferCommandState(bank.MIDDLE)

			this.wait_q.Push(transfer_command)

			channel_id := transfer_command.ChannelId()
			rank_id := transfer_command.RankId()
			dpu_id := transfer_command.DpuId()
			mram_address := transfer_command.MramAddress()
			size := transfer_command.Size()
			byte_stream := transfer_command.ByteStream()

			dma_command := new(dram.DmaCommand)
			dma_command.InitWriteToMram(mram_address, size, byte_stream)

			this.scoreboard[dma_command] = transfer_command

			channel_command := new(channel.ChannelCommand)
			channel_command.Init(channel_id, rank_id, dpu_id, dma_command)

			transfer_command.AppendDmaCommand(dma_command)

			this.channel_command_q.Push(channel_command)
		}
	}
}

func (this *MemoryController) ServiceWaitQ() {
	for i := 0; i < this.wait_q.Length(); i++ {
		transfer_command, _ := this.wait_q.Front(i)

		if transfer_command.IsReady() && transfer_command.IsVmReady() &&
			transfer_command.TransferCommandState() == bank.MIDDLE &&
			this.ready_q.CanPush(1) {
			this.wait_q.Remove(i)

			transfer_command.SetTransferCommandState(bank.END)

			this.ready_q.Push(transfer_command)
		} else if transfer_command.IsReady() {
			this.wait_q.Remove(i)

			transfer_command.SetTransferCommandState(bank.MIDDLE)

			this.vm_wait_q.Push(transfer_command)

			segments := this.memory_mapping.Map(transfer_command.VmAddress(), transfer_command.Size())

			for _, segment := range segments {
				dma_command := new(bank.DmaCommand)
				dma_command.InitWrite(segment, transfer_command)

				transfer_command.AppendVmDmaCommand(dma_command)

				this.vm_dma_command_q.Push(dma_command)
			}
		}
	}
}
