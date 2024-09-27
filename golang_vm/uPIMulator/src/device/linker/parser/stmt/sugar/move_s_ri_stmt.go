package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type MoveSRiStmt struct {
	suffix *expr.Expr
	dc     *lexer.Token
	imm    *expr.Expr
}

func (this *MoveSRiStmt) Init(suffix *expr.Expr, dc *lexer.Token, imm *expr.Expr) {
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

	this.suffix = suffix
	this.dc = dc
	this.imm = imm
}

func (this *MoveSRiStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *MoveSRiStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *MoveSRiStmt) Imm() *expr.Expr {
	return this.imm
}
