package parser

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
)

type Table struct {
	rules []*Rule

	stack *Stack
}

func (this *Table) Init(stack *Stack) {
	this.rules = make([]*Rule, 0)

	this.stack = stack
}

func (this *Table) AddRule(rule *Rule) {
	this.rules = append(this.rules, rule)
}

func (this *Table) IsReducible(token *lexer.Token) bool {
	for i := 0; i < this.stack.Length(); i++ {
		stack_items := this.stack.Front(this.stack.Length() - i)

		for _, rule := range this.rules {
			if rule.IsReducible(stack_items, token) {
				return true
			}
		}
	}
	return false
}

func (this *Table) Reduce(token *lexer.Token) {
	for i := 0; i < this.stack.Length(); i++ {
		stack_items := this.stack.Front(this.stack.Length() - i)

		for _, rule := range this.rules {
			if rule.IsReducible(stack_items, token) {
				stack_item := rule.Reduce(stack_items, token)

				this.stack.Pop(len(stack_items))
				this.stack.Push(stack_item)

				return
			}
		}
	}

	err := errors.New("stack is not reducible")
	panic(err)
}
