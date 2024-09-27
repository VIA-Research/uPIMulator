package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type RiciStmt struct {
	op_code   *expr.Expr
	ra        *expr.Expr
	imm       *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *RiciStmt) Init(
	op_code *expr.Expr,
	ra *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	if op_code.ExprType() != expr.RICI_OP_CODE {
		err := errors.New("op code is not an RICI op code")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
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

	this.op_code = op_code
	this.ra = ra
	this.imm = imm
	this.condition = condition
	this.pc = pc
}

func (this *RiciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *RiciStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *RiciStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *RiciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *RiciStmt) Pc() *expr.Expr {
	return this.pc
}
