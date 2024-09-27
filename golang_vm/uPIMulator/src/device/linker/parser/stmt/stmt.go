package stmt

import (
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
	"uPIMulator/src/device/linker/parser/stmt/directive"
	"uPIMulator/src/device/linker/parser/stmt/instruction"
	"uPIMulator/src/device/linker/parser/stmt/sugar"
)

type StmtType int

const (
	ADDRSIG StmtType = iota
	ADDRSIG_SYM
	ASCII
	ASCIZ
	BYTE
	CFI_DEF_CFA_OFFSET
	CFI_ENDPROC
	CFI_OFFSET
	CFI_SECTIONS
	CFI_STARTPROC
	FILE_NUMBER
	FILE_STRING
	GLOBAL
	LOC_IS_STMT
	LOC_NUMBER
	LOC_PROLOGUE_END
	LONG_PROGRAM_COUNTER
	LONG_SECTION_NAME
	P2_ALIGN
	QUAD
	SECTION_IDENTIFIER_NUMBER
	SECTION_IDENTIFIER
	SECTION_STACK_SIZES
	SECTION_STRING_NUMBER
	SECTION_STRING
	SET
	SHORT
	SIZE
	TEXT
	TYPE
	WEAK
	ZERO_SINGLE_NUMBER
	ZERO_DOUBLE_NUMBER

	RICI
	RRI
	RRIC
	RRICI
	RRR
	RRRC
	RRRCI

	S_RRI
	S_RRIC
	S_RRICI
	S_RRR
	S_RRRC
	S_RRRCI

	RR
	RRC
	RRCI

	S_RR
	S_RRC
	S_RRCI

	DRDICI

	RRRI
	RRRICI

	S_RRRI
	S_RRRICI

	RIR
	RIRC
	RIRCI

	S_RIRC
	S_RIRCI

	R
	RCI

	S_R
	S_RCI

	CI
	I

	DDCI

	ERRI
	EDRI
	S_ERRI

	ERII
	ERIR
	ERID

	DMA_RRI

	NOP

	MOVE_RI
	MOVE_RICI
	MOVE_S_RI
	MOVE_S_RICI

	JEQ_RII
	JEQ_RRI
	JNZ_RI
	JUMP_I
	JUMP_R

	DIV_STEP_DRDI
	BOOT_RI
	STOP
	CALL_RI
	CALL_RR
	BKP
	MOVD_DD
	TIME_CFG_R
	LBS_RRI
	LBS_S_RRI
	LD_DRI
	SB_RIR
	SB_ID_RII
	SB_ID_RI
	SD_RID

	LABEL
)

type Stmt struct {
	stmt_type StmtType

	addrsig_stmt                   *directive.AddrsigStmt
	addrsig_sym_stmt               *directive.AddrsigSymStmt
	ascii_stmt                     *directive.AsciiStmt
	asciz_stmt                     *directive.AscizStmt
	byte_stmt                      *directive.ByteStmt
	cfi_def_cfa_offset_stmt        *directive.CfiDefCfaOffsetStmt
	cfi_endproc_stmt               *directive.CfiEndprocStmt
	cfi_offset_stmt                *directive.CfiOffsetStmt
	cfi_sections_stmt              *directive.CfiSectionsStmt
	cfi_startproc_stmt             *directive.CfiStartprocStmt
	file_number_stmt               *directive.FileNumberStmt
	file_string_stmt               *directive.FileStringStmt
	global_stmt                    *directive.GlobalStmt
	loc_is_stmt_stmt               *directive.LocIsStmtStmt
	loc_number_stmt                *directive.LocNumberStmt
	loc_prologue_end_stmt          *directive.LocPrologueEndStmt
	long_program_counter_stmt      *directive.LongProgramCounterStmt
	long_section_name_stmt         *directive.LongSectionNameStmt
	p2_align_stmt                  *directive.P2AlignStmt
	quad_stmt                      *directive.QuadStmt
	section_identifier_number_stmt *directive.SectionIdentifierNumberStmt
	section_identifier_stmt        *directive.SectionIdentifierStmt
	section_stack_sizes_stmt       *directive.SectionStackSizesStmt
	section_string_number_stmt     *directive.SectionStringNumberStmt
	section_string_stmt            *directive.SectionStringStmt
	set_stmt                       *directive.SetStmt
	short_stmt                     *directive.ShortStmt
	size_stmt                      *directive.SizeStmt
	text_stmt                      *directive.TextStmt
	type_stmt                      *directive.TypeStmt
	weak_stmt                      *directive.WeakStmt
	zero_single_number_stmt        *directive.ZeroSingleNumberStmt
	zero_double_number_stmt        *directive.ZeroDoubleNumberStmt

	ci_stmt      *instruction.CiStmt
	ddci_stmt    *instruction.DdciStmt
	dma_rri_stmt *instruction.DmaRriStmt
	drdici_stmt  *instruction.DrdiciStmt
	edri_stmt    *instruction.EdriStmt
	erid_stmt    *instruction.EridStmt
	erii_stmt    *instruction.EriiStmt
	erir_stmt    *instruction.ErirStmt
	erri_stmt    *instruction.ErriStmt
	i_stmt       *instruction.IStmt
	nop_stmt     *instruction.NopStmt
	r_stmt       *instruction.RStmt
	rci_stmt     *instruction.RciStmt
	rici_stmt    *instruction.RiciStmt
	rir_stmt     *instruction.RirStmt
	rirc_stmt    *instruction.RircStmt
	rirci_stmt   *instruction.RirciStmt
	rr_stmt      *instruction.RrStmt
	rrc_stmt     *instruction.RrcStmt
	rrci_stmt    *instruction.RrciStmt
	rri_stmt     *instruction.RriStmt
	rric_stmt    *instruction.RricStmt
	rrici_stmt   *instruction.RriciStmt
	rrr_stmt     *instruction.RrrStmt
	rrrc_stmt    *instruction.RrrcStmt
	rrrci_stmt   *instruction.RrrciStmt
	rrri_stmt    *instruction.RrriStmt
	rrrici_stmt  *instruction.RrriciStmt

	s_erri_stmt   *instruction.SErriStmt
	s_r_stmt      *instruction.SRStmt
	s_rci_stmt    *instruction.SRciStmt
	s_rirc_stmt   *instruction.SRircStmt
	s_rirci_stmt  *instruction.SRirciStmt
	s_rr_stmt     *instruction.SRrStmt
	s_rrc_stmt    *instruction.SRrcStmt
	s_rrci_stmt   *instruction.SRrciStmt
	s_rri_stmt    *instruction.SRriStmt
	s_rric_stmt   *instruction.SRricStmt
	s_rrici_stmt  *instruction.SRriciStmt
	s_rrr_stmt    *instruction.SRrrStmt
	s_rrrc_stmt   *instruction.SRrrcStmt
	s_rrrci_stmt  *instruction.SRrrciStmt
	s_rrri_stmt   *instruction.SRrriStmt
	s_rrrici_stmt *instruction.SRrriciStmt

	bkp_stmt           *sugar.BkpStmt
	boot_ri_stmt       *sugar.BootRiStmt
	call_ri_stmt       *sugar.CallRiStmt
	call_rr_stmt       *sugar.CallRrStmt
	div_step_drdi_stmt *sugar.DivStepDrdiStmt
	jeq_rii_stmt       *sugar.JeqRiiStmt
	jeq_rri_stmt       *sugar.JeqRriStmt
	jnz_ri_stmt        *sugar.JnzRiStmt
	jump_i_stmt        *sugar.JumpIStmt
	jump_r_stmt        *sugar.JumpRStmt
	lbs_rri_stmt       *sugar.LbsRriStmt
	lbs_s_rri_stmt     *sugar.LbsSRriStmt
	ld_dri_stmt        *sugar.LdDriStmt
	movd_dd_stmt       *sugar.MovdDdStmt
	move_ri_stmt       *sugar.MoveRiStmt
	move_rici_stmt     *sugar.MoveRiciStmt
	move_s_ri_stmt     *sugar.MoveSRiStmt
	move_s_rici_stmt   *sugar.MoveSRiciStmt
	sb_id_ri_stmt      *sugar.SbIdRiStmt
	sb_id_rii_stmt     *sugar.SbIdRiiStmt
	sb_rir_stmt        *sugar.SbRirStmt
	sd_rid_stmt        *sugar.SdRidStmt
	stop_stmt          *sugar.StopStmt
	time_cfg_r_stmt    *sugar.TimeCfgRStmt

	label_stmt *LabelStmt
}

