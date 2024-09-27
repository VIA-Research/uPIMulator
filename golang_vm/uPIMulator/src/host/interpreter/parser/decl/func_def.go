package decl

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/param_list"
	"uPIMulator/src/host/interpreter/parser/stmt"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type FuncDef struct {
	type_specifier *type_specifier.TypeSpecifier
	identifier     *lexer.Token
	param_list     *param_list.ParamList
	body           *stmt.Stmt
}

func (this *FuncDef) Init(
	type_specifier *type_specifier.TypeSpecifier,
	identifier *lexer.Token,
	param_list *param_list.ParamList,
	body *stmt.Stmt,
) {
	if identifier.TokenType() != lexer.IDENTIFIER {
		err := errors.New("identifier's token type is not identifier")
		panic(err)
	} else if body.StmtType() != stmt.BLOCK {
		err := errors.New("body's stmt type is not block")
		panic(err)
	}

	this.type_specifier = type_specifier
	this.identifier = identifier
	this.param_list = param_list
	this.body = body
}

func (this *FuncDef) TypeSpecifier() *type_specifier.TypeSpecifier {
	return this.type_specifier
}

func (this *FuncDef) Identifier() *lexer.Token {
	return this.identifier
}

func (this *FuncDef) ParamList() *param_list.ParamList {
	return this.param_list
}

func (this *FuncDef) Body() *stmt.Stmt {
	return this.body
}
