package misc

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type CommandLineParser struct {
	command_line_options map[string]*CommandLineOption
	args                 map[string]bool
}

func (this *CommandLineParser) Init() {
	this.command_line_options = make(map[string]*CommandLineOption, 0)
	this.args = make(map[string]bool, 0)
}

func (this *CommandLineParser) AddOption(
	command_line_option_type CommandLineOptionType,
	option string,
	default_parameter string,
	help_msg string,
) {
	if _, found := this.command_line_options[option]; found {
		err_msg := fmt.Sprintf("option (%s) is already added to the parser")
		err := errors.New(err_msg)
		panic(err)
	}

	command_line_option := new(CommandLineOption)
	command_line_option.Init(command_line_option_type, option, default_parameter, help_msg)
	this.command_line_options[option] = command_line_option
}

func (this *CommandLineParser) Parse(os_args []string) {
	for i := 1; i < len(os_args); i++ {
		os_arg := os_args[i]

		if os_arg[0:2] == "--" {
			option := os_arg[2:]
			custom_parameter := os_args[i+1]

			this.command_line_options[option].SetCustomParameter(custom_parameter)

			i++
		} else if os_arg[0:1] == "-" {
			arg := os_arg[1:]

			if _, found := this.args[arg]; found {
				err_msg := fmt.Sprintf("arg (%s) is already set", arg)
				err := errors.New(err_msg)
				panic(err)
			}

			this.args[arg] = true
		} else {
			err := errors.New("command line options are corrupted")
			panic(err)
		}
	}
}

func (this *CommandLineParser) BoolParameter(option string) bool {
	if _, found := this.command_line_options[option]; !found {
		err_msg := fmt.Sprintf("option (%s) is not found", option)
		err := errors.New(err_msg)
		panic(err)
	}

	command_line_option := this.command_line_options[option]
	return command_line_option.BoolParameter()
}

func (this *CommandLineParser) IntParameter(option string) int64 {
	if _, found := this.command_line_options[option]; !found {
		err_msg := fmt.Sprintf("option (%s) is not found", option)
		err := errors.New(err_msg)
		panic(err)
	}

	command_line_option := this.command_line_options[option]
	return command_line_option.IntParameter()
}

func (this *CommandLineParser) StringParameter(option string) string {
	if _, found := this.command_line_options[option]; !found {
		err_msg := fmt.Sprintf("option (%s) is not found", option)
		err := errors.New(err_msg)
		panic(err)
	}

	command_line_option := this.command_line_options[option]
	return command_line_option.StringParameter()
}

func (this *CommandLineParser) DataPrepParams() []int {
	string_params := strings.Split(this.StringParameter("data_prep_params"), ",")

	data_prep_params := make([]int, 0)
	for _, string_param := range string_params {
		int_param, err := strconv.Atoi(string_param)
		if err != nil {
			panic(err)
		}

		data_prep_params = append(data_prep_params, int_param)
	}
	return data_prep_params
}

func (this *CommandLineParser) IsArgSet(arg string) bool {
	if _, found := this.args[arg]; found {
		return true
	} else {
		return false
	}
}

func (this *CommandLineParser) Options() []string {
	options := make([]string, 0)
	for option := range this.command_line_options {
		options = append(options, option)
	}

	slices.Sort(options)
	return options
}

func (this *CommandLineParser) Args() []string {
	args := make([]string, 0)
	for arg := range this.args {
		args = append(args, arg)
	}

	slices.Sort(args)
	return args
}

func (this *CommandLineParser) StringifyOptions() string {
	str := "OPTIONS\n"

	for option, command_line_option := range this.command_line_options {
		str += option + "  -->  " + command_line_option.Parameter() + "\n"
	}

	return str
}

func (this *CommandLineParser) StringifyArgs() string {
	str := "ARGS\n"

	for arg := range this.args {
		str += arg + "\n"
	}

	return str
}

func (this *CommandLineParser) StringifyHelpMsgs() string {
	str := "HELP_MSGS\n"

	for option, command_line_option := range this.command_line_options {
		str += option + "  -->  " + command_line_option.HelpMsg() + "\n"
	}

	return str
}
