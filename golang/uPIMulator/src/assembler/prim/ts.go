package prim

import (
	"errors"
	"math"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type Ts struct {
	num_dpus       int
	num_tasklets   int
	num_executions int

	ts_size        int64
	query_length   int64
	query_mean     int64
	query_std      int64
	slice_per_dpu  int64
	exclusion_zone int64
	kernel         int64

	min_vals [][]int64
	min_idxs [][]int64
	max_vals [][]int64
	max_idxs [][]int64

	query_buffer    []int64
	t_series_buffer []int64
	amean_buffer    []int64
	asigma_buffer   []int64

	block_size int64
	dotpip     int64
	elem_size  int64
}

func (this *Ts) Init(command_line_parser *misc.CommandLineParser) {
	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))

	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank
	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.num_executions = 1

	this.ts_size = int64(command_line_parser.DataPrepParams()[0])
	this.query_length = 64

	if this.ts_size%(int64(this.num_dpus)*int64(this.num_tasklets)*this.query_length) != 0 {
		this.ts_size += int64(this.num_dpus)*int64(this.num_tasklets)*this.query_length - this.ts_size%(int64(this.num_dpus)*int64(this.num_tasklets)*this.query_length)
	}

	this.t_series_buffer = make([]int64, 0)
	for i := int64(0); i < this.ts_size; i++ {
		this.t_series_buffer = append(this.t_series_buffer, i%127)
	}

	this.query_buffer = make([]int64, 0)
	for i := int64(0); i < this.query_length; i++ {
		this.query_buffer = append(this.query_buffer, i%127)
	}

	this.asigma_buffer = make([]int64, 0)
	for i := int64(0); i < this.ts_size; i++ {
		this.asigma_buffer = append(this.asigma_buffer, 0)
	}

	this.amean_buffer = make([]int64, 0)
	for i := int64(0); i < this.ts_size; i++ {
		this.amean_buffer = append(this.amean_buffer, 0)
	}

	acum_sum_buffer := make([]int64, 0)
	for i := int64(0); i < this.ts_size; i++ {
		if i == 0 {
			acum_sum_buffer = append(acum_sum_buffer, this.t_series_buffer[i])
		} else {
			acum_sum_buffer = append(acum_sum_buffer, this.t_series_buffer[i]+acum_sum_buffer[i-1])
		}
	}

	asqcum_sum_buffer := make([]int64, 0)
	for i := int64(0); i < this.ts_size; i++ {
		if i == 0 {
			asqcum_sum_buffer = append(
				asqcum_sum_buffer,
				this.t_series_buffer[i]*this.t_series_buffer[i],
			)
		} else {
			asqcum_sum_buffer = append(asqcum_sum_buffer, this.t_series_buffer[i]*this.t_series_buffer[i]+asqcum_sum_buffer[i-1])
		}
	}

	asum_buffer := make([]int64, 0)
	for i := int64(0); i < this.ts_size-this.query_length+1; i++ {
		if i == 0 {
			asum_buffer = append(asum_buffer, acum_sum_buffer[this.query_length-i-1])
		} else {
			asum_buffer = append(asum_buffer, acum_sum_buffer[this.query_length+i-1]-acum_sum_buffer[i-1])
		}
	}

	asum_sq_buffer := make([]int64, 0)
	for i := int64(0); i < this.ts_size-this.query_length+1; i++ {
		if i == 0 {
			asum_sq_buffer = append(asum_sq_buffer, asqcum_sum_buffer[this.query_length+i-1])
		} else {
			asum_sq_buffer = append(asum_sq_buffer, asqcum_sum_buffer[this.query_length+i-1]-asqcum_sum_buffer[i-1])
		}
	}

	amean_temp_buffer := make([]int64, 0)
	for i := int64(0); i < this.ts_size-this.query_length; i++ {
		amean_temp_buffer = append(amean_temp_buffer, asum_buffer[i]/this.query_length)
	}

	asigma_sq_buffer := make([]int64, 0)
	for i := int64(0); i < this.ts_size-this.query_length; i++ {
		asigma_sq_buffer = append(
			asigma_sq_buffer,
			asum_sq_buffer[i]/this.query_length-this.amean_buffer[i]*this.amean_buffer[i],
		)
	}

	this.asigma_buffer = make([]int64, 0)
	for i := int64(0); i < this.ts_size-this.query_length; i++ {
		this.asigma_buffer = append(
			this.asigma_buffer,
			int64(math.Sqrt(float64(asigma_sq_buffer[i]))),
		)
	}

	this.amean_buffer = make([]int64, 0)
	for i := int64(0); i < this.ts_size-this.query_length; i++ {
		this.amean_buffer = append(this.amean_buffer, amean_temp_buffer[i])
	}

	for i := int64(0); i < this.query_length; i++ {
		this.t_series_buffer = append(this.t_series_buffer, 0)
	}

	for i := int64(0); i < this.query_length*2; i++ {
		this.asigma_buffer = append(this.asigma_buffer, 0)
		this.amean_buffer = append(this.amean_buffer, 0)
	}

	if len(this.t_series_buffer) != len(this.amean_buffer) {
		err := errors.New("t series buffer's length != amean buffer's length")
		panic(err)
	}

	this.query_mean = 0
	for i := int64(0); i < this.query_length; i++ {
		this.query_mean += this.query_buffer[i]
	}
	this.query_mean /= this.query_length

	query_variance := int64(0)
	for i := int64(0); i < this.query_length; i++ {
		query_variance += (this.query_buffer[i] - this.query_mean) * (this.query_buffer[i] - this.query_mean)
	}
	query_variance /= this.query_length

	query_std_dev := math.Sqrt(float64(query_variance))
	this.query_std = int64(query_std_dev)

	this.slice_per_dpu = this.ts_size / int64(this.num_dpus)

	this.min_vals = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.min_vals = append(this.min_vals, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			this.min_vals[i] = append(this.min_vals[i], 0x7FFFFFFF)
		}
	}

	this.min_idxs = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.min_idxs = append(this.min_idxs, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			this.min_idxs[i] = append(this.min_idxs[i], 0)
		}
	}

	this.max_vals = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.max_vals = append(this.max_vals, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			this.max_vals[i] = append(this.max_vals[i], 0)
		}
	}

	this.max_idxs = make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		this.max_idxs = append(this.max_idxs, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			this.max_idxs[i] = append(this.max_idxs[i], 0)
		}
	}

	my_start_elems := make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		my_start_elems = append(my_start_elems, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			my_start_elems[i] = append(
				my_start_elems[i],
				this.slice_per_dpu*int64(i)+int64(j)*(this.slice_per_dpu/int64(this.num_tasklets)),
			)
		}
	}

	my_end_elems := make([][]int64, 0)
	for i := 0; i < this.num_dpus; i++ {
		my_end_elems = append(my_end_elems, make([]int64, 0))

		for j := 0; j < this.num_tasklets; j++ {
			my_end_elems[i] = append(
				my_end_elems[i],
				my_start_elems[i][j]+(this.slice_per_dpu/int64(this.num_tasklets))-1,
			)
		}
	}

	for i := 0; i < this.num_dpus; i++ {
		for j := 0; j < this.num_tasklets; j++ {
			if my_end_elems[i][j] > this.slice_per_dpu*int64(i+1)-this.query_length {
				my_end_elems[i][j] = this.slice_per_dpu*int64(i+1) - this.query_length
			}
		}
	}

	this.block_size = 256
	this.elem_size = 4
	increment := this.block_size / this.elem_size
	this.dotpip = this.block_size / this.elem_size
	iter := int64(0)
	for i := 0; i < this.num_dpus; i++ {
		iter = 0
		for j := 0; j < this.num_tasklets; j++ {
			for k := my_start_elems[i][j]; k < my_end_elems[i][j]; k += increment {
				cache_dotprods := make([]int64, 0)
				for l := int64(0); l < this.dotpip; l++ {
					cache_dotprods = append(cache_dotprods, 0)
				}

				for l := int64(0); l < this.query_length/increment; l++ {
					cache_dotprods = this.DotProduct(this.t_series_buffer[k:k+increment],
						this.t_series_buffer[k+increment:k+2*increment],
						this.query_buffer[l*increment:(l+1)*increment],
						cache_dotprods,
					)
				}

				for l := int64(0); l < increment; l++ {
					distance := 2 * (this.query_length - (cache_dotprods[l]-this.query_length*this.amean_buffer[l+iter*increment+int64(i)*this.slice_per_dpu]*this.query_mean)/(this.asigma_buffer[l+iter*increment+int64(i)*this.slice_per_dpu]*this.query_std))

					if distance < this.min_vals[i][j] {
						this.min_vals[i][j] = distance
						this.min_idxs[i][j] = int64(k) + l - (int64(i) * this.slice_per_dpu)
					}
				}

				iter++
			}
		}
	}

	this.exclusion_zone = 0
	this.kernel = 0
}

