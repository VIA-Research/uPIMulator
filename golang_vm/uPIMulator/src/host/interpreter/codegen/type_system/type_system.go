package type_system

import (
	"errors"
	"fmt"
)

type TypeSystem struct {
	methods map[string]*Method
}

func (this *TypeSystem) Init() {
	this.methods = make(map[string]*Method)
}

func (this *TypeSystem) HasMethod(method_name string) bool {
	_, found := this.methods[method_name]
	return found
}

func (this *TypeSystem) Method(method_name string) *Method {
	if !this.HasMethod(method_name) {
		err_msg := fmt.Sprintf("method (%s) is not found", method_name)
		err := errors.New(err_msg)
		panic(err)
	}

	return this.methods[method_name]
}

func (this *TypeSystem) AddMethod(method *Method) {
	this.methods[method.Symbol().Name()] = method
}
