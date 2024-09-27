package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type DivStepDrdiStmt struct {
	op_code *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	db      *lexer.Token
	imm     *expr.Expr
}

func (this *DivStepDrdiStmt) Init(
	op_code *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	db *lexer.Token,
	imm *expr.Expr,
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

	this.op_code = op_code
	this.dc = dc
	this.ra = ra
	this.db = db
	this.imm = imm
}

func (this *DivStepDrdiStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *DivStepDrdiStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *DivStepDrdiStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *DivStepDrdiStmt) Db() *lexer.Token {
	return this.db
}

func (this *DivStepDrdiStmt) Imm() *expr.Expr {
	return this.imm
}
