package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type RiciOpCodeExpr struct {
	token *lexer.Token
}

func (this *RiciOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.ACQUIRE &&
		token_type != lexer.RELEASE &&
		token_type != lexer.BOOT &&
		token_type != lexer.RESUME {
		err := errors.New("token type is not an RICI op code")
		panic(err)
	}

	this.token = token
}

func (this *RiciOpCodeExpr) Token() *lexer.Token {
	return this.token
}
