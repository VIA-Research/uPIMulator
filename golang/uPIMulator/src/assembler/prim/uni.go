package prim

import (
	"errors"
	"math"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Uni struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	buffer_a             []int64
	buffer_c             [][]int64
	input_size_dpu_round int64
	pos                  int64
	input_sizes_dpu      []int64
	kernels              []int64
	t_counts             [][]int64
	firsts               [][]int64
	lasts                [][]int64
}

func (this *Uni) Init(command_line_parser *misc.CommandLineParser) {
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
			if i%2 == 0 {
				a = i
			} else {
				a = i + 1
			}
		} else {
			a = this.buffer_a[input_size-1]
		}

		this.buffer_a = append(this.buffer_a, a)
	}

	this.buffer_c = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.buffer_c = append(this.buffer_c, make([]int64, 0))

		for j := int64(0); j < this.input_size_dpu_round; j++ {
			this.buffer_c[i] = append(this.buffer_c[i], 0)
		}
	}

	for i := 0; i < this.num_dpus; i++ {
		start_elem := this.input_size_dpu_round * int64(i)

		this.buffer_c[i][0] = this.buffer_a[start_elem]

		this.pos = 1
		for j := int64(1); j < this.input_size_dpu_round; j++ {
			if this.buffer_a[start_elem+j] != this.buffer_a[start_elem+j-1] {
				this.buffer_c[i][this.pos] = this.buffer_a[start_elem+j]
				this.pos++
			}
		}
	}

	this.input_sizes_dpu = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.input_sizes_dpu = append(this.input_sizes_dpu, this.input_size_dpu_round*elem_size)
	}

	this.kernels = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.kernels = append(this.kernels, 0)
	}

	this.t_counts = make([][]int64, 0)
	this.firsts = make([][]int64, 0)
	this.lasts = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		start_elem := this.input_size_dpu_round * int64(i)

		this.t_counts = append(this.t_counts, make([]int64, 0))
		this.firsts = append(this.firsts, make([]int64, 0))
		this.lasts = append(this.lasts, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			var t_count int64
			var first int64
			var last int64

			if j == 0 && j != this.num_tasklets-1 {
				t_count = 0
				first = this.buffer_a[start_elem]
				last = 0
			} else if j == this.num_tasklets-1 {
				t_count = this.pos
				first = 0
				last = this.buffer_c[i][this.pos-1]
			} else {
				t_count = 0
				first = 0
				last = 0
			}

			this.t_counts[i] = append(this.t_counts[i], t_count)
			this.firsts[i] = append(this.firsts[i], first)
			this.lasts[i] = append(this.lasts[i], last)
		}
	}
}

func (this *Uni) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_input_arguments_byte_stream := new(encoding.ByteStream)
	dpu_input_arguments_byte_stream.Init()

	input_size_dpu_word := new(word.Word)
	input_size_dpu_word.Init(32)
	input_size_dpu_word.SetValue(this.input_sizes_dpu[dpu_id])
	dpu_input_arguments_byte_stream.Merge(input_size_dpu_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernels[dpu_id])
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Uni) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_results_byte_stream := new(encoding.ByteStream)
	dpu_results_byte_stream.Init()

	if len(this.t_counts[dpu_id]) != len(this.firsts[dpu_id]) {
		err := errors.New("t counts' length != firsts' length")
		panic(err)
	} else if len(this.t_counts[dpu_id]) != len(this.lasts[dpu_id]) {
		err := errors.New("t counts' length != lasts' length")
		panic(err)
	}

	for i := 0; i < len(this.t_counts[dpu_id]); i++ {
		t_count := this.t_counts[dpu_id][i]
		first := this.firsts[dpu_id][i]
		last := this.lasts[dpu_id][i]

		t_count_word := new(word.Word)
		t_count_word.Init(64)
		t_count_word.SetValue(t_count)
		dpu_results_byte_stream.Merge(t_count_word.ToByteStream())

		first_word := new(word.Word)
		first_word.Init(64)
		first_word.SetValue(first)
		dpu_results_byte_stream.Merge(first_word.ToByteStream())

		last_word := new(word.Word)
		last_word.Init(64)
		last_word.SetValue(last)
		dpu_results_byte_stream.Merge(last_word.ToByteStream())
	}

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_RESULTS"] = dpu_results_byte_stream

	return dpu_host
}

func (this *Uni) InputDpuMramHeapPointerName(
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

func (this *Uni) OutputDpuMramHeapPointerName(
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
		element_word.Init(64)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return this.input_size_dpu_round * 8, byte_stream
}

func (this *Uni) NumExecutions() int {
	return this.num_executions
}
