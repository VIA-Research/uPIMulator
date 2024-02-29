package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type DmaRriStmt struct {
	op_code *expr.Expr
	ra      *expr.Expr
	rb      *expr.Expr
	imm     *expr.Expr
}

func (this *DmaRriStmt) Init(op_code *expr.Expr, ra *expr.Expr, rb *expr.Expr, imm *expr.Expr) {
	if op_code.ExprType() != expr.DMA_RRI_OP_CODE {
		err := errors.New("op code is not a DMA_RRI op code")
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

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.ra = ra
	this.rb = rb
	this.imm = imm
}

func (this *DmaRriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *DmaRriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *DmaRriStmt) Rb() *expr.Expr {
	return this.rb
}

func (this *DmaRriStmt) Imm() *expr.Expr {
	return this.imm
}
