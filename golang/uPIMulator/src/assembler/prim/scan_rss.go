package prim

import (
	"errors"
	"math"
	"math/rand"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type ScanRss struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	input_size_dpu_round int64
	buffer_a             []int64
	buffer_c             []int64
	last_result_values   []int64
	result_t_counts      [][]int64
	dpu_arg_size         int64
	kernels              []int64
	t_counts             [][]int64
}

func (this *ScanRss) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.num_executions = 2

	size := int64(command_line_parser.DataPrepParams()[0])

	elem_size := int64(8)

	regs := int64(128)

	is_strong_scaling := true

	var input_size int64
	if is_strong_scaling {
		input_size = size
	} else {
		input_size = size * int64(this.num_dpus)
	}

	input_size_dpu := (input_size-1)/int64(this.num_dpus) + 1

	if input_size_dpu%(int64(this.num_tasklets)*regs) == 0 {
		this.input_size_dpu_round = input_size_dpu
	} else {
		this.input_size_dpu_round = int64(math.Ceil(float64(input_size_dpu)/float64(int64(this.num_tasklets)*regs))) * int64(this.num_tasklets) * regs
	}

	this.buffer_a = make([]int64, 0)
	for i := int64(0); i < this.input_size_dpu_round*int64(this.num_dpus); i++ {
		var a int64
		if i < input_size {
			a = int64(rand.Intn(100))
		} else {
			a = 0
		}

		this.buffer_a = append(this.buffer_a, a)
	}

	this.buffer_c = make([]int64, 0)
	for i := int64(0); i < this.input_size_dpu_round*int64(this.num_dpus); i++ {
		var c int64
		if i == 0 {
			c = this.buffer_a[i]
		} else {
			c = this.buffer_c[i-1] + this.buffer_a[i]
		}

		this.buffer_c = append(this.buffer_c, c)
	}

	this.last_result_values = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.last_result_values = append(this.last_result_values, 0)
	}

	this.result_t_counts = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.result_t_counts = append(this.result_t_counts, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			this.result_t_counts[i] = append(this.result_t_counts[i], 0)
		}
	}

	for i := 0; i < this.num_dpus; i++ {
		start_elem := this.input_size_dpu_round * int64(i)
		end_elem := this.input_size_dpu_round * int64(i+1)

		this.last_result_values[i] = this.Sum1D(this.buffer_a[start_elem:end_elem])
		this.result_t_counts[i][0] = this.last_result_values[i]
	}

	this.dpu_arg_size = this.input_size_dpu_round * elem_size
	this.kernels = []int64{0, 1}

	this.t_counts = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.t_counts = append(this.t_counts, make([]int64, 0))

		for j := 0; j < 2; j++ {
			if i == 0 {
				this.t_counts[i] = []int64{
					0,
					0,
				}
			} else {
				this.t_counts[i] = []int64{
					0,
					this.Sum2D(this.result_t_counts[0:i]),
				}
			}
		}
	}
}

func (this *ScanRss) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
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
	dpu_arg_size_word.SetValue(this.dpu_arg_size)
	dpu_input_arguments_byte_stream.Merge(dpu_arg_size_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernels[execution])
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	t_count_word := new(word.Word)
	t_count_word.Init(64)
	t_count_word.SetValue(this.t_counts[dpu_id][execution])
	dpu_input_arguments_byte_stream.Merge(t_count_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *ScanRss) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_results_byte_stream := new(encoding.ByteStream)
	dpu_results_byte_stream.Init()

	if execution == 0 {
		for _, result_t_count := range this.result_t_counts[dpu_id] {
			result_t_count_word := new(word.Word)
			result_t_count_word.Init(64)
			result_t_count_word.SetValue(result_t_count)
			dpu_results_byte_stream.Merge(result_t_count_word.ToByteStream())
		}
	}

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_RESULTS"] = dpu_results_byte_stream

	return dpu_host
}

func (this *ScanRss) InputDpuMramHeapPointerName(
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

	if execution == 0 {
		start_elem := this.input_size_dpu_round * int64(dpu_id)
		end_elem := this.input_size_dpu_round * int64(dpu_id+1)

		for _, element := range this.buffer_a[start_elem:end_elem] {
			element_word := new(word.Word)
			element_word.Init(64)
			element_word.SetValue(element)
			byte_stream.Merge(element_word.ToByteStream())
		}
	}

	return 0, byte_stream
}

func (this *ScanRss) OutputDpuMramHeapPointerName(
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

	if execution == this.num_executions-1 {
		start_elem := this.input_size_dpu_round * int64(dpu_id)
		end_elem := this.input_size_dpu_round * int64(dpu_id+1)

		for _, element := range this.buffer_c[start_elem:end_elem] {
			element_word := new(word.Word)
			element_word.Init(64)
			element_word.SetValue(element)
			byte_stream.Merge(element_word.ToByteStream())
		}
	}

	return int64(this.input_size_dpu_round * 8), byte_stream
}

func (this *ScanRss) NumExecutions() int {
	return this.num_executions
}

func (this *ScanRss) Sum1D(s []int64) int64 {
	sum := int64(0)
	for _, element := range s {
		sum += element
	}
	return sum
}

func (this *ScanRss) Sum2D(s [][]int64) int64 {
	sum := int64(0)
	for _, elements := range s {
		for _, element := range elements {
			sum += element
		}
	}
	return sum
}

func (this *ScanRss) Pow2(exponent int) int {
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
