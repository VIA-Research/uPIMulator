package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type JumpIStmt struct {
	pc *expr.Expr
}

func (this *JumpIStmt) Init(pc *expr.Expr) {
	if pc.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("pc is not a program counter")
		panic(err)
	}

	this.pc = pc
}

func (this *JumpIStmt) Pc() *expr.Expr {
	return this.pc
}
