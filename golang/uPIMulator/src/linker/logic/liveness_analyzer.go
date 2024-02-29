package logic

import (
	"errors"
	"uPIMulator/src/linker/kernel"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser"
	"uPIMulator/src/linker/parser/expr"
	"uPIMulator/src/linker/parser/stmt"
)

type LivenessAnalyzer struct {
	liveness *kernel.Liveness
	walker   *parser.Walker
}

func (this *LivenessAnalyzer) Init() {
	this.liveness = new(kernel.Liveness)
	this.liveness.Init()

	this.walker = new(parser.Walker)
	this.walker.Init()

	this.walker.RegisterExprCallback(expr.PRIMARY, this.WalkPrimaryExpr)
	this.walker.RegisterStmtCallback(stmt.GLOBAL, this.WalkGlobalStmt)
	this.walker.RegisterStmtCallback(stmt.SET, this.WalkSetStmt)
	this.walker.RegisterStmtCallback(stmt.LABEL, this.WalkLabelStmt)
}

func (this *LivenessAnalyzer) Analyze(relocatable *kernel.Relocatable) *kernel.Liveness {
	this.walker.Walk(relocatable.Ast())
	return this.liveness
}

func (this *LivenessAnalyzer) WalkPrimaryExpr(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PRIMARY {
		err := errors.New("expr type is not primary")
		panic(err)
	}

	primary_expr := expr_.PrimaryExpr()

	token := primary_expr.Token()
	if token.TokenType() == lexer.IDENTIFIER {
		this.liveness.AddUse(token.Attribute())
	}
}

func (this *LivenessAnalyzer) WalkGlobalStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.GLOBAL {
		err := errors.New("stmt type is not global")
		panic(err)
	}

	global_stmt := stmt_.GlobalStmt()

	program_counter_expr := global_stmt.Expr().ProgramCounterExpr()
	primary_expr := program_counter_expr.Expr().PrimaryExpr()

	token := primary_expr.Token()

	if token.TokenType() != lexer.IDENTIFIER {
		err := errors.New("token type is not identifier")
		panic(err)
	}

	attribute := token.Attribute()
	if attribute != "__sys_used_mram_end" {
		this.liveness.AddGlobalSymbol(attribute)
	}
}

func (this *LivenessAnalyzer) WalkSetStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SET {
		err := errors.New("stmt type is not set")
		panic(err)
	}

	set_stmt := stmt_.SetStmt()

	program_counter_expr1 := set_stmt.Expr1().ProgramCounterExpr()
	program_counter_expr2 := set_stmt.Expr2().ProgramCounterExpr()

	primary_expr1 := program_counter_expr1.Expr().PrimaryExpr()
	primary_expr2 := program_counter_expr2.Expr().PrimaryExpr()

	token1 := primary_expr1.Token()
	token2 := primary_expr2.Token()

	if token1.TokenType() != lexer.IDENTIFIER {
		err := errors.New("token1 type is not identifier")
		panic(err)
	}

	if token2.TokenType() != lexer.IDENTIFIER {
		err := errors.New("token2 type is not identifier")
		panic(err)
	}

	this.liveness.AddDef(token1.Attribute())
	this.liveness.AddUse(token2.Attribute())
}

func (this *LivenessAnalyzer) WalkLabelStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.LABEL {
		err := errors.New("stmt type is not label")
		panic(err)
	}

	label_stmt := stmt_.LabelStmt()

	program_counter_expr := label_stmt.Expr().ProgramCounterExpr()
	primary_expr := program_counter_expr.Expr().PrimaryExpr()
	token := primary_expr.Token()

	if token.TokenType() != lexer.IDENTIFIER {
		err := errors.New("token type is not identifier")
		panic(err)
	}

	attribute := token.Attribute()
	if attribute != "__sys_used_mram_end" {
		this.liveness.AddDef(attribute)
	}
}
