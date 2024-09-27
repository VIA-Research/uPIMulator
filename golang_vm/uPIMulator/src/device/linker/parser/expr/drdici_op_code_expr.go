package expr

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type DrdiciOpCodeExpr struct {
	token *lexer.Token
}

func (this *DrdiciOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.DIV_STEP && token_type != lexer.MUL_STEP {
		err := errors.New("token type is not a DRDICI op code")
		panic(err)
	}

	this.token = token
}

func (this *DrdiciOpCodeExpr) Token() *lexer.Token {
	return this.token
}
