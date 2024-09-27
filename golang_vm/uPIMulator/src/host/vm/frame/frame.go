package frame

import (
	"uPIMulator/src/host/abi"
	"uPIMulator/src/host/vm/pc"
	"uPIMulator/src/host/vm/stack"
	"uPIMulator/src/host/vm/symbol"
)

type Frame struct {
	stack            *stack.Stack
	return_stack     *stack.ReturnStack
	fast_scope_chain *symbol.ScopeChain
	arg_scope        *symbol.Scope
	pc               *pc.Pc
}

func (this *Frame) Init(label *abi.Label) {
	this.stack = new(stack.Stack)
	this.stack.Init()

	this.return_stack = new(stack.ReturnStack)
	this.return_stack.Init()

	this.fast_scope_chain = new(symbol.ScopeChain)
	this.fast_scope_chain.Init()

	this.arg_scope = new(symbol.Scope)
	this.arg_scope.Init()

	this.pc = new(pc.Pc)
	this.pc.Init()
	this.pc.Jump(label)
}

func (this *Frame) Stack() *stack.Stack {
	return this.stack
}

func (this *Frame) ReturnStack() *stack.ReturnStack {
	return this.return_stack
}

func (this *Frame) FastScopeChain() *symbol.ScopeChain {
	return this.fast_scope_chain
}

func (this *Frame) ArgScope() *symbol.Scope {
	return this.arg_scope
}

func (this *Frame) Pc() *pc.Pc {
	return this.pc
}

func (this *Frame) Symbols() []*symbol.Symbol {
	symbols := make([]*symbol.Symbol, 0)
	symbols = append(symbols, this.fast_scope_chain.Symbols()...)
	symbols = append(symbols, this.arg_scope.Symbols()...)
	return symbols
}
