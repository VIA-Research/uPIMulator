package logic

import (
	"errors"
	"strconv"
	"uPIMulator/src/linker/kernel"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser"
	"uPIMulator/src/linker/parser/expr"
	"uPIMulator/src/linker/parser/stmt"
	"uPIMulator/src/misc"
)

type LabelAssigner struct {
	executable *kernel.Executable
	walker     *parser.Walker
}

func (this *LabelAssigner) Init() {
	this.walker = new(parser.Walker)
	this.walker.Init()

	this.walker.RegisterStmtCallback(stmt.ASCII, this.WalkAsciiStmt)
	this.walker.RegisterStmtCallback(stmt.ASCIZ, this.WalkAscizStmt)
	this.walker.RegisterStmtCallback(stmt.BYTE, this.WalkByteStmt)
	this.walker.RegisterStmtCallback(stmt.LONG_PROGRAM_COUNTER, this.WalkLongProgramCounterStmt)
	this.walker.RegisterStmtCallback(stmt.LONG_SECTION_NAME, this.WalkLongSectionNameStmt)
	this.walker.RegisterStmtCallback(stmt.QUAD, this.WalkQuadStmt)
	this.walker.RegisterStmtCallback(
		stmt.SECTION_IDENTIFIER_NUMBER,
		this.WalkSectionIdentifierNumberStmt,
	)
	this.walker.RegisterStmtCallback(stmt.SECTION_IDENTIFIER, this.WalkSectionIdentifierStmt)
	this.walker.RegisterStmtCallback(stmt.SECTION_STACK_SIZES, this.WalkSectionStackSizes)
	this.walker.RegisterStmtCallback(stmt.SECTION_STRING_NUMBER, this.WalkSectionStringNumberStmt)
	this.walker.RegisterStmtCallback(stmt.SECTION_STRING, this.WalkSectionStringStmt)
	this.walker.RegisterStmtCallback(stmt.SHORT, this.WalkShortStmt)
	this.walker.RegisterStmtCallback(stmt.TEXT, this.WalkTextStmt)
	this.walker.RegisterStmtCallback(stmt.ZERO_DOUBLE_NUMBER, this.WalkZeroDoubleNumberStmt)
	this.walker.RegisterStmtCallback(stmt.ZERO_SINGLE_NUMBER, this.WalkZeroSingleNumberStmt)

	this.walker.RegisterStmtCallback(stmt.CI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.DDCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.DMA_RRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.DRDICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.EDRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.ERID, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.ERII, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.ERIR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.ERRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.I, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.NOP, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RIRCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RIRC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RIR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRIC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRRCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRRC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRRICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RRR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.RR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.R, this.WalkInstructionStmt)

	this.walker.RegisterStmtCallback(stmt.S_ERRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RIRCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RIRC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRIC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRCI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRC, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_RR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.S_R, this.WalkInstructionStmt)

	this.walker.RegisterStmtCallback(stmt.BKP, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.BOOT_RI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.CALL_RI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.CALL_RR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.DIV_STEP_DRDI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.JEQ_RII, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.JEQ_RRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.JNZ_RI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.JUMP_I, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.JUMP_R, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.LBS_RRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.LBS_S_RRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.LD_DRI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.MOVD_DD, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_RICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_RI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_S_RICI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_S_RI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.SB_ID_RII, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.SB_ID_RI, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.SB_RIR, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.SD_RID, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.STOP, this.WalkInstructionStmt)
	this.walker.RegisterStmtCallback(stmt.TIME_CFG_R, this.WalkInstructionStmt)

	this.walker.RegisterStmtCallback(stmt.LABEL, this.WalkLabelStmt)
}

func (this *LabelAssigner) Assign(executable *kernel.Executable) {
	this.executable = executable
	this.walker.Walk(executable.Ast())
}

func (this *LabelAssigner) WalkAsciiStmt(stmt_ *stmt.Stmt) {
	ascii_stmt := stmt_.AsciiStmt()
	token := ascii_stmt.Token()
	attribute := token.Attribute()

	// TODO(bongjoon.hyun@gmail.com): decode octal code

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + int64(len(attribute)) - 2)
}

