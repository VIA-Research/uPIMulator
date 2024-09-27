package logic

import (
	"errors"
	"uPIMulator/src/device/linker/kernel"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser"
	"uPIMulator/src/device/linker/parser/expr"
	"uPIMulator/src/device/linker/parser/stmt"
)

type SetAssigner struct {
	executable *kernel.Executable
	walker     *parser.Walker
}

func (this *SetAssigner) Init() {
	this.walker = new(parser.Walker)
	this.walker.Init()

	this.walker.RegisterStmtCallback(
		stmt.SECTION_IDENTIFIER_NUMBER,
		this.WalkSectionIdentifierNumberStmt,
	)
	this.walker.RegisterStmtCallback(stmt.SECTION_IDENTIFIER, this.WalkSectionIdentifierStmt)
	this.walker.RegisterStmtCallback(stmt.SECTION_STACK_SIZES, this.WalkSectionStackSizes)
	this.walker.RegisterStmtCallback(stmt.SECTION_STRING_NUMBER, this.WalkSectionStringNumberStmt)
	this.walker.RegisterStmtCallback(stmt.SECTION_STRING, this.WalkSectionStringStmt)
	this.walker.RegisterStmtCallback(stmt.TEXT, this.WalkTextStmt)
}

func (this *SetAssigner) Assign(executable *kernel.Executable) {
	this.executable = executable
	this.walker.Walk(executable.Ast())
}

func (this *SetAssigner) WalkSectionIdentifierNumberStmt(stmt_ *stmt.Stmt) {
	section_identifier_number_stmt := stmt_.SectionIdentifierNumberStmt()

	section_name := this.ConvertSectionName(section_identifier_number_stmt.Expr1())
	name := this.ConvertName(section_identifier_number_stmt.Expr2())

	this.executable.CheckoutSection(section_name, name)
}

func (this *SetAssigner) WalkSectionIdentifierStmt(stmt_ *stmt.Stmt) {
	section_identifier_stmt := stmt_.SectionIdentifierStmt()

	section_name := this.ConvertSectionName(section_identifier_stmt.Expr1())
	name := this.ConvertName(section_identifier_stmt.Expr2())

	this.executable.CheckoutSection(section_name, name)
}

func (this *SetAssigner) WalkSectionStackSizes(stmt_ *stmt.Stmt) {
	section_stack_sizes_stmt := stmt_.SectionStackSizesStmt()

	section_name := kernel.STACK_SIZES

	section_name_expr := section_stack_sizes_stmt.Expr2().SectionNameExpr()
	token := section_name_expr.Token()
	token_type := token.TokenType()

	name := ""
	if token_type == lexer.ATOMIC {
		name += ".atomic."
	} else if token_type == lexer.BSS {
		name += ".bss."
	} else if token_type == lexer.DATA {
		name += ".data."
	} else if token_type == lexer.DEBUG_ABBREV {
		name += ".debug_abbrev."
	} else if token_type == lexer.DEBUG_FRAME {
		name += ".debug_frame."
	} else if token_type == lexer.DEBUG_INFO {
		name += ".debug_info."
	} else if token_type == lexer.DEBUG_LINE {
		name += ".debug_line."
	} else if token_type == lexer.DEBUG_LOC {
		name += ".debug_loc."
	} else if token_type == lexer.DEBUG_RANGES {
		name += ".debug_ranges."
	} else if token_type == lexer.DEBUG_STR {
		name += ".debug_str."
	} else if token_type == lexer.DPU_HOST {
		name += ".dpu_host."
	} else if token_type == lexer.MRAM {
		name += ".mram."
	} else if token_type == lexer.RODATA {
		name += ".rodata."
	} else if token_type == lexer.STACK_SIZES {
		name += ".stack_sizes."
	} else if token_type == lexer.TEXT {
		name += ".text."
	} else {
		err := errors.New("section name is not valid")
		panic(err)
	}

	name += this.ConvertName(section_stack_sizes_stmt.Expr3())

	this.executable.CheckoutSection(section_name, name)
}

