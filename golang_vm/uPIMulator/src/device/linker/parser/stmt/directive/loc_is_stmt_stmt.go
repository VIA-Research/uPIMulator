package directive

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type LocIsStmtStmt struct {
	expr1 *expr.Expr
	expr2 *expr.Expr
	expr3 *expr.Expr
	expr4 *expr.Expr
}

func (this *LocIsStmtStmt) Init(
	expr1 *expr.Expr,
	expr2 *expr.Expr,
	expr3 *expr.Expr,
	expr4 *expr.Expr,
) {
	if expr1.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr1 is not a program counter")
		panic(err)
	}

	if expr2.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr2 is not a program counter")
		panic(err)
	}

	if expr3.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr3 is not a program counter")
		panic(err)
	}

	if expr4.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr4 is not a program counter")
		panic(err)
	}

	this.expr1 = expr1
	this.expr2 = expr2
	this.expr3 = expr3
	this.expr4 = expr4
}

func (this *LocIsStmtStmt) Expr1() *expr.Expr {
	return this.expr1
}

func (this *LocIsStmtStmt) Expr2() *expr.Expr {
	return this.expr2
}

func (this *LocIsStmtStmt) Expr3() *expr.Expr {
	return this.expr3
}

func (this *LocIsStmtStmt) Expr4() *expr.Expr {
	return this.expr4
}
