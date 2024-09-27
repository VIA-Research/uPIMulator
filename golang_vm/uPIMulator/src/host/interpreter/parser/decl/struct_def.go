package decl

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/stmt"
)

type StructDef struct {
	identifier *lexer.Token
	body       *stmt.Stmt
}

func (this *StructDef) Init(identifier *lexer.Token, body *stmt.Stmt) {
	if identifier.TokenType() != lexer.IDENTIFIER {
		err := errors.New("identifier's token type is not identifier")
		panic(err)
	} else if body.StmtType() != stmt.BLOCK {
		err := errors.New("body's stmt type is not block")
		panic(err)
	}

	this.identifier = identifier
	this.body = body
}

func (this *StructDef) Identifier() *lexer.Token {
	return this.identifier
}

func (this *StructDef) Body() *stmt.Stmt {
	return this.body
}