func (this *Stmt) InitAddrsigStmt() {
	this.stmt_type = ADDRSIG

	this.addrsig_stmt = new(directive.AddrsigStmt)
	this.addrsig_stmt.Init()
}

func (this *Stmt) InitAddrsigSymStmt(expr_ *expr.Expr) {
	this.stmt_type = ADDRSIG_SYM

	this.addrsig_sym_stmt = new(directive.AddrsigSymStmt)
	this.addrsig_sym_stmt.Init(expr_)
}

func (this *Stmt) InitAsciiStmt(token *lexer.Token) {
	this.stmt_type = ASCII

	this.ascii_stmt = new(directive.AsciiStmt)
	this.ascii_stmt.Init(token)
}

func (this *Stmt) InitAscizStmt(token *lexer.Token) {
	this.stmt_type = ASCIZ

	this.asciz_stmt = new(directive.AscizStmt)
	this.asciz_stmt.Init(token)
}

func (this *Stmt) InitByteStmt(expr_ *expr.Expr) {
	this.stmt_type = BYTE

	this.byte_stmt = new(directive.ByteStmt)
	this.byte_stmt.Init(expr_)
}

func (this *Stmt) InitCfiDefCfaOffsetStmt(expr_ *expr.Expr) {
	this.stmt_type = CFI_DEF_CFA_OFFSET

	this.cfi_def_cfa_offset_stmt = new(directive.CfiDefCfaOffsetStmt)
	this.cfi_def_cfa_offset_stmt.Init(expr_)
}

func (this *Stmt) InitCfiEndprocStmt() {
	this.stmt_type = CFI_ENDPROC

	this.cfi_endproc_stmt = new(directive.CfiEndprocStmt)
	this.cfi_endproc_stmt.Init()
}

func (this *Stmt) InitCfiOffsetStmt(expr1 *expr.Expr, expr2 *expr.Expr) {
	this.stmt_type = CFI_OFFSET

	this.cfi_offset_stmt = new(directive.CfiOffsetStmt)
	this.cfi_offset_stmt.Init(expr1, expr2)
}

func (this *Stmt) InitCfiSectionsStmt(expr_ *expr.Expr) {
	this.stmt_type = CFI_SECTIONS

	this.cfi_sections_stmt = new(directive.CfiSectionsStmt)
	this.cfi_sections_stmt.Init(expr_)
}

func (this *Stmt) InitCfiStartprocStmt() {
	this.stmt_type = CFI_STARTPROC

	this.cfi_startproc_stmt = new(directive.CfiStartprocStmt)
	this.cfi_startproc_stmt.Init()
}

func (this *Stmt) InitFileNumberStmt(expr_ *expr.Expr, token1 *lexer.Token, token2 *lexer.Token) {
	this.stmt_type = FILE_NUMBER

	this.file_number_stmt = new(directive.FileNumberStmt)
	this.file_number_stmt.Init(expr_, token1, token2)
}

func (this *Stmt) InitFileStringStmt(token *lexer.Token) {
	this.stmt_type = FILE_STRING

	this.file_string_stmt = new(directive.FileStringStmt)
	this.file_string_stmt.Init(token)
}

func (this *Stmt) InitGlobalStmt(expr_ *expr.Expr) {
	this.stmt_type = GLOBAL

	this.global_stmt = new(directive.GlobalStmt)
	this.global_stmt.Init(expr_)
}

func (this *Stmt) InitLocIsStmtStmt(
	expr1 *expr.Expr,
	expr2 *expr.Expr,
	expr3 *expr.Expr,
	expr4 *expr.Expr,
) {
	this.stmt_type = LOC_IS_STMT

	this.loc_is_stmt_stmt = new(directive.LocIsStmtStmt)
	this.loc_is_stmt_stmt.Init(expr1, expr2, expr3, expr4)
}

