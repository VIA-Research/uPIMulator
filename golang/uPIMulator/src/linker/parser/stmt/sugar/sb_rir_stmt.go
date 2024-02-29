package sugar

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type SbRirStmt struct {
	op_code *expr.Expr
	ra      *expr.Expr
	off     *expr.Expr
	rb      *expr.Expr
}

func (this *SbRirStmt) Init(op_code *expr.Expr, ra *expr.Expr, off *expr.Expr, rb *expr.Expr) {
	if op_code.ExprType() != expr.STORE_OP_CODE {
		err := errors.New("op code is not a store op code")
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

	if rb.ExprType() != expr.SRC_REG {
		err := errors.New("rb is not a src reg")
		panic(err)
	}

	this.op_code = op_code
	this.ra = ra
	this.off = off
	this.rb = rb
}

func (this *SbRirStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SbRirStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SbRirStmt) Off() *expr.Expr {
	return this.off
}

func (this *SbRirStmt) Rb() *expr.Expr {
	return this.rb
}
