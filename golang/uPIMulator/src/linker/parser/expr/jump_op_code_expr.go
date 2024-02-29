package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type JumpOpCodeExpr struct {
	token *lexer.Token
}

func (this *JumpOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.JEQ &&
		token_type != lexer.JNEQ &&
		token_type != lexer.JZ &&
		token_type != lexer.JNZ &&
		token_type != lexer.JLTU &&
		token_type != lexer.JGTU &&
		token_type != lexer.JLEU &&
		token_type != lexer.JGEU &&
		token_type != lexer.JLTS &&
		token_type != lexer.JGTS &&
		token_type != lexer.JLES &&
		token_type != lexer.JGES &&
		token_type != lexer.JUMP {
		err := errors.New("token type is not a jump op code")
		panic(err)
	}

	this.token = token
}

func (this *JumpOpCodeExpr) Token() *lexer.Token {
	return this.token
}
