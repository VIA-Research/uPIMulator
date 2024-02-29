package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type RrriOpCodeExpr struct {
	token *lexer.Token
}

func (this *RrriOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.LSL_ADD &&
		token_type != lexer.LSL_SUB &&
		token_type != lexer.LSR_ADD &&
		token_type != lexer.ROL_ADD {
		err := errors.New("token type is not an RRRI op code")
		panic(err)
	}

	this.token = token
}

func (this *RrriOpCodeExpr) Token() *lexer.Token {
	return this.token
}
