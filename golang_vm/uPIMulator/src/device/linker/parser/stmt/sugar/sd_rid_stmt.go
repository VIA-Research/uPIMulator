package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type SdRidStmt struct {
	op_code *expr.Expr
	ra      *expr.Expr
	off     *expr.Expr
	db      *lexer.Token
}

func (this *SdRidStmt) Init(op_code *expr.Expr, ra *expr.Expr, off *expr.Expr, db *lexer.Token) {
	if op_code.ExprType() != expr.STORE_OP_CODE {
		err := errors.New("op code is not a store op code")
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

	if db.TokenType() != lexer.PAIR_REG {
		err := errors.New("dc is not a pair reg")
		panic(err)
	}

	this.op_code = op_code
	this.ra = ra
	this.off = off
	this.db = db
}

func (this *SdRidStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SdRidStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SdRidStmt) Off() *expr.Expr {
	return this.off
}

func (this *SdRidStmt) Db() *lexer.Token {
	return this.db
}
