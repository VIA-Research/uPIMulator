package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type EdriStmt struct {
	op_code *expr.Expr
	endian  *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	off     *expr.Expr
}

func (this *EdriStmt) Init(
	op_code *expr.Expr,
	endian *expr.Expr,
	dc *lexer.Token,
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

	if dc.TokenType() != lexer.PAIR_REG {
		err := errors.New("dc is not a pair reg")
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
	this.dc = dc
	this.ra = ra
	this.off = off
}

func (this *EdriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *EdriStmt) Endian() *expr.Expr {
	return this.endian
}

func (this *EdriStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *EdriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *EdriStmt) Off() *expr.Expr {
	return this.off
}
