package parser

import (
	"errors"
	"uPIMulator/src/linker/parser/expr"
	"uPIMulator/src/linker/parser/stmt"
)

type ExprCallback func(*expr.Expr)
type StmtCallback func(*stmt.Stmt)

type Walker struct {
	expr_callbacks map[expr.ExprType]ExprCallback
	stmt_callbacks map[stmt.StmtType]StmtCallback
}

func (this *Walker) Init() {
	this.expr_callbacks = make(map[expr.ExprType]ExprCallback, 0)
	this.stmt_callbacks = make(map[stmt.StmtType]StmtCallback, 0)
}

func (this *Walker) RegisterExprCallback(expr_type expr.ExprType, expr_callback ExprCallback) {
	if _, found := this.expr_callbacks[expr_type]; found {
		err := errors.New("expr callbak is already registered")
		panic(err)
	}

	this.expr_callbacks[expr_type] = expr_callback
}

func (this *Walker) RegisterStmtCallback(stmt_type stmt.StmtType, stmt_callback StmtCallback) {
	if _, found := this.stmt_callbacks[stmt_type]; found {
		err := errors.New("stmt callbak is already registered")
		panic(err)
	}

	this.stmt_callbacks[stmt_type] = stmt_callback
}

func (this *Walker) Walk(ast *Ast) {
	for i := 0; i < ast.Size(); i++ {
		stmt_ := ast.Get(i)

		stmt_type := stmt_.StmtType()
		if stmt_type == stmt.ASCII {
			this.WalkAsciiStmt(stmt_)
		} else if stmt_type == stmt.ASCIZ {
			this.WalkAscizStmt(stmt_)
		} else if stmt_type == stmt.BYTE {
			this.WalkByteStmt(stmt_)
		} else if stmt_type == stmt.GLOBAL {
			this.WalkGlobalStmt(stmt_)
		} else if stmt_type == stmt.LONG_PROGRAM_COUNTER {
			this.WalkLongProgramCounterStmt(stmt_)
		} else if stmt_type == stmt.LONG_SECTION_NAME {
			this.WalkLongSectionNameStmt(stmt_)
		} else if stmt_type == stmt.P2_ALIGN {
			this.WalkP2AlignStmt(stmt_)
		} else if stmt_type == stmt.QUAD {
			this.WalkQuadStmt(stmt_)
		} else if stmt_type == stmt.SECTION_IDENTIFIER_NUMBER {
			this.WalkSectionIdentifierNumberStmt(stmt_)
		} else if stmt_type == stmt.SECTION_IDENTIFIER {
			this.WalkSectionIdentifierStmt(stmt_)
		} else if stmt_type == stmt.SECTION_STACK_SIZES {
			this.WalkSectionStackSizesStmt(stmt_)
		} else if stmt_type == stmt.SECTION_STRING_NUMBER {
			this.WalkSectionStringNumberStmt(stmt_)
		} else if stmt_type == stmt.SECTION_STRING {
			this.WalkSectionStringStmt(stmt_)
		} else if stmt_type == stmt.SET {
			this.WalkSetStmt(stmt_)
		} else if stmt_type == stmt.SHORT {
			this.WalkShortStmt(stmt_)
		} else if stmt_type == stmt.SIZE {
			this.WalkSizeStmt(stmt_)
		} else if stmt_type == stmt.TEXT {
			this.WalkTextStmt(stmt_)
		} else if stmt_type == stmt.ZERO_DOUBLE_NUMBER {
			this.WalkZeroDoubleNumberStmt(stmt_)
		} else if stmt_type == stmt.ZERO_SINGLE_NUMBER {
			this.WalkZeroSingleNumberStmt(stmt_)
		} else if stmt_type == stmt.CI {
			this.WalkCiStmt(stmt_)
		} else if stmt_type == stmt.DDCI {
			this.WalkDdciStmt(stmt_)
		} else if stmt_type == stmt.DMA_RRI {
			this.WalkDmaRriStmt(stmt_)
		} else if stmt_type == stmt.DRDICI {
			this.WalkDrdiciStmt(stmt_)
		} else if stmt_type == stmt.EDRI {
			this.WalkEdriStmt(stmt_)
		} else if stmt_type == stmt.ERID {
			this.WalkEridStmt(stmt_)
		} else if stmt_type == stmt.ERII {
			this.WalkEriiStmt(stmt_)
		} else if stmt_type == stmt.ERIR {
			this.WalkErirStmt(stmt_)
		} else if stmt_type == stmt.ERRI {
			this.WalkErriStmt(stmt_)
		} else if stmt_type == stmt.I {
			this.WalkIStmt(stmt_)
		} else if stmt_type == stmt.NOP {
			this.WalkNopStmt(stmt_)
		} else if stmt_type == stmt.RCI {
			this.WalkRciStmt(stmt_)
		} else if stmt_type == stmt.RICI {
			this.WalkRiciStmt(stmt_)
		} else if stmt_type == stmt.RIRCI {
			this.WalkRirciStmt(stmt_)
		} else if stmt_type == stmt.RIRC {
			this.WalkRircStmt(stmt_)
		} else if stmt_type == stmt.RIR {
			this.WalkRirStmt(stmt_)
		} else if stmt_type == stmt.RRCI {
			this.WalkRrciStmt(stmt_)
		} else if stmt_type == stmt.RRC {
			this.WalkRrcStmt(stmt_)
		} else if stmt_type == stmt.RRICI {
			this.WalkRriciStmt(stmt_)
		} else if stmt_type == stmt.RRIC {
			this.WalkRricStmt(stmt_)
		} else if stmt_type == stmt.RRI {
			this.WalkRriStmt(stmt_)
		} else if stmt_type == stmt.RRRCI {
			this.WalkRrrciStmt(stmt_)
		} else if stmt_type == stmt.RRRC {
			this.WalkRrrcStmt(stmt_)
		} else if stmt_type == stmt.RRRICI {
			this.WalkRrriciStmt(stmt_)
		} else if stmt_type == stmt.RRRI {
			this.WalkRrriStmt(stmt_)
		} else if stmt_type == stmt.RRR {
			this.WalkRrrStmt(stmt_)
		} else if stmt_type == stmt.RR {
			this.WalkRrStmt(stmt_)
		} else if stmt_type == stmt.R {
			this.WalkRStmt(stmt_)
		} else if stmt_type == stmt.S_ERRI {
			this.WalkSErriStmt(stmt_)
		} else if stmt_type == stmt.S_RCI {
			this.WalkSRciStmt(stmt_)
		} else if stmt_type == stmt.S_RIRCI {
			this.WalkSRirciStmt(stmt_)
		} else if stmt_type == stmt.S_RIRC {
			this.WalkSRircStmt(stmt_)
		} else if stmt_type == stmt.S_RRCI {
			this.WalkSRrciStmt(stmt_)
		} else if stmt_type == stmt.S_RRC {
			this.WalkSRrcStmt(stmt_)
		} else if stmt_type == stmt.S_RRICI {
			this.WalkSRriciStmt(stmt_)
		} else if stmt_type == stmt.S_RRIC {
			this.WalkSRricStmt(stmt_)
		} else if stmt_type == stmt.S_RRI {
			this.WalkSRriStmt(stmt_)
		} else if stmt_type == stmt.S_RRRCI {
			this.WalkSRrrciStmt(stmt_)
		} else if stmt_type == stmt.S_RRRC {
			this.WalkSRrrcStmt(stmt_)
		} else if stmt_type == stmt.S_RRRICI {
			this.WalkSRrriciStmt(stmt_)
		} else if stmt_type == stmt.S_RRRI {
			this.WalkSRrriStmt(stmt_)
		} else if stmt_type == stmt.S_RRR {
			this.WalkSRrrStmt(stmt_)
		} else if stmt_type == stmt.S_RR {
			this.WalkSRrStmt(stmt_)
		} else if stmt_type == stmt.S_R {
			this.WalkSRStmt(stmt_)
		} else if stmt_type == stmt.BKP {
			this.WalkBkpStmt(stmt_)
		} else if stmt_type == stmt.BOOT_RI {
			this.WalkBootRiStmt(stmt_)
		} else if stmt_type == stmt.CALL_RI {
			this.WalkCallRiStmt(stmt_)
		} else if stmt_type == stmt.CALL_RR {
			this.WalkCallRrStmt(stmt_)
		} else if stmt_type == stmt.DIV_STEP_DRDI {
			this.WalkDivStepDrdiStmt(stmt_)
		} else if stmt_type == stmt.JEQ_RII {
			this.WalkJeqRiiStmt(stmt_)
		} else if stmt_type == stmt.JEQ_RRI {
			this.WalkJeqRriStmt(stmt_)
		} else if stmt_type == stmt.JNZ_RI {
			this.WalkJnzRiStmt(stmt_)
		} else if stmt_type == stmt.JUMP_I {
			this.WalkJumpIStmt(stmt_)
		} else if stmt_type == stmt.JUMP_R {
			this.WalkJumpRStmt(stmt_)
		} else if stmt_type == stmt.LBS_RRI {
			this.WalkLbsRriStmt(stmt_)
		} else if stmt_type == stmt.LBS_S_RRI {
			this.WalkLbsSRriStmt(stmt_)
		} else if stmt_type == stmt.LD_DRI {
			this.WalkLdDriStmt(stmt_)
		} else if stmt_type == stmt.MOVD_DD {
			this.WalkMovdDdStmt(stmt_)
		} else if stmt_type == stmt.MOVE_RICI {
			this.WalkMoveRiciStmt(stmt_)
		} else if stmt_type == stmt.MOVE_RI {
			this.WalkMoveRiStmt(stmt_)
		} else if stmt_type == stmt.MOVE_S_RICI {
			this.WalkMoveSRiciStmt(stmt_)
		} else if stmt_type == stmt.MOVE_S_RI {
			this.WalkMoveSRiStmt(stmt_)
		} else if stmt_type == stmt.SB_ID_RII {
			this.WalkSbIdRiiStmt(stmt_)
		} else if stmt_type == stmt.SB_ID_RI {
			this.WalkSbIdRiStmt(stmt_)
		} else if stmt_type == stmt.SB_RIR {
			this.WalkSbRirStmt(stmt_)
		} else if stmt_type == stmt.SD_RID {
			this.WalkSdRidStmt(stmt_)
		} else if stmt_type == stmt.STOP {
			this.WalkStopStmt(stmt_)
		} else if stmt_type == stmt.TIME_CFG_R {
			this.WalkTimeCfgRStmt(stmt_)
		} else if stmt_type == stmt.LABEL {
			this.WalkLabelStmt(stmt_)
		} else {
			continue
		}
	}
}

