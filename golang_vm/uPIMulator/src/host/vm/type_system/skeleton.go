package type_system

import (
	"errors"
	"fmt"
)

type Skeleton struct {
	name   string
	fields []*Field
}

func (this *Skeleton) Init(name string) {
	this.name = name
}

func (this *Skeleton) Name() string {
	return this.name
}

func (this *Skeleton) HasField(field_name string) bool {
	for _, field := range this.fields {
		if field.Name() == field_name {
			return true
		}
	}
	return false
}

func (this *Skeleton) Field(field_name string) *Field {
	if !this.HasField(field_name) {
		err_msg := fmt.Sprintf("skeleton (%s) does not have the field (%s)", this.name, field_name)
		err := errors.New(err_msg)
		panic(err)
	}

	for _, field := range this.fields {
		if field.Name() == field_name {
			return field
		}
	}
	return nil
}

func (this *Skeleton) Fields() []*Field {
	return this.fields
}

func (this *Skeleton) Length() int {
	return len(this.fields)
}

func (this *Skeleton) Get(pos int) *Field {
	return this.fields[pos]
}

func (this *Skeleton) Append(field *Field) {
	this.fields = append(this.fields, field)
}
