package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type RrrcStmt struct {
	op_code   *expr.Expr
	rc        *expr.Expr
	ra        *expr.Expr
	rb        *expr.Expr
	condition *expr.Expr
}

func (this *RrrcStmt) Init(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	rb *expr.Expr,
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

	if rb.ExprType() != expr.SRC_REG {
		err := errors.New("rb is not a src reg")
		panic(err)
	}

	if condition.ExprType() != expr.CONDITION {
		err := errors.New("condition is not a condition")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
	this.ra = ra
	this.rb = rb
	this.condition = condition
}

func (this *RrrcStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RrrcStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RrrcStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RrrcStmt) Rb() *expr.Expr {
	return this.rb
}

func (this *RrrcStmt) Condition() *expr.Expr {
	return this.condition
}
