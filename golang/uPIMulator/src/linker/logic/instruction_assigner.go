package logic

import (
	"errors"
	"strconv"
	"uPIMulator/src/linker/kernel"
	"uPIMulator/src/linker/kernel/directive"
	"uPIMulator/src/linker/kernel/instruction"
	"uPIMulator/src/linker/kernel/instruction/cc"
	"uPIMulator/src/linker/kernel/instruction/reg_descriptor"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser"
	"uPIMulator/src/linker/parser/expr"
	"uPIMulator/src/linker/parser/stmt"
	"uPIMulator/src/misc"
)

type InstructionAssigner struct {
	executable    *kernel.Executable
	walker        *parser.Walker
	linker_script *LinkerScript
}

func (this *InstructionAssigner) Init(linker_script *LinkerScript) {
	this.linker_script = linker_script

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

	this.walker.RegisterStmtCallback(stmt.CI, this.WalkCiStmt)
	this.walker.RegisterStmtCallback(stmt.DMA_RRI, this.WalkDmaRriStmt)
	this.walker.RegisterStmtCallback(stmt.DRDICI, this.WalkDrdiciStmt)
	this.walker.RegisterStmtCallback(stmt.EDRI, this.WalkEdriStmt)
	this.walker.RegisterStmtCallback(stmt.ERID, this.WalkEridStmt)
	this.walker.RegisterStmtCallback(stmt.ERII, this.WalkEriiStmt)
	this.walker.RegisterStmtCallback(stmt.ERIR, this.WalkErirStmt)
	this.walker.RegisterStmtCallback(stmt.ERRI, this.WalkErriStmt)
	this.walker.RegisterStmtCallback(stmt.I, this.WalkIStmt)
	this.walker.RegisterStmtCallback(stmt.NOP, this.WalkNopStmt)
	this.walker.RegisterStmtCallback(stmt.RCI, this.WalkRciStmt)
	this.walker.RegisterStmtCallback(stmt.RICI, this.WalkRiciStmt)
	this.walker.RegisterStmtCallback(stmt.RIRCI, this.WalkRirciStmt)
	this.walker.RegisterStmtCallback(stmt.RIRC, this.WalkRircStmt)
	this.walker.RegisterStmtCallback(stmt.RIR, this.WalkRirStmt)
	this.walker.RegisterStmtCallback(stmt.RRCI, this.WalkRrciStmt)
	this.walker.RegisterStmtCallback(stmt.RRC, this.WalkRrcStmt)
	this.walker.RegisterStmtCallback(stmt.RRICI, this.WalkRriciStmt)
	this.walker.RegisterStmtCallback(stmt.RRIC, this.WalkRricStmt)
	this.walker.RegisterStmtCallback(stmt.RRI, this.WalkRriStmt)
	this.walker.RegisterStmtCallback(stmt.RRRCI, this.WalkRrrciStmt)
	this.walker.RegisterStmtCallback(stmt.RRRC, this.WalkRrrcStmt)
	this.walker.RegisterStmtCallback(stmt.RRRICI, this.WalkRrriciStmt)
	this.walker.RegisterStmtCallback(stmt.RRRI, this.WalkRrriStmt)
	this.walker.RegisterStmtCallback(stmt.RRR, this.WalkRrrStmt)
	this.walker.RegisterStmtCallback(stmt.RR, this.WalkRrStmt)
	this.walker.RegisterStmtCallback(stmt.R, this.WalkRStmt)

	this.walker.RegisterStmtCallback(stmt.S_ERRI, this.WalkSErriStmt)
	this.walker.RegisterStmtCallback(stmt.S_RCI, this.WalkSRciStmt)
	this.walker.RegisterStmtCallback(stmt.S_RIRCI, this.WalkSRirciStmt)
	this.walker.RegisterStmtCallback(stmt.S_RIRC, this.WalkSRircStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRCI, this.WalkSRrciStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRC, this.WalkSRrcStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRICI, this.WalkSRriciStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRIC, this.WalkSRricStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRI, this.WalkSRriStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRCI, this.WalkSRrrciStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRC, this.WalkSRrrcStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRICI, this.WalkSRrriciStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRRI, this.WalkSRrriStmt)
	this.walker.RegisterStmtCallback(stmt.S_RRR, this.WalkSRrrStmt)
	this.walker.RegisterStmtCallback(stmt.S_RR, this.WalkSRrStmt)
	this.walker.RegisterStmtCallback(stmt.S_R, this.WalkSRStmt)

	this.walker.RegisterStmtCallback(stmt.BKP, this.WalkBkpStmt)
	this.walker.RegisterStmtCallback(stmt.BOOT_RI, this.WalkBootRiStmt)
	this.walker.RegisterStmtCallback(stmt.CALL_RI, this.WalkCallRiStmt)
	this.walker.RegisterStmtCallback(stmt.CALL_RR, this.WalkCallRrStmt)
	this.walker.RegisterStmtCallback(stmt.DIV_STEP_DRDI, this.WalkDivStepDrdiStmt)
	this.walker.RegisterStmtCallback(stmt.JEQ_RII, this.WalkJeqRiiStmt)
	this.walker.RegisterStmtCallback(stmt.JEQ_RRI, this.WalkJeqRriStmt)
	this.walker.RegisterStmtCallback(stmt.JNZ_RI, this.WalkJnzRiStmt)
	this.walker.RegisterStmtCallback(stmt.JUMP_I, this.WalkJumpIStmt)
	this.walker.RegisterStmtCallback(stmt.JUMP_R, this.WalkJumpRStmt)
	this.walker.RegisterStmtCallback(stmt.LBS_RRI, this.WalkLbsRriStmt)
	this.walker.RegisterStmtCallback(stmt.LBS_S_RRI, this.WalkLbsSRriStmt)
	this.walker.RegisterStmtCallback(stmt.LD_DRI, this.WalkLdDriStmt)
	this.walker.RegisterStmtCallback(stmt.MOVD_DD, this.WalkMovdDdStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_RICI, this.WalkMoveRiciStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_RI, this.WalkMoveRiStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_S_RICI, this.WalkMoveSRiciStmt)
	this.walker.RegisterStmtCallback(stmt.MOVE_S_RI, this.WalkMoveSRiStmt)
	this.walker.RegisterStmtCallback(stmt.SB_ID_RII, this.WalkSbIdRiiStmt)
	this.walker.RegisterStmtCallback(stmt.SB_ID_RI, this.WalkSbIdRiStmt)
	this.walker.RegisterStmtCallback(stmt.SB_RIR, this.WalkSbRirStmt)
	this.walker.RegisterStmtCallback(stmt.SD_RID, this.WalkSdRidStmt)
	this.walker.RegisterStmtCallback(stmt.STOP, this.WalkStopStmt)
	this.walker.RegisterStmtCallback(stmt.TIME_CFG_R, this.WalkTimeCfgRStmt)

	this.walker.RegisterStmtCallback(stmt.LABEL, this.WalkLabelStmt)
}

func (this *InstructionAssigner) Assign(executable *kernel.Executable) {
	this.executable = executable
	this.walker.Walk(executable.Ast())
}

func (this *InstructionAssigner) WalkAsciiStmt(stmt_ *stmt.Stmt) {
	ascii_stmt := stmt_.AsciiStmt()
	token := ascii_stmt.Token()
	attribute := token.Attribute()
	characters := attribute[1 : len(attribute)-1]

	// TODO(bongjoon.hyun@gmail.com): decode octal code

	ascii_directive := new(directive.AsciiDirective)
	ascii_directive.Init(characters)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(ascii_directive)
}

func (this *InstructionAssigner) WalkAscizStmt(stmt_ *stmt.Stmt) {
	asciz_stmt := stmt_.AscizStmt()
	token := asciz_stmt.Token()
	attribute := token.Attribute()
	characters := attribute[1 : len(attribute)-1]

	asciz_directive := new(directive.AscizDirective)
	asciz_directive.Init(characters)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(asciz_directive)
}

func (this *InstructionAssigner) WalkByteStmt(stmt_ *stmt.Stmt) {
	byte_stmt := stmt_.ByteStmt()

	value := this.EvaluateProgramCounter(byte_stmt.Expr())

	byte_directive := new(directive.ByteDirective)
	byte_directive.Init(value)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(byte_directive)
}

func (this *InstructionAssigner) WalkLongProgramCounterStmt(stmt_ *stmt.Stmt) {
	long_program_counter_stmt := stmt_.LongProgramCounterStmt()

	value := this.EvaluateProgramCounter(long_program_counter_stmt.Expr())

	long_directive := new(directive.LongDirective)
	long_directive.Init(value)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(long_directive)
}

func (this *InstructionAssigner) WalkLongSectionNameStmt(stmt_ *stmt.Stmt) {
	long_section_name_stmt := stmt_.LongSectionNameStmt()

	value := this.EvaluateSectionName(long_section_name_stmt.Expr())

	long_directive := new(directive.LongDirective)
	long_directive.Init(value)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(long_directive)
}

func (this *InstructionAssigner) WalkQuadStmt(stmt_ *stmt.Stmt) {
	quad_stmt := stmt_.QuadStmt()

	value := this.EvaluateProgramCounter(quad_stmt.Expr())

	quad_directive := new(directive.QuadDirective)
	quad_directive.Init(value)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(quad_directive)
}

func (this *InstructionAssigner) WalkSectionIdentifierNumberStmt(stmt_ *stmt.Stmt) {
	section_identifier_number_stmt := stmt_.SectionIdentifierNumberStmt()

	section_name := this.ConvertSectionName(section_identifier_number_stmt.Expr1())
	name := this.ConvertName(section_identifier_number_stmt.Expr2())

	this.executable.CheckoutSection(section_name, name)
}

func (this *InstructionAssigner) WalkSectionIdentifierStmt(stmt_ *stmt.Stmt) {
	section_identifier_stmt := stmt_.SectionIdentifierStmt()

	section_name := this.ConvertSectionName(section_identifier_stmt.Expr1())
	name := this.ConvertName(section_identifier_stmt.Expr2())

	this.executable.CheckoutSection(section_name, name)
}

