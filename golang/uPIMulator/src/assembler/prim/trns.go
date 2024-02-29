package prim

import (
	"errors"
	"math/rand"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Trns struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	num_active_dpus_at_begining int64
	n                           int64
	M                           int64
	m                           int64
	N                           int64

	buffer_a [][][]int64
	buffer_c [][][]int64
	dones    []int64

	kernels []int64
}

func (this *Trns) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	//NOTE(dongjae.lee@kaist.ac.kr): different parameter given if single dpu is simulated
	if num_channels == 1 && num_ranks_per_channel == 1 && num_dpus_per_rank == 1 {
		this.N = 1
		this.n = 4
		this.M = int64(command_line_parser.DataPrepParams()[0])
		this.m = 16
	} else {
		this.N = 64
		this.n = 8
		this.M = int64(command_line_parser.DataPrepParams()[0])
		this.m = 4
	}

	is_strong_scaling := true
	if !is_strong_scaling {
		this.N *= int64(this.num_dpus)
	}

	if this.N > int64(this.num_dpus) {
		this.num_active_dpus_at_begining = int64(this.num_dpus)
	} else {
		this.num_active_dpus_at_begining = this.N
	}

	if is_strong_scaling {
		this.num_executions = int(2 * (this.N / this.num_active_dpus_at_begining))
	} else {
		this.num_executions = int(2 * this.N)
	}

	this.buffer_a = make([][][]int64, 0)
	for i := int64(0); i < this.N; i++ {
		this.buffer_a = append(this.buffer_a, make([][]int64, 0))

		for j := int64(0); j < this.M*this.m; j++ {
			this.buffer_a[i] = append(this.buffer_a[i], make([]int64, 0))

			for k := int64(0); k < this.n; k++ {
				this.buffer_a[i][j] = append(this.buffer_a[i][j], int64(rand.Intn(100)))
			}
		}
	}

	this.dones = make([]int64, 0)
	if (this.M*this.n)/8 == 0 {
		for i := 0; i < 8; i++ {
			this.dones = append(this.dones, 0)
		}
	} else {
		for i := int64(0); i < this.M*this.n; i++ {
			this.dones = append(this.dones, 0)
		}
	}

	this.buffer_c = make([][][]int64, 0)
	for i := int64(0); i < this.N; i++ {
		this.buffer_c = append(this.buffer_c, this.Transpose(this.buffer_a[i]))
	}

	this.kernels = []int64{0, 1}
}

func (this *Trns) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_input_arguments_byte_stream := new(encoding.ByteStream)
	dpu_input_arguments_byte_stream.Init()

	m_word := new(word.Word)
	m_word.Init(32)
	m_word.SetValue(this.m)
	dpu_input_arguments_byte_stream.Merge(m_word.ToByteStream())

	n_word := new(word.Word)
	n_word.Init(32)
	n_word.SetValue(this.n)
	dpu_input_arguments_byte_stream.Merge(n_word.ToByteStream())

	M_word := new(word.Word)
	M_word.Init(32)
	M_word.SetValue(this.M)
	dpu_input_arguments_byte_stream.Merge(M_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernels[execution%2])
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Trns) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	return make(map[string]*encoding.ByteStream, 0)
}

func (this *Trns) InputDpuMramHeapPointerName(
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

	if execution%2 == 0 {
		for _, row := range this.buffer_a[(this.num_active_dpus_at_begining*(int64(execution)/2))+int64(dpu_id)] {
			for _, element := range row {
				element_word := new(word.Word)
				element_word.Init(64)
				element_word.SetValue(element)
				byte_stream.Merge(element_word.ToByteStream())
			}
		}

		for _, done := range this.dones {
			done_word := new(word.Word)
			done_word.Init(8)
			done_word.SetValue(done)
			byte_stream.Merge(done_word.ToByteStream())
		}
	}

	return 0, byte_stream
}

func (this *Trns) OutputDpuMramHeapPointerName(
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

	if execution%2 == 1 {
		for _, row := range this.buffer_c[(this.num_active_dpus_at_begining*(int64(execution)/2))+int64(dpu_id)] {
			for _, element := range row {
				element_word := new(word.Word)
				element_word.Init(64)
				element_word.SetValue(element)
				byte_stream.Merge(element_word.ToByteStream())
			}
		}
	}

	return 0, byte_stream
}

func (this *Trns) NumExecutions() int {
	return this.num_executions
}

func (this *Trns) Transpose(s [][]int64) [][]int64 {
	xl := len(s[0])
	yl := len(s)

	result := make([][]int64, xl)

	for i, _ := range result {
		result[i] = make([]int64, yl)
	}

	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = s[j][i]
		}
	}
	return result
}
