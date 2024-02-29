package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SRrriciStmt struct {
	op_code   *expr.Expr
	suffix    *expr.Expr
	dc        *lexer.Token
	ra        *expr.Expr
	rb        *expr.Expr
	imm       *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *SRrriciStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
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

	if condition.ExprType() != expr.CONDITION {
		err := errors.New("condition is not a condition")
		panic(err)
	}

	if pc.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("pc is not a program counter")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb
	this.imm = imm
	this.condition = condition
	this.pc = pc
}

func (this *SRrriciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRrriciStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRrriciStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRrriciStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SRrriciStmt) Rb() *expr.Expr {
	return this.rb
}

func (this *SRrriciStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *SRrriciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *SRrriciStmt) Pc() *expr.Expr {
	return this.pc
}
