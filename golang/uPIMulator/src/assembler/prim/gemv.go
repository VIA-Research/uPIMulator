package prim

import (
	"errors"
	"math/rand"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Gemv struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	m_size      int64
	n_size      int64
	n_size_pads []int64
	nr_rows     []int64
	max_rows    []int64
	buffer_a    [][]int64
	buffer_b    []int64
	buffer_c    [][]int64
}

func (this *Gemv) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.m_size = int64(command_line_parser.DataPrepParams()[0])
	this.n_size = 64

	this.num_executions = 1

	if this.m_size%int64(this.num_dpus) != 0 {
		err := errors.New("m size % num dpus != 0")
		panic(err)
	}

	this.n_size_pads = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		if this.n_size%2 == 0 {
			this.n_size_pads = append(this.n_size_pads, this.n_size)
		} else {
			this.n_size_pads = append(this.n_size_pads, this.n_size+1)
		}
	}

	this.nr_rows = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.nr_rows = append(this.nr_rows, this.m_size/int64(this.num_dpus))
	}

	this.max_rows = make([]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		if this.nr_rows[i]%2 == 0 {
			this.max_rows = append(this.max_rows, this.nr_rows[i])
		} else {
			this.max_rows = append(this.max_rows, this.nr_rows[i]+1)
		}
	}

	this.buffer_a = make([][]int64, 0)
	for i := int64(0); i < this.m_size; i++ {
		this.buffer_a = append(this.buffer_a, make([]int64, 0))

		for j := int64(0); j < this.n_size; j++ {
			this.buffer_a[i] = append(this.buffer_a[i], int64(rand.Intn(50)))
		}
	}

	this.buffer_b = make([]int64, 0)
	for i := int64(0); i < this.n_size; i++ {
		this.buffer_b = append(this.buffer_b, int64(rand.Intn(50)))
	}

	this.buffer_c = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.buffer_c = append(this.buffer_c, make([]int64, 0))
	}

	for i := 0; i < this.num_dpus; i++ {
		start_row := this.nr_rows[i] * int64(i)
		end_row := this.nr_rows[i] * int64(i+1)

		this.buffer_c[i] = this.MatMul(this.buffer_a[start_row:end_row], this.buffer_b)
	}
}

func (this *Gemv) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_input_arguments_byte_stream := new(encoding.ByteStream)
	dpu_input_arguments_byte_stream.Init()

	n_size_word := new(word.Word)
	n_size_word.Init(32)
	n_size_word.SetValue(this.n_size)
	dpu_input_arguments_byte_stream.Merge(n_size_word.ToByteStream())

	n_size_pad_word := new(word.Word)
	n_size_pad_word.Init(32)
	n_size_pad_word.SetValue(this.n_size_pads[dpu_id])
	dpu_input_arguments_byte_stream.Merge(n_size_pad_word.ToByteStream())

	nr_row_word := new(word.Word)
	nr_row_word.Init(32)
	nr_row_word.SetValue(this.nr_rows[dpu_id])
	dpu_input_arguments_byte_stream.Merge(nr_row_word.ToByteStream())

	max_row_word := new(word.Word)
	max_row_word.Init(32)
	max_row_word.SetValue(this.max_rows[dpu_id])
	dpu_input_arguments_byte_stream.Merge(max_row_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Gemv) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	return make(map[string]*encoding.ByteStream, 0)
}

func (this *Gemv) InputDpuMramHeapPointerName(
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

	start_row := this.nr_rows[dpu_id] * int64(dpu_id)
	end_row := this.nr_rows[dpu_id] * int64(dpu_id+1)

	for _, row := range this.buffer_a[start_row:end_row] {
		for _, element := range row {
			element_word := new(word.Word)
			element_word.Init(32)
			element_word.SetValue(element)
			byte_stream.Merge(element_word.ToByteStream())
		}
	}

	for _, element := range this.buffer_b {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return 0, byte_stream
}

func (this *Gemv) OutputDpuMramHeapPointerName(
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

	start_row := this.nr_rows[dpu_id] * int64(dpu_id)
	end_row := this.nr_rows[dpu_id] * int64(dpu_id+1)

	offset := int64(0)
	for i := 0; i < len(this.buffer_a[start_row:end_row]); i++ {
		for j := 0; j < len(this.buffer_a[i]); j++ {
			offset += 4
		}
	}

	for i := 0; i < len(this.buffer_b); i++ {
		offset += 4
	}

	for _, element := range this.buffer_c[dpu_id] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return offset, byte_stream
}

func (this *Gemv) NumExecutions() int {
	return this.num_executions
}

func (this *Gemv) MatMul(x [][]int64, y []int64) []int64 {
	results := make([]int64, 0)
	for i := 0; i < len(x); i++ {
		results = append(results, 0)
	}

	for i := 0; i < len(x); i++ {
		for j := 0; j < len(y); j++ {
			results[i] += x[i][j] * y[j]
		}
	}

	return results
}
