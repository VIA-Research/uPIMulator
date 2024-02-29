package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RrriciStmt struct {
	op_code   *expr.Expr
	rc        *expr.Expr
	ra        *expr.Expr
	rb        *expr.Expr
	imm       *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *RrriciStmt) Init(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	rb *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	if op_code.ExprType() != expr.RRRI_OP_CODE {
		err := errors.New("op code is not an RRRI op code")
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

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a progrcm counter")
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
	this.imm = imm
	this.condition = condition
	this.pc = pc
}

func (this *RrriciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RrriciStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RrriciStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RrriciStmt) Rb() *expr.Expr {
	return this.rb
}

func (this *RrriciStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *RrriciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *RrriciStmt) Pc() *expr.Expr {
	return this.pc
}
