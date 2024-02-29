package parser

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type Table struct {
	expr_rules []*Rule
	stmt_rules []*Rule

	stack *Stack
}

func (this *Table) Init(stack *Stack) {
	this.expr_rules = make([]*Rule, 0)
	this.stmt_rules = make([]*Rule, 0)

	this.stack = stack
}

func (this *Table) AddExprRule(rule *Rule) {
	this.expr_rules = append(this.expr_rules, rule)
}

func (this *Table) AddStmtRule(rule *Rule) {
	this.stmt_rules = append(this.stmt_rules, rule)
}

func (this *Table) FindReducibleExprRule(token *lexer.Token) (*Rule, []*StackItem) {
	for num := this.stack.NonStmtSize(); num > 0; num-- {
		stack_items := this.stack.Front(num)

		for _, expr_rule := range this.expr_rules {
			if expr_rule.IsReducible(stack_items, token) {
				return expr_rule, stack_items
			}
		}
	}
	return nil, []*StackItem{}
}

func (this *Table) FindReducibleStmtRule(token *lexer.Token) (*Rule, []*StackItem) {
	if token.TokenType() != lexer.NEW_LINE {
		err := errors.New("token is not a new line")
		panic(err)
	}

	num := this.stack.NonStmtSize()
	stack_items := this.stack.Front(num)

	for _, stmt_rule := range this.stmt_rules {
		if stmt_rule.IsReducible(stack_items, token) {
			return stmt_rule, stack_items
		}
	}
	return nil, []*StackItem{}
}
