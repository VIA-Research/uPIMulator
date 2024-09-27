package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type MovdDdStmt struct {
	op_code *expr.Expr
	dc      *lexer.Token
	db      *lexer.Token
}

func (this *MovdDdStmt) Init(op_code *expr.Expr, dc *lexer.Token, db *lexer.Token) {
	if op_code.ExprType() != expr.DDCI_OP_CODE {
		err := errors.New("op code is not a DDCI op code")
		panic(err)
	}

	if dc.TokenType() != lexer.PAIR_REG {
		err := errors.New("dc is not a pair reg")
		panic(err)
	}

	if db.TokenType() != lexer.PAIR_REG {
		err := errors.New("db is not a pair reg")
		panic(err)
	}

	this.op_code = op_code
	this.dc = dc
	this.db = db
}

func (this *MovdDdStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *MovdDdStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *MovdDdStmt) Db() *lexer.Token {
	return this.db
}
