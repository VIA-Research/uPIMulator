package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type PrimaryExpr struct {
	token *lexer.Token
}

func (this *PrimaryExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.POSITIVIE_NUMBER && token_type != lexer.HEX_NUMBER &&
		token_type != lexer.IDENTIFIER {
		err := errors.New("token type is not positive number, hex number, nor identifier")
		panic(err)
	}

	this.token = token
}

func (this *PrimaryExpr) Token() *lexer.Token {
	return this.token
}