func (this *LabelAssigner) WalkAscizStmt(stmt_ *stmt.Stmt) {
	asciz_stmt := stmt_.AscizStmt()
	token := asciz_stmt.Token()
	attribute := token.Attribute()

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + int64(len(attribute)) - 1)
}

func (this *LabelAssigner) WalkByteStmt(stmt_ *stmt.Stmt) {
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + 1)
}

func (this *LabelAssigner) WalkLongProgramCounterStmt(stmt_ *stmt.Stmt) {
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + 4)
}

func (this *LabelAssigner) WalkLongSectionNameStmt(stmt_ *stmt.Stmt) {
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + 4)
}

func (this *LabelAssigner) WalkQuadStmt(stmt_ *stmt.Stmt) {
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + 8)
}

func (this *LabelAssigner) WalkSectionIdentifierNumberStmt(stmt_ *stmt.Stmt) {
	section_identifier_number_stmt := stmt_.SectionIdentifierNumberStmt()

	section_name := this.ConvertSectionName(section_identifier_number_stmt.Expr1())
	name := this.ConvertName(section_identifier_number_stmt.Expr2())
	section_flags := this.ConvertSectionFlags(section_identifier_number_stmt.Token())
	section_type := this.ConvertSectionType(section_identifier_number_stmt.Expr3())

	if this.executable.Section(section_name, name) == nil {
		this.executable.AddSection(section_name, name, section_flags, section_type)
	}

	this.executable.CheckoutSection(section_name, name)
}

func (this *LabelAssigner) WalkSectionIdentifierStmt(stmt_ *stmt.Stmt) {
	section_identifier_stmt := stmt_.SectionIdentifierStmt()

	section_name := this.ConvertSectionName(section_identifier_stmt.Expr1())
	name := this.ConvertName(section_identifier_stmt.Expr2())
	section_flags := this.ConvertSectionFlags(section_identifier_stmt.Token())
	section_type := this.ConvertSectionType(section_identifier_stmt.Expr3())

	if this.executable.Section(section_name, name) == nil {
		this.executable.AddSection(section_name, name, section_flags, section_type)
	}

	this.executable.CheckoutSection(section_name, name)
}

func (this *LabelAssigner) WalkSectionStackSizes(stmt_ *stmt.Stmt) {
	section_stack_sizes_stmt := stmt_.SectionStackSizesStmt()

	section_name := kernel.STACK_SIZES
	section_flags := this.ConvertSectionFlags(section_stack_sizes_stmt.Token())
	section_type := this.ConvertSectionType(section_stack_sizes_stmt.Expr1())

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

	if this.executable.Section(section_name, name) == nil {
		this.executable.AddSection(section_name, name, section_flags, section_type)
	}

	this.executable.CheckoutSection(section_name, name)
}

func (this *LabelAssigner) WalkSectionStringNumberStmt(stmt_ *stmt.Stmt) {
	section_string_number_stmt := stmt_.SectionStringNumberStmt()

	section_name := this.ConvertSectionName(section_string_number_stmt.Expr1())
	name := ""
	section_flags := this.ConvertSectionFlags(section_string_number_stmt.Token())
	section_type := this.ConvertSectionType(section_string_number_stmt.Expr2())

	if this.executable.Section(section_name, name) == nil {
		this.executable.AddSection(section_name, name, section_flags, section_type)
	}

	this.executable.CheckoutSection(section_name, name)
}

func (this *LabelAssigner) WalkSectionStringStmt(stmt_ *stmt.Stmt) {
	section_string_stmt := stmt_.SectionStringStmt()

	section_name := this.ConvertSectionName(section_string_stmt.Expr1())
	name := ""
	section_flags := this.ConvertSectionFlags(section_string_stmt.Token())
	section_type := this.ConvertSectionType(section_string_stmt.Expr2())

	if this.executable.Section(section_name, name) == nil {
		this.executable.AddSection(section_name, name, section_flags, section_type)
	}

	this.executable.CheckoutSection(section_name, name)
}

func (this *LabelAssigner) WalkShortStmt(stmt_ *stmt.Stmt) {
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + 2)
}

