package abi

import (
	"errors"
	"fmt"
)

type Bytecode struct {
	op_code OpCode
	arg1    *int64
	arg2    *int64
	str1    *string
	str2    *string
}

func (this *Bytecode) Init(op_code OpCode, args []int64, strs []string) {
	this.op_code = op_code

	if len(args) == 0 {
		this.arg1 = nil
		this.arg2 = nil
	} else if len(args) == 1 {
		this.arg1 = new(int64)
		*this.arg1 = args[0]

		this.arg2 = nil
	} else if len(args) == 2 {
		this.arg1 = new(int64)
		*this.arg1 = args[0]

		this.arg2 = new(int64)
		*this.arg2 = args[1]
	} else {
		err := errors.New("len(args) > 2")
		panic(err)
	}

	if len(strs) == 0 {
		this.str1 = nil
		this.str2 = nil
	} else if len(strs) == 1 {
		this.str1 = new(string)
		*this.str1 = strs[0]

		this.str2 = nil
	} else if len(strs) == 2 {
		this.str1 = new(string)
		*this.str1 = strs[0]

		this.str2 = new(string)
		*this.str2 = strs[1]
	} else {
		err := errors.New("len(strs) > 2")
		panic(err)
	}
}

func (this *Bytecode) OpCode() OpCode {
	return this.op_code
}

func (this *Bytecode) Arg1() int64 {
	if this.arg1 == nil {
		err := errors.New("arg1 == nil")
		panic(err)
	}

	return *this.arg1
}

func (this *Bytecode) Arg2() int64 {
	if this.arg2 == nil {
		err := errors.New("arg2 == nil")
		panic(err)
	}

	return *this.arg2
}

func (this *Bytecode) Str1() string {
	if this.str1 == nil {
		err := errors.New("str1 == nil")
		panic(err)
	}

	return *this.str1
}

func (this *Bytecode) Str2() string {
	if this.str2 == nil {
		err := errors.New("str2 == nil")
		panic(err)
	}

	return *this.str2
}

