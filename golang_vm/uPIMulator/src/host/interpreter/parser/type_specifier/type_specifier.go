package type_specifier

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
)

type TypeSpecifierType int

const (
	VOID TypeSpecifierType = iota
	CHAR
	SHORT
	INT
	LONG
	STRUCT
)

type TypeSpecifier struct {
	type_specifier_type TypeSpecifierType
	struct_identifier   *lexer.Token
	num_stars           int
}

func (this *TypeSpecifier) InitPrimitive(type_specifier_type TypeSpecifierType) {
	if type_specifier_type == STRUCT {
		err := errors.New("type specifier type is a struct")
		panic(err)
	}

	this.type_specifier_type = type_specifier_type
	this.struct_identifier = nil
	this.num_stars = 0
}

func (this *TypeSpecifier) InitStruct(
	type_specifier_type TypeSpecifierType,
	struct_identifier *lexer.Token,
) {
	if type_specifier_type != STRUCT {
		err := errors.New("type specifier type is not a struct")
		panic(err)
	}

	this.type_specifier_type = type_specifier_type
	this.struct_identifier = struct_identifier
	this.num_stars = 0
}

func (this *TypeSpecifier) TypeSpecifierType() TypeSpecifierType {
	return this.type_specifier_type
}

func (this *TypeSpecifier) StructIdentifier() *lexer.Token {
	if this.struct_identifier == nil {
		err := errors.New("struct identifier == nil")
		panic(err)
	}

	return this.struct_identifier
}

func (this *TypeSpecifier) NumStars() int {
	return this.num_stars
}

func (this *TypeSpecifier) AddStar() {
	this.num_stars++
}
