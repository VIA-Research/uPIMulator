package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RirStmt struct {
	op_code *expr.Expr
	rc      *expr.Expr
	imm     *expr.Expr
	ra      *expr.Expr
}

func (this *RirStmt) Init(op_code *expr.Expr, rc *expr.Expr, imm *expr.Expr, ra *expr.Expr) {
	if op_code.ExprType() != expr.RRI_OP_CODE {
		err := errors.New("op code is not an RRI op code")
		panic(err)
	}

	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
		panic(err)
	}

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
	this.imm = imm
	this.ra = ra
}

func (this *RirStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RirStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RirStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *RirStmt) Ra() *expr.Expr {
	return this.ra
}
