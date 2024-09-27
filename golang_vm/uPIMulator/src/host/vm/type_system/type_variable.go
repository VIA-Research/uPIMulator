package type_system

import (
	"errors"
)

type TypeVariableType int

const (
	VOID TypeVariableType = iota
	CHAR
	SHORT
	INT
	LONG
	STRING
	STRUCT
)

type TypeVariable struct {
	type_variable_type TypeVariableType
	struct_name        *string
	num_stars          int
}

func (this *TypeVariable) InitPrimitive(type_variable_type TypeVariableType, num_stars int) {
	if type_variable_type == STRUCT {
		err := errors.New("type variable type is struct")
		panic(err)
	}

	this.type_variable_type = type_variable_type
	this.struct_name = nil
	this.num_stars = num_stars
}

func (this *TypeVariable) InitStruct(
	type_variable_type TypeVariableType,
	struct_name string,
	num_stars int,
) {
	if type_variable_type != STRUCT {
		err := errors.New("type variable type isn't struct")
		panic(err)
	}

	this.type_variable_type = type_variable_type

	this.struct_name = new(string)
	*this.struct_name = struct_name

	this.num_stars = num_stars
}

func (this *TypeVariable) TypeVariableType() TypeVariableType {
	return this.type_variable_type
}

func (this *TypeVariable) StructName() string {
	if this.struct_name == nil {
		err := errors.New("struct name == nil")
		panic(err)
	}

	return *this.struct_name
}

func (this *TypeVariable) NumStars() int {
	return this.num_stars
}
