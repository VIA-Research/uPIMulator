package expr

import (
	"errors"
)

type ProgramCounterExpr struct {
	expr *Expr
}

func (this *ProgramCounterExpr) Init(expr *Expr) {
	expr_type := expr.ExprType()

	if expr_type != PRIMARY && expr_type != NEGATIVE_NUMBER && expr_type != BINARY_ADD &&
		expr_type != BINARY_SUB {
		err := errors.New("expr type is not primary, negative number, binary add, nor binary sub")
		panic(err)
	}

	this.expr = expr
}

func (this *ProgramCounterExpr) Expr() *Expr {
	return this.expr
}
