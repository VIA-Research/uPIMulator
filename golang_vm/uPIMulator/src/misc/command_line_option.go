package misc

import (
	"errors"
	"fmt"
	"strconv"
)

type CommandLineOptionType int

const (
	BOOL = iota
	INT
	STRING
)

type CommandLineOption struct {
	command_line_option_type CommandLineOptionType
	option                   string
	default_parameter        string
	custom_parameter         string
	help_msg                 string
}

func (this *CommandLineOption) Init(
	command_line_option_type CommandLineOptionType,
	option string,
	default_parameter string,
	help_msg string,
) {
	this.command_line_option_type = command_line_option_type
	this.option = option
	this.default_parameter = default_parameter
	this.custom_parameter = ""
	this.help_msg = help_msg
}

func (this *CommandLineOption) CommandLineOptionType() CommandLineOptionType {
	return this.command_line_option_type
}

func (this *CommandLineOption) Option() string {
	return this.option
}

func (this *CommandLineOption) Parameter() string {
	if this.custom_parameter == "" {
		return this.default_parameter
	} else {
		return this.custom_parameter
	}
}

func (this *CommandLineOption) HelpMsg() string {
	return this.help_msg
}

func (this *CommandLineOption) SetCustomParameter(custom_parameter string) {
	if this.custom_parameter != "" {
		err_msg := fmt.Sprintf("custom parameter (%s) is already set", custom_parameter)
		err := errors.New(err_msg)
		panic(err)
	}

	this.custom_parameter = custom_parameter
}

func (this *CommandLineOption) BoolParameter() bool {
	if this.Parameter() == "true" {
		return true
	} else if this.Parameter() == "false" {
		return false
	} else {
		err_msg := fmt.Sprintf("parameter (%s) is not true or false", this.Parameter())
		err := errors.New(err_msg)
		panic(err)
	}
}

func (this *CommandLineOption) IntParameter() int64 {
	int_parameter, err := strconv.ParseInt(this.Parameter(), 10, 64)

	if err != nil {
		panic(err)
	}

	return int_parameter
}

func (this *CommandLineOption) StringParameter() string {
	return this.Parameter()
}