func (this *Walker) WalkBinaryAddExpr(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.BINARY_ADD {
		err := errors.New("expr type is not binary add expr")
		panic(err)
	}

	if expr_callback, found := this.expr_callbacks[expr.BINARY_ADD]; found {
		expr_callback(expr_)
	}

	binary_add_expr := expr_.BinaryAddExpr()

	operand1 := binary_add_expr.Operand1()
	operand2 := binary_add_expr.Operand2()

	this.WalkPrimaryExpr(operand1)
	this.WalkPrimaryExpr(operand2)
}

func (this *Walker) WalkBinarySubExpr(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.BINARY_SUB {
		err := errors.New("expr type is not binary sub expr")
		panic(err)
	}

	if expr_callback, found := this.expr_callbacks[expr.BINARY_SUB]; found {
		expr_callback(expr_)
	}

	binary_sub_expr := expr_.BinarySubExpr()

	operand1 := binary_sub_expr.Operand1()
	operand2 := binary_sub_expr.Operand2()

	this.WalkPrimaryExpr(operand1)
	this.WalkPrimaryExpr(operand2)
}

func (this *Walker) WalkNegativeNumberExpr(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.NEGATIVE_NUMBER {
		err := errors.New("expr type is not negative number expr")
		panic(err)
	}

	if expr_callback, found := this.expr_callbacks[expr.NEGATIVE_NUMBER]; found {
		expr_callback(expr_)
	}
}

