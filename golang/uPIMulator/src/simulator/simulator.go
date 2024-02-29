package simulator

import (
	"fmt"
	"path/filepath"
	"uPIMulator/src/core"
	"uPIMulator/src/misc"
	"uPIMulator/src/simulator/channel"
	"uPIMulator/src/simulator/host"
)

type Simulator struct {
	host     *host.Host
	channels []*channel.Channel

	bin_dirpath            string
	num_simulation_threads int
	execution              int

	verbose int
}

func (this *Simulator) Init(command_line_parser *misc.CommandLineParser) {
	this.host = new(host.Host)
	this.host.Init(command_line_parser)

	this.verbose = int(command_line_parser.IntParameter("verbose"))

	num_channels := int(command_line_parser.IntParameter("num_channels"))
	this.channels = make([]*channel.Channel, 0)
	for i := 0; i < num_channels; i++ {
		channel_ := new(channel.Channel)
		channel_.Init(i, command_line_parser)

		this.channels = append(this.channels, channel_)
	}

	this.host.ConnectChannels(this.channels)

	this.bin_dirpath = command_line_parser.StringParameter("bin_dirpath")
	this.num_simulation_threads = int(command_line_parser.IntParameter("num_simulation_threads"))
	this.execution = 0

	this.host.Load()
	this.host.Schedule(this.execution)
	this.host.Launch()
}

func (this *Simulator) Fini() {
	this.host.Fini()

	for _, channel_ := range this.channels {
		channel_.Fini()
	}
}

func (this *Simulator) IsFinished() bool {
	return this.execution == this.host.NumExecutions()
}

func (this *Simulator) Cycle() {
	this.host.Cycle()

	thread_pool := new(core.ThreadPool)
	thread_pool.Init(this.num_simulation_threads)

	dpus := this.host.Dpus()
	for _, dpu_ := range dpus {
		cycle_job := new(CycleJob)
		cycle_job.Init(dpu_)

		thread_pool.Enque(cycle_job)
	}

	thread_pool.Start()

	if this.host.IsZombie() {
		fmt.Printf("execution (%d) is finished...\n", this.execution)

		this.host.Check(this.execution)
		this.execution++

		if !this.IsFinished() {
			this.host.Schedule(this.execution)
			this.host.Launch()
		}
	}

	if this.verbose >= 1 {
		fmt.Println("system is cycling...")
	}
}

func (this *Simulator) Dump() {
	file_dumper := new(misc.FileDumper)
	file_dumper.Init(filepath.Join(this.bin_dirpath, "log.txt"))

	lines := make([]string, 0)

	dpus := this.host.Dpus()
	for _, dpu_ := range dpus {
		lines = append(lines, dpu_.StatFactory().ToLines()...)
		lines = append(lines, dpu_.ThreadScheduler().StatFactory().ToLines()...)
		lines = append(lines, dpu_.Logic().StatFactory().ToLines()...)
		lines = append(lines, dpu_.Logic().CycleRule().StatFactory().ToLines()...)
		lines = append(lines, dpu_.MemoryController().StatFactory().ToLines()...)
		lines = append(lines, dpu_.MemoryController().MemoryScheduler().StatFactory().ToLines()...)
		lines = append(lines, dpu_.MemoryController().RowBuffer().StatFactory().ToLines()...)
	}

	file_dumper.WriteLines(lines)
}
