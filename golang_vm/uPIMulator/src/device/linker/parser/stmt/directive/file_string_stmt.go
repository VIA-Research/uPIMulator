package directive

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
)

type FileStringStmt struct {
	token *lexer.Token
}

func (this *FileStringStmt) Init(token *lexer.Token) {
	if token.TokenType() != lexer.STRING {
		err := errors.New("token is not a string")
		panic(err)
	}

	this.token = token
}

func (this *FileStringStmt) Token() *lexer.Token {
	return this.token
}
