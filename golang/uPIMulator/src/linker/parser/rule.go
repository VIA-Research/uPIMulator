package parser

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type Reducible func([]*StackItem) bool
type Reduce func([]*StackItem) *StackItem

type Rule struct {
	precedence map[lexer.TokenType]bool

	reducible Reducible
	reduce    Reduce
}

func (this *Rule) Init(precedence map[lexer.TokenType]bool, reducible Reducible, reduce Reduce) {
	this.precedence = precedence
	this.reducible = reducible
	this.reduce = reduce
}

func (this *Rule) IsReducible(stack_items []*StackItem, token *lexer.Token) bool {
	if _, found := this.precedence[token.TokenType()]; found {
		return false
	} else {
		return this.reducible(stack_items)
	}
}

func (this *Rule) Reduce(stack_items []*StackItem, token *lexer.Token) *StackItem {
	if !this.IsReducible(stack_items, token) {
		err := errors.New("stack items are not reducible")
		panic(err)
	}

	return this.reduce(stack_items)
}
