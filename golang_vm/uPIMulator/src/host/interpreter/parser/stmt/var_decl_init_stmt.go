package stmt

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/expr"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type VarDeclInitStmt struct {
	type_specifier *type_specifier.TypeSpecifier
	identifier     *lexer.Token
	expr           *expr.Expr
}

func (this *VarDeclInitStmt) Init(
	type_specifier *type_specifier.TypeSpecifier,
	identifier *lexer.Token,
	expr_ *expr.Expr,
) {
	if identifier.TokenType() != lexer.IDENTIFIER {
		err := errors.New("identifier's token type is not identifier")
		panic(err)
	}

	this.type_specifier = type_specifier
	this.identifier = identifier
	this.expr = expr_
}

func (this *VarDeclInitStmt) TypeSpecifier() *type_specifier.TypeSpecifier {
	return this.type_specifier
}

func (this *VarDeclInitStmt) Identifier() *lexer.Token {
	return this.identifier
}

func (this *VarDeclInitStmt) Expr() *expr.Expr {
	return this.expr
}
