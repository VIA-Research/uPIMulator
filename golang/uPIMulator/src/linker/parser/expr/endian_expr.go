package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type EndianExpr struct {
	token *lexer.Token
}

func (this *EndianExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.LITTLE &&
		token_type != lexer.BIG {
		err := errors.New("token type is not an endian")
		panic(err)
	}

	this.token = token
}

func (this *EndianExpr) Token() *lexer.Token {
	return this.token
}
