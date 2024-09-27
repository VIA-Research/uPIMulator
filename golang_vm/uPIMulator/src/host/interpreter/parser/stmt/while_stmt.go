package stmt

import (
	"errors"
	"uPIMulator/src/host/interpreter/parser/expr"
)

type WhileStmt struct {
	condition *expr.Expr
	body      *Stmt
}

func (this *WhileStmt) Init(condition *expr.Expr, body *Stmt) {
	if body.StmtType() != BLOCK {
		err := errors.New("body's stmt type is not block")
		panic(err)
	}

	this.condition = condition
	this.body = body
}

func (this *WhileStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *WhileStmt) Body() *Stmt {
	return this.body
}
