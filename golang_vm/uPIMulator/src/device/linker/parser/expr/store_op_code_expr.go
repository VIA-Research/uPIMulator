package expr

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type StoreOpCodeExpr struct {
	token *lexer.Token
}

func (this *StoreOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.SB &&
		token_type != lexer.SB_ID &&
		token_type != lexer.SD &&
		token_type != lexer.SD_ID &&
		token_type != lexer.SH &&
		token_type != lexer.SH_ID &&
		token_type != lexer.SW &&
		token_type != lexer.SW_ID {
		err := errors.New("token type is not a store op code")
		panic(err)
	}

	this.token = token
}

func (this *StoreOpCodeExpr) Token() *lexer.Token {
	return this.token
}
