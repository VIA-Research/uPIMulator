package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type RirciStmt struct {
	op_code   *expr.Expr
	rc        *expr.Expr
	imm       *expr.Expr
	ra        *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *RirciStmt) Init(
	op_code *expr.Expr,
	rc *expr.Expr,
	imm *expr.Expr,
	ra *expr.Expr,
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

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
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

	if pc.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("pc is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.rc = rc
	this.imm = imm
	this.ra = ra
	this.condition = condition
	this.pc = pc
}

func (this *RirciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RirciStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RirciStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *RirciStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RirciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *RirciStmt) Pc() *expr.Expr {
	return this.pc
}
