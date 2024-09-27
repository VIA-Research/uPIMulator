package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type MoveSRiciStmt struct {
	suffix    *expr.Expr
	dc        *lexer.Token
	imm       *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *MoveSRiciStmt) Init(
	suffix *expr.Expr,
	dc *lexer.Token,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	if suffix.ExprType() != expr.SUFFIX {
		err := errors.New("suffix is not a suffix")
		panic(err)
	}

	if dc.TokenType() != lexer.PAIR_REG {
		err := errors.New("dc is not a pair reg")
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

	this.suffix = suffix
	this.dc = dc
	this.imm = imm
	this.condition = condition
	this.pc = pc
}

func (this *MoveSRiciStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *MoveSRiciStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *MoveSRiciStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *MoveSRiciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *MoveSRiciStmt) Pc() *expr.Expr {
	return this.pc
}
