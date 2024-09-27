package stmt

import (
	"errors"
	"uPIMulator/src/host/interpreter/parser/expr"
)

type IfStmt struct {
	if_condition *expr.Expr
	if_body      *Stmt

	else_if_conditions []*expr.Expr
	else_if_bodies     []*Stmt

	else_body *Stmt
}

func (this *IfStmt) Init(if_condition *expr.Expr, if_body *Stmt) {
	if if_body.StmtType() != BLOCK {
		err := errors.New("if body's stmt type is not block")
		panic(err)
	}

	this.if_condition = if_condition
	this.if_body = if_body

	this.else_if_conditions = make([]*expr.Expr, 0)
	this.else_if_bodies = make([]*Stmt, 0)

	this.else_body = nil
}

func (this *IfStmt) IfCondition() *expr.Expr {
	return this.if_condition
}

func (this *IfStmt) IfBody() *Stmt {
	return this.if_body
}

func (this *IfStmt) NumElseIfs() int {
	if len(this.else_if_conditions) != len(this.else_if_bodies) {
		err := errors.New("lengths of else if conditions and bodies are different")
		panic(err)
	}

	return len(this.else_if_conditions)
}

func (this *IfStmt) ElseIfCondition(pos int) *expr.Expr {
	return this.else_if_conditions[pos]
}

func (this *IfStmt) ElseIfBody(pos int) *Stmt {
	return this.else_if_bodies[pos]
}

func (this *IfStmt) AppendElseIf(condition *expr.Expr, body *Stmt) {
	if body.StmtType() != BLOCK {
		err := errors.New("body's stmt type is not block")
		panic(err)
	}

	this.else_if_conditions = append(this.else_if_conditions, condition)
	this.else_if_bodies = append(this.else_if_bodies, body)
}

func (this *IfStmt) HasElseBody() bool {
	return this.else_body != nil
}

func (this *IfStmt) ElseBody() *Stmt {
	if !this.HasElseBody() {
		err := errors.New("else body does not exist")
		panic(err)
	}

	return this.else_body
}

func (this *IfStmt) SetElseBody(body *Stmt) {
	if this.HasElseBody() {
		err := errors.New("else body already exists")
		panic(err)
	} else if body.StmtType() != BLOCK {
		err := errors.New("body's stmt type is not block")
		panic(err)
	}

	this.else_body = body
}
