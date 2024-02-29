package expr

import (
	"errors"
	"uPIMulator/src/linker/lexer"
)

type ConditionExpr struct {
	token *lexer.Token
}

func (this *ConditionExpr) Init(token *lexer.Token) {
	token_type := token.TokenType()

	if token_type != lexer.TRUE &&
		token_type != lexer.FALSE &&
		token_type != lexer.Z &&
		token_type != lexer.NZ &&
		token_type != lexer.E &&
		token_type != lexer.O &&
		token_type != lexer.PL &&
		token_type != lexer.MI &&
		token_type != lexer.OV &&
		token_type != lexer.NOV &&
		token_type != lexer.C &&
		token_type != lexer.NC &&
		token_type != lexer.SZ &&
		token_type != lexer.SNZ &&
		token_type != lexer.SPL &&
		token_type != lexer.SMI &&
		token_type != lexer.SO &&
		token_type != lexer.SE &&
		token_type != lexer.NC5 &&
		token_type != lexer.NC6 &&
		token_type != lexer.NC7 &&
		token_type != lexer.NC8 &&
		token_type != lexer.NC9 &&
		token_type != lexer.NC10 &&
		token_type != lexer.NC11 &&
		token_type != lexer.NC12 &&
		token_type != lexer.NC13 &&
		token_type != lexer.NC14 &&
		token_type != lexer.MAX &&
		token_type != lexer.NMAX &&
		token_type != lexer.SH32 &&
		token_type != lexer.NSH32 &&
		token_type != lexer.EQ &&
		token_type != lexer.NEQ &&
		token_type != lexer.LTU &&
		token_type != lexer.LEU &&
		token_type != lexer.GTU &&
		token_type != lexer.GEU &&
		token_type != lexer.LTS &&
		token_type != lexer.LES &&
		token_type != lexer.GTS &&
		token_type != lexer.GES &&
		token_type != lexer.XZ &&
		token_type != lexer.XNZ &&
		token_type != lexer.XLEU &&
		token_type != lexer.XGTU &&
		token_type != lexer.XLES &&
		token_type != lexer.XGTS &&
		token_type != lexer.SMALL &&
		token_type != lexer.LARGE {
		err := errors.New("token type is not a condition")
		panic(err)
	}

	this.token = token
}

func (this *ConditionExpr) Token() *lexer.Token {
	return this.token
}
