package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type RriOpCodeExpr struct {
	token *lexer.Token
}

func (this *RriOpCodeExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.ADD &&
		token_type != lexer.ADDC &&
		token_type != lexer.AND &&
		token_type != lexer.ANDN &&
		token_type != lexer.ASR &&
		token_type != lexer.CMPB4 &&
		token_type != lexer.LSL &&
		token_type != lexer.LSL1 &&
		token_type != lexer.LSL1X &&
		token_type != lexer.LSLX &&
		token_type != lexer.LSR &&
		token_type != lexer.LSR1 &&
		token_type != lexer.LSR1X &&
		token_type != lexer.LSRX &&
		token_type != lexer.MUL_SH_SH &&
		token_type != lexer.MUL_SH_SL &&
		token_type != lexer.MUL_SH_UH &&
		token_type != lexer.MUL_SH_UL &&
		token_type != lexer.MUL_SL_SH &&
		token_type != lexer.MUL_SL_SL &&
		token_type != lexer.MUL_SL_UH &&
		token_type != lexer.MUL_SL_UL &&
		token_type != lexer.MUL_UH_UH &&
		token_type != lexer.MUL_UH_UL &&
		token_type != lexer.MUL_UL_UH &&
		token_type != lexer.MUL_UL_UL &&
		token_type != lexer.NAND &&
		token_type != lexer.NOR &&
		token_type != lexer.NXOR &&
		token_type != lexer.OR &&
		token_type != lexer.ORN &&
		token_type != lexer.ROL &&
		token_type != lexer.ROR &&
		token_type != lexer.RSUB &&
		token_type != lexer.RSUBC &&
		token_type != lexer.SUB &&
		token_type != lexer.SUBC &&
		token_type != lexer.XOR &&
		token_type != lexer.CALL {
		err := errors.New("token type is not an RRI op code")
		panic(err)
	}

	this.token = token
}

func (this *RriOpCodeExpr) Token() *lexer.Token {
	return this.token
}
