package symbol

import (
	"errors"
	"fmt"
)

type ScopeChain struct {
	scopes []*Scope
}

func (this *ScopeChain) Init() {
	this.scopes = make([]*Scope, 0)
}

func (this *ScopeChain) NewScope() {
	scope := new(Scope)
	scope.Init()

	this.scopes = append(this.scopes, scope)
}

func (this *ScopeChain) DeleteScope() {
	this.scopes = this.scopes[:len(this.scopes)-1]
}

func (this *ScopeChain) LastScope() *Scope {
	return this.scopes[len(this.scopes)-1]
}

func (this *ScopeChain) HasSymbol(symbol_name string) bool {
	for i := len(this.scopes) - 1; i >= 0; i-- {
		scope := this.scopes[i]

		if scope.HasSymbol(symbol_name) {
			return true
		}
	}
	return false
}

func (this *ScopeChain) Symbol(symbol_name string) *Symbol {
	if !this.HasSymbol(symbol_name) {
		err_msg := fmt.Sprintf("symbol (%s) is not found", symbol_name)
		err := errors.New(err_msg)
		panic(err)
	}

	for i := len(this.scopes) - 1; i >= 0; i-- {
		scope := this.scopes[i]

		if scope.HasSymbol(symbol_name) {
			return scope.Symbol(symbol_name)
		}
	}
	return nil
}

func (this *ScopeChain) Symbols() []*Symbol {
	symbols := make([]*Symbol, 0)

	for i := 0; i < len(this.scopes); i++ {
		scope := this.scopes[i]

		symbols = append(symbols, scope.Symbols()...)
	}

	return symbols
}

func (this *ScopeChain) AddSymbol(symbol *Symbol) {
	if this.scopes[len(this.scopes)-1].HasSymbol(symbol.Name()) {
		err_msg := fmt.Sprintf("symbol (%s) already exists", symbol.Name())
		err := errors.New(err_msg)
		panic(err)
	}

	scope := this.scopes[len(this.scopes)-1]
	scope.AddSymbol(symbol)
}

func (this *ScopeChain) HasObject(address int64) bool {
	for _, scope := range this.scopes {
		if scope.HasObject(address) {
			return true
		}
	}
	return false
}
