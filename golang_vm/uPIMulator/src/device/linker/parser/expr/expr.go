package expr

import (
	"uPIMulator/src/device/linker/lexer"
)

type ExprType int

const (
	CI_OP_CODE ExprType = iota
	DDCI_OP_CODE
	DMA_RRI_OP_CODE
	DRDICI_OP_CODE
	I_OP_CODE
	JUMP_OP_CODE
	LOAD_OP_CODE
	R_OP_CODE
	RICI_OP_CODE
	RR_OP_CODE
	RRI_OP_CODE
	RRRI_OP_CODE
	STORE_OP_CODE

	SUFFIX
	CONDITION
	ENDIAN

	SECTION_NAME
	SECTION_TYPE

	SYMBOL_TYPE

	NEGATIVE_NUMBER
	PRIMARY
	BINARY_ADD
	BINARY_SUB
	PROGRAM_COUNTER

	SRC_REG
)

type Expr struct {
	expr_type ExprType

	ci_op_code_expr      *CiOpCodeExpr
	ddci_op_code_expr    *DdciOpCodeExpr
	dma_rri_op_code_expr *DmaRriOpCodeExpr
	drdici_op_code_expr  *DrdiciOpCodeExpr
	i_op_code_expr       *IOpCodeExpr
	jump_op_code_expr    *JumpOpCodeExpr
	load_op_code_expr    *LoadOpCodeExpr
	r_op_code_expr       *ROpCodeExpr
	rici_op_code_expr    *RiciOpCodeExpr
	rr_op_code_expr      *RrOpCodeExpr
	rri_op_code_expr     *RriOpCodeExpr
	rrri_op_code_expr    *RrriOpCodeExpr
	store_op_code_expr   *StoreOpCodeExpr

	suffix_expr    *SuffixExpr
	condition_expr *ConditionExpr
	endian_expr    *EndianExpr

	section_name_expr *SectionNameExpr
	section_type_expr *SectionTypeExpr

	symbol_type_expr *SymbolTypeExpr

	negative_number_expr *NegativeNumberExpr
	primary_expr         *PrimaryExpr
	binary_add_expr      *BinaryAddExpr
	binary_sub_expr      *BinarySubExpr
	program_counter_expr *ProgramCounterExpr

	src_reg_expr *SrcRegExpr
}

func (this *Expr) InitCiOpCodeExpr(token *lexer.Token) {
	this.expr_type = CI_OP_CODE

	this.ci_op_code_expr = new(CiOpCodeExpr)
	this.ci_op_code_expr.Init(token)
}

func (this *Expr) InitDdciOpCodeExpr(token *lexer.Token) {
	this.expr_type = DDCI_OP_CODE

	this.ddci_op_code_expr = new(DdciOpCodeExpr)
	this.ddci_op_code_expr.Init(token)
}

func (this *Expr) InitDmaRriOpCodeExpr(token *lexer.Token) {
	this.expr_type = DMA_RRI_OP_CODE

	this.dma_rri_op_code_expr = new(DmaRriOpCodeExpr)
	this.dma_rri_op_code_expr.Init(token)
}

func (this *Expr) InitDrdiciOpCodeExpr(token *lexer.Token) {
	this.expr_type = DRDICI_OP_CODE

	this.drdici_op_code_expr = new(DrdiciOpCodeExpr)
	this.drdici_op_code_expr.Init(token)
}

func (this *Expr) InitIOpCodeExpr(token *lexer.Token) {
	this.expr_type = I_OP_CODE

	this.i_op_code_expr = new(IOpCodeExpr)
	this.i_op_code_expr.Init(token)
}

func (this *Expr) InitJumpOpCodeExpr(token *lexer.Token) {
	this.expr_type = JUMP_OP_CODE

	this.jump_op_code_expr = new(JumpOpCodeExpr)
	this.jump_op_code_expr.Init(token)
}

func (this *Expr) InitLoadOpCodeExpr(token *lexer.Token) {
	this.expr_type = LOAD_OP_CODE

	this.load_op_code_expr = new(LoadOpCodeExpr)
	this.load_op_code_expr.Init(token)
}

func (this *Expr) InitROpCodeExpr(token *lexer.Token) {
	this.expr_type = R_OP_CODE

	this.r_op_code_expr = new(ROpCodeExpr)
	this.r_op_code_expr.Init(token)
}

func (this *Expr) InitRiciOpCodeExpr(token *lexer.Token) {
	this.expr_type = RICI_OP_CODE

	this.rici_op_code_expr = new(RiciOpCodeExpr)
	this.rici_op_code_expr.Init(token)
}

func (this *Expr) InitRrOpCodeExpr(token *lexer.Token) {
	this.expr_type = RR_OP_CODE

	this.rr_op_code_expr = new(RrOpCodeExpr)
	this.rr_op_code_expr.Init(token)
}

func (this *Expr) InitRriOpCodeExpr(token *lexer.Token) {
	this.expr_type = RRI_OP_CODE

	this.rri_op_code_expr = new(RriOpCodeExpr)
	this.rri_op_code_expr.Init(token)
}

func (this *Expr) InitRrriOpCodeExpr(token *lexer.Token) {
	this.expr_type = RRRI_OP_CODE

	this.rrri_op_code_expr = new(RrriOpCodeExpr)
	this.rrri_op_code_expr.Init(token)
}

func (this *Expr) InitStoreOpCodeExpr(token *lexer.Token) {
	this.expr_type = STORE_OP_CODE

	this.store_op_code_expr = new(StoreOpCodeExpr)
	this.store_op_code_expr.Init(token)
}

