package sugar

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type CallRrStmt struct {
	rc *expr.Expr
	ra *expr.Expr
}

func (this *CallRrStmt) Init(rc *expr.Expr, ra *expr.Expr) {
	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	this.rc = rc
	this.ra = ra
}

func (this *CallRrStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *CallRrStmt) Ra() *expr.Expr {
	return this.ra
}