func (this *Walker) WalkPrimaryExpr(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PRIMARY {
		err := errors.New("expr type is not primary")
		panic(err)
	}

	if expr_callback, found := this.expr_callbacks[expr.PRIMARY]; found {
		expr_callback(expr_)
	}
}

func (this *Walker) WalkProgramCounterExpr(expr_ *expr.Expr) {
	if expr_.ExprType() != expr.PROGRAM_COUNTER {
		err := errors.New("expr type is not program counter expr")
		panic(err)
	}

	if expr_callback, found := this.expr_callbacks[expr.PROGRAM_COUNTER]; found {
		expr_callback(expr_)
	}

	program_counter_expr := expr_.ProgramCounterExpr()

	child_expr := program_counter_expr.Expr()
	child_expr_type := child_expr.ExprType()

	if child_expr_type == expr.PRIMARY {
		this.WalkPrimaryExpr(child_expr)
	} else if child_expr_type == expr.NEGATIVE_NUMBER {
		this.WalkNegativeNumberExpr(child_expr)
	} else if child_expr_type == expr.BINARY_ADD {
		this.WalkBinaryAddExpr(child_expr)
	} else if child_expr_type == expr.BINARY_SUB {
		this.WalkBinarySubExpr(child_expr)
	} else {
		err := errors.New("child expr is not valid for a program counter expr")
		panic(err)
	}
}

