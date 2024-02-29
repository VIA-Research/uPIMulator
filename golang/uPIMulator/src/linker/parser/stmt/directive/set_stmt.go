package directive

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type SetStmt struct {
	expr1 *expr.Expr
	expr2 *expr.Expr
}

func (this *SetStmt) Init(expr1 *expr.Expr, expr2 *expr.Expr) {
	if expr1.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr1 is not a program counter")
		panic(err)
	}

	if expr2.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr2 is not a program counter")
		panic(err)
	}

	this.expr1 = expr1
	this.expr2 = expr2
}

func (this *SetStmt) Expr1() *expr.Expr {
	return this.expr1
}

func (this *SetStmt) Expr2() *expr.Expr {
	return this.expr2
}
