package compiler

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"uPIMulator/src/misc"
)

type Compiler struct {
	root_dirpath string
	benchmark    string

	num_dpus     int
	num_tasklets int
}

func (this *Compiler) Init(command_line_parser *misc.CommandLineParser) {
	this.root_dirpath = command_line_parser.StringParameter("root_dirpath")
	this.benchmark = command_line_parser.StringParameter("benchmark")

	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))
	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank

	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.Build()
}

func (this *Compiler) Build() {
	docker_dirpath := filepath.Join(this.root_dirpath, "docker")

	command := exec.Command("docker", "build", "-t", "bongjoonhyun/upimulator", docker_dirpath)

	err := command.Run()

	if err != nil {
		panic(err)
	}
}

func (this *Compiler) Compile() {
	this.CompileBenchmark()
	this.CompileSdk()
}

func (this *Compiler) CompileBenchmark() {
	command := exec.Command(
		"docker",
		"run",
		"--privileged",
		"--rm",
		"-v",
		this.root_dirpath+":/root/uPIMulator",
		"bongjoonhyun/upimulator",
		"python3",
		"/root/uPIMulator/benchmark/build.py",
		"--num_dpus",
		strconv.Itoa(this.num_dpus),
		"--num_tasklets",
		strconv.Itoa(this.num_tasklets),
	)

	err := command.Run()

	if err != nil {
		panic(err)
	}
}

func (this *Compiler) CompileSdk() {
	command := exec.Command(
		"docker",
		"run",
		"--privileged",
		"--rm",
		"-v",
		this.root_dirpath+":/root/uPIMulator",
		"bongjoonhyun/upimulator",
		"python3",
		"/root/uPIMulator/sdk/build.py",
		"--num_tasklets",
		strconv.Itoa(this.num_tasklets),
	)

	err := command.Run()

	if err != nil {
		panic(err)
	}
}
