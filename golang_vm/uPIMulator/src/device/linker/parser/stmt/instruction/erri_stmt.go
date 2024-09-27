package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type ErriStmt struct {
	op_code *expr.Expr
	endian  *expr.Expr
	rc      *expr.Expr
	ra      *expr.Expr
	off     *expr.Expr
}

func (this *ErriStmt) Init(
	op_code *expr.Expr,
	endian *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
) {
	if op_code.ExprType() != expr.LOAD_OP_CODE {
		err := errors.New("op code is not a load op code")
		panic(err)
	}

	if endian.ExprType() != expr.ENDIAN {
		err := errors.New("endian is not an endian")
		panic(err)
	}

	if rc.ExprType() != expr.SRC_REG {
		err := errors.New("rc is not a src reg")
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

	this.op_code = op_code
	this.endian = endian
	this.rc = rc
	this.ra = ra
	this.off = off
}

func (this *ErriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *ErriStmt) Endian() *expr.Expr {
	return this.endian
}

func (this *ErriStmt) Rc() *expr.Expr {
	return this.rc
}

func (this *ErriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *ErriStmt) Off() *expr.Expr {
	return this.off
}
