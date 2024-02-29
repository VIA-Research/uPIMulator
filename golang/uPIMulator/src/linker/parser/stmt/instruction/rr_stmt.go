package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RrStmt struct {
	op_code *expr.Expr
	rc      *expr.Expr
	ra      *expr.Expr
}

func (this *RrStmt) Init(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr) {
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

	this.op_code = op_code
	this.rc = rc
	this.ra = ra
}

func (this *RrStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RrStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RrStmt) Ra() *expr.Expr {
	return this.ra
}
