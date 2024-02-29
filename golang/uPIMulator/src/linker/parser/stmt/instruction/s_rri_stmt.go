package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SRriStmt struct {
	op_code *expr.Expr
	suffix  *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	imm     *expr.Expr
}

func (this *SRriStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	imm *expr.Expr,
) {
	if op_code.ExprType() != expr.RRI_OP_CODE {
		err := errors.New("op code is not an RRI op code")
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

	if imm.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("imm is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.imm = imm
}

func (this *SRriStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRriStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRriStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRriStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SRriStmt) Imm() *expr.Expr {
	return this.imm
}