func (this *Expr) InitSuffixExpr(token *lexer.Token) {
	this.expr_type = SUFFIX

	this.suffix_expr = new(SuffixExpr)
	this.suffix_expr.Init(token)
}

func (this *Expr) InitConditionExpr(token *lexer.Token) {
	this.expr_type = CONDITION

	this.condition_expr = new(ConditionExpr)
	this.condition_expr.Init(token)
}

func (this *Expr) InitEndianExpr(token *lexer.Token) {
	this.expr_type = ENDIAN

	this.endian_expr = new(EndianExpr)
	this.endian_expr.Init(token)
}

func (this *Expr) InitSectionNameExpr(token *lexer.Token) {
	this.expr_type = SECTION_NAME

	this.section_name_expr = new(SectionNameExpr)
	this.section_name_expr.Init(token)
}

func (this *Expr) InitSectionTypeExpr(token *lexer.Token) {
	this.expr_type = SECTION_TYPE

	this.section_type_expr = new(SectionTypeExpr)
	this.section_type_expr.Init(token)
}

func (this *Expr) InitSymbolTypeExpr(token *lexer.Token) {
	this.expr_type = SYMBOL_TYPE

	this.symbol_type_expr = new(SymbolTypeExpr)
	this.symbol_type_expr.Init(token)
}

func (this *Expr) InitNegativeNumberExpr(token *lexer.Token) {
	this.expr_type = NEGATIVE_NUMBER

	this.negative_number_expr = new(NegativeNumberExpr)
	this.negative_number_expr.Init(token)
}

func (this *Expr) InitPrimaryExpr(token *lexer.Token) {
	this.expr_type = PRIMARY

	this.primary_expr = new(PrimaryExpr)
	this.primary_expr.Init(token)
}

func (this *Expr) InitBinaryAddExpr(operand1 *Expr, operand2 *Expr) {
	this.expr_type = BINARY_ADD

	this.binary_add_expr = new(BinaryAddExpr)
	this.binary_add_expr.Init(operand1, operand2)
}

func (this *Expr) InitBinarySubExpr(operand1 *Expr, operand2 *Expr) {
	this.expr_type = BINARY_SUB

	this.binary_sub_expr = new(BinarySubExpr)
	this.binary_sub_expr.Init(operand1, operand2)
}

func (this *Expr) InitProgramCounterExpr(expr *Expr) {
	this.expr_type = PROGRAM_COUNTER

	this.program_counter_expr = new(ProgramCounterExpr)
	this.program_counter_expr.Init(expr)
}

func (this *Expr) InitSrcRegExpr(token *lexer.Token) {
	this.expr_type = SRC_REG

	this.src_reg_expr = new(SrcRegExpr)
	this.src_reg_expr.Init(token)
}

func (this *Expr) ExprType() ExprType {
	return this.expr_type
}

func (this *Expr) CiOpCodeExpr() *CiOpCodeExpr {
	return this.ci_op_code_expr
}

func (this *Expr) DdciOpCodeExpr() *DdciOpCodeExpr {
	return this.ddci_op_code_expr
}

func (this *Expr) DmaRriOpCodeExpr() *DmaRriOpCodeExpr {
	return this.dma_rri_op_code_expr
}

func (this *Expr) DrdiciOpCodeExpr() *DrdiciOpCodeExpr {
	return this.drdici_op_code_expr
}

func (this *Expr) IOpCodeExpr() *IOpCodeExpr {
	return this.i_op_code_expr
}

func (this *Expr) JumpOpCodeExpr() *JumpOpCodeExpr {
	return this.jump_op_code_expr
}

func (this *Expr) LoadOpCodeExpr() *LoadOpCodeExpr {
	return this.load_op_code_expr
}

func (this *Expr) ROpCodeExpr() *ROpCodeExpr {
	return this.r_op_code_expr
}

func (this *Expr) RiciOpCodeExpr() *RiciOpCodeExpr {
	return this.rici_op_code_expr
}

func (this *Expr) RrOpCodeExpr() *RrOpCodeExpr {
	return this.rr_op_code_expr
}

func (this *Expr) RriOpCodeExpr() *RriOpCodeExpr {
	return this.rri_op_code_expr
}

func (this *Expr) RrriOpCodeExpr() *RrriOpCodeExpr {
	return this.rrri_op_code_expr
}

func (this *Expr) StoreOpCodeExpr() *StoreOpCodeExpr {
	return this.store_op_code_expr
}

func (this *Expr) SuffixExpr() *SuffixExpr {
	return this.suffix_expr
}

func (this *Expr) ConditionExpr() *ConditionExpr {
	return this.condition_expr
}

func (this *Expr) EndianExpr() *EndianExpr {
	return this.endian_expr
}

func (this *Expr) SectionNameExpr() *SectionNameExpr {
	return this.section_name_expr
}

func (this *Expr) SectionTypeExpr() *SectionTypeExpr {
	return this.section_type_expr
}

func (this *Expr) SymbolTypeExpr() *SymbolTypeExpr {
	return this.symbol_type_expr
}

func (this *Expr) NegativeNumberExpr() *NegativeNumberExpr {
	return this.negative_number_expr
}

func (this *Expr) PrimaryExpr() *PrimaryExpr {
	return this.primary_expr
}

func (this *Expr) BinaryAddExpr() *BinaryAddExpr {
	return this.binary_add_expr
}

func (this *Expr) BinarySubExpr() *BinarySubExpr {
	return this.binary_sub_expr
}

func (this *Expr) ProgramCounterExpr() *ProgramCounterExpr {
	return this.program_counter_expr
}

func (this *Expr) SrcRegExpr() *SrcRegExpr {
	return this.src_reg_expr
}