func (this *Stmt) InitLocNumberStmt(expr1 *expr.Expr, expr2 *expr.Expr, expr3 *expr.Expr) {
	this.stmt_type = LOC_NUMBER

	this.loc_number_stmt = new(directive.LocNumberStmt)
	this.loc_number_stmt.Init(expr1, expr2, expr3)
}

func (this *Stmt) InitLocPrologueEndStmt(expr1 *expr.Expr, expr2 *expr.Expr, expr3 *expr.Expr) {
	this.stmt_type = LOC_PROLOGUE_END

	this.loc_prologue_end_stmt = new(directive.LocPrologueEndStmt)
	this.loc_prologue_end_stmt.Init(expr1, expr2, expr3)
}

func (this *Stmt) InitLongProgramCounterStmt(expr_ *expr.Expr) {
	this.stmt_type = LONG_PROGRAM_COUNTER

	this.long_program_counter_stmt = new(directive.LongProgramCounterStmt)
	this.long_program_counter_stmt.Init(expr_)
}

func (this *Stmt) InitLongSectionNameStmt(expr_ *expr.Expr) {
	this.stmt_type = LONG_SECTION_NAME

	this.long_section_name_stmt = new(directive.LongSectionNameStmt)
	this.long_section_name_stmt.Init(expr_)
}

func (this *Stmt) InitP2AlignStmt(expr_ *expr.Expr) {
	this.stmt_type = P2_ALIGN

	this.p2_align_stmt = new(directive.P2AlignStmt)
	this.p2_align_stmt.Init(expr_)
}

func (this *Stmt) InitQuadStmt(expr_ *expr.Expr) {
	this.stmt_type = QUAD

	this.quad_stmt = new(directive.QuadStmt)
	this.quad_stmt.Init(expr_)
}

func (this *Stmt) InitSectionIdentifierNumberStmt(
	expr1 *expr.Expr,
	expr2 *expr.Expr,
	token *lexer.Token,
	expr3 *expr.Expr,
	expr4 *expr.Expr,
) {
	this.stmt_type = SECTION_IDENTIFIER_NUMBER

	this.section_identifier_number_stmt = new(directive.SectionIdentifierNumberStmt)
	this.section_identifier_number_stmt.Init(expr1, expr2, token, expr3, expr4)
}

func (this *Stmt) InitSectionIdentifierStmt(
	expr1 *expr.Expr,
	expr2 *expr.Expr,
	token *lexer.Token,
	expr3 *expr.Expr,
) {
	this.stmt_type = SECTION_IDENTIFIER

	this.section_identifier_stmt = new(directive.SectionIdentifierStmt)
	this.section_identifier_stmt.Init(expr1, expr2, token, expr3)
}

func (this *Stmt) InitSectionStackSizesStmt(
	token *lexer.Token,
	expr1 *expr.Expr,
	expr2 *expr.Expr,
	expr3 *expr.Expr,
) {
	this.stmt_type = SECTION_STACK_SIZES

	this.section_stack_sizes_stmt = new(directive.SectionStackSizesStmt)
	this.section_stack_sizes_stmt.Init(token, expr1, expr2, expr3)
}

func (this *Stmt) InitSectionStringNumberStmt(
	expr1 *expr.Expr,
	token *lexer.Token,
	expr2 *expr.Expr,
	expr3 *expr.Expr,
) {
	this.stmt_type = SECTION_STRING_NUMBER

	this.section_string_number_stmt = new(directive.SectionStringNumberStmt)
	this.section_string_number_stmt.Init(expr1, token, expr2, expr3)
}

func (this *Stmt) InitSectionStringStmt(expr1 *expr.Expr, token *lexer.Token, expr2 *expr.Expr) {
	this.stmt_type = SECTION_STRING

	this.section_string_stmt = new(directive.SectionStringStmt)
	this.section_string_stmt.Init(expr1, token, expr2)
}

func (this *Stmt) InitSetStmt(expr1 *expr.Expr, expr2 *expr.Expr) {
	this.stmt_type = SET

	this.set_stmt = new(directive.SetStmt)
	this.set_stmt.Init(expr1, expr2)
}

func (this *Stmt) InitShortStmt(expr_ *expr.Expr) {
	this.stmt_type = SHORT

	this.short_stmt = new(directive.ShortStmt)
	this.short_stmt.Init(expr_)
}

func (this *Stmt) InitSizeStmt(expr1 *expr.Expr, expr2 *expr.Expr) {
	this.stmt_type = SIZE

	this.size_stmt = new(directive.SizeStmt)
	this.size_stmt.Init(expr1, expr2)
}

func (this *Stmt) InitTextStmt() {
	this.stmt_type = TEXT

	this.text_stmt = new(directive.TextStmt)
	this.text_stmt.Init()
}

func (this *Stmt) InitTypeStmt(expr1 *expr.Expr, expr2 *expr.Expr) {
	this.stmt_type = TYPE

	this.type_stmt = new(directive.TypeStmt)
	this.type_stmt.Init(expr1, expr2)
}

func (this *Stmt) InitWeakStmt(expr_ *expr.Expr) {
	this.stmt_type = WEAK

	this.weak_stmt = new(directive.WeakStmt)
	this.weak_stmt.Init(expr_)
}

func (this *Stmt) InitZeroSingleNumberStmt(expr_ *expr.Expr) {
	this.stmt_type = ZERO_SINGLE_NUMBER

	this.zero_single_number_stmt = new(directive.ZeroSingleNumberStmt)
	this.zero_single_number_stmt.Init(expr_)
}

func (this *Stmt) InitZeroDoubleNumberStmt(expr1 *expr.Expr, expr2 *expr.Expr) {
	this.stmt_type = ZERO_DOUBLE_NUMBER

	this.zero_double_number_stmt = new(directive.ZeroDoubleNumberStmt)
	this.zero_double_number_stmt.Init(expr1, expr2)
}

func (this *Stmt) InitCiStmt(expr1 *expr.Expr, expr2 *expr.Expr, expr3 *expr.Expr) {
	this.stmt_type = CI

	this.ci_stmt = new(instruction.CiStmt)
	this.ci_stmt.Init(expr1, expr2, expr3)
}

