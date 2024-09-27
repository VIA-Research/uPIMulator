package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type JeqRriStmt struct {
	op_code *expr.Expr
	ra      *expr.Expr
	rb      *expr.Expr
	pc      *expr.Expr
}

func (this *JeqRriStmt) Init(op_code *expr.Expr, ra *expr.Expr, rb *expr.Expr, pc *expr.Expr) {
	if op_code.ExprType() != expr.JUMP_OP_CODE {
		err := errors.New("op code is not a jump op code")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	if rb.ExprType() != expr.SRC_REG {
		err := errors.New("rb is not a src reg")
		panic(err)
	}

	if pc.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("pc is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.ra = ra
	this.rb = rb
	this.pc = pc
}

func (this *JeqRriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *JeqRriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *JeqRriStmt) Rb() *expr.Expr {
	return this.rb
}

func (this *JeqRriStmt) Pc() *expr.Expr {
	return this.pc
}
