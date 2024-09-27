package stmt

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type LabelStmt struct {
	expr *expr.Expr
}

func (this *LabelStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("pc is not a program counter")
		panic(err)
	}

	this.expr = expr_
}

func (this *LabelStmt) Expr() *expr.Expr {
	return this.expr
}
