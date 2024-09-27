package directive

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/expr"
)

type DefineDirective struct {
	lvalue *lexer.Token
	rvalue *expr.Expr
}

func (this *DefineDirective) Init(lvalue *lexer.Token, rvalue *expr.Expr) {
	if lvalue.TokenType() != lexer.IDENTIFIER {
		err := errors.New("lvalue's token type is not identifier")
		panic(err)
	}

	this.lvalue = lvalue
	this.rvalue = rvalue
}

func (this *DefineDirective) Lvalue() *lexer.Token {
	return this.lvalue
}

func (this *DefineDirective) Rvalue() *expr.Expr {
	return this.rvalue
}
