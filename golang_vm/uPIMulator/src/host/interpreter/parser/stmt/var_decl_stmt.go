package stmt

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type VarDeclStmt struct {
	type_specifier *type_specifier.TypeSpecifier
	identifier     *lexer.Token
}

func (this *VarDeclStmt) Init(
	type_specifier *type_specifier.TypeSpecifier,
	identifier *lexer.Token,
) {
	if identifier.TokenType() != lexer.IDENTIFIER {
		err := errors.New("identifier's token type is not identifier")
		panic(err)
	}

	this.type_specifier = type_specifier
	this.identifier = identifier
}

func (this *VarDeclStmt) TypeSpecifier() *type_specifier.TypeSpecifier {
	return this.type_specifier
}

func (this *VarDeclStmt) Identifier() *lexer.Token {
	return this.identifier
}
