package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type ErirStmt struct {
	op_code *expr.Expr
	endian  *expr.Expr
	ra      *expr.Expr
	off     *expr.Expr
	rb      *expr.Expr
}

func (this *ErirStmt) Init(
	op_code *expr.Expr,
	endian *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
	rb *expr.Expr,
) {
	if op_code.ExprType() != expr.STORE_OP_CODE {
		err := errors.New("op code is not a store op code")
		panic(err)
	}

	if endian.ExprType() != expr.ENDIAN {
		err := errors.New("endian is not an endian")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	if off.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("off is not a program counter")
		panic(err)
	}

	if rb.ExprType() != expr.SRC_REG {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.endian = endian
	this.ra = ra
	this.off = off
	this.rb = rb
}

func (this *ErirStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *ErirStmt) Endian() *expr.Expr {
	return this.endian
}

func (this *ErirStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *ErirStmt) Off() *expr.Expr {
	return this.off
}

func (this *ErirStmt) Rb() *expr.Expr {
	return this.rb
}
