package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type JumpRStmt struct {
	ra *expr.Expr
}

func (this *JumpRStmt) Init(ra *expr.Expr) {
	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	this.ra = ra
}

func (this *JumpRStmt) Ra() *expr.Expr {
	return this.ra
}
