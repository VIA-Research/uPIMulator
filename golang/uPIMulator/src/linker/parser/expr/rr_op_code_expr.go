package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type RrOpCodeExpr struct {
	token *lexer.Token
}

func (this *RrOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.CAO &&
		token_type != lexer.CLO &&
		token_type != lexer.CLS &&
		token_type != lexer.CLZ &&
		token_type != lexer.EXTSB &&
		token_type != lexer.EXTSH &&
		token_type != lexer.EXTUB &&
		token_type != lexer.EXTUH &&
		token_type != lexer.SATS &&
		token_type != lexer.TIME_CFG &&
		token_type != lexer.MOVE &&
		token_type != lexer.NEG &&
		token_type != lexer.NOT {
		err := errors.New("token type is not an RR op code")
		panic(err)
	}

	this.token = token
}

func (this *RrOpCodeExpr) Token() *lexer.Token {
	return this.token
}
