package prim

import (
	"errors"
	"math"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Sel struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	buffer_a             []int64
	buffer_c             [][]int64
	input_size_dpu_round int64
	pos                  []int64
	dpu_arg_sizes        []int64
	kernels              []int64
	results              [][]int64
}

func (this *Sel) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.num_executions = 1

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
			a = i + 1
		} else {
			a = 0
		}

		this.buffer_a = append(this.buffer_a, a)
	}

	this.buffer_c = make([][]int64, 0)
	this.pos = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.buffer_c = append(this.buffer_c, make([]int64, 0))

		start_elem := this.input_size_dpu_round * int64(i)
		end_elem := this.input_size_dpu_round * int64(i+1)

		for _, a := range this.buffer_a[start_elem:end_elem] {
			if a%2 != 0 {
				this.buffer_c[i] = append(this.buffer_c[i], a)
			}
		}

		this.pos = append(this.pos, int64(len(this.buffer_c[i])))

		for int64(len(this.buffer_c[i])) < this.input_size_dpu_round {
			this.buffer_c[i] = append(this.buffer_c[i], 0)
		}
	}

	this.dpu_arg_sizes = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.dpu_arg_sizes = append(this.dpu_arg_sizes, this.input_size_dpu_round*elem_size)
	}

	this.kernels = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.kernels = append(this.kernels, 0)
	}

	this.results = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.results = append(this.results, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			if j != this.num_tasklets-1 {
				this.results[i] = append(this.results[i], 0)
			} else {
				this.results[i] = append(this.results[i], this.pos[i])
			}
		}
	}
}

func (this *Sel) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
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

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Sel) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_results_byte_stream := new(encoding.ByteStream)
	dpu_results_byte_stream.Init()

	for _, result := range this.results[dpu_id] {
		result_word := new(word.Word)
		result_word.Init(32)
		result_word.SetValue(result)
		dpu_results_byte_stream.Merge(result_word.ToByteStream())
	}

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_RESULTS"] = dpu_results_byte_stream

	return dpu_host
}

func (this *Sel) InputDpuMramHeapPointerName(
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

	start_elem := this.input_size_dpu_round * int64(dpu_id)

	for i := int64(0); i < this.input_size_dpu_round; i++ {
		element_word := new(word.Word)
		element_word.Init(64)
		element_word.SetValue(this.buffer_a[start_elem+i])
		byte_stream.Merge(element_word.ToByteStream())
	}

	return 0, byte_stream
}

func (this *Sel) OutputDpuMramHeapPointerName(
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

	for i, element := range this.buffer_c[dpu_id] {
		if int64(i) >= this.pos[dpu_id] {
			break
		}

		element_word := new(word.Word)
		element_word.Init(64)
		element_word.SetValue(int64(element))
		byte_stream.Merge(element_word.ToByteStream())
	}

	return this.input_size_dpu_round * 8, byte_stream
}

func (this *Sel) NumExecutions() int {
	return this.num_executions
}
