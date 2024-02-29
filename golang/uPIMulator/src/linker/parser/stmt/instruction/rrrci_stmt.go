package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RrrciStmt struct {
	op_code   *expr.Expr
	rc        *expr.Expr
	ra        *expr.Expr
	rb        *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *RrrciStmt) Init(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	rb *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
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

	if pc.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("pc is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
	this.ra = ra
	this.rb = rb
	this.condition = condition
	this.pc = pc
}

func (this *RrrciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RrrciStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RrrciStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RrrciStmt) Rb() *expr.Expr {
	return this.rb
}

func (this *RrrciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *RrrciStmt) Pc() *expr.Expr {
	return this.pc
}
