package expr

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type NegativeNumberExpr struct {
	token *lexer.Token
}

func (this *NegativeNumberExpr) Init(token *lexer.Token) {
	if token.TokenType() != lexer.POSITIVIE_NUMBER {
		err := errors.New("token type is not positive number")
		panic(err)
	}

	this.token = token
}

func (this *NegativeNumberExpr) Token() *lexer.Token {
	return this.token
}
