package type_system

import (
	"errors"
	"fmt"
)

type Method struct {
	symbol *Symbol
	params []*Symbol
}

func (this *Method) Init(symbol *Symbol) {
	this.symbol = symbol
	this.params = make([]*Symbol, 0)
}

func (this *Method) Symbol() *Symbol {
	return this.symbol
}

func (this *Method) HasParam(param_name string) bool {
	for _, param := range this.params {
		if param.Name() == param_name {
			return true
		}
	}
	return false
}

func (this *Method) Param(param_name string) *Symbol {
	if !this.HasParam(param_name) {
		err_msg := fmt.Sprintf("param (%s) is not found", param_name)
		err := errors.New(err_msg)
		panic(err)
	}

	for _, param := range this.params {
		if param.Name() == param_name {
			return param
		}
	}
	return nil
}

func (this *Method) Params() []*Symbol {
	return this.params
}

func (this *Method) AppendParam(param *Symbol) {
	this.params = append(this.params, param)
}
