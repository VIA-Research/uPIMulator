package directive

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SectionStackSizesStmt struct {
	token *lexer.Token
	expr1 *expr.Expr
	expr2 *expr.Expr
	expr3 *expr.Expr
}

func (this *SectionStackSizesStmt) Init(
	token *lexer.Token,
	expr1 *expr.Expr,
	expr2 *expr.Expr,
	expr3 *expr.Expr,
) {
	if token.TokenType() != lexer.STRING {
		err := errors.New("token is not a string")
		panic(err)
	}

	if expr1.ExprType() != expr.SECTION_TYPE {
		err := errors.New("expr1 is not a section type")
		panic(err)
	}

	if expr2.ExprType() != expr.SECTION_NAME {
		err := errors.New("expr2 is not a section name")
		panic(err)
	}

	if expr3.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr3 is not a program counter")
		panic(err)
	}

	this.token = token
	this.expr1 = expr1
	this.expr2 = expr2
	this.expr3 = expr3
}

func (this *SectionStackSizesStmt) Token() *lexer.Token {
	return this.token
}

func (this *SectionStackSizesStmt) Expr1() *expr.Expr {
	return this.expr1
}

func (this *SectionStackSizesStmt) Expr2() *expr.Expr {
	return this.expr2
}

func (this *SectionStackSizesStmt) Expr3() *expr.Expr {
	return this.expr3
}
