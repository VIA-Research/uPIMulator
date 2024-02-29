package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RricStmt struct {
	op_code   *expr.Expr
	rc        *expr.Expr
	ra        *expr.Expr
	imm       *expr.Expr
	condition *expr.Expr
}

func (this *RricStmt) Init(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
) {
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

	if condition.ExprType() != expr.CONDITION {
		err := errors.New("condition is not a condition")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
	this.ra = ra
	this.imm = imm
	this.condition = condition
}

func (this *RricStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RricStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RricStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RricStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *RricStmt) Condition() *expr.Expr {
	return this.condition
}
