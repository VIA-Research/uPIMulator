package directive

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type LongSectionNameStmt struct {
	expr *expr.Expr
}

func (this *LongSectionNameStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.SECTION_NAME {
		err := errors.New("expr is not a section name")
		panic(err)
	}

	this.expr = expr_
}

func (this *LongSectionNameStmt) Expr() *expr.Expr {
	return this.expr
}
