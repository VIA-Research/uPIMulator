package expr

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type SymbolTypeExpr struct {
	token *lexer.Token
}

func (this *SymbolTypeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.FUNCTION && token_type != lexer.OBJECT {
		err := errors.New("token type is not a symbol type")
		panic(err)
	}

	this.token = token
}

func (this *SymbolTypeExpr) Token() *lexer.Token {
	return this.token
}
