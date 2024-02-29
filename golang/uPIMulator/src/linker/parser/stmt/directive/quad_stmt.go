package directive

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type QuadStmt struct {
	expr *expr.Expr
}

func (this *QuadStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	this.expr = expr_
}

func (this *QuadStmt) Expr() *expr.Expr {
	return this.expr
}
