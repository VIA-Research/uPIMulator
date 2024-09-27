package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type DdciStmt struct {
	op_code   *expr.Expr
	dc        *lexer.Token
	db        *lexer.Token
	condition *expr.Expr
	pc        *expr.Expr
}

func (this *DdciStmt) Init(
	op_code *expr.Expr,
	dc *lexer.Token,
	db *lexer.Token,
	condition *expr.Expr,
	pc *expr.Expr,
) {
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
	this.db = db
	this.condition = condition
	this.pc = pc
}

func (this *DdciStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *DdciStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *DdciStmt) Db() *lexer.Token {
	return this.db
}

func (this *DdciStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *DdciStmt) Pc() *expr.Expr {
	return this.pc
}
