package directive

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
)

type IncludeDirective struct {
	header *lexer.Token
}

func (this *IncludeDirective) Init(header *lexer.Token) {
	if header.TokenType() != lexer.HEADER {
		err := errors.New("header's token type is not header")
		panic(err)
	}

	this.header = header
}

func (this *IncludeDirective) Header() *lexer.Token {
	return this.header
}
