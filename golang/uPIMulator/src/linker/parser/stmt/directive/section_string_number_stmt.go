package directive

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SectionStringNumberStmt struct {
	expr1 *expr.Expr
	token *lexer.Token
	expr2 *expr.Expr
	expr3 *expr.Expr
}

func (this *SectionStringNumberStmt) Init(
	expr1 *expr.Expr,
	token *lexer.Token,
	expr2 *expr.Expr,
	expr3 *expr.Expr,
) {
	if expr1.ExprType() != expr.SECTION_NAME {
		err := errors.New("expr1 is not a section name")
		panic(err)
	}

	if token.TokenType() != lexer.STRING {
		err := errors.New("token is not a string")
		panic(err)
	}

	if expr2.ExprType() != expr.SECTION_TYPE {
		err := errors.New("expr2 is not a section type")
		panic(err)
	}

	if expr3.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr3 is not a program counter")
		panic(err)
	}

	this.expr1 = expr1
	this.token = token
	this.expr2 = expr2
	this.expr3 = expr3
}

func (this *SectionStringNumberStmt) Expr1() *expr.Expr {
	return this.expr1
}

func (this *SectionStringNumberStmt) Token() *lexer.Token {
	return this.token
}

func (this *SectionStringNumberStmt) Expr2() *expr.Expr {
	return this.expr2
}

func (this *SectionStringNumberStmt) Expr3() *expr.Expr {
	return this.expr3
}
