package stmt

type BlockStmt struct {
	stmts []*Stmt
}

func (this *BlockStmt) Init() {
	this.stmts = make([]*Stmt, 0)
}

func (this *BlockStmt) Length() int {
	return len(this.stmts)
}

func (this *BlockStmt) Get(pos int) *Stmt {
	return this.stmts[pos]
}

func (this *BlockStmt) Append(stmt *Stmt) {
	this.stmts = append(this.stmts, stmt)
}
