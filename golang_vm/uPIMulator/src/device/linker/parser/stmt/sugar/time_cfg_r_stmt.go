package sugar

import (
	"errors"
	"uPIMulator/src/device/linker/parser/expr"
)

type TimeCfgRStmt struct {
	ra *expr.Expr
}

func (this *TimeCfgRStmt) Init(ra *expr.Expr) {
	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	this.ra = ra
}

func (this *TimeCfgRStmt) Ra() *expr.Expr {
	return this.ra
}
