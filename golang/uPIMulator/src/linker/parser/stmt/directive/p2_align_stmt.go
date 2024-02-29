package directive

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type P2AlignStmt struct {
	expr *expr.Expr
}

func (this *P2AlignStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	this.expr = expr_
}

func (this *P2AlignStmt) Expr() *expr.Expr {
	return this.expr
}
