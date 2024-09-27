package type_system

import (
	"errors"
)

type SymbolType int

const (
	VOID SymbolType = iota
	CHAR
	SHORT
	INT
	LONG
	STRING
	STRUCT
)

type Symbol struct {
	symbol_type SymbolType
	struct_name *string
	num_stars   int

	name string
}

func (this *Symbol) InitPrimitive(symbol_type SymbolType, num_stars int, name string) {
	if symbol_type == STRUCT {
		err := errors.New("symbol type is struct")
		panic(err)
	}

	this.symbol_type = symbol_type
	this.struct_name = nil
	this.num_stars = num_stars
	this.name = name
}

func (this *Symbol) InitStruct(
	symbol_type SymbolType,
	struct_name string,
	num_stars int,
	name string,
) {
	if symbol_type != STRUCT {
		err := errors.New("symbol type is not struct")
		panic(err)
	}

	this.symbol_type = symbol_type

	this.struct_name = new(string)
	*this.struct_name = struct_name

	this.num_stars = num_stars
	this.name = name
}

func (this *Symbol) SymbolType() SymbolType {
	return this.symbol_type
}

func (this *Symbol) StructName() string {
	if this.struct_name == nil {
		err := errors.New("struct name == nil")
		panic(err)
	}

	return *this.struct_name
}

func (this *Symbol) NumStars() int {
	return this.num_stars
}

func (this *Symbol) Name() string {
	return this.name
}
