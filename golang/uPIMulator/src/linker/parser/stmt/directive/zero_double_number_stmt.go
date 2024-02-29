package directive

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type ZeroDoubleNumberStmt struct {
	expr1 *expr.Expr
	expr2 *expr.Expr
}

func (this *ZeroDoubleNumberStmt) Init(expr1 *expr.Expr, expr2 *expr.Expr) {
	if expr1.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	if expr2.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	this.expr1 = expr1
	this.expr2 = expr2
}

func (this *ZeroDoubleNumberStmt) Expr1() *expr.Expr {
	return this.expr1
}

func (this *ZeroDoubleNumberStmt) Expr2() *expr.Expr {
	return this.expr2
}
