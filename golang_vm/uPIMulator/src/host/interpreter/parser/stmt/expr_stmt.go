package stmt

import (
	"uPIMulator/src/host/interpreter/parser/expr"
)

type ExprStmt struct {
	expr *expr.Expr
}

func (this *ExprStmt) Init(expr_ *expr.Expr) {
	this.expr = expr_
}

func (this *ExprStmt) Expr() *expr.Expr {
	return this.expr
}
