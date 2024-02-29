package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RriStmt struct {
	op_code *expr.Expr
	rc      *expr.Expr
	ra      *expr.Expr
	imm     *expr.Expr
}

func (this *RriStmt) Init(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr, imm *expr.Expr) {
	if op_code.ExprType() != expr.RRI_OP_CODE {
		err := errors.New("op code is not an RRI op code")
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

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a progrcm counter")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
	this.ra = ra
	this.imm = imm
}

func (this *RriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RriStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RriStmt) Imm() *expr.Expr {
	return this.imm
}
