package decl

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/param_list"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type FuncDecl struct {
	type_specifier *type_specifier.TypeSpecifier
	identifier     *lexer.Token
	param_list     *param_list.ParamList
}

func (this *FuncDecl) Init(
	type_specifier *type_specifier.TypeSpecifier,
	identifier *lexer.Token,
	param_list *param_list.ParamList,
) {
	if identifier.TokenType() != lexer.IDENTIFIER {
		err := errors.New("identifier's token type is not identifier")
		panic(err)
	}

	this.type_specifier = type_specifier
	this.identifier = identifier
	this.param_list = param_list
}

func (this *FuncDecl) TypeSpecifier() *type_specifier.TypeSpecifier {
	return this.type_specifier
}

func (this *FuncDecl) Identifier() *lexer.Token {
	return this.identifier
}

func (this *FuncDecl) ParamList() *param_list.ParamList {
	return this.param_list
}
