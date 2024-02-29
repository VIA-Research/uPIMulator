package prim

import (
	"errors"
	"math"
	"math/rand"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Va struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	input_size_dpu_8bytes int64
	buffer_a              []int64
	buffer_b              []int64
	buffer_c              []int64
	sizes                 []int64
	transfer_sizes        []int64
	kernels               []int64
}

func (this *Va) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.num_executions = 1

	buffer_size := int64(command_line_parser.DataPrepParams()[0])

	elem_size := int64(4)

	is_strong_scaling := true

	var input_size int64
	if is_strong_scaling {
		input_size = buffer_size
	} else {
		input_size = buffer_size * int64(this.num_dpus)
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
	this.buffer_b = make([]int64, 0)
	this.buffer_c = make([]int64, 0)
	for i := int64(0); i < this.input_size_dpu_8bytes*int64(this.num_dpus); i++ {
		a := int64(rand.Intn(this.Pow2(31)))
		b := int64(rand.Intn(this.Pow2(31)))

		c := a + b

		this.buffer_a = append(this.buffer_a, a)
		this.buffer_b = append(this.buffer_b, b)
		this.buffer_c = append(this.buffer_c, c)
	}

	this.sizes = make([]int64, 0)
	for i := 0; i < this.num_dpus-1; i++ {
		this.sizes = append(this.sizes, this.input_size_dpu_8bytes*elem_size)
	}
	size := (input_size_8bytes - this.input_size_dpu_8bytes*int64(this.num_dpus-1)) * elem_size
	this.sizes = append(this.sizes, size)

	this.transfer_sizes = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.transfer_sizes = append(this.transfer_sizes, this.input_size_dpu_8bytes*elem_size)
	}

	this.kernels = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.kernels = append(this.kernels, 0)
	}
}

func (this *Va) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_input_arguments_byte_stream := new(encoding.ByteStream)
	dpu_input_arguments_byte_stream.Init()

	size_word := new(word.Word)
	size_word.Init(32)
	size_word.SetValue(this.sizes[dpu_id])
	dpu_input_arguments_byte_stream.Merge(size_word.ToByteStream())

	transfer_size_word := new(word.Word)
	transfer_size_word.Init(32)
	transfer_size_word.SetValue(this.transfer_sizes[dpu_id])
	dpu_input_arguments_byte_stream.Merge(transfer_size_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernels[dpu_id])
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Va) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	return make(map[string]*encoding.ByteStream, 0)
}

func (this *Va) InputDpuMramHeapPointerName(
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
		element_word.Init(32)
		element_word.SetValue(this.buffer_a[start_elem+i])
		byte_stream.Merge(element_word.ToByteStream())
	}

	for i := int64(0); i < this.input_size_dpu_8bytes; i++ {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(this.buffer_b[start_elem+i])
		byte_stream.Merge(element_word.ToByteStream())
	}

	return 0, byte_stream
}

func (this *Va) OutputDpuMramHeapPointerName(
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
		element_word.Init(32)
		element_word.SetValue(this.buffer_c[start_elem+i])
		byte_stream.Merge(element_word.ToByteStream())
	}

	return this.input_size_dpu_8bytes * 4, byte_stream
}

func (this *Va) NumExecutions() int {
	return this.num_executions
}

func (this *Va) Pow2(exponent int) int {
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
