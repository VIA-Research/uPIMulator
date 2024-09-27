package program

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"uPIMulator/src/encoding"
	"uPIMulator/src/misc"
)

type Task struct {
	bin_dirpath string

	benchmark    string
	num_dpus     int
	num_tasklets int

	addresses map[string]int64
	values    map[string]int64

	atomic *encoding.ByteStream
	iram   *encoding.ByteStream
	wram   *encoding.ByteStream
	mram   *encoding.ByteStream
}

func (this *Task) Init(command_line_parser *misc.CommandLineParser) {
	this.bin_dirpath = command_line_parser.StringParameter("bin_dirpath")

	this.benchmark = command_line_parser.StringParameter("benchmark")

	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))
	this.num_tasklets = num_channels * num_ranks_per_channel * num_dpus_per_rank

	this.InitAddresses()
	this.InitValues()
	this.InitAtomic()
	this.InitIram()
	this.InitWram()
	this.InitMram()
}

func (this *Task) InitAddresses() {
	path := filepath.Join(this.bin_dirpath, "addresses.txt")

	file_scanner := new(misc.FileScanner)
	file_scanner.Init(path)

	lines := file_scanner.ReadLines()

	this.addresses = make(map[string]int64)

	for _, line := range lines {
		words := strings.Split(line, ":")

		name := words[0]
		address, err := strconv.ParseInt(words[1][1:], 10, 64)

		if err != nil {
			panic(err)
		}

		this.addresses[name] = address
	}
}

func (this *Task) InitValues() {
	path := filepath.Join(this.bin_dirpath, "values.txt")

	file_scanner := new(misc.FileScanner)
	file_scanner.Init(path)

	lines := file_scanner.ReadLines()

	this.values = make(map[string]int64)

	for _, line := range lines {
		words := strings.Split(line, ":")

		name := words[0]
		address, err := strconv.ParseInt(words[1][1:], 10, 64)

		if err != nil {
			panic(err)
		}

		this.values[name] = address
	}
}

func (this *Task) InitAtomic() {
	path := filepath.Join(this.bin_dirpath, "atomic.bin")

	this.atomic = this.InitByteStream(path)
}

func (this *Task) InitIram() {
	path := filepath.Join(this.bin_dirpath, "iram.bin")

	this.iram = this.InitByteStream(path)
}

func (this *Task) InitWram() {
	path := filepath.Join(this.bin_dirpath, "wram.bin")

	this.wram = this.InitByteStream(path)
}

func (this *Task) InitMram() {
	path := filepath.Join(this.bin_dirpath, "mram.bin")

	this.mram = this.InitByteStream(path)
}

func (this *Task) InitByteStream(path string) *encoding.ByteStream {
	file_scanner := new(misc.FileScanner)
	file_scanner.Init(path)
	lines := file_scanner.ReadLines()

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for _, line := range lines {
		value, err := strconv.Atoi(line)

		if err != nil {
			panic(err)
		}

		byte_stream.Append(uint8(value))
	}

	return byte_stream
}

func (this *Task) Atomic() *encoding.ByteStream {
	return this.atomic
}

func (this *Task) Iram() *encoding.ByteStream {
	return this.iram
}

func (this *Task) Wram() *encoding.ByteStream {
	return this.wram
}

func (this *Task) Mram() *encoding.ByteStream {
	return this.mram
}

func (this *Task) Addresses() map[string]int64 {
	return this.addresses
}

func (this *Task) Values() map[string]int64 {
	return this.values
}

func (this *Task) SysUsedMramEnd() int64 {
	if _, found := this.values["__sys_used_mram_end"]; !found {
		err := errors.New("__sys_used_mram_end is not found")
		panic(err)
	}

	return this.values["__sys_used_mram_end"]
}

func (this *Task) SysEnd() int64 {
	if _, found := this.addresses["__sys_end"]; !found {
		err := errors.New("__sys_end is not found")
		panic(err)
	}

	return this.addresses["__sys_end"]
}
