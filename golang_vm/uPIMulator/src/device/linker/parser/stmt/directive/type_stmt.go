package directive

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type TypeStmt struct {
	expr1 *expr.Expr
	expr2 *expr.Expr
}

func (this *TypeStmt) Init(expr1 *expr.Expr, expr2 *expr.Expr) {
	if expr1.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	if expr2.ExprType() != expr.SYMBOL_TYPE {
		err := errors.New("expr is not a symbol type")
		panic(err)
	}

	this.expr1 = expr2
}

func (this *TypeStmt) Expr1() *expr.Expr {
	return this.expr1
}

func (this *TypeStmt) Expr2() *expr.Expr {
	return this.expr2
}
