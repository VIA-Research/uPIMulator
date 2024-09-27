package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type RrcStmt struct {
	op_code   *expr.Expr
	rc        *expr.Expr
	ra        *expr.Expr
	condition *expr.Expr
}

func (this *RrcStmt) Init(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr, condition *expr.Expr) {
	if op_code.ExprType() != expr.RR_OP_CODE {
		err := errors.New("op code is not an RR op code")
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

	if condition.ExprType() != expr.CONDITION {
		err := errors.New("condition is not a condition")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
	this.ra = ra
	this.condition = condition
}

func (this *RrcStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RrcStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RrcStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RrcStmt) Condition() *expr.Expr {
	return this.condition
}
