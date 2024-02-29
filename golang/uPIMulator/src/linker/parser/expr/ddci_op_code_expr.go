package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type DdciOpCodeExpr struct {
	token *lexer.Token
}

func (this *DdciOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.MOVD && token_type != lexer.SWAPD {
		err := errors.New("token type is not a DDCI op code")
		panic(err)
	}

	this.token = token
}

func (this *DdciOpCodeExpr) Token() *lexer.Token {
	return this.token
}
