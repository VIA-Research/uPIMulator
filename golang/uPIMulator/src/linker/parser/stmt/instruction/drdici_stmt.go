package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type DrdiciStmt struct {
	op_code   *expr.Expr
	dc        *lexer.Token
	ra        *expr.Expr
	db        *lexer.Token
	imm       *expr.Expr
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *DrdiciStmt) Init(
	op_code *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	db *lexer.Token,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	if op_code.ExprType() != expr.DRDICI_OP_CODE {
		err := errors.New("op code is not a DRDICI op code")
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

	if db.TokenType() != lexer.PAIR_REG {
		err := errors.New("db is not a pair reg")
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
	this.dc = dc
	this.ra = ra
	this.db = db
	this.imm = imm
	this.condition = condition
	this.pc = pc
}

func (this *DrdiciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *DrdiciStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *DrdiciStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *DrdiciStmt) Db() *lexer.Token {
	return this.db
}

func (this *DrdiciStmt) Imm() *expr.Expr {
	return this.imm
}

func (this *DrdiciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *DrdiciStmt) Pc() *expr.Expr {
	return this.pc
}
