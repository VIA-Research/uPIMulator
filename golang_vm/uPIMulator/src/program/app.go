package program

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"uPIMulator/src/host/abi"
	"uPIMulator/src/misc"
)

type App struct {
	bin_dirpath string

	benchmark    string
	num_dpus     int
	num_tasklets int

	labels []*abi.Label
}

func (this *App) Init(command_line_parser *misc.CommandLineParser) {
	this.bin_dirpath = command_line_parser.StringParameter("bin_dirpath")

	this.benchmark = command_line_parser.StringParameter("benchmark")

	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))
	this.num_tasklets = num_channels * num_ranks_per_channel * num_dpus_per_rank

	this.labels = make([]*abi.Label, 0)
	this.LoadLabels()
}

func (this *App) HasLabel(label_name string) bool {
	for _, label := range this.labels {
		if label.Name() == label_name {
			return true
		}
	}
	return false
}

func (this *App) Label(label_name string) *abi.Label {
	if !this.HasLabel(label_name) {
		err_msg := fmt.Sprintf("label (%s) is not found", label_name)
		err := errors.New(err_msg)
		panic(err)
	}

	for _, label := range this.labels {
		if label.Name() == label_name {
			return label
		}
	}

	return nil
}

func (this *App) LoadLabels() {
	path := filepath.Join(this.bin_dirpath, "bytecode.txt")

	file_scanner := new(misc.FileScanner)
	file_scanner.Init(path)

	for _, line := range file_scanner.ReadLines() {
		if strings.Contains(line, ":") {
			label_name := line[:len(line)-1]

			label := new(abi.Label)
			label.Init(label_name)

			this.labels = append(this.labels, label)
		} else {
			words := strings.Split(line, " ")

			var op_code abi.OpCode
			args := make([]int64, 0)
			strs := make([]string, 0)
			for i, word := range words {
				if i == 0 {
					op_code = this.ConvertToOpCode(word[1:])
				} else {
					if arg, err := strconv.ParseInt(word, 10, 64); err == nil {
						args = append(args, arg)
					} else {
						strs = append(strs, word)
					}
				}
			}

			bytecode := new(abi.Bytecode)
			bytecode.Init(op_code, args, strs)

			this.labels[len(this.labels)-1].Append(bytecode)
		}
	}
}

