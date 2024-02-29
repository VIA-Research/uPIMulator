package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type CiStmt struct {
	op_code   *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *CiStmt) Init(op_code *expr.Expr, condition *expr.Expr, pc *expr.Expr) {
	if op_code.ExprType() != expr.CI_OP_CODE {
		err := errors.New("op code is not a CI op code")
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
	this.condition = condition
	this.pc = pc
}

func (this *CiStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *CiStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *CiStmt) Pc() *expr.Expr {
	return this.pc
}
