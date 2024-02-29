package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type SectionTypeExpr struct {
	token *lexer.Token
}

func (this *SectionTypeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.PROGBITS && token_type != lexer.NOBITS {
		err := errors.New("token type is not a section type")
		panic(err)
	}

	this.token = token
}

func (this *SectionTypeExpr) Token() *lexer.Token {
	return this.token
}
