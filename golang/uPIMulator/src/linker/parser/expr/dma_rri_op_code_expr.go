package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type DmaRriOpCodeExpr struct {
	token *lexer.Token
}

func (this *DmaRriOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.LDMA && token_type != lexer.LDMAI && token_type != lexer.SDMA {
		err := errors.New("token type is not a DMA_RRI op code")
		panic(err)
	}

	this.token = token
}

func (this *DmaRriOpCodeExpr) Token() *lexer.Token {
	return this.token
}
