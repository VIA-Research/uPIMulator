package symbol

import (
	"errors"
	"fmt"
)

type Scope struct {
	symbols map[string]*Symbol
}

func (this *Scope) Init() {
	this.symbols = make(map[string]*Symbol)
}

func (this *Scope) HasSymbol(symbol_name string) bool {
	_, found := this.symbols[symbol_name]
	return found
}

func (this *Scope) Symbol(symbol_name string) *Symbol {
	if !this.HasSymbol(symbol_name) {
		err_msg := fmt.Sprintf("symbol (%s) is not found", symbol_name)
		err := errors.New(err_msg)
		panic(err)
	}

	return this.symbols[symbol_name]
}

func (this *Scope) Symbols() []*Symbol {
	symbols := make([]*Symbol, 0)
	for _, symbol := range this.symbols {
		symbols = append(symbols, symbol)
	}
	return symbols
}

func (this *Scope) AddSymbol(symbol *Symbol) {
	if this.HasSymbol(symbol.Name()) {
		err_msg := fmt.Sprintf("symbol (%s) already exists", symbol.Name())
		err := errors.New(err_msg)
		panic(err)
	}

	this.symbols[symbol.Name()] = symbol
}

func (this *Scope) HasObject(address int64) bool {
	for _, symbol := range this.symbols {
		if symbol.Object().Address() == address {
			return true
		}
	}
	return false
}

func (this *Scope) Clear() {
	this.symbols = make(map[string]*Symbol)
}
