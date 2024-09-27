package misc

import (
	"errors"
	"fmt"
	"os"
)

type CommandLineValidator struct {
	command_line_parser *CommandLineParser
}

func (this *CommandLineValidator) Init(command_line_parser *CommandLineParser) {
	this.command_line_parser = command_line_parser
}

func (this *CommandLineValidator) Validate() {
	if this.command_line_parser.IntParameter("num_channels") <= 0 {
		err := errors.New("num_channels <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_ranks_per_channel") <= 0 {
		err := errors.New("num_ranks <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_dpus_per_rank") <= 0 {
		err := errors.New("num_dpus <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_vm_channels") <= 0 {
		err := errors.New("num_vm_channels <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_vm_ranks_per_channel") <= 0 {
		err := errors.New("num_vm_ranks_per_channel <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_vm_banks_per_rank") <= 0 {
		err := errors.New("num_vm_banks_per_rank <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_tasklets") <= 0 {
		err := errors.New("num_tasklets <= 0")
		panic(err)
	}

	if _, stat_err := os.Stat(this.command_line_parser.StringParameter("root_dirpath")); os.IsNotExist(
		stat_err,
	) {
		err_msg := fmt.Sprintf(
			"root_dirpath (%s) does not exist",
			this.command_line_parser.StringParameter("root_dirpath"),
		)
		err := errors.New(err_msg)
		panic(err)
	}

	if this.command_line_parser.IntParameter("logic_frequency") <= 0 {
		err := errors.New("logic_frequency <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("memory_frequency") <= 0 {
		err := errors.New("memory_frequency <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_pipeline_stages") <= 0 {
		err := errors.New("num_pipeline_stages <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("num_revolver_scheduling_cycles") < 0 {
		err := errors.New("num_revolver_scheduling_cycles < 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("wordline_size") <= 0 {
		err := errors.New("wordline_size <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("min_access_granularity") <= 0 {
		err := errors.New("min_access_granularity <= 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("t_rcd") < 0 {
		err := errors.New("t_rcd < 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("t_ras") < 0 {
		err := errors.New("t_ras < 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("t_rp") < 0 {
		err := errors.New("t_rp < 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("t_cl") < 0 {
		err := errors.New("t_cl < 0")
		panic(err)
	}

	if this.command_line_parser.IntParameter("t_bl") < 0 {
		err := errors.New("t_bl < 0")
		panic(err)
	}
}