func (this *InstructionAssigner) WalkSectionStackSizes(stmt_ *stmt.Stmt) {
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

func (this *InstructionAssigner) WalkSectionStringNumberStmt(stmt_ *stmt.Stmt) {
	section_string_number_stmt := stmt_.SectionStringNumberStmt()

	section_name := this.ConvertSectionName(section_string_number_stmt.Expr1())
	name := ""

	this.executable.CheckoutSection(section_name, name)
}

func (this *InstructionAssigner) WalkSectionStringStmt(stmt_ *stmt.Stmt) {
	section_string_stmt := stmt_.SectionStringStmt()

	section_name := this.ConvertSectionName(section_string_stmt.Expr1())
	name := ""

	this.executable.CheckoutSection(section_name, name)
}

func (this *InstructionAssigner) WalkShortStmt(stmt_ *stmt.Stmt) {
	short_stmt := stmt_.ShortStmt()

	value := this.EvaluateProgramCounter(short_stmt.Expr())

	short_directive := new(directive.ShortDirective)
	short_directive.Init(value)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(short_directive)
}

func (this *InstructionAssigner) WalkTextStmt(stmt_ *stmt.Stmt) {
	section_name := kernel.TEXT
	name := ""

	this.executable.CheckoutSection(section_name, name)
}

func (this *InstructionAssigner) WalkZeroDoubleNumberStmt(stmt_ *stmt.Stmt) {
	zero_double_number_stmt := stmt_.ZeroDoubleNumberStmt()

	size := this.EvaluateProgramCounter(zero_double_number_stmt.Expr1())
	value := this.EvaluateProgramCounter(zero_double_number_stmt.Expr2())

	zero_directive := new(directive.ZeroDirective)
	zero_directive.Init(size, value)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(zero_directive)
}

func (this *InstructionAssigner) WalkZeroSingleNumberStmt(stmt_ *stmt.Stmt) {
	zero_single_number_stmt := stmt_.ZeroSingleNumberStmt()

	size := this.EvaluateProgramCounter(zero_single_number_stmt.Expr())

	zero_directive := new(directive.ZeroDirective)
	zero_directive.Init(size, 0)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(zero_directive)
}

func (this *InstructionAssigner) WalkCiStmt(stmt_ *stmt.Stmt) {
	ci_stmt := stmt_.CiStmt()

	op_code := this.ConvertCiOpCode(ci_stmt.OpCode())
	condition := this.ConvertCondition(ci_stmt.Condition())
	pc := this.EvaluateProgramCounter(ci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitCi(op_code, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkDdciStmt(stmt_ *stmt.Stmt) {
	ddci_stmt := stmt_.DdciStmt()

	op_code := this.ConvertDdciOpCode(ddci_stmt.OpCode())
	dc := this.ConvertPairReg(ddci_stmt.Dc())
	db := this.ConvertPairReg(ddci_stmt.Db())
	condition := this.ConvertCondition(ddci_stmt.Condition())
	pc := this.EvaluateProgramCounter(ddci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitDdci(op_code, dc, db, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkDmaRriStmt(stmt_ *stmt.Stmt) {
	dma_rri_stmt := stmt_.DmaRriStmt()

	op_code := this.ConvertDmaRriOpCode(dma_rri_stmt.OpCode())
	ra := this.ConvertSrcReg(dma_rri_stmt.Ra())
	rb := this.ConvertSrcReg(dma_rri_stmt.Rb())
	imm := this.EvaluateProgramCounter(dma_rri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitDmaRri(op_code, ra, rb, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkDrdiciStmt(stmt_ *stmt.Stmt) {
	drdici_stmt := stmt_.DrdiciStmt()

	op_code := this.ConvertDrdiciOpCode(drdici_stmt.OpCode())
	dc := this.ConvertPairReg(drdici_stmt.Dc())
	ra := this.ConvertSrcReg(drdici_stmt.Ra())
	db := this.ConvertPairReg(drdici_stmt.Db())
	imm := this.EvaluateProgramCounter(drdici_stmt.Imm())
	condition := this.ConvertCondition(drdici_stmt.Condition())
	pc := this.EvaluateProgramCounter(drdici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitDrdici(op_code, dc, ra, db, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkEdriStmt(stmt_ *stmt.Stmt) {
	edri_stmt := stmt_.EdriStmt()

	op_code := this.ConvertLoadOpCode(edri_stmt.OpCode())
	endian := this.ConvertEndian(edri_stmt.Endian())
	dc := this.ConvertPairReg(edri_stmt.Dc())
	ra := this.ConvertSrcReg(edri_stmt.Ra())
	off := this.EvaluateProgramCounter(edri_stmt.Off())

	instruction_ := new(instruction.Instruction)
	instruction_.InitEdri(op_code, endian, dc, ra, off)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkEridStmt(stmt_ *stmt.Stmt) {
	erid_stmt := stmt_.EridStmt()

	op_code := this.ConvertStoreOpCode(erid_stmt.OpCode())
	endian := this.ConvertEndian(erid_stmt.Endian())
	ra := this.ConvertSrcReg(erid_stmt.Ra())
	off := this.EvaluateProgramCounter(erid_stmt.Off())
	db := this.ConvertPairReg(erid_stmt.Db())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErid(op_code, endian, ra, off, db)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkEriiStmt(stmt_ *stmt.Stmt) {
	erii_stmt := stmt_.EriiStmt()

	op_code := this.ConvertStoreOpCode(erii_stmt.OpCode())
	endian := this.ConvertEndian(erii_stmt.Endian())
	ra := this.ConvertSrcReg(erii_stmt.Ra())
	off := this.EvaluateProgramCounter(erii_stmt.Off())
	imm := this.EvaluateProgramCounter(erii_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErii(op_code, endian, ra, off, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkErirStmt(stmt_ *stmt.Stmt) {
	erir_stmt := stmt_.ErirStmt()

	op_code := this.ConvertStoreOpCode(erir_stmt.OpCode())
	endian := this.ConvertEndian(erir_stmt.Endian())
	ra := this.ConvertSrcReg(erir_stmt.Ra())
	off := this.EvaluateProgramCounter(erir_stmt.Off())
	rb := this.ConvertSrcReg(erir_stmt.Rb())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErir(op_code, endian, ra, off, rb)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkErriStmt(stmt_ *stmt.Stmt) {
	erri_stmt := stmt_.ErriStmt()

	op_code := this.ConvertLoadOpCode(erri_stmt.OpCode())
	endian := this.ConvertEndian(erri_stmt.Endian())
	rc := this.ConvertGpReg(erri_stmt.Rc())
	ra := this.ConvertSrcReg(erri_stmt.Ra())
	off := this.EvaluateProgramCounter(erri_stmt.Off())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErri(op_code, endian, rc, ra, off)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkIStmt(stmt_ *stmt.Stmt) {
	i_stmt := stmt_.IStmt()

	op_code := this.ConvertIOpCode(i_stmt.OpCode())
	imm := this.EvaluateProgramCounter(i_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitI(op_code, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkNopStmt(stmt_ *stmt.Stmt) {
	nop_stmt := stmt_.NopStmt()

	op_code := this.ConvertROpCode(nop_stmt.OpCode())

	instruction_ := new(instruction.Instruction)
	instruction_.InitZ(op_code)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRciStmt(stmt_ *stmt.Stmt) {
	rci_stmt := stmt_.RciStmt()

	is_zero_reg := this.IsZeroReg(rci_stmt.Rc())

	op_code := this.ConvertROpCode(rci_stmt.OpCode())
	condition := this.ConvertCondition(rci_stmt.Condition())
	pc := this.EvaluateProgramCounter(rci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rci_stmt.Rc())

		instruction_.InitRci(op_code, rc, condition, pc)
	} else {
		instruction_.InitZci(op_code, condition, pc)
	}

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRiciStmt(stmt_ *stmt.Stmt) {
	rici_stmt := stmt_.RiciStmt()

	op_code := this.ConvertRiciOpCode(rici_stmt.OpCode())
	ra := this.ConvertSrcReg(rici_stmt.Ra())
	imm := this.EvaluateProgramCounter(rici_stmt.Imm())
	condition := this.ConvertCondition(rici_stmt.Condition())
	pc := this.EvaluateProgramCounter(rici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitRici(op_code, ra, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRirciStmt(stmt_ *stmt.Stmt) {
	rirci_stmt := stmt_.RirciStmt()

	is_zero_reg := this.IsZeroReg(rirci_stmt.Rc())

	op_code := this.ConvertRriOpCode(rirci_stmt.OpCode())
	imm := this.EvaluateProgramCounter(rirci_stmt.Imm())
	ra := this.ConvertSrcReg(rirci_stmt.Ra())
	condition := this.ConvertCondition(rirci_stmt.Condition())
	pc := this.EvaluateProgramCounter(rirci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rirci_stmt.Rc())

		instruction_.InitRirci(op_code, rc, imm, ra, condition, pc)
	} else {
		instruction_.InitZirci(op_code, imm, ra, condition, pc)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRircStmt(stmt_ *stmt.Stmt) {
	rirc_stmt := stmt_.RircStmt()

	is_zero_reg := this.IsZeroReg(rirc_stmt.Rc())

	op_code := this.ConvertRriOpCode(rirc_stmt.OpCode())
	imm := this.EvaluateProgramCounter(rirc_stmt.Imm())
	ra := this.ConvertSrcReg(rirc_stmt.Ra())
	condition := this.ConvertCondition(rirc_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rirc_stmt.Rc())

		instruction_.InitRirc(op_code, rc, imm, ra, condition)
	} else {
		instruction_.InitZirc(op_code, imm, ra, condition)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRirStmt(stmt_ *stmt.Stmt) {
	rir_stmt := stmt_.RirStmt()

	is_zero_reg := this.IsZeroReg(rir_stmt.Rc())

	op_code := this.ConvertRriOpCode(rir_stmt.OpCode())
	imm := this.EvaluateProgramCounter(rir_stmt.Imm())
	ra := this.ConvertSrcReg(rir_stmt.Ra())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rir_stmt.Rc())

		instruction_.InitRir(op_code, rc, imm, ra)
	} else {
		instruction_.InitZir(op_code, imm, ra)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrciStmt(stmt_ *stmt.Stmt) {
	rrci_stmt := stmt_.RrciStmt()

	op_code := this.ConvertRrOpCode(rrci_stmt.OpCode())
	ra := this.ConvertSrcReg(rrci_stmt.Ra())
	condition := this.ConvertCondition(rrci_stmt.Condition())
	pc := this.EvaluateProgramCounter(rrci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	if op_code == instruction.OR {
		rc := this.ConvertGpReg(rrci_stmt.Rc())
		imm := int64(0)

		instruction_.InitRrici(op_code, rc, ra, imm, condition, pc)
	} else if op_code == instruction.SUB {
		rc := this.ConvertGpReg(rrci_stmt.Rc())
		imm := int64(0)

		instruction_.InitRirci(op_code, rc, imm, ra, condition, pc)
	} else if op_code == instruction.XOR {
		is_zero_reg := this.IsZeroReg(rrci_stmt.Rc())

		if !is_zero_reg {
			rc := this.ConvertGpReg(rrci_stmt.Rc())
			imm := int64(-1)

			instruction_.InitRrici(op_code, rc, ra, imm, condition, pc)
		} else {
			imm := int64(-1)

			instruction_.InitZrici(op_code, ra, imm, condition, pc)
		}
	} else {
		is_zero_reg := this.IsZeroReg(rrci_stmt.Rc())

		if !is_zero_reg {
			rc := this.ConvertGpReg(rrci_stmt.Rc())

			instruction_.InitRrci(op_code, rc, ra, condition, pc)
		} else {
			instruction_.InitZrci(op_code, ra, condition, pc)
		}
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrcStmt(stmt_ *stmt.Stmt) {
	rrc_stmt := stmt_.RrcStmt()

	is_zero_reg := this.IsZeroReg(rrc_stmt.Rc())

	op_code := this.ConvertRrOpCode(rrc_stmt.OpCode())
	ra := this.ConvertSrcReg(rrc_stmt.Ra())
	condition := this.ConvertCondition(rrc_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rrc_stmt.Rc())

		instruction_.InitRrc(op_code, rc, ra, condition)
	} else {
		instruction_.InitZrc(op_code, ra, condition)
	}

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRriciStmt(stmt_ *stmt.Stmt) {
	rrici_stmt := stmt_.RriciStmt()

	is_zero_reg := this.IsZeroReg(rrici_stmt.Rc())

	op_code := this.ConvertRriOpCode(rrici_stmt.OpCode())
	ra := this.ConvertSrcReg(rrici_stmt.Ra())
	imm := this.EvaluateProgramCounter(rrici_stmt.Imm())
	condition := this.ConvertCondition(rrici_stmt.Condition())
	pc := this.EvaluateProgramCounter(rrici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rrici_stmt.Rc())

		instruction_.InitRrici(op_code, rc, ra, imm, condition, pc)
	} else {
		instruction_.InitZrici(op_code, ra, imm, condition, pc)
	}

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRricStmt(stmt_ *stmt.Stmt) {
	rric_stmt := stmt_.RricStmt()

	is_zero_reg := this.IsZeroReg(rric_stmt.Rc())

	op_code := this.ConvertRriOpCode(rric_stmt.OpCode())
	ra := this.ConvertSrcReg(rric_stmt.Ra())
	imm := this.EvaluateProgramCounter(rric_stmt.Imm())
	condition := this.ConvertCondition(rric_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rric_stmt.Rc())

		if condition != cc.FALSE {
			instruction_.InitRric(op_code, rc, ra, imm, condition)
		} else {
			instruction_.InitRrif(op_code, rc, ra, imm, condition)
		}
	} else {
		if condition != cc.FALSE {
			instruction_.InitZric(op_code, ra, imm, condition)
		} else {
			instruction_.InitZrif(op_code, ra, imm, condition)
		}
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRriStmt(stmt_ *stmt.Stmt) {
	rri_stmt := stmt_.RriStmt()

	is_zero_reg := this.IsZeroReg(rri_stmt.Rc())

	op_code := this.ConvertRriOpCode(rri_stmt.OpCode())
	ra := this.ConvertSrcReg(rri_stmt.Ra())
	imm := this.EvaluateProgramCounter(rri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	if op_code == instruction.ANDN || op_code == instruction.NAND || op_code == instruction.NOR ||
		op_code == instruction.NXOR ||
		op_code == instruction.ORN {
		rc := this.ConvertGpReg(rri_stmt.Rc())
		condition := cc.FALSE

		instruction_.InitRrif(op_code, rc, ra, imm, condition)
	} else if !is_zero_reg {
		rc := this.ConvertGpReg(rri_stmt.Rc())

		instruction_.InitRri(op_code, rc, ra, imm)
	} else {
		instruction_.InitZri(op_code, ra, imm)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrrciStmt(stmt_ *stmt.Stmt) {
	rrrci_stmt := stmt_.RrrciStmt()

	is_zero_reg := this.IsZeroReg(rrrci_stmt.Rc())

	op_code := this.ConvertRriOpCode(rrrci_stmt.OpCode())
	ra := this.ConvertSrcReg(rrrci_stmt.Ra())
	rb := this.ConvertSrcReg(rrrci_stmt.Rb())
	condition := this.ConvertCondition(rrrci_stmt.Condition())
	pc := this.EvaluateProgramCounter(rrrci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rrrci_stmt.Rc())

		instruction_.InitRrrci(op_code, rc, ra, rb, condition, pc)
	} else {
		instruction_.InitZrrci(op_code, ra, rb, condition, pc)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrrcStmt(stmt_ *stmt.Stmt) {
	rrrc_stmt := stmt_.RrrcStmt()

	is_zero_reg := this.IsZeroReg(rrrc_stmt.Rc())

	op_code := this.ConvertRriOpCode(rrrc_stmt.OpCode())
	ra := this.ConvertSrcReg(rrrc_stmt.Ra())
	rb := this.ConvertSrcReg(rrrc_stmt.Rb())
	condition := this.ConvertCondition(rrrc_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rrrc_stmt.Rc())

		instruction_.InitRrrc(op_code, rc, ra, rb, condition)
	} else {
		instruction_.InitZrrc(op_code, ra, rb, condition)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrriciStmt(stmt_ *stmt.Stmt) {
	rrrici_stmt := stmt_.RrriciStmt()

	is_zero_reg := this.IsZeroReg(rrrici_stmt.Rc())

	op_code := this.ConvertRrriOpCode(rrrici_stmt.OpCode())
	ra := this.ConvertSrcReg(rrrici_stmt.Ra())
	rb := this.ConvertSrcReg(rrrici_stmt.Rb())
	imm := this.EvaluateProgramCounter(rrrici_stmt.Imm())
	condition := this.ConvertCondition(rrrici_stmt.Condition())
	pc := this.EvaluateProgramCounter(rrrici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rrrici_stmt.Rc())

		instruction_.InitRrrici(op_code, rc, ra, rb, imm, condition, pc)
	} else {
		instruction_.InitZrrici(op_code, ra, rb, imm, condition, pc)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrriStmt(stmt_ *stmt.Stmt) {
	rrri_stmt := stmt_.RrriStmt()

	is_zero_reg := this.IsZeroReg(rrri_stmt.Rc())

	op_code := this.ConvertRrriOpCode(rrri_stmt.OpCode())
	ra := this.ConvertSrcReg(rrri_stmt.Ra())
	rb := this.ConvertSrcReg(rrri_stmt.Rb())
	imm := this.EvaluateProgramCounter(rrri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rrri_stmt.Rc())

		instruction_.InitRrri(op_code, rc, ra, rb, imm)
	} else {
		instruction_.InitZrri(op_code, ra, rb, imm)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrrStmt(stmt_ *stmt.Stmt) {
	rrr_stmt := stmt_.RrrStmt()

	is_zero_reg := this.IsZeroReg(rrr_stmt.Rc())

	op_code := this.ConvertRriOpCode(rrr_stmt.OpCode())
	ra := this.ConvertSrcReg(rrr_stmt.Ra())
	rb := this.ConvertSrcReg(rrr_stmt.Rb())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(rrr_stmt.Rc())

		instruction_.InitRrr(op_code, rc, ra, rb)
	} else {
		instruction_.InitZrr(op_code, ra, rb)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRrStmt(stmt_ *stmt.Stmt) {
	rr_stmt := stmt_.RrStmt()

	op_code := this.ConvertRrOpCode(rr_stmt.OpCode())
	ra := this.ConvertSrcReg(rr_stmt.Ra())

	instruction_ := new(instruction.Instruction)
	if op_code == instruction.OR {
		rc := this.ConvertGpReg(rr_stmt.Rc())
		imm := int64(0)
		condition := cc.FALSE

		instruction_.InitRrif(op_code, rc, ra, imm, condition)
	} else if op_code == instruction.SUB {
		rc := this.ConvertGpReg(rr_stmt.Rc())
		imm := int64(0)

		instruction_.InitRir(op_code, rc, imm, ra)
	} else if op_code == instruction.XOR {
		is_zero_reg := this.IsZeroReg(rr_stmt.Rc())

		if !is_zero_reg {
			rc := this.ConvertGpReg(rr_stmt.Rc())
			imm := int64(-1)

			instruction_.InitRri(op_code, rc, ra, imm)
		} else {
			imm := int64(-1)

			instruction_.InitZri(op_code, ra, imm)
		}
	} else {
		is_zero_reg := this.IsZeroReg(rr_stmt.Rc())

		if !is_zero_reg {
			rc := this.ConvertGpReg(rr_stmt.Rc())

			instruction_.InitRr(op_code, rc, ra)
		} else {
			instruction_.InitZr(op_code, ra)
		}
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkRStmt(stmt_ *stmt.Stmt) {
	r_stmt := stmt_.RStmt()

	is_zero_reg := this.IsZeroReg(r_stmt.Rc())

	op_code := this.ConvertROpCode(r_stmt.OpCode())

	instruction_ := new(instruction.Instruction)
	if !is_zero_reg {
		rc := this.ConvertGpReg(r_stmt.Rc())

		instruction_.InitR(op_code, rc)
	} else {
		instruction_.InitZ(op_code)
	}

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSErriStmt(stmt_ *stmt.Stmt) {
	s_erri_stmt := stmt_.SErriStmt()

	op_code := this.ConvertLoadOpCode(s_erri_stmt.OpCode())
	suffix := this.ConvertSuffix(s_erri_stmt.Suffix(), instruction.ERRI)
	endian := this.ConvertEndian(s_erri_stmt.Endian())
	dc := this.ConvertPairReg(s_erri_stmt.Dc())
	ra := this.ConvertSrcReg(s_erri_stmt.Ra())
	off := this.EvaluateProgramCounter(s_erri_stmt.Off())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSErri(op_code, suffix, endian, dc, ra, off)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRciStmt(stmt_ *stmt.Stmt) {
	s_rci_stmt := stmt_.SRciStmt()

	op_code := this.ConvertROpCode(s_rci_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rci_stmt.Suffix(), instruction.RCI)
	dc := this.ConvertPairReg(s_rci_stmt.Dc())
	condition := this.ConvertCondition(s_rci_stmt.Condition())
	pc := this.EvaluateProgramCounter(s_rci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRci(op_code, suffix, dc, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRirciStmt(stmt_ *stmt.Stmt) {
	s_rirci_stmt := stmt_.SRirciStmt()

	op_code := this.ConvertRriOpCode(s_rirci_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rirci_stmt.Suffix(), instruction.RIRCI)
	dc := this.ConvertPairReg(s_rirci_stmt.Dc())
	imm := this.EvaluateProgramCounter(s_rirci_stmt.Imm())
	ra := this.ConvertSrcReg(s_rirci_stmt.Ra())
	condition := this.ConvertCondition(s_rirci_stmt.Condition())
	pc := this.EvaluateProgramCounter(s_rirci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRirci(op_code, suffix, dc, imm, ra, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRircStmt(stmt_ *stmt.Stmt) {
	s_rirc_stmt := stmt_.SRircStmt()

	op_code := this.ConvertRriOpCode(s_rirc_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rirc_stmt.Suffix(), instruction.RIRC)
	dc := this.ConvertPairReg(s_rirc_stmt.Dc())
	imm := this.EvaluateProgramCounter(s_rirc_stmt.Imm())
	ra := this.ConvertSrcReg(s_rirc_stmt.Ra())
	condition := this.ConvertCondition(s_rirc_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRirc(op_code, suffix, dc, imm, ra, condition)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrciStmt(stmt_ *stmt.Stmt) {
	s_rrci_stmt := stmt_.SRrciStmt()

	op_code := this.ConvertRrOpCode(s_rrci_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrci_stmt.Suffix(), instruction.RRCI)
	dc := this.ConvertPairReg(s_rrci_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrci_stmt.Ra())
	condition := this.ConvertCondition(s_rrci_stmt.Condition())
	pc := this.EvaluateProgramCounter(s_rrci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	if op_code == instruction.OR {
		if suffix == instruction.S_RRCI {
			suffix = instruction.S_RRICI
			imm := int64(0)

			instruction_.InitSRrici(op_code, suffix, dc, ra, imm, condition, pc)
		} else {
			suffix = instruction.U_RRICI
			imm := int64(0)

			instruction_.InitSRrici(op_code, suffix, dc, ra, imm, condition, pc)
		}
	} else {
		instruction_.InitSRrci(op_code, suffix, dc, ra, condition, pc)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrcStmt(stmt_ *stmt.Stmt) {
	s_rrc_stmt := stmt_.SRrcStmt()

	op_code := this.ConvertRriOpCode(s_rrc_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrc_stmt.Suffix(), instruction.RRC)
	dc := this.ConvertPairReg(s_rrc_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrc_stmt.Ra())
	condition := this.ConvertCondition(s_rrc_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrc(op_code, suffix, dc, ra, condition)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRriciStmt(stmt_ *stmt.Stmt) {
	s_rrici_stmt := stmt_.SRriciStmt()

	op_code := this.ConvertRriOpCode(s_rrici_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrici_stmt.Suffix(), instruction.RRICI)
	dc := this.ConvertPairReg(s_rrici_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrici_stmt.Ra())
	imm := this.EvaluateProgramCounter(s_rrici_stmt.Imm())
	condition := this.ConvertCondition(s_rrici_stmt.Condition())
	pc := this.EvaluateProgramCounter(s_rrici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrici(op_code, suffix, dc, ra, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRricStmt(stmt_ *stmt.Stmt) {
	s_rric_stmt := stmt_.SRricStmt()

	op_code := this.ConvertRriOpCode(s_rric_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rric_stmt.Suffix(), instruction.RRIC)
	dc := this.ConvertPairReg(s_rric_stmt.Dc())
	ra := this.ConvertSrcReg(s_rric_stmt.Ra())
	imm := this.EvaluateProgramCounter(s_rric_stmt.Imm())
	condition := this.ConvertCondition(s_rric_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	if condition != cc.FALSE {
		instruction_.InitSRric(op_code, suffix, dc, ra, imm, condition)
	} else {
		instruction_.InitSRrif(op_code, suffix, dc, ra, imm, condition)
	}

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRriStmt(stmt_ *stmt.Stmt) {
	s_rri_stmt := stmt_.SRriStmt()

	op_code := this.ConvertRriOpCode(s_rri_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rri_stmt.Suffix(), instruction.RRI)
	dc := this.ConvertPairReg(s_rri_stmt.Dc())
	ra := this.ConvertSrcReg(s_rri_stmt.Ra())
	imm := this.EvaluateProgramCounter(s_rri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRri(op_code, suffix, dc, ra, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrrciStmt(stmt_ *stmt.Stmt) {
	s_rrrci_stmt := stmt_.SRrrciStmt()

	op_code := this.ConvertRriOpCode(s_rrrci_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrrci_stmt.Suffix(), instruction.RRRCI)
	dc := this.ConvertPairReg(s_rrrci_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrrci_stmt.Ra())
	rb := this.ConvertSrcReg(s_rrrci_stmt.Rb())
	condition := this.ConvertCondition(s_rrrci_stmt.Condition())
	pc := this.EvaluateProgramCounter(s_rrrci_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrrci(op_code, suffix, dc, ra, rb, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrrcStmt(stmt_ *stmt.Stmt) {
	s_rrrc_stmt := stmt_.SRrrcStmt()

	op_code := this.ConvertRriOpCode(s_rrrc_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrrc_stmt.Suffix(), instruction.RRRC)
	dc := this.ConvertPairReg(s_rrrc_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrrc_stmt.Ra())
	rb := this.ConvertSrcReg(s_rrrc_stmt.Rb())
	condition := this.ConvertCondition(s_rrrc_stmt.Condition())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrrc(op_code, suffix, dc, ra, rb, condition)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrriciStmt(stmt_ *stmt.Stmt) {
	s_rrrici_stmt := stmt_.SRrriciStmt()

	op_code := this.ConvertRriOpCode(s_rrrici_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrrici_stmt.Suffix(), instruction.RRRICI)
	dc := this.ConvertPairReg(s_rrrici_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrrici_stmt.Ra())
	rb := this.ConvertSrcReg(s_rrrici_stmt.Rb())
	imm := this.EvaluateProgramCounter(s_rrrici_stmt.Imm())
	condition := this.ConvertCondition(s_rrrici_stmt.Condition())
	pc := this.EvaluateProgramCounter(s_rrrici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrrici(op_code, suffix, dc, ra, rb, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrriStmt(stmt_ *stmt.Stmt) {
	s_rrri_stmt := stmt_.SRrriStmt()

	op_code := this.ConvertRriOpCode(s_rrri_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrri_stmt.Suffix(), instruction.RRRI)
	dc := this.ConvertPairReg(s_rrri_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrri_stmt.Ra())
	rb := this.ConvertSrcReg(s_rrri_stmt.Rb())
	imm := this.EvaluateProgramCounter(s_rrri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrri(op_code, suffix, dc, ra, rb, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrrStmt(stmt_ *stmt.Stmt) {
	s_rrr_stmt := stmt_.SRrrStmt()

	op_code := this.ConvertRriOpCode(s_rrr_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rrr_stmt.Suffix(), instruction.RRR)
	dc := this.ConvertPairReg(s_rrr_stmt.Dc())
	ra := this.ConvertSrcReg(s_rrr_stmt.Ra())
	rb := this.ConvertSrcReg(s_rrr_stmt.Rb())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrr(op_code, suffix, dc, ra, rb)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRrStmt(stmt_ *stmt.Stmt) {
	s_rr_stmt := stmt_.SRrStmt()

	op_code := this.ConvertRrOpCode(s_rr_stmt.OpCode())
	suffix := this.ConvertSuffix(s_rr_stmt.Suffix(), instruction.RR)
	dc := this.ConvertPairReg(s_rr_stmt.Dc())
	ra := this.ConvertSrcReg(s_rr_stmt.Ra())

	instruction_ := new(instruction.Instruction)
	if op_code == instruction.OR {
		if suffix == instruction.S_RR {
			imm := int64(0)
			condition := cc.FALSE

			instruction_.InitSRrif(op_code, instruction.S_RRIF, dc, ra, imm, condition)
		} else if suffix == instruction.U_RR {
			imm := int64(0)
			condition := cc.FALSE

			instruction_.InitSRrif(op_code, instruction.U_RRIF, dc, ra, imm, condition)
		} else {
			err := errors.New("suffix is not valid")
			panic(err)
		}
	} else {
		instruction_.InitSRr(op_code, suffix, dc, ra)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSRStmt(stmt_ *stmt.Stmt) {
	s_r_stmt := stmt_.SRStmt()

	op_code := this.ConvertRriOpCode(s_r_stmt.OpCode())
	suffix := this.ConvertSuffix(s_r_stmt.Suffix(), instruction.R)
	dc := this.ConvertPairReg(s_r_stmt.Dc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSR(op_code, suffix, dc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkBkpStmt(stmt_ *stmt.Stmt) {
	op_code := instruction.FAULT
	imm := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitI(op_code, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkBootRiStmt(stmt_ *stmt.Stmt) {
	boot_ri_stmt := stmt_.BootRiStmt()

	op_code := this.ConvertRiciOpCode(boot_ri_stmt.OpCode())
	ra := this.ConvertSrcReg(boot_ri_stmt.Ra())
	imm := this.EvaluateProgramCounter(boot_ri_stmt.Imm())
	condition := cc.FALSE
	pc := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitRici(op_code, ra, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkCallRiStmt(stmt_ *stmt.Stmt) {
	call_ri_stmt := stmt_.CallRiStmt()

	op_code := instruction.CALL
	rc := this.ConvertGpReg(call_ri_stmt.Rc())

	zero_reg := new(reg_descriptor.SpRegDescriptor)
	*zero_reg = reg_descriptor.ZERO
	ra := new(reg_descriptor.SrcRegDescriptor)
	ra.InitSpRegDescriptor(zero_reg)

	imm := this.EvaluateProgramCounter(call_ri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitRri(op_code, rc, ra, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkCallRrStmt(stmt_ *stmt.Stmt) {
	call_rr_stmt := stmt_.CallRrStmt()

	op_code := instruction.CALL
	rc := this.ConvertGpReg(call_rr_stmt.Rc())
	ra := this.ConvertSrcReg(call_rr_stmt.Ra())
	imm := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitRri(op_code, rc, ra, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkDivStepDrdiStmt(stmt_ *stmt.Stmt) {
	div_step_drdi_stmt := stmt_.DivStepDrdiStmt()

	op_code := this.ConvertDrdiciOpCode(div_step_drdi_stmt.OpCode())
	dc := this.ConvertPairReg(div_step_drdi_stmt.Dc())
	ra := this.ConvertSrcReg(div_step_drdi_stmt.Ra())
	db := this.ConvertPairReg(div_step_drdi_stmt.Db())
	imm := this.EvaluateProgramCounter(div_step_drdi_stmt.Imm())
	condition := cc.FALSE
	pc := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitDrdici(op_code, dc, ra, db, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkJeqRiiStmt(stmt_ *stmt.Stmt) {
	jeq_rii_stmt := stmt_.JeqRiiStmt()

	op_code := this.ConvertJumpOpCode(jeq_rii_stmt.OpCode())
	ra := this.ConvertSrcReg(jeq_rii_stmt.Ra())
	imm := this.EvaluateProgramCounter(jeq_rii_stmt.Imm())
	condition := this.ConvertJumpCondition(jeq_rii_stmt.OpCode())
	pc := this.EvaluateProgramCounter(jeq_rii_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitZrici(op_code, ra, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkJeqRriStmt(stmt_ *stmt.Stmt) {
	jeq_rri_stmt := stmt_.JeqRriStmt()

	op_code := this.ConvertJumpOpCode(jeq_rri_stmt.OpCode())
	ra := this.ConvertSrcReg(jeq_rri_stmt.Ra())
	rb := this.ConvertSrcReg(jeq_rri_stmt.Rb())
	condition := this.ConvertJumpCondition(jeq_rri_stmt.OpCode())
	pc := this.EvaluateProgramCounter(jeq_rri_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitZrrci(op_code, ra, rb, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkJnzRiStmt(stmt_ *stmt.Stmt) {
	jnz_ri_stmt := stmt_.JnzRiStmt()

	op_code := this.ConvertJumpOpCode(jnz_ri_stmt.OpCode())
	ra := this.ConvertSrcReg(jnz_ri_stmt.Ra())

	instruction_ := new(instruction.Instruction)
	if op_code == instruction.SUB {
		imm := int64(0)
		condition := this.ConvertJumpCondition(jnz_ri_stmt.OpCode())
		pc := this.EvaluateProgramCounter(jnz_ri_stmt.Pc())

		instruction_.InitZrici(op_code, ra, imm, condition, pc)
	} else if op_code == instruction.CALL {
		imm := this.EvaluateProgramCounter(jnz_ri_stmt.Pc())

		instruction_.InitZri(op_code, ra, imm)
	} else {
		err := errors.New("op code is not valid")
		panic(err)
	}
	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkJumpIStmt(stmt_ *stmt.Stmt) {
	jump_i_stmt := stmt_.JumpIStmt()

	op_code := instruction.CALL

	zero_reg := new(reg_descriptor.SpRegDescriptor)
	*zero_reg = reg_descriptor.ZERO
	ra := new(reg_descriptor.SrcRegDescriptor)
	ra.InitSpRegDescriptor(zero_reg)

	imm := this.EvaluateProgramCounter(jump_i_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitZri(op_code, ra, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkJumpRStmt(stmt_ *stmt.Stmt) {
	jump_r_stmt := stmt_.JumpRStmt()

	op_code := instruction.CALL
	ra := this.ConvertSrcReg(jump_r_stmt.Ra())
	imm := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitZri(op_code, ra, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkLbsRriStmt(stmt_ *stmt.Stmt) {
	lbs_rri_stmt := stmt_.LbsRriStmt()

	op_code := this.ConvertLoadOpCode(lbs_rri_stmt.OpCode())
	endian := instruction.LITTLE
	rc := this.ConvertGpReg(lbs_rri_stmt.Rc())
	ra := this.ConvertSrcReg(lbs_rri_stmt.Ra())
	off := this.EvaluateProgramCounter(lbs_rri_stmt.Off())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErri(op_code, endian, rc, ra, off)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkLbsSRriStmt(stmt_ *stmt.Stmt) {
	lbs_s_rri_stmt := stmt_.LbsSRriStmt()

	op_code := this.ConvertLoadOpCode(lbs_s_rri_stmt.OpCode())
	suffix := this.ConvertSuffix(lbs_s_rri_stmt.Suffix(), instruction.ERRI)
	endian := instruction.LITTLE
	dc := this.ConvertPairReg(lbs_s_rri_stmt.Dc())
	ra := this.ConvertSrcReg(lbs_s_rri_stmt.Ra())
	off := this.EvaluateProgramCounter(lbs_s_rri_stmt.Off())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSErri(op_code, suffix, endian, dc, ra, off)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkLdDriStmt(stmt_ *stmt.Stmt) {
	ld_dri_stmt := stmt_.LdDriStmt()

	op_code := this.ConvertLoadOpCode(ld_dri_stmt.OpCode())
	endian := instruction.LITTLE
	dc := this.ConvertPairReg(ld_dri_stmt.Dc())
	ra := this.ConvertSrcReg(ld_dri_stmt.Ra())
	off := this.EvaluateProgramCounter(ld_dri_stmt.Off())

	instruction_ := new(instruction.Instruction)
	instruction_.InitEdri(op_code, endian, dc, ra, off)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkMovdDdStmt(stmt_ *stmt.Stmt) {
	movd_dd_stmt := stmt_.MovdDdStmt()

	op_code := this.ConvertDdciOpCode(movd_dd_stmt.OpCode())
	dc := this.ConvertPairReg(movd_dd_stmt.Dc())
	db := this.ConvertPairReg(movd_dd_stmt.Db())
	condition := cc.FALSE
	pc := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitDdci(op_code, dc, db, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkMoveRiciStmt(stmt_ *stmt.Stmt) {
	move_rici_stmt := stmt_.MoveRiciStmt()

	op_code := instruction.OR
	rc := this.ConvertGpReg(move_rici_stmt.Rc())

	zero_reg := new(reg_descriptor.SpRegDescriptor)
	*zero_reg = reg_descriptor.ZERO
	ra := new(reg_descriptor.SrcRegDescriptor)
	ra.InitSpRegDescriptor(zero_reg)

	imm := this.EvaluateProgramCounter(move_rici_stmt.Imm())
	condition := this.ConvertCondition(move_rici_stmt.Condition())
	pc := this.EvaluateProgramCounter(move_rici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitRrici(op_code, rc, ra, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkMoveRiStmt(stmt_ *stmt.Stmt) {
	move_ri_stmt := stmt_.MoveRiStmt()

	op_code := instruction.OR
	rc := this.ConvertGpReg(move_ri_stmt.Rc())

	zero_reg := new(reg_descriptor.SpRegDescriptor)
	*zero_reg = reg_descriptor.ZERO
	ra := new(reg_descriptor.SrcRegDescriptor)
	ra.InitSpRegDescriptor(zero_reg)

	imm := this.EvaluateProgramCounter(move_ri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitRri(op_code, rc, ra, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkMoveSRiciStmt(stmt_ *stmt.Stmt) {
	move_s_rici_stmt := stmt_.MoveSRiciStmt()

	op_code := instruction.OR
	suffix := this.ConvertSuffix(move_s_rici_stmt.Suffix(), instruction.RRICI)
	dc := this.ConvertPairReg(move_s_rici_stmt.Dc())

	zero_reg := new(reg_descriptor.SpRegDescriptor)
	*zero_reg = reg_descriptor.ZERO
	ra := new(reg_descriptor.SrcRegDescriptor)
	ra.InitSpRegDescriptor(zero_reg)

	imm := this.EvaluateProgramCounter(move_s_rici_stmt.Imm())
	condition := this.ConvertCondition(move_s_rici_stmt.Condition())
	pc := this.EvaluateProgramCounter(move_s_rici_stmt.Pc())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRrici(op_code, suffix, dc, ra, imm, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkMoveSRiStmt(stmt_ *stmt.Stmt) {
	move_s_ri_stmt := stmt_.MoveSRiStmt()

	// NOTE(bongjoon.hyun@gmail.com): move.s is implemented by using and.s:rki
	op_code := instruction.AND
	suffix := this.ConvertSuffix(move_s_ri_stmt.Suffix(), instruction.RRI)
	dc := this.ConvertPairReg(move_s_ri_stmt.Dc())

	lneg_reg := new(reg_descriptor.SpRegDescriptor)
	*lneg_reg = reg_descriptor.LNEG
	ra := new(reg_descriptor.SrcRegDescriptor)
	ra.InitSpRegDescriptor(lneg_reg)

	imm := this.EvaluateProgramCounter(move_s_ri_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitSRri(op_code, suffix, dc, ra, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSbIdRiiStmt(stmt_ *stmt.Stmt) {
	sb_id_rii_stmt := stmt_.SbIdRiiStmt()

	op_code := this.ConvertStoreOpCode(sb_id_rii_stmt.OpCode())
	endian := instruction.LITTLE
	ra := this.ConvertSrcReg(sb_id_rii_stmt.Ra())
	off := this.EvaluateProgramCounter(sb_id_rii_stmt.Off())
	imm := this.EvaluateProgramCounter(sb_id_rii_stmt.Imm())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErii(op_code, endian, ra, off, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSbIdRiStmt(stmt_ *stmt.Stmt) {
	sb_id_ri_stmt := stmt_.SbIdRiStmt()

	op_code := this.ConvertStoreOpCode(sb_id_ri_stmt.OpCode())
	endian := instruction.LITTLE
	ra := this.ConvertSrcReg(sb_id_ri_stmt.Ra())
	off := this.EvaluateProgramCounter(sb_id_ri_stmt.Off())
	imm := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitErii(op_code, endian, ra, off, imm)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSbRirStmt(stmt_ *stmt.Stmt) {
	sb_rir_stmt := stmt_.SbRirStmt()

	op_code := this.ConvertStoreOpCode(sb_rir_stmt.OpCode())
	endian := instruction.LITTLE
	ra := this.ConvertSrcReg(sb_rir_stmt.Ra())
	off := this.EvaluateProgramCounter(sb_rir_stmt.Off())
	rb := this.ConvertSrcReg(sb_rir_stmt.Rb())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErir(op_code, endian, ra, off, rb)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkSdRidStmt(stmt_ *stmt.Stmt) {
	sb_rid_stmt := stmt_.SdRidStmt()

	op_code := this.ConvertStoreOpCode(sb_rid_stmt.OpCode())
	endian := instruction.LITTLE
	ra := this.ConvertSrcReg(sb_rid_stmt.Ra())
	off := this.EvaluateProgramCounter(sb_rid_stmt.Off())
	db := this.ConvertPairReg(sb_rid_stmt.Db())

	instruction_ := new(instruction.Instruction)
	instruction_.InitErid(op_code, endian, ra, off, db)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkStopStmt(stmt_ *stmt.Stmt) {
	op_code := instruction.STOP
	condition := cc.FALSE
	pc := int64(0)

	instruction_ := new(instruction.Instruction)
	instruction_.InitCi(op_code, condition, pc)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkTimeCfgRStmt(stmt_ *stmt.Stmt) {
	time_cfg_r_stmt := stmt_.TimeCfgRStmt()

	op_code := instruction.TIME_CFG
	ra := this.ConvertSrcReg(time_cfg_r_stmt.Ra())

	instruction_ := new(instruction.Instruction)
	instruction_.InitZr(op_code, ra)

	cur_label := this.executable.CurSection().CurLabel()
	cur_label.Append(instruction_)
}

func (this *InstructionAssigner) WalkLabelStmt(stmt_ *stmt.Stmt) {
	label_stmt := stmt_.LabelStmt()

	program_counter_expr := label_stmt.Expr().ProgramCounterExpr()
	primary_expr := program_counter_expr.Expr().PrimaryExpr()
	token := primary_expr.Token()
	label_name := token.Attribute()

	if label_name != "__sys_used_mram_end" {
		this.executable.CurSection().CheckoutLabel(label_name)
	}
}

func (this *InstructionAssigner) ConvertSectionName(expr_ *expr.Expr) kernel.SectionName {
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

func (this *InstructionAssigner) ConvertName(expr_ *expr.Expr) string {
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

func (this *InstructionAssigner) ConvertCiOpCode(op_code *expr.Expr) instruction.OpCode {
	ci_op_code_expr := op_code.CiOpCodeExpr()

	token_type := ci_op_code_expr.Token().TokenType()
	if token_type == lexer.STOP {
		return instruction.STOP
	} else {
		err := errors.New("CI op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertDdciOpCode(op_code *expr.Expr) instruction.OpCode {
	ddci_op_code_expr := op_code.DdciOpCodeExpr()

	token_type := ddci_op_code_expr.Token().TokenType()
	if token_type == lexer.MOVD {
		return instruction.MOVD
	} else if token_type == lexer.SWAPD {
		return instruction.SWAPD
	} else {
		err := errors.New("DDCI op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertDmaRriOpCode(op_code *expr.Expr) instruction.OpCode {
	dma_rri_op_code_expr := op_code.DmaRriOpCodeExpr()

	token_type := dma_rri_op_code_expr.Token().TokenType()
	if token_type == lexer.LDMA {
		return instruction.LDMA
	} else if token_type == lexer.LDMAI {
		return instruction.LDMAI
	} else if token_type == lexer.SDMA {
		return instruction.SDMA
	} else {
		err := errors.New("DMA_RRI op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertDrdiciOpCode(op_code *expr.Expr) instruction.OpCode {
	drdici_op_code_expr := op_code.DrdiciOpCodeExpr()

	token_type := drdici_op_code_expr.Token().TokenType()
	if token_type == lexer.DIV_STEP {
		return instruction.DIV_STEP
	} else if token_type == lexer.MUL_STEP {
		return instruction.MUL_STEP
	} else {
		err := errors.New("DRDICI op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertIOpCode(op_code *expr.Expr) instruction.OpCode {
	i_op_code_expr := op_code.IOpCodeExpr()

	token_type := i_op_code_expr.Token().TokenType()
	if token_type == lexer.FAULT {
		return instruction.FAULT
	} else {
		err := errors.New("I op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertRiciOpCode(op_code *expr.Expr) instruction.OpCode {
	rici_op_code_expr := op_code.RiciOpCodeExpr()

	token_type := rici_op_code_expr.Token().TokenType()
	if token_type == lexer.ACQUIRE {
		return instruction.ACQUIRE
	} else if token_type == lexer.RELEASE {
		return instruction.RELEASE
	} else if token_type == lexer.BOOT {
		return instruction.BOOT
	} else if token_type == lexer.RESUME {
		return instruction.RESUME
	} else {
		err := errors.New("I op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertROpCode(op_code *expr.Expr) instruction.OpCode {
	r_op_code_expr := op_code.ROpCodeExpr()

	token_type := r_op_code_expr.Token().TokenType()
	if token_type == lexer.TIME {
		return instruction.TIME
	} else if token_type == lexer.NOP {
		return instruction.NOP
	} else {
		err := errors.New("R op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertRrOpCode(op_code *expr.Expr) instruction.OpCode {
	rr_op_code_expr := op_code.RrOpCodeExpr()

	token_type := rr_op_code_expr.Token().TokenType()
	if token_type == lexer.CAO {
		return instruction.CAO
	} else if token_type == lexer.CLO {
		return instruction.CLO
	} else if token_type == lexer.CLS {
		return instruction.CLS
	} else if token_type == lexer.CLZ {
		return instruction.CLZ
	} else if token_type == lexer.EXTSB {
		return instruction.EXTSB
	} else if token_type == lexer.EXTSH {
		return instruction.EXTSH
	} else if token_type == lexer.EXTUB {
		return instruction.EXTUB
	} else if token_type == lexer.EXTUH {
		return instruction.EXTUH
	} else if token_type == lexer.SATS {
		return instruction.SATS
	} else if token_type == lexer.TIME_CFG {
		return instruction.TIME_CFG
	} else if token_type == lexer.MOVE {
		return instruction.OR
	} else if token_type == lexer.NEG {
		return instruction.SUB
	} else if token_type == lexer.NOT {
		return instruction.XOR
	} else {
		err := errors.New("RR op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertRriOpCode(op_code *expr.Expr) instruction.OpCode {
	rri_op_code_expr := op_code.RriOpCodeExpr()

	token_type := rri_op_code_expr.Token().TokenType()
	if token_type == lexer.ADD {
		return instruction.ADD
	} else if token_type == lexer.ADDC {
		return instruction.ADDC
	} else if token_type == lexer.AND {
		return instruction.AND
	} else if token_type == lexer.ANDN {
		return instruction.ANDN
	} else if token_type == lexer.ASR {
		return instruction.ASR
	} else if token_type == lexer.CMPB4 {
		return instruction.CMPB4
	} else if token_type == lexer.LSL {
		return instruction.LSL
	} else if token_type == lexer.LSL1 {
		return instruction.LSL1
	} else if token_type == lexer.LSL1X {
		return instruction.LSL1X
	} else if token_type == lexer.LSLX {
		return instruction.LSLX
	} else if token_type == lexer.LSR {
		return instruction.LSR
	} else if token_type == lexer.LSR1 {
		return instruction.LSR1
	} else if token_type == lexer.LSR1X {
		return instruction.LSR1X
	} else if token_type == lexer.LSRX {
		return instruction.LSRX
	} else if token_type == lexer.MUL_SH_SH {
		return instruction.MUL_SH_SH
	} else if token_type == lexer.MUL_SH_SL {
		return instruction.MUL_SH_SL
	} else if token_type == lexer.MUL_SH_UH {
		return instruction.MUL_SH_UH
	} else if token_type == lexer.MUL_SH_UL {
		return instruction.MUL_SH_UL
	} else if token_type == lexer.MUL_SL_SH {
		return instruction.MUL_SL_SH
	} else if token_type == lexer.MUL_SL_SL {
		return instruction.MUL_SL_SL
	} else if token_type == lexer.MUL_SL_UH {
		return instruction.MUL_SL_UH
	} else if token_type == lexer.MUL_SL_UL {
		return instruction.MUL_SL_UL
	} else if token_type == lexer.MUL_UH_UH {
		return instruction.MUL_UH_UH
	} else if token_type == lexer.MUL_UH_UL {
		return instruction.MUL_UH_UL
	} else if token_type == lexer.MUL_UL_UH {
		return instruction.MUL_UL_UH
	} else if token_type == lexer.MUL_UL_UL {
		return instruction.MUL_UL_UL
	} else if token_type == lexer.NAND {
		return instruction.NAND
	} else if token_type == lexer.NOR {
		return instruction.NOR
	} else if token_type == lexer.NXOR {
		return instruction.NXOR
	} else if token_type == lexer.OR {
		return instruction.OR
	} else if token_type == lexer.ORN {
		return instruction.ORN
	} else if token_type == lexer.ROL {
		return instruction.ROL
	} else if token_type == lexer.ROR {
		return instruction.ROR
	} else if token_type == lexer.RSUB {
		return instruction.RSUB
	} else if token_type == lexer.RSUBC {
		return instruction.RSUBC
	} else if token_type == lexer.SUB {
		return instruction.SUB
	} else if token_type == lexer.SUBC {
		return instruction.SUBC
	} else if token_type == lexer.XOR {
		return instruction.XOR
	} else if token_type == lexer.CALL {
		return instruction.CALL
	} else if token_type == lexer.HASH {
		return instruction.HASH
	} else {
		err := errors.New("RRI op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertRrriOpCode(op_code *expr.Expr) instruction.OpCode {
	rrri_op_code_expr := op_code.RrriOpCodeExpr()

	token_type := rrri_op_code_expr.Token().TokenType()
	if token_type == lexer.LSL_ADD {
		return instruction.LSL_ADD
	} else if token_type == lexer.LSL_SUB {
		return instruction.LSL_SUB
	} else if token_type == lexer.LSR_ADD {
		return instruction.LSR_ADD
	} else if token_type == lexer.ROL_ADD {
		return instruction.ROL_ADD
	} else {
		err := errors.New("RRRI op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertLoadOpCode(op_code *expr.Expr) instruction.OpCode {
	load_op_code_expr := op_code.LoadOpCodeExpr()

	token_type := load_op_code_expr.Token().TokenType()
	if token_type == lexer.LBS {
		return instruction.LBS
	} else if token_type == lexer.LBU {
		return instruction.LBU
	} else if token_type == lexer.LD {
		return instruction.LD
	} else if token_type == lexer.LHS {
		return instruction.LHS
	} else if token_type == lexer.LHU {
		return instruction.LHU
	} else if token_type == lexer.LW {
		return instruction.LW
	} else {
		err := errors.New("load op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertStoreOpCode(op_code *expr.Expr) instruction.OpCode {
	store_op_code_expr := op_code.StoreOpCodeExpr()

	token_type := store_op_code_expr.Token().TokenType()
	if token_type == lexer.SB {
		return instruction.SB
	} else if token_type == lexer.SB_ID {
		return instruction.SB_ID
	} else if token_type == lexer.SD {
		return instruction.SD
	} else if token_type == lexer.SD_ID {
		return instruction.SD_ID
	} else if token_type == lexer.SH {
		return instruction.SH
	} else if token_type == lexer.SH_ID {
		return instruction.SH_ID
	} else if token_type == lexer.SW {
		return instruction.SW
	} else if token_type == lexer.SW_ID {
		return instruction.SW_ID
	} else {
		err := errors.New("store op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertJumpOpCode(op_code *expr.Expr) instruction.OpCode {
	jump_op_code_expr := op_code.JumpOpCodeExpr()

	token_type := jump_op_code_expr.Token().TokenType()
	if token_type == lexer.JEQ {
		return instruction.SUB
	} else if token_type == lexer.JGES {
		return instruction.SUB
	} else if token_type == lexer.JGEU {
		return instruction.SUB
	} else if token_type == lexer.JGTS {
		return instruction.SUB
	} else if token_type == lexer.JGTU {
		return instruction.SUB
	} else if token_type == lexer.JLES {
		return instruction.SUB
	} else if token_type == lexer.JLEU {
		return instruction.SUB
	} else if token_type == lexer.JLTS {
		return instruction.SUB
	} else if token_type == lexer.JLTU {
		return instruction.SUB
	} else if token_type == lexer.JNEQ {
		return instruction.SUB
	} else if token_type == lexer.JNZ {
		return instruction.SUB
	} else if token_type == lexer.JUMP {
		return instruction.CALL
	} else if token_type == lexer.JZ {
		return instruction.SUB
	} else {
		err := errors.New("jump op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertSuffix(
	suffix *expr.Expr,
	base instruction.Suffix,
) instruction.Suffix {
	suffix_expr := suffix.SuffixExpr()

	token_type := suffix_expr.Token().TokenType()
	if token_type == lexer.S {
		if base == instruction.ERRI {
			return instruction.S_ERRI
		} else if base == instruction.RCI {
			return instruction.S_RCI
		} else if base == instruction.RIRCI {
			return instruction.S_RIRCI
		} else if base == instruction.RIRC {
			return instruction.S_RIRC
		} else if base == instruction.RRCI {
			return instruction.S_RRCI
		} else if base == instruction.RRC {
			return instruction.S_RRC
		} else if base == instruction.RRICI {
			return instruction.S_RRICI
		} else if base == instruction.RRIC {
			return instruction.S_RRIC
		} else if base == instruction.RRI {
			return instruction.S_RRI
		} else if base == instruction.RRRCI {
			return instruction.S_RRRCI
		} else if base == instruction.RRRC {
			return instruction.S_RRRC
		} else if base == instruction.RRRICI {
			return instruction.S_RRRICI
		} else if base == instruction.RRRI {
			return instruction.S_RRRI
		} else if base == instruction.RRR {
			return instruction.S_RRR
		} else if base == instruction.RR {
			return instruction.S_RR
		} else if base == instruction.R {
			return instruction.S_R
		} else {
			err := errors.New("base is not valid")
			panic(err)
		}
	} else if token_type == lexer.U {
		if base == instruction.ERRI {
			return instruction.U_ERRI
		} else if base == instruction.RCI {
			return instruction.U_RCI
		} else if base == instruction.RIRCI {
			return instruction.U_RIRCI
		} else if base == instruction.RIRC {
			return instruction.U_RIRC
		} else if base == instruction.RRCI {
			return instruction.U_RRCI
		} else if base == instruction.RRC {
			return instruction.U_RRC
		} else if base == instruction.RRICI {
			return instruction.U_RRICI
		} else if base == instruction.RRIC {
			return instruction.U_RRIC
		} else if base == instruction.RRI {
			return instruction.U_RRI
		} else if base == instruction.RRRCI {
			return instruction.U_RRRCI
		} else if base == instruction.RRRC {
			return instruction.U_RRRC
		} else if base == instruction.RRRICI {
			return instruction.U_RRRICI
		} else if base == instruction.RRRI {
			return instruction.U_RRRI
		} else if base == instruction.RRR {
			return instruction.U_RRR
		} else if base == instruction.RR {
			return instruction.U_RR
		} else if base == instruction.R {
			return instruction.U_R
		} else {
			err := errors.New("base is not valid")
			panic(err)
		}
	} else {
		err := errors.New("suffix is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertGpReg(expr_ *expr.Expr) *reg_descriptor.GpRegDescriptor {
	src_reg_expr := expr_.SrcRegExpr()

	index, err := strconv.Atoi(src_reg_expr.Token().Attribute()[1:])

	if err != nil {
		panic(err)
	}

	gp_reg_descriptor := new(reg_descriptor.GpRegDescriptor)
	gp_reg_descriptor.Init(index)
	return gp_reg_descriptor
}

func (this *InstructionAssigner) ConvertSrcReg(expr_ *expr.Expr) *reg_descriptor.SrcRegDescriptor {
	src_reg_expr := expr_.SrcRegExpr()

	token := src_reg_expr.Token()
	token_type := token.TokenType()
	if token_type == lexer.GP_REG {
		gp_reg_descriptor := this.ConvertGpReg(expr_)

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitGpRegDescriptor(gp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.ZERO_REG {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.ZERO

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.ONE {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.ONE

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.ID {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.ID

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.ID2 {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.ID2

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.ID4 {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.ID4

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.ID8 {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.ID8

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.LNEG {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.LNEG

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else if token_type == lexer.MNEG {
		sp_reg_descriptor := new(reg_descriptor.SpRegDescriptor)
		*sp_reg_descriptor = reg_descriptor.MNEG

		src_reg_descriptor := new(reg_descriptor.SrcRegDescriptor)
		src_reg_descriptor.InitSpRegDescriptor(sp_reg_descriptor)
		return src_reg_descriptor
	} else {
		err := errors.New("src reg is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertPairReg(
	token *lexer.Token,
) *reg_descriptor.PairRegDescriptor {
	index, err := strconv.Atoi(token.Attribute()[1:])

	if err != nil {
		panic(err)
	}

	pair_reg_descriptor := new(reg_descriptor.PairRegDescriptor)
	pair_reg_descriptor.Init(index)
	return pair_reg_descriptor
}

func (this *InstructionAssigner) ConvertCondition(expr_ *expr.Expr) cc.Condition {
	condition_expr := expr_.ConditionExpr()

	token := condition_expr.Token()
	token_type := token.TokenType()
	if token_type == lexer.TRUE {
		return cc.TRUE
	} else if token_type == lexer.FALSE {
		return cc.FALSE
	} else if token_type == lexer.Z {
		return cc.Z
	} else if token_type == lexer.NZ {
		return cc.NZ
	} else if token_type == lexer.E {
		return cc.E
	} else if token_type == lexer.O {
		return cc.O
	} else if token_type == lexer.PL {
		return cc.PL
	} else if token_type == lexer.MI {
		return cc.MI
	} else if token_type == lexer.OV {
		return cc.OV
	} else if token_type == lexer.NOV {
		return cc.NOV
	} else if token_type == lexer.C {
		return cc.C
	} else if token_type == lexer.NC {
		return cc.NC
	} else if token_type == lexer.SZ {
		return cc.SZ
	} else if token_type == lexer.SNZ {
		return cc.SNZ
	} else if token_type == lexer.SPL {
		return cc.SPL
	} else if token_type == lexer.SMI {
		return cc.SMI
	} else if token_type == lexer.SO {
		return cc.SO
	} else if token_type == lexer.SE {
		return cc.SE
	} else if token_type == lexer.NC5 {
		return cc.NC5
	} else if token_type == lexer.NC6 {
		return cc.NC6
	} else if token_type == lexer.NC7 {
		return cc.NC7
	} else if token_type == lexer.NC8 {
		return cc.NC8
	} else if token_type == lexer.NC9 {
		return cc.NC9
	} else if token_type == lexer.NC10 {
		return cc.NC10
	} else if token_type == lexer.NC11 {
		return cc.NC11
	} else if token_type == lexer.NC12 {
		return cc.NC12
	} else if token_type == lexer.NC13 {
		return cc.NC13
	} else if token_type == lexer.NC14 {
		return cc.NC14
	} else if token_type == lexer.MAX {
		return cc.MAX
	} else if token_type == lexer.NMAX {
		return cc.NMAX
	} else if token_type == lexer.SH32 {
		return cc.SH32
	} else if token_type == lexer.NSH32 {
		return cc.NSH32
	} else if token_type == lexer.EQ {
		return cc.EQ
	} else if token_type == lexer.NEQ {
		return cc.NEQ
	} else if token_type == lexer.LTU {
		return cc.LTU
	} else if token_type == lexer.LEU {
		return cc.LEU
	} else if token_type == lexer.GTU {
		return cc.GTU
	} else if token_type == lexer.GEU {
		return cc.GEU
	} else if token_type == lexer.LTS {
		return cc.LTS
	} else if token_type == lexer.LES {
		return cc.LES
	} else if token_type == lexer.GTS {
		return cc.GTS
	} else if token_type == lexer.GES {
		return cc.GES
	} else if token_type == lexer.XZ {
		return cc.XZ
	} else if token_type == lexer.XNZ {
		return cc.XNZ
	} else if token_type == lexer.XLEU {
		return cc.XLEU
	} else if token_type == lexer.XGTU {
		return cc.XGTU
	} else if token_type == lexer.XLES {
		return cc.XLES
	} else if token_type == lexer.XGTS {
		return cc.XGTS
	} else if token_type == lexer.SMALL {
		return cc.SMALL
	} else if token_type == lexer.LARGE {
		return cc.LARGE
	} else {
		err := errors.New("condition is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertEndian(expr_ *expr.Expr) instruction.Endian {
	endian_expr := expr_.EndianExpr()

	token := endian_expr.Token()
	token_type := token.TokenType()
	if token_type == lexer.LITTLE {
		return instruction.LITTLE
	} else if token_type == lexer.BIG {
		return instruction.BIG
	} else {
		err := errors.New("endian is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) ConvertJumpCondition(op_code *expr.Expr) cc.Condition {
	jump_op_code_expr := op_code.JumpOpCodeExpr()

	token_type := jump_op_code_expr.Token().TokenType()
	if token_type == lexer.JEQ {
		return cc.Z
	} else if token_type == lexer.JGES {
		return cc.GES
	} else if token_type == lexer.JGEU {
		return cc.GEU
	} else if token_type == lexer.JGTS {
		return cc.GTS
	} else if token_type == lexer.JGTU {
		return cc.GTU
	} else if token_type == lexer.JLES {
		return cc.LES
	} else if token_type == lexer.JLEU {
		return cc.LEU
	} else if token_type == lexer.JLTS {
		return cc.LTS
	} else if token_type == lexer.JLTU {
		return cc.LTU
	} else if token_type == lexer.JNEQ {
		return cc.NZ
	} else if token_type == lexer.JNZ {
		return cc.NZ
	} else if token_type == lexer.JZ {
		return cc.Z
	} else {
		err := errors.New("jump op code is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) EvaluateProgramCounter(expr_ *expr.Expr) int64 {
	program_counter_expr := expr_.ProgramCounterExpr()

	child_expr := program_counter_expr.Expr()
	child_expr_type := child_expr.ExprType()
	if child_expr_type == expr.PRIMARY {
		return this.EvaluatePrimary(child_expr)
	} else if child_expr_type == expr.NEGATIVE_NUMBER {
		return this.EvaluateNegativeNumber(child_expr)
	} else if child_expr_type == expr.BINARY_ADD {
		return this.EvaluateBinaryAdd(child_expr)
	} else if child_expr_type == expr.BINARY_SUB {
		return this.EvaluateBinarySub(child_expr)
	} else {
		err := errors.New("program counter expr is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) EvaluatePrimary(expr_ *expr.Expr) int64 {
	primary_expr := expr_.PrimaryExpr()

	token := primary_expr.Token()
	token_type := token.TokenType()
	if token_type == lexer.POSITIVIE_NUMBER {
		return this.EvaluatePositiveNumber(token)
	} else if token_type == lexer.HEX_NUMBER {
		return this.EvaluateHexNumber(token)
	} else if token_type == lexer.IDENTIFIER {
		return this.EvaluateIdentifier(token)
	} else {
		err := errors.New("primary expr is not valid")
		panic(err)
	}
}

func (this *InstructionAssigner) EvaluatePositiveNumber(token *lexer.Token) int64 {
	value, err := strconv.ParseInt(token.Attribute(), 10, 64)

	if err != nil {
		panic(err)
	}

	return value
}

func (this *InstructionAssigner) EvaluateHexNumber(token *lexer.Token) int64 {
	attribute := token.Attribute()
	var value int64
	var err error
	if attribute[:2] == "0x" {
		value, err = strconv.ParseInt(attribute[2:], 16, 64)
	} else {
		value, err = strconv.ParseInt(attribute, 16, 64)
	}

	if err != nil {
		panic(err)
	}

	return value
}

func (this *InstructionAssigner) EvaluateIdentifier(token *lexer.Token) int64 {
	name := token.Attribute()

	label := this.executable.Label(name)
	linker_constant := this.linker_script.LinkerConstant(name)

	if label != nil {
		if linker_constant != nil {
			err := errors.New("label and linker constant both exist")
			panic(err)
		}

		return label.Address()
	} else if linker_constant != nil {
		return linker_constant.Value()
	} else {
		err := errors.New("label and linker constant do not exist")
		panic(err)
	}
}

func (this *InstructionAssigner) EvaluateNegativeNumber(expr_ *expr.Expr) int64 {
	negative_number_expr := expr_.NegativeNumberExpr()

	return -this.EvaluatePositiveNumber(negative_number_expr.Token())
}

func (this *InstructionAssigner) EvaluateBinaryAdd(expr_ *expr.Expr) int64 {
	binary_add_expr := expr_.BinaryAddExpr()

	is_operand1_nr_tasklets := binary_add_expr.Operand1().
		PrimaryExpr().
		Token().
		Attribute() ==
		"NR_TASKLETS"
	is_operand2_number := binary_add_expr.Operand2().
		PrimaryExpr().
		Token().
		TokenType() !=
		lexer.IDENTIFIER
	if !is_operand1_nr_tasklets && is_operand2_number {
		config_loader := new(misc.ConfigLoader)
		config_loader.Init()

		iram_data_size := int64(config_loader.IramDataWidth() / 8)

		return this.EvaluatePrimary(
			binary_add_expr.Operand1(),
		) + iram_data_size*this.EvaluatePrimary(
			binary_add_expr.Operand2(),
		)
	} else {
		return this.EvaluatePrimary(
			binary_add_expr.Operand1(),
		) + this.EvaluatePrimary(
			binary_add_expr.Operand2(),
		)
	}
}

func (this *InstructionAssigner) EvaluateBinarySub(expr_ *expr.Expr) int64 {
	binary_sub_expr := expr_.BinarySubExpr()

	is_operand1_nr_tasklets := binary_sub_expr.Operand1().
		PrimaryExpr().
		Token().
		Attribute() ==
		"NR_TASKLETS"
	is_operand2_number := binary_sub_expr.Operand2().
		PrimaryExpr().
		Token().
		TokenType() !=
		lexer.IDENTIFIER
	if !is_operand1_nr_tasklets && is_operand2_number {
		config_loader := new(misc.ConfigLoader)
		config_loader.Init()

		iram_data_size := int64(config_loader.IramDataWidth() / 8)

		return this.EvaluatePrimary(
			binary_sub_expr.Operand1(),
		) - iram_data_size*this.EvaluatePrimary(
			binary_sub_expr.Operand2(),
		)
	} else {
		return this.EvaluatePrimary(
			binary_sub_expr.Operand1(),
		) - this.EvaluatePrimary(
			binary_sub_expr.Operand2(),
		)
	}
}

func (this *InstructionAssigner) EvaluateSectionName(expr_ *expr.Expr) int64 {
	section_name := this.ConvertSectionName(expr_)
	name := ""

	return this.executable.Section(section_name, name).Address()
}

func (this *InstructionAssigner) IsZeroReg(expr_ *expr.Expr) bool {
	src_reg_expr := expr_.SrcRegExpr()

	return src_reg_expr.Token().TokenType() == lexer.ZERO_REG
}
