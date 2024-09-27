package stmt

import (
	"errors"
	"uPIMulator/src/host/interpreter/parser/expr"
)

type ForStmt struct {
	initialization *Stmt
	condition      *expr.Expr
	update         *Stmt
	body           *Stmt
}

func (this *ForStmt) Init(initialization *Stmt, condition *expr.Expr, update *Stmt, body *Stmt) {
	if body.StmtType() != BLOCK {
		err := errors.New("body's stmt type is not block")
		panic(err)
	}

	this.initialization = initialization
	this.condition = condition
	this.update = update
	this.body = body
}

func (this *ForStmt) Initialization() *Stmt {
	return this.initialization
}

func (this *ForStmt) Condition() *expr.Expr {
	return this.condition
}

func (this *ForStmt) Update() *Stmt {
	return this.update
}

func (this *ForStmt) Body() *Stmt {
	return this.body
}
