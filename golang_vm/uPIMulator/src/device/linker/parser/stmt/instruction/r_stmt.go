package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type RStmt struct {
	op_code *expr.Expr
	rc      *expr.Expr
}

func (this *RStmt) Init(op_code *expr.Expr, rc *expr.Expr) {
	if op_code.ExprType() != expr.R_OP_CODE {
		err := errors.New("op code is not an R op code")
		panic(err)
	}

	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
}

func (this *RStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RStmt) Rc() *expr.Expr {
	return this.rc
}
