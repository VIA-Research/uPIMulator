package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RrrStmt struct {
	op_code *expr.Expr
	rc      *expr.Expr
	ra      *expr.Expr
	rb      *expr.Expr
}

func (this *RrrStmt) Init(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr, rb *expr.Expr) {
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

	this.op_code = op_code
	this.rc = rc
	this.ra = ra
	this.rb = rb
}

func (this *RrrStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RrrStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RrrStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RrrStmt) Rb() *expr.Expr {
	return this.rb
}
