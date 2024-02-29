package parser

import (
	"uPIMulator/src/linker/parser/stmt"
)

type Ast struct {
	stmts []*stmt.Stmt
}

func (this *Ast) Init(stmts []*stmt.Stmt) {
	this.stmts = stmts
}

func (this *Ast) Size() int {
	return len(this.stmts)
}

func (this *Ast) Get(pos int) *stmt.Stmt {
	return this.stmts[pos]
}