func (this *Bytecode) Stringify() string {
	if this.op_code == NEW_SCOPE {
		return "NEW_SCOPE"
	} else if this.op_code == DELETE_SCOPE {
		return "DELETE_SCOPE"
	} else if this.op_code == PUSH_CHAR {
		return fmt.Sprintf("PUSH_CHAR %d", *this.arg1)
	} else if this.op_code == PUSH_SHORT {
		return fmt.Sprintf("PUSH_SHORT %d", *this.arg1)
	} else if this.op_code == PUSH_INT {
		return fmt.Sprintf("PUSH_INT %d", *this.arg1)
	} else if this.op_code == PUSH_LONG {
		return fmt.Sprintf("PUSH_LONG %d", *this.arg1)
	} else if this.op_code == PUSH_STRING {
		return fmt.Sprintf("PUSH_STRING %s", *this.str1)
	} else if this.op_code == POP {
		return "POP"
	} else if this.op_code == BEGIN_STRUCT {
		return fmt.Sprintf("BEGIN_STRUCT %s", *this.str1)
	} else if this.op_code == APPEND_VOID {
		return fmt.Sprintf("APPEND_VOID %d %s", *this.arg1, *this.str1)
	} else if this.op_code == APPEND_CHAR {
		return fmt.Sprintf("APPEND_VOID %d %s", *this.arg1, *this.str1)
	} else if this.op_code == APPEND_SHORT {
		return fmt.Sprintf("APPEND_SHORT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == APPEND_INT {
		return fmt.Sprintf("APPEND_INT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == APPEND_LONG {
		return fmt.Sprintf("APPEND_LONG %d %s", *this.arg1, *this.str1)
	} else if this.op_code == APPEND_STRUCT {
		return fmt.Sprintf("APPEND_STRUCT %d %s %s", *this.arg1, *this.str1, *this.str2)
	} else if this.op_code == END_STRUCT {
		return "END_STRUCT"
	} else if this.op_code == NEW_GLOBAL_VOID {
		return fmt.Sprintf("NEW_GLOBAL_VOID %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_GLOBAL_CHAR {
		return fmt.Sprintf("NEW_GLOBAL_CHAR %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_GLOBAL_SHORT {
		return fmt.Sprintf("NEW_GLOBAL_SHORT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_GLOBAL_INT {
		return fmt.Sprintf("NEW_GLOBAL_INT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_GLOBAL_LONG {
		return fmt.Sprintf("NEW_GLOBAL_LONG %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_FAST_VOID {
		return fmt.Sprintf("NEW_FAST_VOID %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_FAST_CHAR {
		return fmt.Sprintf("NEW_FAST_CHAR %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_FAST_SHORT {
		return fmt.Sprintf("NEW_FAST_SHORT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_FAST_INT {
		return fmt.Sprintf("NEW_FAST_INT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_FAST_LONG {
		return fmt.Sprintf("NEW_FAST_LONG %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_FAST_STRUCT {
		return fmt.Sprintf("NEW_FAST_STRUCT %d %s %s", *this.arg1, *this.str1, *this.str2)
	} else if this.op_code == NEW_ARG_VOID {
		return fmt.Sprintf("NEW_ARG_VOID %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_ARG_CHAR {
		return fmt.Sprintf("NEW_ARG_CHAR %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_ARG_SHORT {
		return fmt.Sprintf("NEW_ARG_SHORT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_ARG_INT {
		return fmt.Sprintf("NEW_ARG_INT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_ARG_LONG {
		return fmt.Sprintf("NEW_ARG_LONG %d %s", *this.arg1, *this.str1)
	} else if this.op_code == NEW_ARG_STRUCT {
		return fmt.Sprintf("NEW_ARG_STRUCT %d %s %s", *this.arg1, *this.str1, *this.str2)
	} else if this.op_code == NEW_RETURN_VOID {
		return fmt.Sprintf("NEW_RETURN_VOID %d", *this.arg1)
	} else if this.op_code == NEW_RETURN_CHAR {
		return fmt.Sprintf("NEW_RETURN_CHAR %d", *this.arg1)
	} else if this.op_code == NEW_RETURN_SHORT {
		return fmt.Sprintf("NEW_RETURN_SHORT %d", *this.arg1)
	} else if this.op_code == NEW_RETURN_INT {
		return fmt.Sprintf("NEW_RETURN_INT %d", *this.arg1)
	} else if this.op_code == NEW_RETURN_LONG {
		return fmt.Sprintf("NEW_RETURN_LONG %d", *this.arg1)
	} else if this.op_code == NEW_RETURN_STRUCT {
		return fmt.Sprintf("NEW_RETURN_STRUCT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == SIZEOF_VOID {
		return fmt.Sprintf("SIZEOF_VOID %d", *this.arg1)
	} else if this.op_code == SIZEOF_CHAR {
		return fmt.Sprintf("SIZEOF_CHAR %d", *this.arg1)
	} else if this.op_code == SIZEOF_SHORT {
		return fmt.Sprintf("SIZEOF_SHORT %d", *this.arg1)
	} else if this.op_code == SIZEOF_INT {
		return fmt.Sprintf("SIZEOF_INT %d", *this.arg1)
	} else if this.op_code == SIZEOF_LONG {
		return fmt.Sprintf("SIZEOF_LONG %d", *this.arg1)
	} else if this.op_code == SIZEOF_STRUCT {
		return fmt.Sprintf("SIZEOF_STRUCT %d %s", *this.arg1, *this.str1)
	} else if this.op_code == GET_IDENTIFIER {
		return fmt.Sprintf("GET_IDENTIFIER %s", *this.str1)
	} else if this.op_code == GET_ARG_IDENTIFIER {
		return fmt.Sprintf("GET_ARG_IDENTIFIER %s", *this.str1)
	} else if this.op_code == GET_SUBSCRIPT {
		return "GET_SUBSCRIPT"
	} else if this.op_code == GET_ACCESS {
		return fmt.Sprintf("GET_ACCESS %s", *this.str1)
	} else if this.op_code == GET_REFERENCE {
		return fmt.Sprintf("GET_REFERENCE %s", *this.str1)
	} else if this.op_code == GET_ADDRESS {
		return "GET_ADDRESS"
	} else if this.op_code == GET_VALUE {
		return "GET_VALUE"
	} else if this.op_code == ALLOC {
		return "ALLOC"
	} else if this.op_code == FREE {
		return "FREE"
	} else if this.op_code == ASSERT {
		return "ASSERT"
	} else if this.op_code == ADD {
		return "ADD"
	} else if this.op_code == SUB {
		return "SUB"
	} else if this.op_code == MUL {
		return "MUL"
	} else if this.op_code == DIV {
		return "DIV"
	} else if this.op_code == MOD {
		return "MOD"
	} else if this.op_code == LSHIFT {
		return "LSHIFT"
	} else if this.op_code == RSHIFT {
		return "RSHIFT"
	} else if this.op_code == NEGATE {
		return "NEGATE"
	} else if this.op_code == TILDE {
		return "TILDE"
	} else if this.op_code == SQRT {
		return "SQRT"
	} else if this.op_code == BITWISE_AND {
		return "BITWISE_AND"
	} else if this.op_code == BITWISE_XOR {
		return "BITWISE_XOR"
	} else if this.op_code == BITWISE_OR {
		return "BITWISE_OR"
	} else if this.op_code == LOGICAL_AND {
		return "LOGICAL_AND"
	} else if this.op_code == LOGICAL_OR {
		return "LOGICAL_OR"
	} else if this.op_code == LOGICAL_NOT {
		return "LOGICAL_NOT"
	} else if this.op_code == EQ {
		return "EQ"
	} else if this.op_code == NOT_EQ {
		return "NOT_EQ"
	} else if this.op_code == LESS {
		return "LESS"
	} else if this.op_code == LESS_EQ {
		return "LESS_EQ"
	} else if this.op_code == GREATER {
		return "GREATER"
	} else if this.op_code == GREATER_EQ {
		return "GREATER_EQ"
	} else if this.op_code == CONDITIONAL {
		return "CONDITIONAL"
	} else if this.op_code == ASSIGN {
		return "ASSIGN"
	} else if this.op_code == ASSIGN_STAR {
		return "ASSIGN_STAR"
	} else if this.op_code == ASSIGN_DIV {
		return "ASSIGN_DIV"
	} else if this.op_code == ASSIGN_MOD {
		return "ASSIGN_MOD"
	} else if this.op_code == ASSIGN_ADD {
		return "ASSIGN_ADD"
	} else if this.op_code == ASSIGN_SUB {
		return "ASSIGN_SUB"
	} else if this.op_code == ASSIGN_LSHIFT {
		return "ASSIGN_LSHIFT"
	} else if this.op_code == ASSIGN_RSHIFT {
		return "ASSIGN_RSHIFT"
	} else if this.op_code == ASSIGN_BITWISE_AND {
		return "ASSIGN_BITWISE_AND"
	} else if this.op_code == ASSIGN_BITWISE_XOR {
		return "ASSIGN_BITWISE_XOR"
	} else if this.op_code == ASSIGN_BITWISE_OR {
		return "ASSIGN_BITWISE_OR"
	} else if this.op_code == ASSIGN_PLUS_PLUS {
		return "ASSIGN_PLUS_PLUS"
	} else if this.op_code == ASSIGN_MINUS_MINUS {
		return "ASSIGN_MINUS_MINUS"
	} else if this.op_code == ASSIGN_RETURN {
		return "ASSIGN_RETURN"
	} else if this.op_code == JUMP {
		return fmt.Sprintf("JUMP %s", *this.str1)
	} else if this.op_code == JUMP_IF_ZERO {
		return fmt.Sprintf("JUMP_IF_ZERO %s", *this.str1)
	} else if this.op_code == JUMP_IF_NONZERO {
		return fmt.Sprintf("JUMP_IF_NONZERO %s", *this.str1)
	} else if this.op_code == CALL {
		return fmt.Sprintf("CALL %s", *this.str1)
	} else if this.op_code == RETURN {
		return "RETURN"
	} else if this.op_code == NOP {
		return "NOP"
	} else if this.op_code == DPU_ALLOC {
		return fmt.Sprintf("DPU_ALLOC %d", *this.arg1)
	} else if this.op_code == DPU_LOAD {
		return fmt.Sprintf("DPU_LOAD %d %s", *this.arg1, *this.str1)
	} else if this.op_code == DPU_PREPARE {
		return "DPU_PREPARE"
	} else if this.op_code == DPU_TRANSFER {
		return "DPU_TRANSFER"
	} else if this.op_code == DPU_COPY_TO {
		return "DPU_COPY_TO"
	} else if this.op_code == DPU_COPY_FROM {
		return "DPU_COPY_FROM"
	} else if this.op_code == DPU_LAUNCH {
		return "DPU_LAUNCH"
	} else if this.op_code == DPU_FREE {
		return "DPU_FREE"
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
}
