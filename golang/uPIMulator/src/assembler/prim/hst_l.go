package prim

import (
	"errors"
	"math"
	"math/rand"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type HstL struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	num_bins              int64
	input_size_dpu_8bytes int64
	buffer_a              []int64
	buffer_c              [][]int64
	sizes                 []int64
	dpu_arg_sizes         []int64
	transfer_size         int64
	kernel                int64
}

func (this *HstL) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.num_executions = 1

	size := int64(command_line_parser.DataPrepParams()[0])
	this.num_bins = 256

	elem_size := int64(4)

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
	for i := int64(0); i < input_size; i++ {
		this.buffer_a = append(this.buffer_a, int64(rand.Intn(4096)))
	}

	depth := int64(12)

	this.buffer_c = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.buffer_c = append(this.buffer_c, make([]int64, 0))

		for j := int64(0); j < this.num_bins; j++ {
			this.buffer_c[i] = append(this.buffer_c[i], 0)
		}
	}

	for i := 0; i < this.num_dpus; i++ {
		start_elem := this.input_size_dpu_8bytes * int64(i)
		end_elem := this.input_size_dpu_8bytes * int64(i+1)

		for _, elem := range this.buffer_a[start_elem:end_elem] {
			this.buffer_c[i][(elem*this.num_bins)>>depth] += 1
		}
	}

	this.dpu_arg_sizes = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		if i != this.num_dpus-1 {
			this.dpu_arg_sizes = append(this.dpu_arg_sizes, this.input_size_dpu_8bytes*elem_size)
		} else {
			this.dpu_arg_sizes = append(this.dpu_arg_sizes, (input_size_8bytes-this.input_size_dpu_8bytes*int64(this.num_dpus-1))*elem_size)
		}
	}

	this.transfer_size = this.input_size_dpu_8bytes * elem_size
	this.kernel = 0
}

func (this *HstL) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
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

	transfer_size_word := new(word.Word)
	transfer_size_word.Init(32)
	transfer_size_word.SetValue(this.transfer_size)
	dpu_input_arguments_byte_stream.Merge(transfer_size_word.ToByteStream())

	num_bins_word := new(word.Word)
	num_bins_word.Init(32)
	num_bins_word.SetValue(this.num_bins)
	dpu_input_arguments_byte_stream.Merge(num_bins_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernel)
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *HstL) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	return make(map[string]*encoding.ByteStream, 0)
}

func (this *HstL) InputDpuMramHeapPointerName(
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
	end_elem := this.input_size_dpu_8bytes * int64(dpu_id+1)

	for _, element := range this.buffer_a[start_elem:end_elem] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return 0, byte_stream
}

func (this *HstL) OutputDpuMramHeapPointerName(
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

	for _, element := range this.buffer_c[dpu_id] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return this.input_size_dpu_8bytes * 4, byte_stream
}

func (this *HstL) NumExecutions() int {
	return this.num_executions
}

func (this *HstL) Pow2(exponent int) int {
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