func (this *Stmt) InitDdciStmt(
	op_code *expr.Expr,
	dc *lexer.Token,
	db *lexer.Token,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = DDCI

	this.ddci_stmt = new(instruction.DdciStmt)
	this.ddci_stmt.Init(op_code, dc, db, condition, pc)
}

func (this *Stmt) InitDmaRriStmt(op_code *expr.Expr, ra *expr.Expr, rb *expr.Expr, imm *expr.Expr) {
	this.stmt_type = DMA_RRI

	this.dma_rri_stmt = new(instruction.DmaRriStmt)
	this.dma_rri_stmt.Init(op_code, ra, rb, imm)
}

func (this *Stmt) InitDrdiciStmt(
	op_code *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	db *lexer.Token,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = DRDICI

	this.drdici_stmt = new(instruction.DrdiciStmt)
	this.drdici_stmt.Init(op_code, dc, ra, db, imm, condition, pc)
}

func (this *Stmt) InitEdriStmt(
	op_code *expr.Expr,
	endian *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	off *expr.Expr,
) {
	this.stmt_type = EDRI

	this.edri_stmt = new(instruction.EdriStmt)
	this.edri_stmt.Init(op_code, endian, dc, ra, off)
}

func (this *Stmt) InitEridStmt(
	op_code *expr.Expr,
	endian *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
	db *lexer.Token,
) {
	this.stmt_type = ERID

	this.erid_stmt = new(instruction.EridStmt)
	this.erid_stmt.Init(op_code, endian, ra, off, db)
}

func (this *Stmt) InitEriiStmt(
	op_code *expr.Expr,
	endian *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
	imm *expr.Expr,
) {
	this.stmt_type = ERII

	this.erii_stmt = new(instruction.EriiStmt)
	this.erii_stmt.Init(op_code, endian, ra, off, imm)
}

func (this *Stmt) InitErirStmt(
	op_code *expr.Expr,
	endian *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
	rb *expr.Expr,
) {
	this.stmt_type = ERIR

	this.erir_stmt = new(instruction.ErirStmt)
	this.erir_stmt.Init(op_code, endian, ra, off, rb)
}

func (this *Stmt) InitErriStmt(
	op_code *expr.Expr,
	endian *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
) {
	this.stmt_type = ERRI

	this.erri_stmt = new(instruction.ErriStmt)
	this.erri_stmt.Init(op_code, endian, rc, ra, off)
}

func (this *Stmt) InitIStmt(op_code *expr.Expr, imm *expr.Expr) {
	this.stmt_type = I

	this.i_stmt = new(instruction.IStmt)
	this.i_stmt.Init(op_code, imm)
}

func (this *Stmt) InitNopStmt(op_code *expr.Expr) {
	this.stmt_type = NOP

	this.nop_stmt = new(instruction.NopStmt)
	this.nop_stmt.Init(op_code)
}

func (this *Stmt) InitRStmt(op_code *expr.Expr, rc *expr.Expr) {
	this.stmt_type = R

	this.r_stmt = new(instruction.RStmt)
	this.r_stmt.Init(op_code, rc)
}

func (this *Stmt) InitRciStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = RCI

	this.rci_stmt = new(instruction.RciStmt)
	this.rci_stmt.Init(op_code, rc, condition, pc)
}

func (this *Stmt) InitRiciStmt(
	op_code *expr.Expr,
	ra *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = RICI

	this.rici_stmt = new(instruction.RiciStmt)
	this.rici_stmt.Init(op_code, ra, imm, condition, pc)
}

func (this *Stmt) InitRirStmt(op_code *expr.Expr, rc *expr.Expr, imm *expr.Expr, ra *expr.Expr) {
	this.stmt_type = RIR

	this.rir_stmt = new(instruction.RirStmt)
	this.rir_stmt.Init(op_code, rc, imm, ra)
}

func (this *Stmt) InitRircStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	imm *expr.Expr,
	ra *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = RIRC

	this.rirc_stmt = new(instruction.RircStmt)
	this.rirc_stmt.Init(op_code, rc, imm, ra, condition)
}

func (this *Stmt) InitRirciStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	imm *expr.Expr,
	ra *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = RIRCI

	this.rirci_stmt = new(instruction.RirciStmt)
	this.rirci_stmt.Init(op_code, rc, imm, ra, condition, pc)
}

func (this *Stmt) InitRrStmt(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr) {
	this.stmt_type = RR

	this.rr_stmt = new(instruction.RrStmt)
	this.rr_stmt.Init(op_code, rc, ra)
}

func (this *Stmt) InitRrcStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = RRC

	this.rrc_stmt = new(instruction.RrcStmt)
	this.rrc_stmt.Init(op_code, rc, ra, condition)
}

func (this *Stmt) InitRrciStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = RRCI

	this.rrci_stmt = new(instruction.RrciStmt)
	this.rrci_stmt.Init(op_code, rc, ra, condition, pc)
}

func (this *Stmt) InitRriStmt(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr, imm *expr.Expr) {
	this.stmt_type = RRI

	this.rri_stmt = new(instruction.RriStmt)
	this.rri_stmt.Init(op_code, rc, ra, imm)
}

func (this *Stmt) InitRricStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = RRIC

	this.rric_stmt = new(instruction.RricStmt)
	this.rric_stmt.Init(op_code, rc, ra, imm, condition)
}

func (this *Stmt) InitRriciStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = RRICI

	this.rrici_stmt = new(instruction.RriciStmt)
	this.rrici_stmt.Init(op_code, rc, ra, imm, condition, pc)
}

func (this *Stmt) InitRrrStmt(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr, rb *expr.Expr) {
	this.stmt_type = RRR

	this.rrr_stmt = new(instruction.RrrStmt)
	this.rrr_stmt.Init(op_code, rc, ra, rb)
}

func (this *Stmt) InitRrrcStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	rb *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = RRRC

	this.rrrc_stmt = new(instruction.RrrcStmt)
	this.rrrc_stmt.Init(op_code, rc, ra, rb, condition)
}