func (this *Ts) InputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_input_arguments_byte_stream := new(encoding.ByteStream)
	dpu_input_arguments_byte_stream.Init()

	ts_size_word := new(word.Word)
	ts_size_word.Init(32)
	ts_size_word.SetValue(this.ts_size)
	dpu_input_arguments_byte_stream.Merge(ts_size_word.ToByteStream())

	query_length_word := new(word.Word)
	query_length_word.Init(32)
	query_length_word.SetValue(this.query_length)
	dpu_input_arguments_byte_stream.Merge(query_length_word.ToByteStream())

	query_mean_word := new(word.Word)
	query_mean_word.Init(32)
	query_mean_word.SetValue(this.query_mean)
	dpu_input_arguments_byte_stream.Merge(query_mean_word.ToByteStream())

	query_std_word := new(word.Word)
	query_std_word.Init(32)
	query_std_word.SetValue(this.query_std)
	dpu_input_arguments_byte_stream.Merge(query_std_word.ToByteStream())

	slice_per_dpu_word := new(word.Word)
	slice_per_dpu_word.Init(32)
	slice_per_dpu_word.SetValue(this.slice_per_dpu)
	dpu_input_arguments_byte_stream.Merge(slice_per_dpu_word.ToByteStream())

	exclusion_zone_word := new(word.Word)
	exclusion_zone_word.Init(32)
	exclusion_zone_word.SetValue(this.exclusion_zone)
	dpu_input_arguments_byte_stream.Merge(exclusion_zone_word.ToByteStream())

	kernel_word := new(word.Word)
	kernel_word.Init(32)
	kernel_word.SetValue(this.kernel)
	dpu_input_arguments_byte_stream.Merge(kernel_word.ToByteStream())

	dpu_host := make(map[string]*encoding.ByteStream, 0)
	dpu_host["DPU_INPUT_ARGUMENTS"] = dpu_input_arguments_byte_stream

	return dpu_host
}

