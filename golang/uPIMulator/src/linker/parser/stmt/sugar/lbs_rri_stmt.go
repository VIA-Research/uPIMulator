package sugar

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type LbsRriStmt struct {
	op_code *expr.Expr
	rc      *expr.Expr
	ra      *expr.Expr
	off     *expr.Expr
}

func (this *LbsRriStmt) Init(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr, off *expr.Expr) {
	if op_code.ExprType() != expr.LOAD_OP_CODE {
		err := errors.New("op code is not a load op code")
		panic(err)
	}

	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
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
	this.rc = rc
	this.ra = ra
	this.off = off
}

func (this *LbsRriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *LbsRriStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *LbsRriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *LbsRriStmt) Off() *expr.Expr {
	return this.off
}
