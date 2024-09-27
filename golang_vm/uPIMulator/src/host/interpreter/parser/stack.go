package parser

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
)

type Stack struct {
	stack_items []*StackItem
}

func (this *Stack) Init() {
	this.stack_items = make([]*StackItem, 0)
}

func (this *Stack) Length() int {
	return len(this.stack_items)
}

func (this *Stack) Push(stack_item *StackItem) {
	this.stack_items = append(this.stack_items, stack_item)
}

func (this *Stack) Pop(num int) {
	this.stack_items = this.stack_items[:len(this.stack_items)-num]
}

func (this *Stack) Front(num int) []*StackItem {
	stack_items := make([]*StackItem, 0)
	for i := 0; i < num; i++ {
		stack_item := this.stack_items[len(this.stack_items)-num+i]
		stack_items = append(stack_items, stack_item)
	}
	return stack_items
}

func (this *Stack) CanAccept() bool {
	for i, stack_item := range this.stack_items {
		if i < len(this.stack_items)-1 {
			if stack_item.StackItemType() != DIRECTIVE && stack_item.StackItemType() != DECL {
				return false
			}
		} else {
			if stack_item.StackItemType() != TOKEN || stack_item.Token().TokenType() != lexer.END_OF_FILE {
				return false
			}
		}
	}
	return true
}

func (this *Stack) Accept() *Ast {
	if !this.CanAccept() {
		err := errors.New("stack cannot be accepted")
		panic(err)
	}

	ast := new(Ast)
	ast.Init(this.stack_items[:len(this.stack_items)-1])

	return ast
}
