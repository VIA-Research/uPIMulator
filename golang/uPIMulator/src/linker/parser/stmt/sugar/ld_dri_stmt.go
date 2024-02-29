package sugar

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type LdDriStmt struct {
	op_code *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	off     *expr.Expr
}

func (this *LdDriStmt) Init(op_code *expr.Expr, dc *lexer.Token, ra *expr.Expr, off *expr.Expr) {
	if op_code.ExprType() != expr.LOAD_OP_CODE {
		err := errors.New("op code is not a load op code")
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

	if off.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("off is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.dc = dc
	this.ra = ra
	this.off = off
}

func (this *LdDriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *LdDriStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *LdDriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *LdDriStmt) Off() *expr.Expr {
	return this.off
}
