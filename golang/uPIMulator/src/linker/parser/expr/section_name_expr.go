package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type SectionNameExpr struct {
	token *lexer.Token
}

func (this *SectionNameExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.ATOMIC &&
		token_type != lexer.BSS &&
		token_type != lexer.DATA &&
		token_type != lexer.DEBUG_ABBREV &&
		token_type != lexer.DEBUG_FRAME &&
		token_type != lexer.DEBUG_INFO &&
		token_type != lexer.DEBUG_LINE &&
		token_type != lexer.DEBUG_LOC &&
		token_type != lexer.DEBUG_RANGES &&
		token_type != lexer.DEBUG_STR &&
		token_type != lexer.DPU_HOST &&
		token_type != lexer.MRAM &&
		token_type != lexer.RODATA &&
		token_type != lexer.STACK_SIZES &&
		token_type != lexer.TEXT {
		err := errors.New("token type is not a section name")
		panic(err)
	}

	this.token = token
}

func (this *SectionNameExpr) Token() *lexer.Token {
	return this.token
}