func (this *SetAssigner) WalkSectionStringNumberStmt(stmt_ *stmt.Stmt) {
	section_string_number_stmt := stmt_.SectionStringNumberStmt()

	section_name := this.ConvertSectionName(section_string_number_stmt.Expr1())
	name := ""

	this.executable.CheckoutSection(section_name, name)
}

func (this *SetAssigner) WalkSectionStringStmt(stmt_ *stmt.Stmt) {
	section_string_stmt := stmt_.SectionStringStmt()

	section_name := this.ConvertSectionName(section_string_stmt.Expr1())
	name := ""

	this.executable.CheckoutSection(section_name, name)
}

func (this *SetAssigner) WalkSetStmt(stmt_ *stmt.Stmt) {
	set_stmt := stmt_.SetStmt()

	program_counter_expr1 := set_stmt.Expr1().ProgramCounterExpr()
	program_counter_expr2 := set_stmt.Expr2().ProgramCounterExpr()

	primary_expr1 := program_counter_expr1.Expr().PrimaryExpr()
	primary_expr2 := program_counter_expr2.Expr().PrimaryExpr()

	token1 := primary_expr1.Token()
	token2 := primary_expr2.Token()

	attribute1 := token1.Attribute()
	attribute2 := token2.Attribute()

	src_label := this.executable.CurSection().Label(attribute1)

	if this.executable.CurSection().Label(attribute2) == nil {
		this.executable.CurSection().AppendLabel(attribute2)
	}

	dst_label := this.executable.CurSection().Label(attribute2)

	dst_label.SetAddress(src_label.Address())
}

func (this *SetAssigner) WalkTextStmt(stmt_ *stmt.Stmt) {
	section_name := kernel.TEXT
	name := ""

	this.executable.CheckoutSection(section_name, name)
}

func (this *SetAssigner) ConvertSectionName(expr_ *expr.Expr) kernel.SectionName {
	section_name_expr := expr_.SectionNameExpr()
	token_type := section_name_expr.Token().TokenType()

	if token_type == lexer.ATOMIC {
		return kernel.ATOMIC
	} else if token_type == lexer.BSS {
		return kernel.BSS
	} else if token_type == lexer.DATA {
		return kernel.DATA
	} else if token_type == lexer.DEBUG_ABBREV {
		return kernel.DEBUG_ABBREV
	} else if token_type == lexer.DEBUG_FRAME {
		return kernel.DEBUG_FRAME
	} else if token_type == lexer.DEBUG_INFO {
		return kernel.DEBUG_INFO
	} else if token_type == lexer.DEBUG_LINE {
		return kernel.DEBUG_LINE
	} else if token_type == lexer.DEBUG_LOC {
		return kernel.DEBUG_LOC
	} else if token_type == lexer.DEBUG_RANGES {
		return kernel.DEBUG_RANGES
	} else if token_type == lexer.DEBUG_STR {
		return kernel.DEBUG_STR
	} else if token_type == lexer.DPU_HOST {
		return kernel.DPU_HOST
	} else if token_type == lexer.MRAM {
		return kernel.MRAM
	} else if token_type == lexer.RODATA {
		return kernel.RODATA
	} else if token_type == lexer.STACK_SIZES {
		return kernel.STACK_SIZES
	} else if token_type == lexer.TEXT {
		return kernel.TEXT
	} else {
		err := errors.New("section name is not valid")
		panic(err)
	}
}

func (this *SetAssigner) ConvertName(expr_ *expr.Expr) string {
	program_counter_expr := expr_.ProgramCounterExpr()
	primary_expr := program_counter_expr.Expr().PrimaryExpr()
	token := primary_expr.Token()

	if token.TokenType() != lexer.IDENTIFIER {
		err := errors.New("token type is not identifier")
		panic(err)
	}

	attribute := token.Attribute()

	if attribute[0] != '.' {
		err := errors.New("attribute does not start with .")
		panic(err)
	}

	return attribute[1:]
}
