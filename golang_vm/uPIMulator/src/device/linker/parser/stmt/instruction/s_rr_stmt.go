package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type SRrStmt struct {
	op_code *expr.Expr
	suffix  *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
}

func (this *SRrStmt) Init(op_code *expr.Expr, suffix *expr.Expr, dc *lexer.Token, ra *expr.Expr) {
	if op_code.ExprType() != expr.RR_OP_CODE {
		err := errors.New("op code is not an RR op code")
		panic(err)
	}

	if suffix.ExprType() != expr.SUFFIX {
		err := errors.New("suffix is not a suffix")
		panic(err)
	}

	if dc.TokenType() != lexer.PAIR_REG {
		err := errors.New("dc is not a pair reg")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
}

func (this *SRrStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRrStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRrStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRrStmt) Ra() *expr.Expr {
	return this.ra
}
