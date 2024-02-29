package prim

import (
	"errors"
	"math"
	"math/rand"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Red struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	input_size_dpu_8bytes int64
	buffer_a              []int64
	counts                []int64
	dpu_arg_sizes         []int64
	kernels               []int64
	input_t_counts        []int64
	cycles                [][]int64
	t_counts              [][]int64
}

func (this *Red) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.num_executions = 1

	size := int64(command_line_parser.DataPrepParams()[0])

	elem_size := int64(8)

	is_strong_scaling := true

	var input_size int64
	if is_strong_scaling {
		input_size = size
	} else {
		input_size = size * int64(this.num_dpus)
	}

	var input_size_8bytes int64
	if (input_size*elem_size)%8 == 0 {
		input_size_8bytes = input_size
	} else {
		input_size_8bytes = int64(math.Ceil(float64(input_size)/float64(8)) * 8)
	}

	input_size_dpu := (input_size-1)/int64(this.num_dpus) + 1

	if (input_size_dpu*elem_size)%8 == 0 {
		this.input_size_dpu_8bytes = input_size_dpu
	} else {
		this.input_size_dpu_8bytes = int64(math.Ceil(float64(input_size_dpu)/float64(8)) * 8)
	}

	this.buffer_a = make([]int64, 0)
	for i := int64(0); i < this.input_size_dpu_8bytes*int64(this.num_dpus); i++ {
		a := int64(rand.Intn(this.Pow2(31)))

		this.buffer_a = append(this.buffer_a, a)
	}

	this.counts = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.counts = append(this.counts, 0)
	}

	this.dpu_arg_sizes = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.dpu_arg_sizes = append(this.dpu_arg_sizes, 0)
	}

	for i := 0; i < this.num_dpus; i++ {
		start_elem := this.input_size_dpu_8bytes * int64(i)

		var end_elem int64
		if i != this.num_dpus-1 {
			end_elem = this.input_size_dpu_8bytes * int64(i+1)
		} else {
			end_elem = input_size
		}

		this.counts[i] = this.Sum(this.buffer_a[start_elem:end_elem])

		if i != this.num_dpus-1 {
			this.dpu_arg_sizes[i] = this.input_size_dpu_8bytes * elem_size
		} else {
			this.dpu_arg_sizes[i] = (input_size_8bytes - this.input_size_dpu_8bytes*int64(this.num_dpus-1)) * elem_size
		}
	}

	this.kernels = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.kernels = append(this.kernels, 0)
	}

	this.input_t_counts = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.input_t_counts = append(this.input_t_counts, 0)
	}

	this.cycles = make([][]int64, 0)
	this.t_counts = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.cycles = append(this.cycles, make([]int64, 0))
		this.t_counts = append(this.t_counts, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			this.cycles[i] = append(this.cycles[i], 0)

			if j == 0 {
				this.t_counts[i] = append(this.t_counts[i], this.counts[i])
			} else {
				this.t_counts[i] = append(this.t_counts[i], 0)
			}
		}
	}
}

func (this *Red) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_input_arguments_byte_stream := new(encoding.ByteStream)
	dpu_input_arguments_byte_stream.Init()

	dpu_arg_size_word := new(word.Word)
	dpu_arg_size_word.Init(32)
	dpu_arg_size_word.SetValue(this.dpu_arg_sizes[dpu_id])
	dpu_input_arguments_byte_stream.Merge(dpu_arg_size_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernels[dpu_id])
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	input_t_count_word := new(word.Word)
	input_t_count_word.Init(32)
	input_t_count_word.SetValue(this.input_t_counts[dpu_id])
	dpu_input_arguments_byte_stream.Merge(input_t_count_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Red) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_results_byte_stream := new(encoding.ByteStream)
	dpu_results_byte_stream.Init()

	if len(this.cycles[dpu_id]) != len(this.t_counts[dpu_id]) {
		err := errors.New("cycles' length != t counts' length")
		panic(err)
	}

	for i := 0; i < len(this.cycles[dpu_id]); i++ {
		cycle_word := new(word.Word)
		cycle_word.Init(64)
		cycle_word.SetValue(this.cycles[dpu_id][i])
		dpu_results_byte_stream.Merge(cycle_word.ToByteStream())

		t_count_word := new(word.Word)
		t_count_word.Init(64)
		t_count_word.SetValue(this.t_counts[dpu_id][i])
		dpu_results_byte_stream.Merge(t_count_word.ToByteStream())
	}

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_RESULTS"] = dpu_results_byte_stream

	return dpu_host
}

func (this *Red) InputDpuMramHeapPointerName(
	execution int,
	dpu_id int,
) (int64, *encoding.ByteStream) {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	start_elem := this.input_size_dpu_8bytes * int64(dpu_id)

	for i := int64(0); i < this.input_size_dpu_8bytes; i++ {
		element_word := new(word.Word)
		element_word.Init(64)
		element_word.SetValue(this.buffer_a[start_elem+i])
		byte_stream.Merge(element_word.ToByteStream())
	}

	return 0, byte_stream
}

func (this *Red) OutputDpuMramHeapPointerName(
	execution int,
	dpu_id int,
) (int64, *encoding.ByteStream) {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	return 0, byte_stream
}

func (this *Red) NumExecutions() int {
	return this.num_executions
}

func (this *Red) Sum(s []int64) int64 {
	sum := int64(0)
	for _, element := range s {
		sum += element
	}
	return sum
}

func (this *Red) Pow2(exponent int) int {
	if exponent < 0 {
		err := errors.New("exponent < 0")
		panic(err)
	}

	value := 1
	for i := 0; i < exponent; i++ {
		value *= 2
	}
	return value
}
