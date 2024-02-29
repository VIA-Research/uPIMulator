package parser

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/stmt"
)

type Stack struct {
	stack_items []*StackItem
}

func (this *Stack) Init() {
	this.stack_items = make([]*StackItem, 0)
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

func (this *Stack) NonStmtSize() int {
	non_stmt_size := 0
	for i := len(this.stack_items) - 1; i >= 0; i-- {
		stack_item := this.stack_items[i]

		if stack_item.StackItemType() != STMT {
			non_stmt_size++
		} else {
			break
		}
	}
	return non_stmt_size
}

func (this *Stack) AreStmts() bool {
	for _, stack_item := range this.stack_items {
		if stack_item.StackItemType() == STMT {
			continue
		} else {
			return false
		}
	}
	return true
}

func (this *Stack) CanAccept() bool {
	for i, stack_item := range this.stack_items {
		if i < len(this.stack_items)-1 {
			if stack_item.StackItemType() == STMT {
				continue
			} else {
				return false
			}
		} else {
			if stack_item.StackItemType() == TOKEN && stack_item.Token().TokenType() == lexer.END_OF_FILE {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (this *Stack) Accept() *Ast {
	if !this.CanAccept() {
		err := errors.New("stack cannot be accepted")
		panic(err)
	}

	stmts := make([]*stmt.Stmt, 0)
	for i := 0; i < len(this.stack_items)-1; i++ {
		stmts = append(stmts, this.stack_items[i].Stmt())
	}

	ast := new(Ast)
	ast.Init(stmts)

	return ast
}
