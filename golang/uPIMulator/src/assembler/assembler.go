package assembler

import (
	"errors"
	"fmt"
	"path/filepath"
	"uPIMulator/src/assembler/prim"
	"uPIMulator/src/misc"
)

type Assembler struct {
	bin_dirpath string

	benchmark string

	num_channels          int
	num_ranks_per_channel int
	num_dpus_per_rank     int
	num_dpus              int

	num_tasklets int

	assemblables map[string]Assemblable
}

func (this *Assembler) Init(command_line_parser *misc.CommandLineParser) {
	this.bin_dirpath = command_line_parser.StringParameter("bin_dirpath")

	this.benchmark = command_line_parser.StringParameter("benchmark")

	this.num_channels = int(command_line_parser.IntParameter("num_channels"))
	this.num_ranks_per_channel = int(command_line_parser.IntParameter("num_ranks_per_channel"))
	this.num_dpus_per_rank = int(command_line_parser.IntParameter("num_dpus_per_rank"))
	this.num_dpus = this.num_channels * this.num_ranks_per_channel * this.num_dpus_per_rank

	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.assemblables = make(map[string]Assemblable, 0)

	this.assemblables["BS"] = new(prim.Bs)
	this.assemblables["GEMV"] = new(prim.Gemv)
	this.assemblables["HST-L"] = new(prim.HstL)
	this.assemblables["HST-S"] = new(prim.HstS)
	this.assemblables["MLP"] = new(prim.Mlp)
	this.assemblables["RED"] = new(prim.Red)
	this.assemblables["SCAN-RSS"] = new(prim.ScanRss)
	this.assemblables["SCAN-SSA"] = new(prim.ScanSsa)
	this.assemblables["SEL"] = new(prim.Sel)
	this.assemblables["TRNS"] = new(prim.Trns)
	this.assemblables["TS"] = new(prim.Ts)
	this.assemblables["UNI"] = new(prim.Uni)
	this.assemblables["VA"] = new(prim.Va)

	if assemblable, found := this.assemblables[this.benchmark]; found {
		assemblable.Init(command_line_parser)
	} else {
		err := errors.New("assemblable is not found")
		panic(err)
	}
}

func (this *Assembler) Assemble() {
	this.AssembleInputDpuHost()
	this.AssembleOutputDpuHost()
	this.AssembleInputDpuMramHeapPointerName()
	this.AssembleOutputDpuMramHeapPointerName()
	this.AssembleNumExecutions()
}

func (this *Assembler) AssembleInputDpuHost() {
	assemblable := this.assemblables[this.benchmark]

	for execution := 0; execution < assemblable.NumExecutions(); execution++ {
		for dpu_id := 0; dpu_id < this.num_dpus; dpu_id++ {
			input_dpu_host := assemblable.InputDpuHost(execution, dpu_id)

			for name, byte_stream := range input_dpu_host {
				filename := fmt.Sprintf("input_%s_%d_%d.bin", name, execution, dpu_id)
				filepath_ := filepath.Join(this.bin_dirpath, filename)

				file_dumper := new(misc.FileDumper)
				file_dumper.Init(filepath_)

				lines := make([]string, 0)
				for i := int64(0); i < byte_stream.Size(); i++ {
					line := fmt.Sprintf("%d", byte_stream.Get(int(i)))

					lines = append(lines, line)
				}

				file_dumper.WriteLines(lines)
			}
		}
	}
}

func (this *Assembler) AssembleOutputDpuHost() {
	assemblable := this.assemblables[this.benchmark]

	for execution := 0; execution < assemblable.NumExecutions(); execution++ {
		for dpu_id := 0; dpu_id < this.num_dpus; dpu_id++ {
			output_dpu_host := assemblable.OutputDpuHost(execution, dpu_id)

			for name, byte_stream := range output_dpu_host {
				filename := fmt.Sprintf("output_%s_%d_%d.bin", name, execution, dpu_id)
				filepath_ := filepath.Join(this.bin_dirpath, filename)

				file_dumper := new(misc.FileDumper)
				file_dumper.Init(filepath_)

				lines := make([]string, 0)
				for i := int64(0); i < byte_stream.Size(); i++ {
					line := fmt.Sprintf("%d", byte_stream.Get(int(i)))

					lines = append(lines, line)
				}

				file_dumper.WriteLines(lines)
			}
		}
	}
}

func (this *Assembler) AssembleInputDpuMramHeapPointerName() {
	assemblable := this.assemblables[this.benchmark]

	for execution := 0; execution < assemblable.NumExecutions(); execution++ {
		for dpu_id := 0; dpu_id < this.num_dpus; dpu_id++ {
			offset, byte_stream := assemblable.InputDpuMramHeapPointerName(execution, dpu_id)

			filename := fmt.Sprintf(
				"input_dpu_mram_heap_pointer_name_%d_%d_%d.bin",
				offset,
				execution,
				dpu_id,
			)
			filepath_ := filepath.Join(this.bin_dirpath, filename)

			file_dumper := new(misc.FileDumper)
			file_dumper.Init(filepath_)

			lines := make([]string, 0)
			for i := int64(0); i < byte_stream.Size(); i++ {
				line := fmt.Sprintf("%d", byte_stream.Get(int(i)))

				lines = append(lines, line)
			}

			file_dumper.WriteLines(lines)
		}
	}
}

func (this *Assembler) AssembleOutputDpuMramHeapPointerName() {
	assemblable := this.assemblables[this.benchmark]

	for execution := 0; execution < assemblable.NumExecutions(); execution++ {
		for dpu_id := 0; dpu_id < this.num_dpus; dpu_id++ {
			offset, byte_stream := assemblable.OutputDpuMramHeapPointerName(execution, dpu_id)

			filename := fmt.Sprintf(
				"output_dpu_mram_heap_pointer_name_%d_%d_%d.bin",
				offset,
				execution,
				dpu_id,
			)
			filepath_ := filepath.Join(this.bin_dirpath, filename)

			file_dumper := new(misc.FileDumper)
			file_dumper.Init(filepath_)

			lines := make([]string, 0)
			for i := int64(0); i < byte_stream.Size(); i++ {
				line := fmt.Sprintf("%d", byte_stream.Get(int(i)))

				lines = append(lines, line)
			}

			file_dumper.WriteLines(lines)
		}
	}
}

func (this *Assembler) AssembleNumExecutions() {
	assemblable := this.assemblables[this.benchmark]

	path := filepath.Join(this.bin_dirpath, "num_executions.txt")

	file_dumper := new(misc.FileDumper)
	file_dumper.Init(path)

	line := fmt.Sprintf("%d", assemblable.NumExecutions())
	lines := make([]string, 0)
	lines = append(lines, line)

	file_dumper.WriteLines(lines)
}
