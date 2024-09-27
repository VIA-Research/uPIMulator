package parser

import (
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
	"uPIMulator/src/device/linker/parser/stmt"
)

type StackItemType int

const (
	TOKEN StackItemType = iota
	EXPR
	STMT
)

type StackItem struct {
	stack_item_type StackItemType

	token *lexer.Token
	expr  *expr.Expr
	stmt  *stmt.Stmt
}

func (this *StackItem) InitToken(token *lexer.Token) {
	this.stack_item_type = TOKEN

	this.token = token
}

func (this *StackItem) InitExpr(expr *expr.Expr) {
	this.stack_item_type = EXPR

	this.expr = expr
}

func (this *StackItem) InitStmt(stmt *stmt.Stmt) {
	this.stack_item_type = STMT

	this.stmt = stmt
}

func (this *StackItem) StackItemType() StackItemType {
	return this.stack_item_type
}

func (this *StackItem) Token() *lexer.Token {
	return this.token
}

func (this *StackItem) Expr() *expr.Expr {
	return this.expr
}

func (this *StackItem) Stmt() *stmt.Stmt {
	return this.stmt
}
