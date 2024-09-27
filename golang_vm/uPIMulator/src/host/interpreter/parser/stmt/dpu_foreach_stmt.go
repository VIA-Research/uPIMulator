package stmt

import (
	"errors"
	"uPIMulator/src/host/interpreter/parser/expr"
)

type DpuForeachStmt struct {
	foreach *expr.ArgList
	body    *Stmt
}

func (this *DpuForeachStmt) Init(foreach *expr.ArgList, body *Stmt) {
	if foreach.Length() != 3 {
		err := errors.New("arg list's length != 3")
		panic(err)
	} else if body.StmtType() != BLOCK {
		err := errors.New("body's stmt type is not block")
		panic(err)
	}

	this.foreach = foreach
	this.body = body
}

func (this *DpuForeachStmt) Foreach() *expr.ArgList {
	return this.foreach
}

func (this *DpuForeachStmt) Body() *Stmt {
	return this.body
}
