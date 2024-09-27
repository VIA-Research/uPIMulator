package symbol

import (
	"uPIMulator/src/host/vm/base"
	"uPIMulator/src/host/vm/type_system"
)

type Symbol struct {
	name          string
	type_variable *type_system.TypeVariable
	object        *base.Object
}

func (this *Symbol) Init(
	name string,
	type_variable *type_system.TypeVariable,
	object *base.Object,
) {
	this.name = name
	this.type_variable = type_variable
	this.object = object
}

func (this *Symbol) Name() string {
	return this.name
}

func (this *Symbol) TypeVariable() *type_system.TypeVariable {
	return this.type_variable
}

func (this *Symbol) Object() *base.Object {
	return this.object
}
