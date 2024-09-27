package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type CallRiStmt struct {
	rc  *expr.Expr
	imm *expr.Expr
}

func (this *CallRiStmt) Init(rc *expr.Expr, imm *expr.Expr) {
	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
		panic(err)
	}

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	this.rc = rc
	this.imm = imm
}

func (this *CallRiStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *CallRiStmt) Imm() *expr.Expr {
	return this.imm
}