func (this *Stmt) InitRrrciStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	rb *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = RRRCI

	this.rrrci_stmt = new(instruction.RrrciStmt)
	this.rrrci_stmt.Init(op_code, rc, ra, rb, condition, pc)
}

func (this *Stmt) InitRrriStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	rb *expr.Expr,
	imm *expr.Expr,
) {
	this.stmt_type = RRRI

	this.rrri_stmt = new(instruction.RrriStmt)
	this.rrri_stmt.Init(op_code, rc, ra, rb, imm)
}

func (this *Stmt) InitRrriciStmt(
	op_code *expr.Expr,
	rc *expr.Expr,
	ra *expr.Expr,
	rb *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = RRRICI

	this.rrrici_stmt = new(instruction.RrriciStmt)
	this.rrrici_stmt.Init(op_code, rc, ra, rb, imm, condition, pc)
}

func (this *Stmt) InitSErriStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	endian *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	off *expr.Expr,
) {
	this.stmt_type = S_ERRI

	this.s_erri_stmt = new(instruction.SErriStmt)
	this.s_erri_stmt.Init(op_code, suffix, endian, dc, ra, off)
}

func (this *Stmt) InitSRStmt(op_code *expr.Expr, suffix *expr.Expr, dc *lexer.Token) {
	this.stmt_type = S_R

	this.s_r_stmt = new(instruction.SRStmt)
	this.s_r_stmt.Init(op_code, suffix, dc)
}

func (this *Stmt) InitSRciStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = S_RCI

	this.s_rci_stmt = new(instruction.SRciStmt)
	this.s_rci_stmt.Init(op_code, suffix, dc, condition, pc)
}

func (this *Stmt) InitSRircStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	imm *expr.Expr,
	ra *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = S_RIRC

	this.s_rirc_stmt = new(instruction.SRircStmt)
	this.s_rirc_stmt.Init(op_code, suffix, dc, imm, ra, condition)
}

func (this *Stmt) InitSRirciStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	imm *expr.Expr,
	ra *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = S_RIRCI

	this.s_rirci_stmt = new(instruction.SRirciStmt)
	this.s_rirci_stmt.Init(op_code, suffix, dc, imm, ra, condition, pc)
}

func (this *Stmt) InitSRrStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
) {
	this.stmt_type = S_RR

	this.s_rr_stmt = new(instruction.SRrStmt)
	this.s_rr_stmt.Init(op_code, suffix, dc, ra)
}

func (this *Stmt) InitSRrcStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = S_RRC

	this.s_rrc_stmt = new(instruction.SRrcStmt)
	this.s_rrc_stmt.Init(op_code, suffix, dc, ra, condition)
}

func (this *Stmt) InitSRrciStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = S_RRCI

	this.s_rrci_stmt = new(instruction.SRrciStmt)
	this.s_rrci_stmt.Init(op_code, suffix, dc, ra, condition, pc)
}

func (this *Stmt) InitSRriStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	imm *expr.Expr,
) {
	this.stmt_type = S_RRI

	this.s_rri_stmt = new(instruction.SRriStmt)
	this.s_rri_stmt.Init(op_code, suffix, dc, ra, imm)
}

func (this *Stmt) InitSRricStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = S_RRIC

	this.s_rric_stmt = new(instruction.SRricStmt)
	this.s_rric_stmt.Init(op_code, suffix, dc, ra, imm, condition)
}

func (this *Stmt) InitSRriciStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = S_RRICI

	this.s_rrici_stmt = new(instruction.SRriciStmt)
	this.s_rrici_stmt.Init(op_code, suffix, dc, ra, imm, condition, pc)
}

func (this *Stmt) InitSRrrStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
) {
	this.stmt_type = S_RRR

	this.s_rrr_stmt = new(instruction.SRrrStmt)
	this.s_rrr_stmt.Init(op_code, suffix, dc, ra, rb)
}

func (this *Stmt) InitSRrrcStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
	condition *expr.Expr,
) {
	this.stmt_type = S_RRRC

	this.s_rrrc_stmt = new(instruction.SRrrcStmt)
	this.s_rrrc_stmt.Init(op_code, suffix, dc, ra, rb, condition)
}

func (this *Stmt) InitSRrrciStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = S_RRRCI

	this.s_rrrci_stmt = new(instruction.SRrrciStmt)
	this.s_rrrci_stmt.Init(op_code, suffix, dc, ra, rb, condition, pc)
}

func (this *Stmt) InitSRrriStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
	imm *expr.Expr,
) {
	this.stmt_type = S_RRRI

	this.s_rrri_stmt = new(instruction.SRrriStmt)
	this.s_rrri_stmt.Init(op_code, suffix, dc, ra, rb, imm)
}

func (this *Stmt) InitSRrriciStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = S_RRRICI

	this.s_rrrici_stmt = new(instruction.SRrriciStmt)
	this.s_rrrici_stmt.Init(op_code, suffix, dc, ra, rb, imm, condition, pc)
}

func (this *Stmt) InitBkpStmt() {
	this.stmt_type = BKP

	this.bkp_stmt = new(sugar.BkpStmt)
	this.bkp_stmt.Init()
}

func (this *Stmt) InitBootRiStmt(op_code *expr.Expr, ra *expr.Expr, imm *expr.Expr) {
	this.stmt_type = BOOT_RI

	this.boot_ri_stmt = new(sugar.BootRiStmt)
	this.boot_ri_stmt.Init(op_code, ra, imm)
}

func (this *Stmt) InitCallRiStmt(rc *expr.Expr, imm *expr.Expr) {
	this.stmt_type = CALL_RI

	this.call_ri_stmt = new(sugar.CallRiStmt)
	this.call_ri_stmt.Init(rc, imm)
}

func (this *Stmt) InitCallRrStmt(rc *expr.Expr, ra *expr.Expr) {
	this.stmt_type = CALL_RR

	this.call_rr_stmt = new(sugar.CallRrStmt)
	this.call_rr_stmt.Init(rc, ra)
}

