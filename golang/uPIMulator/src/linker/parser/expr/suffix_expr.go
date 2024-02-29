package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type SuffixExpr struct {
	token *lexer.Token
}

func (this *SuffixExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.S && token_type != lexer.U {
		err := errors.New("token type is not a suffix")
		panic(err)
	}

	this.token = token
}

func (this *SuffixExpr) Token() *lexer.Token {
	return this.token
}
