package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type RciStmt struct {
	op_code   *expr.Expr
	rc        *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *RciStmt) Init(op_code *expr.Expr, rc *expr.Expr, condition *expr.Expr, pc *expr.Expr) {
	if op_code.ExprType() != expr.R_OP_CODE {
		err := errors.New("op code is not an R op code")
		panic(err)
	}

	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
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
	this.condition = condition
	this.pc = pc
}

func (this *RciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RciStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *RciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *RciStmt) Pc() *expr.Expr {
	return this.pc
}
