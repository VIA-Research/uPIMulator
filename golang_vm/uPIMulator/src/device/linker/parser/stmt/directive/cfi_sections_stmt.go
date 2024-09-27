package directive

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type CfiSectionsStmt struct {
	expr *expr.Expr
}

func (this *CfiSectionsStmt) Init(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.SECTION_NAME {
		err := errors.New("expr is not a section name")
		panic(err)
	}

	this.expr = expr_
}

func (this *CfiSectionsStmt) Expr() *expr.Expr {
	return this.expr
}