func (this *Stmt) InitDivStepDrdiStmt(
	op_code *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	db *lexer.Token,
	imm *expr.Expr,
) {
	this.stmt_type = DIV_STEP_DRDI

	this.div_step_drdi_stmt = new(sugar.DivStepDrdiStmt)
	this.div_step_drdi_stmt.Init(op_code, dc, ra, db, imm)
}

func (this *Stmt) InitJeqRiiStmt(op_code *expr.Expr, ra *expr.Expr, imm *expr.Expr, pc *expr.Expr) {
	this.stmt_type = JEQ_RII

	this.jeq_rii_stmt = new(sugar.JeqRiiStmt)
	this.jeq_rii_stmt.Init(op_code, ra, imm, pc)
}

func (this *Stmt) InitJeqRriStmt(op_code *expr.Expr, ra *expr.Expr, rb *expr.Expr, pc *expr.Expr) {
	this.stmt_type = JEQ_RRI

	this.jeq_rri_stmt = new(sugar.JeqRriStmt)
	this.jeq_rri_stmt.Init(op_code, ra, rb, pc)
}

func (this *Stmt) InitJnzRiStmt(op_code *expr.Expr, ra *expr.Expr, pc *expr.Expr) {
	this.stmt_type = JNZ_RI

	this.jnz_ri_stmt = new(sugar.JnzRiStmt)
	this.jnz_ri_stmt.Init(op_code, ra, pc)
}

func (this *Stmt) InitJumpIStmt(pc *expr.Expr) {
	this.stmt_type = JUMP_I

	this.jump_i_stmt = new(sugar.JumpIStmt)
	this.jump_i_stmt.Init(pc)
}

func (this *Stmt) InitJumpRStmt(ra *expr.Expr) {
	this.stmt_type = JUMP_R

	this.jump_r_stmt = new(sugar.JumpRStmt)
	this.jump_r_stmt.Init(ra)
}

func (this *Stmt) InitLbsRriStmt(op_code *expr.Expr, rc *expr.Expr, ra *expr.Expr, off *expr.Expr) {
	this.stmt_type = LBS_RRI

	this.lbs_rri_stmt = new(sugar.LbsRriStmt)
	this.lbs_rri_stmt.Init(op_code, rc, ra, off)
}

func (this *Stmt) InitLbsSRriStmt(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	off *expr.Expr,
) {
	this.stmt_type = LBS_S_RRI

	this.lbs_s_rri_stmt = new(sugar.LbsSRriStmt)
	this.lbs_s_rri_stmt.Init(op_code, suffix, dc, ra, off)
}

func (this *Stmt) InitLdDriStmt(
	op_code *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	off *expr.Expr,
) {
	this.stmt_type = LD_DRI

	this.ld_dri_stmt = new(sugar.LdDriStmt)
	this.ld_dri_stmt.Init(op_code, dc, ra, off)
}

func (this *Stmt) InitMovdDdStmt(op_code *expr.Expr, dc *lexer.Token, db *lexer.Token) {
	this.stmt_type = MOVD_DD

	this.movd_dd_stmt = new(sugar.MovdDdStmt)
	this.movd_dd_stmt.Init(op_code, dc, db)
}

func (this *Stmt) InitMoveRiStmt(rc *expr.Expr, imm *expr.Expr) {
	this.stmt_type = MOVE_RI

	this.move_ri_stmt = new(sugar.MoveRiStmt)
	this.move_ri_stmt.Init(rc, imm)
}

func (this *Stmt) InitMoveRiciStmt(
	rc *expr.Expr,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = MOVE_RICI

	this.move_rici_stmt = new(sugar.MoveRiciStmt)
	this.move_rici_stmt.Init(rc, imm, condition, pc)
}

func (this *Stmt) InitMoveSRiStmt(suffix *expr.Expr, dc *lexer.Token, imm *expr.Expr) {
	this.stmt_type = MOVE_S_RI

	this.move_s_ri_stmt = new(sugar.MoveSRiStmt)
	this.move_s_ri_stmt.Init(suffix, dc, imm)
}

func (this *Stmt) InitMoveSRiciStmt(
	suffix *expr.Expr,
	dc *lexer.Token,
	imm *expr.Expr,
	condition *expr.Expr,
	pc *expr.Expr,
) {
	this.stmt_type = MOVE_S_RICI

	this.move_s_rici_stmt = new(sugar.MoveSRiciStmt)
	this.move_s_rici_stmt.Init(suffix, dc, imm, condition, pc)
}

func (this *Stmt) InitSbIdRiStmt(op_code *expr.Expr, ra *expr.Expr, off *expr.Expr) {
	this.stmt_type = SB_ID_RI

	this.sb_id_ri_stmt = new(sugar.SbIdRiStmt)
	this.sb_id_ri_stmt.Init(op_code, ra, off)
}

func (this *Stmt) InitSbIdRiiStmt(
	op_code *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
	imm *expr.Expr,
) {
	this.stmt_type = SB_ID_RII

	this.sb_id_rii_stmt = new(sugar.SbIdRiiStmt)
	this.sb_id_rii_stmt.Init(op_code, ra, off, imm)
}

func (this *Stmt) InitSbRirStmt(op_code *expr.Expr, ra *expr.Expr, off *expr.Expr, rb *expr.Expr) {
	this.stmt_type = SB_RIR

	this.sb_rir_stmt = new(sugar.SbRirStmt)
	this.sb_rir_stmt.Init(op_code, ra, off, rb)
}

func (this *Stmt) InitSdRidStmt(
	op_code *expr.Expr,
	ra *expr.Expr,
	off *expr.Expr,
	db *lexer.Token,
) {
	this.stmt_type = SD_RID

	this.sd_rid_stmt = new(sugar.SdRidStmt)
	this.sd_rid_stmt.Init(op_code, ra, off, db)
}

func (this *Stmt) InitStopStmt() {
	this.stmt_type = STOP

	this.stop_stmt = new(sugar.StopStmt)
	this.stop_stmt.Init()
}

func (this *Stmt) InitTimeCfgRStmt(ra *expr.Expr) {
	this.stmt_type = TIME_CFG_R

	this.time_cfg_r_stmt = new(sugar.TimeCfgRStmt)
	this.time_cfg_r_stmt.Init(ra)
}

