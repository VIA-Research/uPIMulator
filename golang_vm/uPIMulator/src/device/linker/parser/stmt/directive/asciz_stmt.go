package directive

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type AscizStmt struct {
	token *lexer.Token
}

func (this *AscizStmt) Init(token *lexer.Token) {
	if token.TokenType() != lexer.STRING {
		err := errors.New("token is not a string")
		panic(err)
	}

	this.token = token
}

func (this *AscizStmt) Token() *lexer.Token {
	return this.token
}
