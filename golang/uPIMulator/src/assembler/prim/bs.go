package prim

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Bs struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	size              int64
	input_buffer      []int64
	query_buffer      []int64
	results           [][]int64
	slice_per_dpu     int64
	query_per_tasklet int64
	kernel            int64
}

func (this *Bs) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.num_executions = 1

	this.size = int64(command_line_parser.DataPrepParams()[0])
	num_queries := int64(command_line_parser.DataPrepParams()[0] / 8)

	if num_queries%int64(this.num_dpus*this.num_tasklets) != 0 {
		num_queries += int64(this.num_dpus*this.num_tasklets) - num_queries%int64(this.num_dpus*this.num_tasklets)
	}

	if num_queries%int64(this.num_dpus*this.num_tasklets) != 0 {
		err := errors.New("num queries % (num dpus * num tasklets) != 0")
		panic(err)
	}

	this.input_buffer = make([]int64, 0)
	for i := int64(0); i < this.size; i++ {
		this.input_buffer = append(this.input_buffer, i+1)
	}

	this.query_buffer = make([]int64, 0)
	for i := int64(0); i < num_queries; i++ {
		this.query_buffer = append(this.query_buffer, i)
	}

	this.results = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.results = append(this.results, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			this.results[i] = append(this.results[i], 0)
		}
	}

	this.slice_per_dpu = num_queries / int64(this.num_dpus)
	this.query_per_tasklet = this.slice_per_dpu / int64(this.num_tasklets)

	for i := 0; i < this.num_dpus; i++ {
		for j := 0; j < this.num_tasklets; j++ {
			for k := 0; k < int(this.query_per_tasklet); k++ {
				is_found := false
				l := int64(0)
				r := this.size - 1

				for l <= r {
					m := l + (r-l)/2

					if this.input_buffer[m] == this.query_buffer[k+j*int(this.query_per_tasklet)+i*int(this.slice_per_dpu)] {
						this.results[i][j] = m
						is_found = true
						break
					}

					if this.input_buffer[m] < this.query_buffer[k+j*int(this.query_per_tasklet)+i*int(this.slice_per_dpu)] {
						l = m + 1
					} else {
						r = m - 1
					}
				}

				if !is_found {
					this.results[i][j] = -1
				}
			}
		}
	}

	this.kernel = 0
}

func (this *Bs) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
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
	size_word.Init(64)
	size_word.SetValue(this.size)
	dpu_input_arguments_byte_stream.Merge(size_word.ToByteStream())

	slice_per_dpu_word := new(word.Word)
	slice_per_dpu_word.Init(64)
	slice_per_dpu_word.SetValue(this.slice_per_dpu)
	dpu_input_arguments_byte_stream.Merge(slice_per_dpu_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernel)
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Bs) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
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
		result_word.Init(64)
		result_word.SetValue(result)
		dpu_results_byte_stream.Merge(result_word.ToByteStream())
	}

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_RESULTS"] = dpu_results_byte_stream

	return dpu_host
}

func (this *Bs) InputDpuMramHeapPointerName(
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

	for _, element := range this.input_buffer {
		element_word := new(word.Word)
		element_word.Init(64)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	start_elem := this.slice_per_dpu * int64(dpu_id)
	end_elem := this.slice_per_dpu * int64(dpu_id+1)

	for _, element := range this.query_buffer[start_elem:end_elem] {
		element_word := new(word.Word)
		element_word.Init(64)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}
	return 0, byte_stream
}

func (this *Bs) OutputDpuMramHeapPointerName(
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

func (this *Bs) NumExecutions() int {
	return this.num_executions
}

func (this *Bs) Sum(s []int64) int64 {
	sum := int64(0)
	for _, element := range s {
		sum += element
	}
	return sum
}

func (this *Bs) Pow2(exponent int) int {
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
