package expr

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type IOpCodeExpr struct {
	token *lexer.Token
}

func (this *IOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.FAULT && token_type != lexer.BKP {
		err := errors.New("token type is not an I op code")
		panic(err)
	}

	this.token = token
}

func (this *IOpCodeExpr) Token() *lexer.Token {
	return this.token
}