func (this *App) ConvertToOpCode(op_code string) abi.OpCode {
	if op_code == "NEW_SCOPE" {
		return abi.NEW_SCOPE
	} else if op_code == "DELETE_SCOPE" {
		return abi.DELETE_SCOPE
	} else if op_code == "PUSH_CHAR" {
		return abi.PUSH_CHAR
	} else if op_code == "PUSH_SHORT" {
		return abi.PUSH_SHORT
	} else if op_code == "PUSH_INT" {
		return abi.PUSH_INT
	} else if op_code == "PUSH_LONG" {
		return abi.PUSH_LONG
	} else if op_code == "PUSH_STRING" {
		return abi.PUSH_STRING
	} else if op_code == "POP" {
		return abi.POP
	} else if op_code == "BEGIN_STRUCT" {
		return abi.BEGIN_STRUCT
	} else if op_code == "APPEND_VOID" {
		return abi.APPEND_VOID
	} else if op_code == "APPEND_CHAR" {
		return abi.APPEND_CHAR
	} else if op_code == "APPEND_SHORT" {
		return abi.APPEND_SHORT
	} else if op_code == "APPEND_INT" {
		return abi.APPEND_INT
	} else if op_code == "APPEND_LONG" {
		return abi.APPEND_LONG
	} else if op_code == "APPEND_STRUCT" {
		return abi.APPEND_STRUCT
	} else if op_code == "END_STRUCT" {
		return abi.END_STRUCT
	} else if op_code == "NEW_GLOBAL_VOID" {
		return abi.NEW_GLOBAL_VOID
	} else if op_code == "NEW_GLOBAL_CHAR" {
		return abi.NEW_GLOBAL_CHAR
	} else if op_code == "NEW_GLOBAL_SHORT" {
		return abi.NEW_GLOBAL_SHORT
	} else if op_code == "NEW_GLOBAL_INT" {
		return abi.NEW_GLOBAL_INT
	} else if op_code == "NEW_GLOBAL_LONG" {
		return abi.NEW_GLOBAL_LONG
	} else if op_code == "NEW_FAST_VOID" {
		return abi.NEW_FAST_VOID
	} else if op_code == "NEW_FAST_CHAR" {
		return abi.NEW_FAST_CHAR
	} else if op_code == "NEW_FAST_SHORT" {
		return abi.NEW_FAST_SHORT
	} else if op_code == "NEW_FAST_INT" {
		return abi.NEW_FAST_INT
	} else if op_code == "NEW_FAST_LONG" {
		return abi.NEW_FAST_LONG
	} else if op_code == "NEW_FAST_STRUCT" {
		return abi.NEW_FAST_STRUCT
	} else if op_code == "NEW_ARG_VOID" {
		return abi.NEW_ARG_VOID
	} else if op_code == "NEW_ARG_CHAR" {
		return abi.NEW_ARG_CHAR
	} else if op_code == "NEW_ARG_SHORT" {
		return abi.NEW_ARG_SHORT
	} else if op_code == "NEW_ARG_INT" {
		return abi.NEW_ARG_INT
	} else if op_code == "NEW_ARG_LONG" {
		return abi.NEW_ARG_LONG
	} else if op_code == "NEW_ARG_STRUCT" {
		return abi.NEW_ARG_STRUCT
	} else if op_code == "NEW_RETURN_VOID" {
		return abi.NEW_RETURN_VOID
	} else if op_code == "NEW_RETURN_CHAR" {
		return abi.NEW_RETURN_CHAR
	} else if op_code == "NEW_RETURN_SHORT" {
		return abi.NEW_RETURN_SHORT
	} else if op_code == "NEW_RETURN_INT" {
		return abi.NEW_RETURN_INT
	} else if op_code == "NEW_RETURN_LONG" {
		return abi.NEW_RETURN_LONG
	} else if op_code == "NEW_RETURN_STRUCT" {
		return abi.NEW_RETURN_STRUCT
	} else if op_code == "SIZEOF_VOID" {
		return abi.SIZEOF_VOID
	} else if op_code == "SIZEOF_CHAR" {
		return abi.SIZEOF_CHAR
	} else if op_code == "SIZEOF_SHORT" {
		return abi.SIZEOF_SHORT
	} else if op_code == "SIZEOF_INT" {
		return abi.SIZEOF_INT
	} else if op_code == "SIZEOF_LONG" {
		return abi.SIZEOF_LONG
	} else if op_code == "SIZEOF_STRUCT" {
		return abi.SIZEOF_STRUCT
	} else if op_code == "GET_IDENTIFIER" {
		return abi.GET_IDENTIFIER
	} else if op_code == "GET_ARG_IDENTIFIER" {
		return abi.GET_ARG_IDENTIFIER
	} else if op_code == "GET_SUBSCRIPT" {
		return abi.GET_SUBSCRIPT
	} else if op_code == "GET_ACCESS" {
		return abi.GET_ACCESS
	} else if op_code == "GET_REFERENCE" {
		return abi.GET_REFERENCE
	} else if op_code == "GET_ADDRESS" {
		return abi.GET_ADDRESS
	} else if op_code == "GET_VALUE" {
		return abi.GET_VALUE
	} else if op_code == "ALLOC" {
		return abi.ALLOC
	} else if op_code == "FREE" {
		return abi.FREE
	} else if op_code == "ASSERT" {
		return abi.ASSERT
	} else if op_code == "ADD" {
		return abi.ADD
	} else if op_code == "SUB" {
		return abi.SUB
	} else if op_code == "MUL" {
		return abi.MUL
	} else if op_code == "DIV" {
		return abi.DIV
	} else if op_code == "MOD" {
		return abi.MOD
	} else if op_code == "LSHIFT" {
		return abi.LSHIFT
	} else if op_code == "RSHIFT" {
		return abi.RSHIFT
	} else if op_code == "NEGATE" {
		return abi.NEGATE
	} else if op_code == "TILDE" {
		return abi.TILDE
	} else if op_code == "SQRT" {
		return abi.SQRT
	} else if op_code == "BITWISE_AND" {
		return abi.BITWISE_AND
	} else if op_code == "BITWISE_XOR" {
		return abi.BITWISE_XOR
	} else if op_code == "BITWISE_OR" {
		return abi.BITWISE_OR
	} else if op_code == "LOGICAL_AND" {
		return abi.LOGICAL_AND
	} else if op_code == "LOGICAL_OR" {
		return abi.LOGICAL_OR
	} else if op_code == "LOGICAL_NOT" {
		return abi.LOGICAL_NOT
	} else if op_code == "EQ" {
		return abi.EQ
	} else if op_code == "NOT_EQ" {
		return abi.NOT_EQ
	} else if op_code == "LESS" {
		return abi.LESS
	} else if op_code == "LESS_EQ" {
		return abi.LESS_EQ
	} else if op_code == "GREATER" {
		return abi.GREATER
	} else if op_code == "GREATER_EQ" {
		return abi.GREATER_EQ
	} else if op_code == "CONDITIONAL" {
		return abi.CONDITIONAL
	} else if op_code == "ASSIGN" {
		return abi.ASSIGN
	} else if op_code == "ASSIGN_STAR" {
		return abi.ASSIGN_STAR
	} else if op_code == "ASSIGN_DIV" {
		return abi.ASSIGN_DIV
	} else if op_code == "ASSIGN_MOD" {
		return abi.ASSIGN_MOD
	} else if op_code == "ASSIGN_ADD" {
		return abi.ASSIGN_ADD
	} else if op_code == "ASSIGN_SUB" {
		return abi.ASSIGN_SUB
	} else if op_code == "ASSIGN_LSHIFT" {
		return abi.ASSIGN_LSHIFT
	} else if op_code == "ASSIGN_RSHIFT" {
		return abi.ASSIGN_RSHIFT
	} else if op_code == "ASSIGN_BITWISE_AND" {
		return abi.ASSIGN_BITWISE_AND
	} else if op_code == "ASSIGN_BITWISE_XOR" {
		return abi.ASSIGN_BITWISE_XOR
	} else if op_code == "ASSIGN_BITWISE_OR" {
		return abi.ASSIGN_BITWISE_OR
	} else if op_code == "ASSIGN_PLUS_PLUS" {
		return abi.ASSIGN_PLUS_PLUS
	} else if op_code == "ASSIGN_MINUS_MINUS" {
		return abi.ASSIGN_MINUS_MINUS
	} else if op_code == "ASSIGN_RETURN" {
		return abi.ASSIGN_RETURN
	} else if op_code == "JUMP" {
		return abi.JUMP
	} else if op_code == "JUMP_IF_ZERO" {
		return abi.JUMP_IF_ZERO
	} else if op_code == "JUMP_IF_NONZERO" {
		return abi.JUMP_IF_NONZERO
	} else if op_code == "CALL" {
		return abi.CALL
	} else if op_code == "RETURN" {
		return abi.RETURN
	} else if op_code == "NOP" {
		return abi.NOP
	} else if op_code == "DPU_ALLOC" {
		return abi.DPU_ALLOC
	} else if op_code == "DPU_LOAD" {
		return abi.DPU_LOAD
	} else if op_code == "DPU_PREPARE" {
		return abi.DPU_PREPARE
	} else if op_code == "DPU_TRANSFER" {
		return abi.DPU_TRANSFER
	} else if op_code == "DPU_COPY_TO" {
		return abi.DPU_COPY_TO
	} else if op_code == "DPU_COPY_FROM" {
		return abi.DPU_COPY_FROM
	} else if op_code == "DPU_LAUNCH" {
		return abi.DPU_LAUNCH
	} else if op_code == "DPU_FREE" {
		return abi.DPU_FREE
	} else {
		err_msg := fmt.Sprintf("op code (%s) is not valid", op_code)
		err := errors.New(err_msg)
		panic(err)
	}
}
