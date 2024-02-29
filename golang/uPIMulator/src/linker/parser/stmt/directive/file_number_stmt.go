package directive

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type FileNumberStmt struct {
	expr   *expr.Expr
	token1 *lexer.Token
	token2 *lexer.Token
}

func (this *FileNumberStmt) Init(expr_ *expr.Expr, token1 *lexer.Token, token2 *lexer.Token) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	if token1.TokenType() != lexer.STRING {
		err := errors.New("token1 is not a string")
		panic(err)
	}

	if token2.TokenType() != lexer.STRING {
		err := errors.New("token2 is not a string")
		panic(err)
	}

	this.expr = expr_
	this.token1 = token1
	this.token2 = token2
}

func (this *FileNumberStmt) Expr() *expr.Expr {
	return this.expr
}

func (this *FileNumberStmt) Token1() *lexer.Token {
	return this.token1
}

func (this *FileNumberStmt) Token2() *lexer.Token {
	return this.token2
}
