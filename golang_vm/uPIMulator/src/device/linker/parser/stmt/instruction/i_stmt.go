package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type IStmt struct {
	op_code *expr.Expr
	imm     *expr.Expr
}

func (this *IStmt) Init(op_code *expr.Expr, imm *expr.Expr) {
	if op_code.ExprType() != expr.I_OP_CODE {
		err := errors.New("op code is not an I op code")
		panic(err)
	}

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.imm = imm
}

func (this *IStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *IStmt) Imm() *expr.Expr {
	return this.imm
}