func (this *Stmt) InitLabelStmt(expr_ *expr.Expr) {
	this.stmt_type = LABEL

	this.label_stmt = new(LabelStmt)
	this.label_stmt.Init(expr_)
}

func (this *Stmt) StmtType() StmtType {
	return this.stmt_type
}

func (this *Stmt) AddrsigStmt() *directive.AddrsigStmt {
	return this.addrsig_stmt
}

func (this *Stmt) AddrsigSymStmt() *directive.AddrsigSymStmt {
	return this.addrsig_sym_stmt
}

func (this *Stmt) AsciiStmt() *directive.AsciiStmt {
	return this.ascii_stmt
}

func (this *Stmt) AscizStmt() *directive.AscizStmt {
	return this.asciz_stmt
}

func (this *Stmt) ByteStmt() *directive.ByteStmt {
	return this.byte_stmt
}

func (this *Stmt) CfiDefCfaOffsetStmt() *directive.CfiDefCfaOffsetStmt {
	return this.cfi_def_cfa_offset_stmt
}

func (this *Stmt) CfiEndprocStmt() *directive.CfiEndprocStmt {
	return this.cfi_endproc_stmt
}

func (this *Stmt) CfiOffsetStmt() *directive.CfiOffsetStmt {
	return this.cfi_offset_stmt
}

func (this *Stmt) CfiSectionsStmt() *directive.CfiSectionsStmt {
	return this.cfi_sections_stmt
}

func (this *Stmt) CfiStartprocStmt() *directive.CfiStartprocStmt {
	return this.cfi_startproc_stmt
}

func (this *Stmt) FileNumberStmt() *directive.FileNumberStmt {
	return this.file_number_stmt
}

func (this *Stmt) FileStringStmt() *directive.FileStringStmt {
	return this.file_string_stmt
}

func (this *Stmt) GlobalStmt() *directive.GlobalStmt {
	return this.global_stmt
}

func (this *Stmt) LocIsStmtStmt() *directive.LocIsStmtStmt {
	return this.loc_is_stmt_stmt
}

func (this *Stmt) LocNumberStmt() *directive.LocNumberStmt {
	return this.loc_number_stmt
}

func (this *Stmt) LocPrologueEndStmt() *directive.LocPrologueEndStmt {
	return this.loc_prologue_end_stmt
}

func (this *Stmt) LongProgramCounterStmt() *directive.LongProgramCounterStmt {
	return this.long_program_counter_stmt
}

func (this *Stmt) LongSectionNameStmt() *directive.LongSectionNameStmt {
	return this.long_section_name_stmt
}

func (this *Stmt) P2AlignStmt() *directive.P2AlignStmt {
	return this.p2_align_stmt
}

func (this *Stmt) QuadStmt() *directive.QuadStmt {
	return this.quad_stmt
}

func (this *Stmt) SectionIdentifierNumberStmt() *directive.SectionIdentifierNumberStmt {
	return this.section_identifier_number_stmt
}

func (this *Stmt) SectionIdentifierStmt() *directive.SectionIdentifierStmt {
	return this.section_identifier_stmt
}

func (this *Stmt) SectionStackSizesStmt() *directive.SectionStackSizesStmt {
	return this.section_stack_sizes_stmt
}

func (this *Stmt) SectionStringNumberStmt() *directive.SectionStringNumberStmt {
	return this.section_string_number_stmt
}

func (this *Stmt) SectionStringStmt() *directive.SectionStringStmt {
	return this.section_string_stmt
}

func (this *Stmt) SetStmt() *directive.SetStmt {
	return this.set_stmt
}

func (this *Stmt) ShortStmt() *directive.ShortStmt {
	return this.short_stmt
}

func (this *Stmt) SizeStmt() *directive.SizeStmt {
	return this.size_stmt
}

func (this *Stmt) TextStmt() *directive.TextStmt {
	return this.text_stmt
}

func (this *Stmt) TypeStmt() *directive.TypeStmt {
	return this.type_stmt
}

func (this *Stmt) WeakStmt() *directive.WeakStmt {
	return this.weak_stmt
}

func (this *Stmt) ZeroSingleNumberStmt() *directive.ZeroSingleNumberStmt {
	return this.zero_single_number_stmt
}

func (this *Stmt) ZeroDoubleNumberStmt() *directive.ZeroDoubleNumberStmt {
	return this.zero_double_number_stmt
}

func (this *Stmt) CiStmt() *instruction.CiStmt {
	return this.ci_stmt
}

func (this *Stmt) DdciStmt() *instruction.DdciStmt {
	return this.ddci_stmt
}

func (this *Stmt) DmaRriStmt() *instruction.DmaRriStmt {
	return this.dma_rri_stmt
}

func (this *Stmt) DrdiciStmt() *instruction.DrdiciStmt {
	return this.drdici_stmt
}

func (this *Stmt) EdriStmt() *instruction.EdriStmt {
	return this.edri_stmt
}

func (this *Stmt) EridStmt() *instruction.EridStmt {
	return this.erid_stmt
}

func (this *Stmt) EriiStmt() *instruction.EriiStmt {
	return this.erii_stmt
}

func (this *Stmt) ErirStmt() *instruction.ErirStmt {
	return this.erir_stmt
}

func (this *Stmt) ErriStmt() *instruction.ErriStmt {
	return this.erri_stmt
}

func (this *Stmt) IStmt() *instruction.IStmt {
	return this.i_stmt
}

func (this *Stmt) NopStmt() *instruction.NopStmt {
	return this.nop_stmt
}

func (this *Stmt) RStmt() *instruction.RStmt {
	return this.r_stmt
}

func (this *Stmt) RciStmt() *instruction.RciStmt {
	return this.rci_stmt
}

func (this *Stmt) RiciStmt() *instruction.RiciStmt {
	return this.rici_stmt
}

func (this *Stmt) RirStmt() *instruction.RirStmt {
	return this.rir_stmt
}

func (this *Stmt) RircStmt() *instruction.RircStmt {
	return this.rirc_stmt
}

