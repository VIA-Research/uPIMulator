package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SRciStmt struct {
	op_code   *expr.Expr
	suffix    *expr.Expr
	dc        *lexer.Token
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *SRciStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	if op_code.ExprType() != expr.R_OP_CODE {
		err := errors.New("op code is not an R op code")
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
	this.condition = condition
	this.pc = pc
}

func (this *SRciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRciStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRciStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *SRciStmt) Pc() *expr.Expr {
	return this.pc
}
