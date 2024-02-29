package directive

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SectionStringStmt struct {
	expr1 *expr.Expr
	token *lexer.Token
	expr2 *expr.Expr
}

func (this *SectionStringStmt) Init(expr1 *expr.Expr, token *lexer.Token, expr2 *expr.Expr) {
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

	this.expr1 = expr1
	this.token = token
	this.expr2 = expr2
}

func (this *SectionStringStmt) Expr1() *expr.Expr {
	return this.expr1
}

func (this *SectionStringStmt) Token() *lexer.Token {
	return this.token
}

func (this *SectionStringStmt) Expr2() *expr.Expr {
	return this.expr2
}