func (this *Stmt) RirciStmt() *instruction.RirciStmt {
	return this.rirci_stmt
}

func (this *Stmt) RrStmt() *instruction.RrStmt {
	return this.rr_stmt
}

func (this *Stmt) RrcStmt() *instruction.RrcStmt {
	return this.rrc_stmt
}

func (this *Stmt) RrciStmt() *instruction.RrciStmt {
	return this.rrci_stmt
}

func (this *Stmt) RriStmt() *instruction.RriStmt {
	return this.rri_stmt
}

func (this *Stmt) RricStmt() *instruction.RricStmt {
	return this.rric_stmt
}

func (this *Stmt) RriciStmt() *instruction.RriciStmt {
	return this.rrici_stmt
}

func (this *Stmt) RrrStmt() *instruction.RrrStmt {
	return this.rrr_stmt
}

func (this *Stmt) RrrcStmt() *instruction.RrrcStmt {
	return this.rrrc_stmt
}

func (this *Stmt) RrrciStmt() *instruction.RrrciStmt {
	return this.rrrci_stmt
}

func (this *Stmt) RrriStmt() *instruction.RrriStmt {
	return this.rrri_stmt
}

func (this *Stmt) RrriciStmt() *instruction.RrriciStmt {
	return this.rrrici_stmt
}

func (this *Stmt) SErriStmt() *instruction.SErriStmt {
	return this.s_erri_stmt
}

func (this *Stmt) SRStmt() *instruction.SRStmt {
	return this.s_r_stmt
}

func (this *Stmt) SRciStmt() *instruction.SRciStmt {
	return this.s_rci_stmt
}

func (this *Stmt) SRircStmt() *instruction.SRircStmt {
	return this.s_rirc_stmt
}

func (this *Stmt) SRirciStmt() *instruction.SRirciStmt {
	return this.s_rirci_stmt
}

func (this *Stmt) SRrStmt() *instruction.SRrStmt {
	return this.s_rr_stmt
}

func (this *Stmt) SRrcStmt() *instruction.SRrcStmt {
	return this.s_rrc_stmt
}

func (this *Stmt) SRrciStmt() *instruction.SRrciStmt {
	return this.s_rrci_stmt
}

func (this *Stmt) SRriStmt() *instruction.SRriStmt {
	return this.s_rri_stmt
}

func (this *Stmt) SRricStmt() *instruction.SRricStmt {
	return this.s_rric_stmt
}

func (this *Stmt) SRriciStmt() *instruction.SRriciStmt {
	return this.s_rrici_stmt
}

func (this *Stmt) SRrrStmt() *instruction.SRrrStmt {
	return this.s_rrr_stmt
}

func (this *Stmt) SRrrcStmt() *instruction.SRrrcStmt {
	return this.s_rrrc_stmt
}

func (this *Stmt) SRrrciStmt() *instruction.SRrrciStmt {
	return this.s_rrrci_stmt
}

func (this *Stmt) SRrriStmt() *instruction.SRrriStmt {
	return this.s_rrri_stmt
}

func (this *Stmt) SRrriciStmt() *instruction.SRrriciStmt {
	return this.s_rrrici_stmt
}

func (this *Stmt) BkpStmt() *sugar.BkpStmt {
	return this.bkp_stmt
}

func (this *Stmt) BootRiStmt() *sugar.BootRiStmt {
	return this.boot_ri_stmt
}

func (this *Stmt) CallRiStmt() *sugar.CallRiStmt {
	return this.call_ri_stmt
}

func (this *Stmt) CallRrStmt() *sugar.CallRrStmt {
	return this.call_rr_stmt
}

func (this *Stmt) DivStepDrdiStmt() *sugar.DivStepDrdiStmt {
	return this.div_step_drdi_stmt
}

func (this *Stmt) JeqRiiStmt() *sugar.JeqRiiStmt {
	return this.jeq_rii_stmt
}

func (this *Stmt) JeqRriStmt() *sugar.JeqRriStmt {
	return this.jeq_rri_stmt
}

func (this *Stmt) JnzRiStmt() *sugar.JnzRiStmt {
	return this.jnz_ri_stmt
}

func (this *Stmt) JumpIStmt() *sugar.JumpIStmt {
	return this.jump_i_stmt
}

func (this *Stmt) JumpRStmt() *sugar.JumpRStmt {
	return this.jump_r_stmt
}

func (this *Stmt) LbsRriStmt() *sugar.LbsRriStmt {
	return this.lbs_rri_stmt
}

func (this *Stmt) LbsSRriStmt() *sugar.LbsSRriStmt {
	return this.lbs_s_rri_stmt
}

func (this *Stmt) LdDriStmt() *sugar.LdDriStmt {
	return this.ld_dri_stmt
}

func (this *Stmt) MovdDdStmt() *sugar.MovdDdStmt {
	return this.movd_dd_stmt
}

func (this *Stmt) MoveRiStmt() *sugar.MoveRiStmt {
	return this.move_ri_stmt
}

func (this *Stmt) MoveRiciStmt() *sugar.MoveRiciStmt {
	return this.move_rici_stmt
}

func (this *Stmt) MoveSRiStmt() *sugar.MoveSRiStmt {
	return this.move_s_ri_stmt
}

func (this *Stmt) MoveSRiciStmt() *sugar.MoveSRiciStmt {
	return this.move_s_rici_stmt
}

func (this *Stmt) SbIdRiStmt() *sugar.SbIdRiStmt {
	return this.sb_id_ri_stmt
}

func (this *Stmt) SbIdRiiStmt() *sugar.SbIdRiiStmt {
	return this.sb_id_rii_stmt
}

func (this *Stmt) SbRirStmt() *sugar.SbRirStmt {
	return this.sb_rir_stmt
}

func (this *Stmt) SdRidStmt() *sugar.SdRidStmt {
	return this.sd_rid_stmt
}

func (this *Stmt) StopStmt() *sugar.StopStmt {
	return this.stop_stmt
}

func (this *Stmt) TimeCfgRStmt() *sugar.TimeCfgRStmt {
	return this.time_cfg_r_stmt
}

func (this *Stmt) LabelStmt() *LabelStmt {
	return this.label_stmt
}
