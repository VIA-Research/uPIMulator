package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SErriStmt struct {
	op_code *expr.Expr
	suffix  *expr.Expr
	endian  *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	off     *expr.Expr
}

func (this *SErriStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	endian *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	off *expr.Expr,
) {
	if op_code.ExprType() != expr.LOAD_OP_CODE {
		err := errors.New("op code is not a load op code")
		panic(err)
	}

	if suffix.ExprType() != expr.SUFFIX {
		err := errors.New("suffix is not a suffix")
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
	this.suffix = suffix
	this.endian = endian
	this.dc = dc
	this.ra = ra
	this.off = off
}

func (this *SErriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SErriStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SErriStmt) Endian() *expr.Expr {
	return this.endian
}

func (this *SErriStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SErriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SErriStmt) Off() *expr.Expr {
	return this.off
}
