package prim

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Mlp struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	m_size      int64
	n_size      int64
	num_layers  int64
	n_size_pads []int64
	nr_rows     []int64
	max_rows    []int64
	buffer_a    [][][]int64
	buffer_b    [][]int64
	buffer_c    [][][]int64
}

func (this *Mlp) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.m_size = int64(command_line_parser.DataPrepParams()[0])
	this.n_size = int64(command_line_parser.DataPrepParams()[0])

	this.num_layers = 3

	this.num_executions = int(this.num_layers)

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

	this.buffer_a = make([][][]int64, 0)
	for i := int64(0); i < this.num_layers; i++ {
		this.buffer_a = append(this.buffer_a, make([][]int64, 0))

		for j := int64(0); j < this.m_size; j++ {
			this.buffer_a[i] = append(this.buffer_a[i], make([]int64, 0))

			for k := int64(0); k < this.n_size; k++ {
				if k%100 < 98 {
					this.buffer_a[i][j] = append(this.buffer_a[i][j], 0)
				} else {
					this.buffer_a[i][j] = append(this.buffer_a[i][j], (i+k)%2)
				}
			}
		}
	}

	this.buffer_b = make([][]int64, 0)
	for i := int64(0); i < this.num_layers; i++ {
		this.buffer_b = append(this.buffer_b, make([]int64, 0))

		for j := int64(0); j < this.n_size; j++ {
			if j%50 < 48 {
				this.buffer_b[i] = append(this.buffer_b[i], 0)
			} else {
				this.buffer_b[i] = append(this.buffer_b[i], j%2)
			}
		}
	}

	this.buffer_c = make([][][]int64, 0)
	for i := int64(0); i < this.num_layers; i++ {
		this.buffer_c = append(this.buffer_c, make([][]int64, 0))

		for j := 0; j < this.num_dpus; j++ {
			this.buffer_c[i] = append(this.buffer_c[i], make([]int64, 0))
		}
	}

	for i := int64(0); i < this.num_layers; i++ {
		for j := 0; j < this.num_dpus; j++ {
			start_row := this.nr_rows[j] * int64(j)
			end_row := this.nr_rows[j] * int64(j+1)

			this.buffer_c[i][j] = this.MatMul(this.buffer_a[i][start_row:end_row], this.buffer_b[i])

			for k := 0; k < len(this.buffer_c[i][j]); k++ {
				if this.buffer_c[i][j][k] < 0 {
					this.buffer_c[i][j][k] = 0
				}
			}

			if i < this.num_layers-1 {
				if len(this.buffer_b[i+1][start_row:end_row]) != len(this.buffer_c[i][j]) {
					err := errors.New("buffer b [i+1][start_row:end_row]'s length != buffer c [i][j]'s length")
					panic(err)
				}

				for k, element := range this.buffer_c[i][j] {
					this.buffer_b[i+1][k] = element
				}
			}
		}
	}
}

func (this *Mlp) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
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

func (this *Mlp) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	return make(map[string]*encoding.ByteStream, 0)
}

func (this *Mlp) InputDpuMramHeapPointerName(
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

	for _, row := range this.buffer_a[execution][start_row:end_row] {
		for _, element := range row {
			element_word := new(word.Word)
			element_word.Init(32)
			element_word.SetValue(element)
			byte_stream.Merge(element_word.ToByteStream())
		}
	}

	for _, element := range this.buffer_b[execution] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return 0, byte_stream
}

func (this *Mlp) OutputDpuMramHeapPointerName(
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
	for i := 0; i < len(this.buffer_a[execution][start_row:end_row]); i++ {
		for j := 0; j < len(this.buffer_a[execution][i]); j++ {
			offset += 4
		}
	}

	for i := 0; i < len(this.buffer_b[execution]); i++ {
		offset += 4
	}

	for _, element := range this.buffer_c[execution][dpu_id] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return offset, byte_stream
}

func (this *Mlp) NumExecutions() int {
	return this.num_executions
}

func (this *Mlp) MatMul(x [][]int64, y []int64) []int64 {
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
