package instruction

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
)

type NopStmt struct {
	op_code *expr.Expr
}

func (this *NopStmt) Init(op_code *expr.Expr) {
	if op_code.ExprType() != expr.R_OP_CODE {
		err := errors.New("op code is not an R op code")
		panic(err)
	}

	this.op_code = op_code
}

func (this *NopStmt) OpCode() *expr.Expr {
	return this.op_code
}