func (this *Walker) WalkAsciiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ASCII {
		err := errors.New("stmt type is not an ASCII stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ASCII]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkAscizStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ASCIZ {
		err := errors.New("stmt type is not an ASCIZ stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ASCIZ]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkByteStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.BYTE {
		err := errors.New("stmt type is not a byte stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.BYTE]; found {
		stmt_callback(stmt_)
	}

	byte_stmt := stmt_.ByteStmt()

	expr_ := byte_stmt.Expr()

	this.WalkProgramCounterExpr(expr_)
}

func (this *Walker) WalkGlobalStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.GLOBAL {
		err := errors.New("stmt type is not a global stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.GLOBAL]; found {
		stmt_callback(stmt_)
	}

	global_stmt := stmt_.GlobalStmt()

	expr_ := global_stmt.Expr()

	this.WalkProgramCounterExpr(expr_)
}

func (this *Walker) WalkLongProgramCounterStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.LONG_PROGRAM_COUNTER {
		err := errors.New("stmt type is not a long program counter stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.LONG_PROGRAM_COUNTER]; found {
		stmt_callback(stmt_)
	}

	long_program_counter_stmt := stmt_.LongProgramCounterStmt()

	expr_ := long_program_counter_stmt.Expr()

	this.WalkProgramCounterExpr(expr_)
}

func (this *Walker) WalkLongSectionNameStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.LONG_SECTION_NAME {
		err := errors.New("stmt type is not a long section name stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.LONG_SECTION_NAME]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkP2AlignStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.P2_ALIGN {
		err := errors.New("stmt type is not a p2align stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.P2_ALIGN]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkQuadStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.QUAD {
		err := errors.New("stmt type is not a quad stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.QUAD]; found {
		stmt_callback(stmt_)
	}

	quad_stmt := stmt_.QuadStmt()

	expr_ := quad_stmt.Expr()

	this.WalkProgramCounterExpr(expr_)
}

func (this *Walker) WalkSectionIdentifierNumberStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SECTION_IDENTIFIER_NUMBER {
		err := errors.New("stmt type is not a section identifier number stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SECTION_IDENTIFIER_NUMBER]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSectionIdentifierStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SECTION_IDENTIFIER {
		err := errors.New("stmt type is not a section identifier stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SECTION_IDENTIFIER]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSectionStackSizesStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SECTION_STACK_SIZES {
		err := errors.New("stmt type is not a section stack sizes stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SECTION_STACK_SIZES]; found {
		stmt_callback(stmt_)
	}
}
func (this *Walker) WalkSectionStringNumberStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SECTION_STRING_NUMBER {
		err := errors.New("stmt type is not a section string number stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SECTION_STRING_NUMBER]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSectionStringStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SECTION_STRING {
		err := errors.New("stmt type is not a section string stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SECTION_STRING]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSetStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SET {
		err := errors.New("stmt type is not a set stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SET]; found {
		stmt_callback(stmt_)
	}

	set_stmt := stmt_.SetStmt()

	expr1 := set_stmt.Expr1()
	expr2 := set_stmt.Expr2()

	this.WalkProgramCounterExpr(expr1)
	this.WalkProgramCounterExpr(expr2)
}

func (this *Walker) WalkShortStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SHORT {
		err := errors.New("stmt type is not a short stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SHORT]; found {
		stmt_callback(stmt_)
	}

	short_stmt := stmt_.ShortStmt()

	expr_ := short_stmt.Expr()

	this.WalkProgramCounterExpr(expr_)
}

func (this *Walker) WalkSizeStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SIZE {
		err := errors.New("stmt type is not a size stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SIZE]; found {
		stmt_callback(stmt_)
	}

	size_stmt := stmt_.SizeStmt()

	expr1 := size_stmt.Expr1()
	expr2 := size_stmt.Expr2()

	this.WalkProgramCounterExpr(expr1)
	this.WalkProgramCounterExpr(expr2)
}

func (this *Walker) WalkTextStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.TEXT {
		err := errors.New("stmt type is not a text stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.TEXT]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkZeroDoubleNumberStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ZERO_DOUBLE_NUMBER {
		err := errors.New("stmt type is not a zero double number stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ZERO_DOUBLE_NUMBER]; found {
		stmt_callback(stmt_)
	}

	zero_double_number_stmt := stmt_.ZeroDoubleNumberStmt()

	expr1 := zero_double_number_stmt.Expr1()
	expr2 := zero_double_number_stmt.Expr2()

	this.WalkProgramCounterExpr(expr1)
	this.WalkProgramCounterExpr(expr2)
}

func (this *Walker) WalkZeroSingleNumberStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ZERO_SINGLE_NUMBER {
		err := errors.New("stmt type is not a zero single number stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ZERO_SINGLE_NUMBER]; found {
		stmt_callback(stmt_)
	}

	zero_single_number_stmt := stmt_.ZeroSingleNumberStmt()

	expr_ := zero_single_number_stmt.Expr()

	this.WalkProgramCounterExpr(expr_)
}

func (this *Walker) WalkCiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.CI {
		err := errors.New("stmt type is not a CI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.CI]; found {
		stmt_callback(stmt_)
	}

	ci_stmt := stmt_.CiStmt()

	pc := ci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkDdciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.DDCI {
		err := errors.New("stmt type is not a DDCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.DDCI]; found {
		stmt_callback(stmt_)
	}

	ddci_stmt := stmt_.DdciStmt()

	pc := ddci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkDmaRriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.DMA_RRI {
		err := errors.New("stmt type is not a DMA_RRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.DMA_RRI]; found {
		stmt_callback(stmt_)
	}

	dma_rri_stmt := stmt_.DmaRriStmt()

	imm := dma_rri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkDrdiciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.DRDICI {
		err := errors.New("stmt type is not a DRDICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.DRDICI]; found {
		stmt_callback(stmt_)
	}

	drdici_stmt := stmt_.DrdiciStmt()

	imm := drdici_stmt.Imm()
	pc := drdici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkEdriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.EDRI {
		err := errors.New("stmt type is not an EDRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.EDRI]; found {
		stmt_callback(stmt_)
	}

	edri_stmt := stmt_.EdriStmt()

	off := edri_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkEridStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ERID {
		err := errors.New("stmt type is not an ERID stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ERID]; found {
		stmt_callback(stmt_)
	}

	erid_stmt := stmt_.EridStmt()

	off := erid_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkEriiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ERII {
		err := errors.New("stmt type is not an ERII stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ERII]; found {
		stmt_callback(stmt_)
	}

	erii_stmt := stmt_.EriiStmt()

	off := erii_stmt.Off()
	imm := erii_stmt.Imm()

	this.WalkProgramCounterExpr(off)
	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkErirStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ERIR {
		err := errors.New("stmt type is not an ERIR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ERIR]; found {
		stmt_callback(stmt_)
	}

	erir_stmt := stmt_.ErirStmt()

	off := erir_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkErriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.ERRI {
		err := errors.New("stmt type is not an ERRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.ERRI]; found {
		stmt_callback(stmt_)
	}

	erri_stmt := stmt_.ErriStmt()

	off := erri_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkIStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.I {
		err := errors.New("stmt type is not an I stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.I]; found {
		stmt_callback(stmt_)
	}

	i_stmt := stmt_.IStmt()

	imm := i_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkNopStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.NOP {
		err := errors.New("stmt type is not a NOP stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.NOP]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkRciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RCI {
		err := errors.New("stmt type is not an RCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RCI]; found {
		stmt_callback(stmt_)
	}

	rci_stmt := stmt_.RciStmt()

	pc := rci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkRiciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RICI {
		err := errors.New("stmt type is not an RICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RICI]; found {
		stmt_callback(stmt_)
	}

	rici_stmt := stmt_.RiciStmt()

	imm := rici_stmt.Imm()
	pc := rici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkRirciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RIRCI {
		err := errors.New("stmt type is not an RIRCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RIRCI]; found {
		stmt_callback(stmt_)
	}

	rirci_stmt := stmt_.RirciStmt()

	imm := rirci_stmt.Imm()
	pc := rirci_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkRircStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RIRC {
		err := errors.New("stmt type is not an RIRC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RIRC]; found {
		stmt_callback(stmt_)
	}

	rirc_stmt := stmt_.RircStmt()

	imm := rirc_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkRirStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RIR {
		err := errors.New("stmt type is not an RIR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RIR]; found {
		stmt_callback(stmt_)
	}

	rir_stmt := stmt_.RirStmt()

	imm := rir_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkRrciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRCI {
		err := errors.New("stmt type is not an RRCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRCI]; found {
		stmt_callback(stmt_)
	}

	rrci_stmt := stmt_.RrciStmt()

	pc := rrci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkRrcStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRC {
		err := errors.New("stmt type is not an RRC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRC]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkRriciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRICI {
		err := errors.New("stmt type is not an RRICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRICI]; found {
		stmt_callback(stmt_)
	}

	rrici_stmt := stmt_.RriciStmt()

	imm := rrici_stmt.Imm()
	pc := rrici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkRricStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRIC {
		err := errors.New("stmt type is not an RRIC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRIC]; found {
		stmt_callback(stmt_)
	}

	rric_stmt := stmt_.RricStmt()

	imm := rric_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkRriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRI {
		err := errors.New("stmt type is not an RRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRI]; found {
		stmt_callback(stmt_)
	}

	rri_stmt := stmt_.RriStmt()

	imm := rri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkRrrciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRRCI {
		err := errors.New("stmt type is not an RRRCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRRCI]; found {
		stmt_callback(stmt_)
	}

	rrrci_stmt := stmt_.RrrciStmt()

	pc := rrrci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkRrrcStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRRC {
		err := errors.New("stmt type is not an RRRC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRRC]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkRrriciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRRICI {
		err := errors.New("stmt type is not an RRRICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRRICI]; found {
		stmt_callback(stmt_)
	}

	rrrici_stmt := stmt_.RrriciStmt()

	imm := rrrici_stmt.Imm()
	pc := rrrici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkRrriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRRI {
		err := errors.New("stmt type is not an RRRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRRI]; found {
		stmt_callback(stmt_)
	}

	rrri_stmt := stmt_.RrriStmt()

	imm := rrri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkRrrStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RRR {
		err := errors.New("stmt type is not an RRR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRR]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkRrStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.RR {
		err := errors.New("stmt type is not an RR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RR]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkRStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.R {
		err := errors.New("stmt type is not an R stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.R]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSErriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_ERRI {
		err := errors.New("stmt type is not a S_ERRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_ERRI]; found {
		stmt_callback(stmt_)
	}

	s_erri_stmt := stmt_.SErriStmt()

	off := s_erri_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkSRciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RCI {
		err := errors.New("stmt type is not a S_RCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RCI]; found {
		stmt_callback(stmt_)
	}

	s_rci_stmt := stmt_.SRciStmt()

	pc := s_rci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkSRirciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RIRCI {
		err := errors.New("stmt type is not a S_RIRCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RIRCI]; found {
		stmt_callback(stmt_)
	}

	s_rirci_stmt := stmt_.SRirciStmt()

	imm := s_rirci_stmt.Imm()
	pc := s_rirci_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkSRircStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RIRC {
		err := errors.New("stmt type is not a S_RIRC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RIRC]; found {
		stmt_callback(stmt_)
	}

	s_rirc_stmt := stmt_.SRircStmt()

	imm := s_rirc_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkSRrciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRCI {
		err := errors.New("stmt type is not a S_RRCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRCI]; found {
		stmt_callback(stmt_)
	}

	s_rrci_stmt := stmt_.SRrciStmt()

	pc := s_rrci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkSRrcStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRC {
		err := errors.New("stmt type is not a S_RRC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRC]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSRriciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRICI {
		err := errors.New("stmt type is not a S_RRICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRICI]; found {
		stmt_callback(stmt_)
	}

	s_rrici_stmt := stmt_.SRriciStmt()

	imm := s_rrici_stmt.Imm()
	pc := s_rrici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkSRricStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRIC {
		err := errors.New("stmt type is not a S_RRIC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRIC]; found {
		stmt_callback(stmt_)
	}

	s_rric_stmt := stmt_.SRricStmt()

	imm := s_rric_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkSRriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRI {
		err := errors.New("stmt type is not a S_RRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRI]; found {
		stmt_callback(stmt_)
	}

	s_rri_stmt := stmt_.SRriStmt()

	imm := s_rri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkSRrrciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRRCI {
		err := errors.New("stmt type is not an S_RRRCI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRRCI]; found {
		stmt_callback(stmt_)
	}

	s_rrrci_stmt := stmt_.SRrrciStmt()

	pc := s_rrrci_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkSRrrcStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRRC {
		err := errors.New("stmt type is not an S_RRRC stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.RRRC]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSRrriciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRRICI {
		err := errors.New("stmt type is not a S_RRRICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRRICI]; found {
		stmt_callback(stmt_)
	}

	s_rrrici_stmt := stmt_.SRrriciStmt()

	imm := s_rrrici_stmt.Imm()
	pc := s_rrrici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkSRrriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRRI {
		err := errors.New("stmt type is not a S_RRRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRRI]; found {
		stmt_callback(stmt_)
	}

	s_rrri_stmt := stmt_.SRrriStmt()

	imm := s_rrri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkSRrrStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RRR {
		err := errors.New("stmt type is not a RRR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RRR]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSRrStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_RR {
		err := errors.New("stmt type is not a S_RR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_RR]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkSRStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.S_R {
		err := errors.New("stmt type is not a S_R stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.S_R]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkBkpStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.BKP {
		err := errors.New("stmt type is not a BKP stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.BKP]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkBootRiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.BOOT_RI {
		err := errors.New("stmt type is not a BOOT_RI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.BOOT_RI]; found {
		stmt_callback(stmt_)
	}

	boot_ri_stmt := stmt_.BootRiStmt()

	imm := boot_ri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkCallRiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.CALL_RI {
		err := errors.New("stmt type is not a CALL_RI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.CALL_RI]; found {
		stmt_callback(stmt_)
	}

	call_ri_stmt := stmt_.CallRiStmt()

	imm := call_ri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkCallRrStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.CALL_RR {
		err := errors.New("stmt type is not a CALL_RR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.CALL_RR]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkDivStepDrdiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.DIV_STEP_DRDI {
		err := errors.New("stmt type is not a DIV_STEP_DRDI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.DIV_STEP_DRDI]; found {
		stmt_callback(stmt_)
	}

	div_step_drdi_stmt := stmt_.DivStepDrdiStmt()

	imm := div_step_drdi_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkJeqRiiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.JEQ_RII {
		err := errors.New("stmt type is not a JEQ_RII stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.JEQ_RII]; found {
		stmt_callback(stmt_)
	}

	jeq_rii_stmt := stmt_.JeqRiiStmt()

	imm := jeq_rii_stmt.Imm()
	pc := jeq_rii_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkJeqRriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.JEQ_RRI {
		err := errors.New("stmt type is not a JEQ_RRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.JEQ_RRI]; found {
		stmt_callback(stmt_)
	}

	jeq_rri_stmt := stmt_.JeqRriStmt()

	pc := jeq_rri_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkJnzRiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.JNZ_RI {
		err := errors.New("stmt type is not a JNZ_RI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.JNZ_RI]; found {
		stmt_callback(stmt_)
	}

	jnz_ri_stmt := stmt_.JnzRiStmt()

	pc := jnz_ri_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkJumpIStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.JUMP_I {
		err := errors.New("stmt type is not a JUMP_I stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.JUMP_I]; found {
		stmt_callback(stmt_)
	}

	jump_i_stmt := stmt_.JumpIStmt()

	pc := jump_i_stmt.Pc()

	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkJumpRStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.JUMP_R {
		err := errors.New("stmt type is not a JUMP_R stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.JUMP_R]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkLbsRriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.LBS_RRI {
		err := errors.New("stmt type is not a LBS_RRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.LBS_RRI]; found {
		stmt_callback(stmt_)
	}

	lbs_rri_stmt := stmt_.LbsRriStmt()

	off := lbs_rri_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkLbsSRriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.LBS_S_RRI {
		err := errors.New("stmt type is not a LBS_S_RRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.LBS_S_RRI]; found {
		stmt_callback(stmt_)
	}

	lbs_s_rri_stmt := stmt_.LbsSRriStmt()

	off := lbs_s_rri_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkLdDriStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.LD_DRI {
		err := errors.New("stmt type is not a LD_DRI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.LD_DRI]; found {
		stmt_callback(stmt_)
	}

	ld_dri_stmt := stmt_.LdDriStmt()

	off := ld_dri_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkMovdDdStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.MOVD_DD {
		err := errors.New("stmt type is not a MOVD_DD stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.MOVD_DD]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkMoveRiciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.MOVE_RICI {
		err := errors.New("stmt type is not a MOVE_RICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.MOVE_RICI]; found {
		stmt_callback(stmt_)
	}

	move_rici_stmt := stmt_.MoveRiciStmt()

	imm := move_rici_stmt.Imm()
	pc := move_rici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkMoveRiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.MOVE_RI {
		err := errors.New("stmt type is not a MOVE_RI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.MOVE_RI]; found {
		stmt_callback(stmt_)
	}

	move_ri_stmt := stmt_.MoveRiStmt()

	imm := move_ri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkMoveSRiciStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.MOVE_S_RICI {
		err := errors.New("stmt type is not a MOVE_S_RICI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.MOVE_S_RICI]; found {
		stmt_callback(stmt_)
	}

	move_s_rici_stmt := stmt_.MoveSRiciStmt()

	imm := move_s_rici_stmt.Imm()
	pc := move_s_rici_stmt.Pc()

	this.WalkProgramCounterExpr(imm)
	this.WalkProgramCounterExpr(pc)
}

func (this *Walker) WalkMoveSRiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.MOVE_S_RI {
		err := errors.New("stmt type is not a MOVE_S_RI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.MOVE_S_RI]; found {
		stmt_callback(stmt_)
	}

	move_s_ri_stmt := stmt_.MoveSRiStmt()

	imm := move_s_ri_stmt.Imm()

	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkSbIdRiiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SB_ID_RII {
		err := errors.New("stmt type is not a SB_ID_RII stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SB_ID_RII]; found {
		stmt_callback(stmt_)
	}

	sb_id_rii_stmt := stmt_.SbIdRiiStmt()

	off := sb_id_rii_stmt.Off()
	imm := sb_id_rii_stmt.Imm()

	this.WalkProgramCounterExpr(off)
	this.WalkProgramCounterExpr(imm)
}

func (this *Walker) WalkSbIdRiStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SB_ID_RI {
		err := errors.New("stmt type is not a SB_ID_RI stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SB_ID_RI]; found {
		stmt_callback(stmt_)
	}

	sb_id_ri_stmt := stmt_.SbIdRiStmt()

	off := sb_id_ri_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkSbRirStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SB_RIR {
		err := errors.New("stmt type is not a SB_RIR stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SB_RIR]; found {
		stmt_callback(stmt_)
	}

	sb_rir_stmt := stmt_.SbRirStmt()

	off := sb_rir_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkSdRidStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.SD_RID {
		err := errors.New("stmt type is not a SD_RID stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.SD_RID]; found {
		stmt_callback(stmt_)
	}

	sd_rid_stmt := stmt_.SdRidStmt()

	off := sd_rid_stmt.Off()

	this.WalkProgramCounterExpr(off)
}

func (this *Walker) WalkStopStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.STOP {
		err := errors.New("stmt type is not a STOP stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.STOP]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkTimeCfgRStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.TIME_CFG_R {
		err := errors.New("stmt type is not a TIME_CFG_R stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.TIME_CFG_R]; found {
		stmt_callback(stmt_)
	}
}

func (this *Walker) WalkLabelStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() != stmt.LABEL {
		err := errors.New("stmt type is not a LABEL stmt")
		panic(err)
	}

	if stmt_callback, found := this.stmt_callbacks[stmt.LABEL]; found {
		stmt_callback(stmt_)
	}

	label_stmt := stmt_.LabelStmt()

	expr_ := label_stmt.Expr()

	this.WalkProgramCounterExpr(expr_)
}
