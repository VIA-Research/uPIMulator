package directive

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type WeakStmt struct {
	expr *expr.Expr
}

func (this *WeakStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	this.expr = expr_
}

func (this *WeakStmt) Expr() *expr.Expr {
	return this.expr
}
