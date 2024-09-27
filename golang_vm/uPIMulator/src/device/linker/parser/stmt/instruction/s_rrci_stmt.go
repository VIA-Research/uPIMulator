package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type SRrciStmt struct {
	op_code   *expr.Expr
	suffix    *expr.Expr
	dc        *lexer.Token
	ra        *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *SRrciStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	if op_code.ExprType() != expr.RR_OP_CODE {
		err := errors.New("op code is not an RR op code")
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
	this.condition = condition
	this.pc = pc
}

func (this *SRrciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRrciStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRrciStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRrciStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SRrciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *SRrciStmt) Pc() *expr.Expr {
	return this.pc
}
