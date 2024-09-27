package directive

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type AsciiStmt struct {
	token *lexer.Token
}

func (this *AsciiStmt) Init(token *lexer.Token) {
	if token.TokenType() != lexer.STRING {
		err := errors.New("token is not a string")
		panic(err)
	}

	this.token = token
}

func (this *AsciiStmt) Token() *lexer.Token {
	return this.token
}
