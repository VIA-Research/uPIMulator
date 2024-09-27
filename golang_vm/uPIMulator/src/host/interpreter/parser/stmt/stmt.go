package stmt

type StmtType int

const (
	EMPTY StmtType = iota
	VAR_DECL
	VAR_DECL_INIT
	FOR
	DPU_FOREACH
	WHILE
	CONTINUE
	BREAK
	IF
	RETURN
	EXPR
	BLOCK
)

type Stmt struct {
	stmt_type StmtType

	empty_stmt         *EmptyStmt
	var_decl_stmt      *VarDeclStmt
	var_decl_init_stmt *VarDeclInitStmt
	for_stmt           *ForStmt
	dpu_foreach_stmt   *DpuForeachStmt
	while_stmt         *WhileStmt
	continue_stmt      *ContinueStmt
	break_stmt         *BreakStmt
	if_stmt            *IfStmt
	return_stmt        *ReturnStmt
	expr_stmt          *ExprStmt
	block_stmt         *BlockStmt
}

func (this *Stmt) InitEmptyStmt(empty_stmt *EmptyStmt) {
	this.stmt_type = EMPTY

	this.empty_stmt = empty_stmt
}

func (this *Stmt) InitVarDeclStmt(var_decl_stmt *VarDeclStmt) {
	this.stmt_type = VAR_DECL

	this.var_decl_stmt = var_decl_stmt
}

func (this *Stmt) InitVarDeclInitStmt(var_decl_init_stmt *VarDeclInitStmt) {
	this.stmt_type = VAR_DECL_INIT

	this.var_decl_init_stmt = var_decl_init_stmt
}

func (this *Stmt) InitForStmt(for_stmt *ForStmt) {
	this.stmt_type = FOR

	this.for_stmt = for_stmt
}

func (this *Stmt) InitDpuForeachStmt(dpu_foreach_stmt *DpuForeachStmt) {
	this.stmt_type = DPU_FOREACH

	this.dpu_foreach_stmt = dpu_foreach_stmt
}

func (this *Stmt) InitWhileStmt(while_stmt *WhileStmt) {
	this.stmt_type = WHILE

	this.while_stmt = while_stmt
}

func (this *Stmt) InitContinueStmt(continue_stmt *ContinueStmt) {
	this.stmt_type = CONTINUE

	this.continue_stmt = continue_stmt
}

func (this *Stmt) InitBreakStmt(break_stmt *BreakStmt) {
	this.stmt_type = BREAK

	this.break_stmt = break_stmt
}

func (this *Stmt) InitIfStmt(if_stmt *IfStmt) {
	this.stmt_type = IF

	this.if_stmt = if_stmt
}

func (this *Stmt) InitReturnStmt(return_stmt *ReturnStmt) {
	this.stmt_type = RETURN

	this.return_stmt = return_stmt
}

func (this *Stmt) InitExprStmt(expr_stmt *ExprStmt) {
	this.stmt_type = EXPR

	this.expr_stmt = expr_stmt
}

func (this *Stmt) InitBlockStmt(block_stmt *BlockStmt) {
	this.stmt_type = BLOCK

	this.block_stmt = block_stmt
}

func (this *Stmt) StmtType() StmtType {
	return this.stmt_type
}

func (this *Stmt) EmptyStmt() *EmptyStmt {
	return this.empty_stmt
}

func (this *Stmt) VarDeclStmt() *VarDeclStmt {
	return this.var_decl_stmt
}

func (this *Stmt) VarDeclInitStmt() *VarDeclInitStmt {
	return this.var_decl_init_stmt
}

func (this *Stmt) ForStmt() *ForStmt {
	return this.for_stmt
}

func (this *Stmt) DpuForeachStmt() *DpuForeachStmt {
	return this.dpu_foreach_stmt
}

func (this *Stmt) WhileStmt() *WhileStmt {
	return this.while_stmt
}

func (this *Stmt) ContinueStmt() *ContinueStmt {
	return this.continue_stmt
}

func (this *Stmt) BreakStmt() *BreakStmt {
	return this.break_stmt
}

func (this *Stmt) IfStmt() *IfStmt {
	return this.if_stmt
}

func (this *Stmt) ReturnStmt() *ReturnStmt {
	return this.return_stmt
}

func (this *Stmt) ExprStmt() *ExprStmt {
	return this.expr_stmt
}

func (this *Stmt) BlockStmt() *BlockStmt {
	return this.block_stmt
}
