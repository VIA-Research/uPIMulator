package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type SrcRegExpr struct {
	token *lexer.Token
}

func (this *SrcRegExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.GP_REG &&
		token_type != lexer.ZERO_REG &&
		token_type != lexer.ONE &&
		token_type != lexer.ID &&
		token_type != lexer.ID2 &&
		token_type != lexer.ID4 &&
		token_type != lexer.ID8 &&
		token_type != lexer.LNEG &&
		token_type != lexer.MNEG {
		err := errors.New("token type is not a src reg")
		panic(err)
	}

	this.token = token
}

func (this *SrcRegExpr) Token() *lexer.Token {
	return this.token
}
