package dpu

import (
	"errors"
	"fmt"
	"uPIMulator/src/device/simulator/dpu/dram"
	"uPIMulator/src/device/simulator/dpu/logic"
	"uPIMulator/src/device/simulator/dpu/sram"
	"uPIMulator/src/misc"
)

type Dpu struct {
	channel_id int
	rank_id    int
	dpu_id     int

	logic_frequency  int64
	memory_frequency int64
	frequency_ratio  float64
	cycles           int64

	control_interface *ControlInterface
	threads           []*logic.Thread
	thread_scheduler  *logic.ThreadScheduler
	atomic            *sram.Atomic
	iram              *sram.Iram
	wram              *sram.Wram
	mram              *dram.Mram
	operand_collector *logic.OperandCollector
	memory_controller *dram.MemoryController
	dma               *logic.Dma
	logic             *logic.Logic

	stat_factory *misc.StatFactory
}

func (this *Dpu) Init(
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

	this.logic_frequency = command_line_parser.IntParameter("logic_frequency")
	this.memory_frequency = command_line_parser.IntParameter("memory_frequency")
	this.frequency_ratio = float64(this.memory_frequency) / float64(this.logic_frequency)
	this.cycles = 0

	this.control_interface = new(ControlInterface)
	this.control_interface.Init()

	this.threads = make([]*logic.Thread, 0)
	num_threads := int(command_line_parser.IntParameter("num_tasklets"))
	for i := 0; i < num_threads; i++ {
		thread := new(logic.Thread)
		thread.Init(i)
		this.threads = append(this.threads, thread)
	}

	this.thread_scheduler = new(logic.ThreadScheduler)
	this.thread_scheduler.Init(channel_id, rank_id, dpu_id, this.threads, command_line_parser)

	this.atomic = new(sram.Atomic)
	this.atomic.Init()

	this.iram = new(sram.Iram)
	this.iram.Init()

	this.wram = new(sram.Wram)
	this.wram.Init()

	this.mram = new(dram.Mram)
	this.mram.Init(command_line_parser)

	this.operand_collector = new(logic.OperandCollector)
	this.operand_collector.Init()
	this.operand_collector.ConnectWram(this.wram)

	this.memory_controller = new(dram.MemoryController)
	this.memory_controller.Init(channel_id, rank_id, dpu_id, command_line_parser)
	this.memory_controller.ConnectMram(this.mram)

	this.dma = new(logic.Dma)
	this.dma.Init()
	this.dma.ConnectAtomic(this.atomic)
	this.dma.ConnectIram(this.iram)
	this.dma.ConnectOperandCollector(this.operand_collector)
	this.dma.ConnectMemoryController(this.memory_controller)

	this.logic = new(logic.Logic)
	this.logic.Init(channel_id, rank_id, dpu_id, command_line_parser)
	this.logic.ConnectThreadScheduler(this.thread_scheduler)
	this.logic.ConnectAtomic(this.atomic)
	this.logic.ConnectIram(this.iram)
	this.logic.ConnectOperandCollector(this.operand_collector)
	this.logic.ConnectDma(this.dma)

	name := fmt.Sprintf("DPU%d-%d-%d", channel_id, rank_id, dpu_id)
	this.stat_factory = new(misc.StatFactory)
	this.stat_factory.Init(name)
}

func (this *Dpu) Fini() {
	for _, thread := range this.threads {
		thread.Fini()
	}

	this.atomic.Fini()
	this.iram.Fini()
	this.wram.Fini()
	this.mram.Fini()

	this.operand_collector.Fini()
	this.memory_controller.Fini()

	this.logic.Fini()
	this.dma.Fini()
}

func (this *Dpu) ChannelId() int {
	return this.channel_id
}

func (this *Dpu) RankId() int {
	return this.rank_id
}

func (this *Dpu) DpuId() int {
	return this.dpu_id
}

func (this *Dpu) ThreadScheduler() *logic.ThreadScheduler {
	return this.thread_scheduler
}

func (this *Dpu) Logic() *logic.Logic {
	return this.logic
}

func (this *Dpu) MemoryController() *dram.MemoryController {
	return this.memory_controller
}

func (this *Dpu) Dma() *logic.Dma {
	return this.dma
}

func (this *Dpu) Threads() []*logic.Thread {
	return this.threads
}

func (this *Dpu) StatFactory() *misc.StatFactory {
	return this.stat_factory
}

func (this *Dpu) Boot() {
	this.control_interface.SetBoot()
	this.thread_scheduler.Boot(0)
}

func (this *Dpu) Unboot() {
	this.control_interface.UnsetBoot()
	this.thread_scheduler.Unboot()
}

func (this *Dpu) IsZombie() bool {
	for _, thread := range this.threads {
		if thread.ThreadState() != logic.ZOMBIE {
			return false
		}
	}
	return this.logic.IsEmpty() && this.memory_controller.IsEmpty()
}

func (this *Dpu) Cycle() {
	this.dma.Cycle()

	num_memory_cycles := int(
		this.frequency_ratio*float64(this.cycles) - this.frequency_ratio*float64(this.cycles-1),
	)
	for i := 0; i < num_memory_cycles; i++ {
		this.memory_controller.Cycle()
	}

	if this.control_interface.Boot() {
		for _, thread := range this.threads {
			thread.IncrementIssueCycle()
		}

		this.thread_scheduler.Cycle()
		this.logic.Cycle()
	}

	this.cycles++
}
