package stack

import (
	"uPIMulator/src/host/vm/type_system"
)

type StackItem struct {
	type_variable *type_system.TypeVariable
	address       int64
	size          int64
}

func (this *StackItem) Init(type_variable *type_system.TypeVariable, address int64, size int64) {
	this.type_variable = type_variable
	this.address = address
	this.size = size
}

func (this *StackItem) TypeVariable() *type_system.TypeVariable {
	return this.type_variable
}

func (this *StackItem) Address() int64 {
	return this.address
}

func (this *StackItem) Size() int64 {
	return this.size
}
