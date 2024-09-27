package param_list

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type Param struct {
	type_specifier *type_specifier.TypeSpecifier
	identifier     *lexer.Token
}

func (this *Param) Init(type_specifier *type_specifier.TypeSpecifier, identifier *lexer.Token) {
	if identifier.TokenType() != lexer.IDENTIFIER {
		err := errors.New("identifier's token type is not identifier")
		panic(err)
	}

	this.type_specifier = type_specifier
	this.identifier = identifier
}

func (this *Param) TypeSpecifier() *type_specifier.TypeSpecifier {
	return this.type_specifier
}

func (this *Param) Identifier() *lexer.Token {
	return this.identifier
}
