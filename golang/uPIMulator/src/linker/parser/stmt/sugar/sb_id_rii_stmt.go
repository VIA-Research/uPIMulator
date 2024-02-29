package sugar

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type SbIdRiiStmt struct {
	op_code *expr.Expr
	ra      *expr.Expr
	off     *expr.Expr
	imm     *expr.Expr
}

func (this *SbIdRiiStmt) Init(op_code *expr.Expr, ra *expr.Expr, off *expr.Expr, imm *expr.Expr) {
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

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.ra = ra
	this.off = off
	this.imm = imm
}

func (this *SbIdRiiStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SbIdRiiStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SbIdRiiStmt) Off() *expr.Expr {
	return this.off
}

func (this *SbIdRiiStmt) Imm() *expr.Expr {
	return this.imm
}
