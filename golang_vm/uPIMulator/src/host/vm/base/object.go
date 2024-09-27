package base

import (
	"errors"
	"uPIMulator/src/host/vm/type_system"
)

type ObjectType int

const (
	TEMPORARY ObjectType = iota
	UNTEMPORARY
)

type Object struct {
	object_type   ObjectType
	type_variable *type_system.TypeVariable

	address int64
	size    int64
}

func (this *Object) Init(object_type ObjectType, address int64, size int64) {
	if address <= 0 {
		err := errors.New("address <= 0")
		panic(err)
	} else if size <= 0 {
		err := errors.New("size <= 0")
		panic(err)
	}

	this.object_type = object_type

	this.address = address
	this.size = size
}

func (this *Object) ObjectType() ObjectType {
	return this.object_type
}

func (this *Object) HasTypeVariable() bool {
	return this.type_variable != nil
}

func (this *Object) TypeVariable() *type_system.TypeVariable {
	return this.type_variable
}

func (this *Object) SetTypeVariable(type_variable *type_system.TypeVariable) {
	if this.type_variable != nil {
		err := errors.New("type variable != nil")
		panic(err)
	}

	this.type_variable = type_variable
}

func (this *Object) Address() int64 {
	return this.address
}

func (this *Object) Size() int64 {
	return this.size
}
