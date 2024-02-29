package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type LoadOpCodeExpr struct {
	token *lexer.Token
}

func (this *LoadOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.LBS &&
		token_type != lexer.LBU &&
		token_type != lexer.LD &&
		token_type != lexer.LHS &&
		token_type != lexer.LHU &&
		token_type != lexer.LW {
		err := errors.New("token type is not a load op code")
		panic(err)
	}

	this.token = token
}

func (this *LoadOpCodeExpr) Token() *lexer.Token {
	return this.token
}
