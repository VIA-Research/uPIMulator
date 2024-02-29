package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type CiOpCodeExpr struct {
	token *lexer.Token
}

func (this *CiOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.STOP {
		err := errors.New("token type is not a CI op code")
		panic(err)
	}

	this.token = token
}

func (this *CiOpCodeExpr) Token() *lexer.Token {
	return this.token
}
