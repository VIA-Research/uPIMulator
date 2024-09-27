package directive

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type ZeroSingleNumberStmt struct {
	expr *expr.Expr
}

func (this *ZeroSingleNumberStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	this.expr = expr_
}

func (this *ZeroSingleNumberStmt) Expr() *expr.Expr {
	return this.expr
}
