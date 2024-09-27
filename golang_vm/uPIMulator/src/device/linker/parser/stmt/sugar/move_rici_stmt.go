package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type MoveRiciStmt struct {
	rc        *expr.Expr
	imm       *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *MoveRiciStmt) Init(rc *expr.Expr, imm *expr.Expr, condition *expr.Expr, pc *expr.Expr) {
	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
		panic(err)
	}

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
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

	this.rc = rc
	this.imm = imm
	this.condition = condition
	this.pc = pc
}

func (this *MoveRiciStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *MoveRiciStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *MoveRiciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *MoveRiciStmt) Pc() *expr.Expr {
	return this.pc
}
