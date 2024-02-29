package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type EridStmt struct {
	op_code *expr.Expr
	endian  *expr.Expr
	ra      *expr.Expr
	off     *expr.Expr
	db      *lexer.Token
}

func (this *EridStmt) Init(
	op_code *expr.Expr,
	endian *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
	db *lexer.Token,
) {
	if op_code.ExprType() != expr.STORE_OP_CODE {
		err := errors.New("op code is not a store op code")
		panic(err)
	}

	if endian.ExprType() != expr.ENDIAN {
		err := errors.New("endian is not an endian")
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
		err := errors.New("db is not a pair reg")
		panic(err)
	}

	this.op_code = op_code
	this.endian = endian
	this.ra = ra
	this.off = off
	this.db = db
}

func (this *EridStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *EridStmt) Endian() *expr.Expr {
	return this.endian
}

func (this *EridStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *EridStmt) Off() *expr.Expr {
	return this.off
}

func (this *EridStmt) Db() *lexer.Token {
	return this.db
}