func (this *LabelAssigner) WalkTextStmt(stmt_ *stmt.Stmt) {
	section_name := kernel.TEXT
	name := ""
	section_flags := make(map[kernel.SectionFlag]bool, 0)
	section_type := kernel.PROGBITS

	if this.executable.Section(section_name, name) == nil {
		this.executable.AddSection(section_name, name, section_flags, section_type)
	}

	this.executable.CheckoutSection(section_name, name)
}

func (this *LabelAssigner) WalkZeroDoubleNumberStmt(stmt_ *stmt.Stmt) {
	zero_double_number_stmt := stmt_.ZeroDoubleNumberStmt()

	program_counter_expr := zero_double_number_stmt.Expr1().ProgramCounterExpr()
	primary_expr := program_counter_expr.Expr().PrimaryExpr()
	token := primary_expr.Token()
	attribute := token.Attribute()

	size, err := strconv.ParseInt(attribute, 10, 64)
	if err != nil {
		panic(err)
	}

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + size)
}

func (this *LabelAssigner) WalkZeroSingleNumberStmt(stmt_ *stmt.Stmt) {
	zero_single_number_stmt := stmt_.ZeroSingleNumberStmt()

	program_counter_expr := zero_single_number_stmt.Expr().ProgramCounterExpr()
	primary_expr := program_counter_expr.Expr().PrimaryExpr()
	token := primary_expr.Token()
	attribute := token.Attribute()

	size, err := strconv.ParseInt(attribute, 10, 64)
	if err != nil {
		panic(err)
	}

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + size)
}

func (this *LabelAssigner) WalkInstructionStmt(stmt_ *stmt.Stmt) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	iram_data_width := config_loader.IramDataWidth()
	instruction_size := int64(iram_data_width / 8)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.SetSize(cur_label.Size() + instruction_size)
}

func (this *LabelAssigner) WalkLabelStmt(stmt_ *stmt.Stmt) {
	label_stmt := stmt_.LabelStmt()

	program_counter_expr := label_stmt.Expr().ProgramCounterExpr()
	primary_expr := program_counter_expr.Expr().PrimaryExpr()
	token := primary_expr.Token()
	label_name := token.Attribute()

	if label_name != "__sys_used_mram_end" {
		if this.executable.CurSection().Label(label_name) == nil {
			this.executable.CurSection().AppendLabel(label_name)
		}

		this.executable.CurSection().CheckoutLabel(label_name)
	}
}

func (this *LabelAssigner) ConvertSectionName(expr_ *expr.Expr) kernel.SectionName {
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

func (this *LabelAssigner) ConvertName(expr_ *expr.Expr) string {
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

func (this *LabelAssigner) ConvertSectionFlags(token *lexer.Token) map[kernel.SectionFlag]bool {
	attribute := token.Attribute()

	section_flags := make(map[kernel.SectionFlag]bool, 0)
	for i := 1; i < len(attribute)-1; i++ {
		if attribute[i] == 'a' {
			section_flags[kernel.ALLOC] = true
		} else if attribute[i] == 'w' {
			section_flags[kernel.WRITE] = true
		} else if attribute[i] == 'x' {
			section_flags[kernel.EXECINSTR] = true
		} else if attribute[i] == 'o' {
			section_flags[kernel.LINK_ORDER] = true
		} else if attribute[i] == 'M' {
			section_flags[kernel.MERGE] = true
		} else if attribute[i] == 'S' {
			section_flags[kernel.STRINGS] = true
		} else {
			err := errors.New("section flag is not valid")
			panic(err)
		}
	}
	return section_flags
}

func (this *LabelAssigner) ConvertSectionType(expr_ *expr.Expr) kernel.SectionType {
	section_type_expr := expr_.SectionTypeExpr()
	token := section_type_expr.Token()
	token_type := token.TokenType()

	if token_type == lexer.PROGBITS {
		return kernel.PROGBITS
	} else if token_type == lexer.NOBITS {
		return kernel.NOBITS
	} else {
		err := errors.New("section type is not valid")
		panic(err)
	}
}
