package sugar

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type LbsSRriStmt struct {
	op_code *expr.Expr
	suffix  *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	off     *expr.Expr
}

func (this *LbsSRriStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	off *expr.Expr,
) {
	if op_code.ExprType() != expr.LOAD_OP_CODE {
		err := errors.New("op code is not a load op code")
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

	if off.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("off is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.off = off
}

func (this *LbsSRriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *LbsSRriStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *LbsSRriStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *LbsSRriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *LbsSRriStmt) Off() *expr.Expr {
	return this.off
}
