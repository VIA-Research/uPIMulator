package expr

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
)

type PrimaryExprType int

const (
	IDENTIFIER PrimaryExprType = iota
	NUMBER
	STRING
	NULLPTR
	PAREN
)

type PrimaryExpr struct {
	primary_expr_type PrimaryExprType
	token             *lexer.Token
	expr              *Expr
}

func (this *PrimaryExpr) InitIdentifier(token *lexer.Token) {
	this.primary_expr_type = IDENTIFIER
	this.token = token
	this.expr = nil
}

func (this *PrimaryExpr) InitNumber(token *lexer.Token) {
	this.primary_expr_type = NUMBER
	this.token = token
	this.expr = nil
}

func (this *PrimaryExpr) InitString(token *lexer.Token) {
	this.primary_expr_type = STRING
	this.token = token
	this.expr = nil
}

func (this *PrimaryExpr) InitNullptr(token *lexer.Token) {
	this.primary_expr_type = NULLPTR
	this.token = token
	this.expr = nil
}

func (this *PrimaryExpr) InitParen(expr *Expr) {
	this.primary_expr_type = PAREN
	this.token = nil
	this.expr = expr
}

func (this *PrimaryExpr) PrimaryExprType() PrimaryExprType {
	return this.primary_expr_type
}

func (this *PrimaryExpr) Token() *lexer.Token {
	if this.token == nil {
		err := errors.New("token == nil")
		panic(err)
	}

	return this.token
}

func (this *PrimaryExpr) Expr() *Expr {
	if this.expr == nil {
		err := errors.New("expr == nil")
		panic(err)
	}

	return this.expr
}
