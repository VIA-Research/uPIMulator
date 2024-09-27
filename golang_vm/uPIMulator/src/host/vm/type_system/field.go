package type_system

type Field struct {
	type_variable *TypeVariable
	name          string
}

func (this *Field) Init(type_variable *TypeVariable, name string) {
	this.type_variable = type_variable
	this.name = name
}

func (this *Field) TypeVariable() *TypeVariable {
	return this.type_variable
}

func (this *Field) Name() string {
	return this.name
}
