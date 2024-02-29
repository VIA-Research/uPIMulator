package directive

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type AddrsigSymStmt struct {
	expr *expr.Expr
}

func (this *AddrsigSymStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr is not a program counter")
		panic(err)
	}

	this.expr = expr_
}

func (this *AddrsigSymStmt) Expr() *expr.Expr {
	return this.expr
}
