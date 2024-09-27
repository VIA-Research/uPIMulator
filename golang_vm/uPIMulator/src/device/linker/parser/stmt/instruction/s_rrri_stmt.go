package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type SRrriStmt struct {
	op_code *expr.Expr
	suffix  *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	rb      *expr.Expr
	imm     *expr.Expr
}

func (this *SRrriStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
	imm *expr.Expr,
) {
	if op_code.ExprType() != expr.RRRI_OP_CODE {
		err := errors.New("op code is not an RRRI op code")
		panic(err)
	}

	if suffix.ExprType() != expr.SUFFIX {
		err := errors.New("suffix is not a suffix")
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

	if rb.ExprType() != expr.SRC_REG {
		err := errors.New("rb is not a src reg")
		panic(err)
	}

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb
	this.imm = imm
}

func (this *SRrriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRrriStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRrriStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRrriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SRrriStmt) Rb() *expr.Expr {
	return this.rb
}

func (this *SRrriStmt) Imm() *expr.Expr {
	return this.imm
}
