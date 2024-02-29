package sugar

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type BootRiStmt struct {
	op_code *expr.Expr
	ra      *expr.Expr
	imm     *expr.Expr
}

func (this *BootRiStmt) Init(op_code *expr.Expr, ra *expr.Expr, imm *expr.Expr) {
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

	this.op_code = op_code
	this.ra = ra
	this.imm = imm
}

func (this *BootRiStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *BootRiStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *BootRiStmt) Imm() *expr.Expr {
	return this.imm
}
