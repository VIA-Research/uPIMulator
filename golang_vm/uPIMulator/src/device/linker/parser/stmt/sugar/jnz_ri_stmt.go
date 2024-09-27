package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type JnzRiStmt struct {
	op_code *expr.Expr
	ra      *expr.Expr
	pc      *expr.Expr
}

func (this *JnzRiStmt) Init(op_code *expr.Expr, ra *expr.Expr, pc *expr.Expr) {
	if op_code.ExprType() != expr.JUMP_OP_CODE {
		err := errors.New("op code is not a jump op code")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	if pc.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("pc is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.ra = ra
	this.pc = pc
}

func (this *JnzRiStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *JnzRiStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *JnzRiStmt) Pc() *expr.Expr {
	return this.pc
}