func (this *Ts) OutputDpuHost(execution int, dpu_id int) map[string]*encoding.ByteStream {
	if execution >= this.num_executions {
		err := errors.New("execution >= num executions")
		panic(err)
	} else if dpu_id >= this.num_dpus {
		err := errors.New("DPU ID >= num DPUs")
		panic(err)
	}

	dpu_results_byte_stream := new(encoding.ByteStream)
	dpu_results_byte_stream.Init()

	for i := 0; i < this.num_tasklets; i++ {
		min_val_word := new(word.Word)
		min_val_word.Init(32)
		min_val_word.SetValue(this.min_vals[dpu_id][i])
		dpu_results_byte_stream.Merge(min_val_word.ToByteStream())

		min_idx_word := new(word.Word)
		min_idx_word.Init(32)
		min_idx_word.SetValue(this.min_idxs[dpu_id][i])
		dpu_results_byte_stream.Merge(min_idx_word.ToByteStream())

		max_val_word := new(word.Word)
		max_val_word.Init(32)
		max_val_word.SetValue(this.max_vals[dpu_id][i])
		dpu_results_byte_stream.Merge(max_val_word.ToByteStream())

		max_idx_word := new(word.Word)
		max_idx_word.Init(32)
		max_idx_word.SetValue(this.max_idxs[dpu_id][i])
		dpu_results_byte_stream.Merge(max_idx_word.ToByteStream())
	}

	dpu_results := make(map[string]*encoding.ByteStream, 0)
	dpu_results["DPU_RESULTS"] = dpu_results_byte_stream

	return dpu_results
}

func (this *Ts) InputDpuMramHeapPointerName(
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

	start_elem := this.slice_per_dpu * int64(dpu_id)
	end_elem := this.slice_per_dpu*int64(dpu_id+1) + this.query_length

	for _, element := range this.query_buffer {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	for _, element := range this.t_series_buffer[start_elem:end_elem] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	for _, element := range this.amean_buffer[start_elem:end_elem] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	for _, element := range this.asigma_buffer[start_elem:end_elem] {
		element_word := new(word.Word)
		element_word.Init(32)
		element_word.SetValue(element)
		byte_stream.Merge(element_word.ToByteStream())
	}

	return 0, byte_stream
}

func (this *Ts) OutputDpuMramHeapPointerName(
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

func (this *Ts) NumExecutions() int {
	return this.num_executions
}

func (this *Ts) DotProduct(a []int64, a_aux []int64, query []int64, result []int64) []int64 {
	for i := int64(0); i < this.block_size/this.elem_size; i++ {
		for j := int64(0); j < this.dotpip; j++ {
			if (j + i) > (this.block_size/this.elem_size)-1 {
				result[j] += a_aux[(j+i)-this.block_size/this.elem_size] * query[i]
			} else {
				result[j] += a[j+i] * query[i]
			}
		}
	}

	return result
}
