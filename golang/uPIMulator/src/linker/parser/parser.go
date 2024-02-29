package parser

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
	"uPIMulator/src/linker/parser/stmt"
)

type Parser struct {
	stack *Stack
	table *Table
}

func (this *Parser) Init() {
	this.stack = new(Stack)
	this.stack.Init()

	this.table = new(Table)
	this.table.Init(this.stack)

	this.RegisterAccessExpr()
	this.RegisterConcatExpr()
	this.RegisterBinaryAddExpr()
	this.RegisterBinarySubExpr()
	this.RegisterCiOpCodeExpr()
	this.RegisterConditionExpr()
	this.RegisterDdciOpCodeExpr()
	this.RegisterDmaRriOpCodeExpr()
	this.RegisterDrdiciOpCodeExpr()
	this.RegisterEndianExpr()
	this.RegisterIOpCodeExpr()
	this.RegisterJumpOpCodeExpr()
	this.RegisterLoadOpCodeExpr()
	this.RegisterNegativeNumberExpr()
	this.RegisterPrimaryExpr()
	this.RegisterProgramCounterExpr()
	this.RegisterROpCodeExpr()
	this.RegisterRiciOpCodeExpr()
	this.RegisterRrOpCodeExpr()
	this.RegisterRriOpCodeExpr()
	this.RegisterRrriOpCodeExpr()
	this.RegisterSectionNameExpr()
	this.RegisterSectionTypeExpr()
	this.RegisterSrcRegExpr()
	this.RegisterStoreOpCodeExpr()
	this.RegisterSuffixExpr()
	this.RegisterSymbolTypeExpr()

	this.RegisterAddrsigStmt()
	this.RegisterAddrsigSymStmt()
	this.RegisterAsciiStmt()
	this.RegisterAscizStmt()
	this.RegisterByteStmt()
	this.RegisterCfiDefCfaOffsetStmt()
	this.RegisterCfiEndprocStmt()
	this.RegisterCfiOffsetStmt()
	this.RegisterCfiSectionsStmt()
	this.RegisterCfiStartprocStmt()
	this.RegisterFileNumberStmt()
	this.RegisterFileStringStmt()
	this.RegisterGlobalStmt()
	this.RegisterLocIsStmtStmt()
	this.RegisterLocNumberStmt()
	this.RegisterLocPrologueEndStmt()
	this.RegisterLongProgramCounterStmt()
	this.RegisterLongSectionNameStmt()
	this.RegisterP2AlignStmt()
	this.RegisterQuadStmt()
	this.RegisterSectionIdentifierNumberStmt()
	this.RegisterSectionIdentifierStmt()
	this.RegisterSectionStackSizesStmt()
	this.RegisterSectionStringNumberStmt()
	this.RegisterSectionStringStmt()
	this.RegisterSetStmt()
	this.RegisterShortStmt()
	this.RegisterSizeStmt()
	this.RegisterTextStmt()
	this.RegisterTypeStmt()
	this.RegisterWeakStmt()
	this.RegisterZeroDoubleNumberStmt()
	this.RegisterZeroSingleNumberStmt()

	this.RegisterCiStmt()
	this.RegisterDdciStmt()
	this.RegisterDmaRriStmt()
	this.RegisterDrdiciStmt()
	this.RegisterEdriStmt()
	this.RegisterEriiStmt()
	this.RegisterErirStmt()
	this.RegisterErriStmt()
	this.RegisterIStmt()
	this.RegisterRciStmt()
	this.RegisterRiciStmt()
	this.RegisterRirciStmt()
	this.RegisterRircStmt()
	this.RegisterRirStmt()
	this.RegisterRrciStmt()
	this.RegisterRrcStmt()
	this.RegisterRriciStmt()
	this.RegisterRricStmt()
	this.RegisterRriStmt()
	this.RegisterRrrciStmt()
	this.RegisterRrrcStmt()
	this.RegisterRrriciStmt()
	this.RegisterRrriStmt()
	this.RegisterRrrStmt()
	this.RegisterRrStmt()
	this.RegisterRStmt()

	this.RegisterSErriStmt()
	this.RegisterSRciStmt()
	this.RegisterSRirciStmt()
	this.RegisterSRrciStmt()
	this.RegisterSRrcStmt()
	this.RegisterSRriciStmt()
	this.RegisterSRricStmt()
	this.RegisterSRriStmt()
	this.RegisterSRrrciStmt()
	this.RegisterSRrrcStmt()
	this.RegisterSRrriciStmt()
	this.RegisterSRrriStmt()
	this.RegisterSRrrStmt()
	this.RegisterSRrStmt()
	this.RegisterSRStmt()

	this.RegisterNopStmt()
	this.RegisterBkpStmt()
	this.RegisterBootRiStmt()
	this.RegisterCallRiStmt()
	this.RegisterCallRrStmt()
	this.RegisterDivStepDrdiStmt()
	this.RegisterJeqRiiStmt()
	this.RegisterJeqRriStmt()
	this.RegisterJnzRiStmt()
	this.RegisterJumpIStmt()
	this.RegisterJumpRStmt()
	this.RegisterLbsRriStmt()
	this.RegisterLbsSRriStmt()
	this.RegisterLdDriStmt()
	this.RegisterMovdDdStmt()
	this.RegisterMoveRiciStmt()
	this.RegisterMoveRiStmt()
	this.RegisterMoveSRiciStmt()
	this.RegisterMoveSRiStmt()
	this.RegisterSbIdRiiStmt()
	this.RegisterSbIdRiStmt()
	this.RegisterSbRirStmt()
	this.RegisterSdRidStmt()
	this.RegisterStopStmt()
	this.RegisterTimeCfgRStmt()

	this.RegisterLabelStmt()
}

func (this *Parser) Parse(token_stream *lexer.TokenStream) *Ast {
	for pos := 0; pos < token_stream.Size(); pos++ {
		token := token_stream.Get(pos)

		this.ReduceExpr(token)

		if token.TokenType() != lexer.NEW_LINE {
			stack_item := new(StackItem)
			stack_item.InitToken(token)

			this.stack.Push(stack_item)
		} else {
			this.ReduceStmt(token)

			if !this.stack.AreStmts() {
				err := errors.New("stack are not stmts")
				panic(err)
			}
		}
	}

	if !this.stack.CanAccept() {
		err := errors.New("stack cannot be accepted")
		panic(err)
	}

	return this.stack.Accept()
}

func (this *Parser) ReduceExpr(token *lexer.Token) {
	for {
		reducible_expr_rule, stack_items := this.table.FindReducibleExprRule(token)

		if reducible_expr_rule != nil {
			stack_item := reducible_expr_rule.Reduce(stack_items, token)

			this.stack.Pop(len(stack_items))
			this.stack.Push(stack_item)
		} else {
			break
		}
	}
}

func (this *Parser) ReduceStmt(token *lexer.Token) {
	for {
		reducible_stmt_rule, stack_items := this.table.FindReducibleStmtRule(token)

		if reducible_stmt_rule != nil {
			stack_item := reducible_stmt_rule.Reduce(stack_items, token)

			this.stack.Pop(len(stack_items))
			this.stack.Push(stack_item)
		} else {
			break
		}
	}
}

func (this *Parser) RegisterAccessExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.IDENTIFIER &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.IDENTIFIER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token1 := stack_items[0].Token()
		token2 := stack_items[1].Token()

		attribute := token1.Attribute() + token2.Attribute()

		token := new(lexer.Token)
		token.Init(lexer.IDENTIFIER, attribute)

		stack_item := new(StackItem)
		stack_item.InitToken(token)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterConcatExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.NEGATIVE_NUMBER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		operand1 := stack_items[0].Expr()
		operand2 := stack_items[1].Expr()

		negative_number_expr := operand2.NegativeNumberExpr()
		token := negative_number_expr.Token()

		primary_expr := new(expr.Expr)
		primary_expr.InitPrimaryExpr(token)

		binary_sub_expr := new(expr.Expr)
		binary_sub_expr.InitBinarySubExpr(operand1, primary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(binary_sub_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterBinaryAddExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PLUS &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.PRIMARY {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		operand1 := stack_items[0].Expr()
		operand2 := stack_items[2].Expr()

		binary_add_expr := new(expr.Expr)
		binary_add_expr.InitBinaryAddExpr(operand1, operand2)

		stack_item := new(StackItem)
		stack_item.InitExpr(binary_add_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterBinarySubExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.MINUS &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.PRIMARY {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		operand1 := stack_items[0].Expr()
		operand2 := stack_items[2].Expr()

		binary_sub_expr := new(expr.Expr)
		binary_sub_expr.InitBinarySubExpr(operand1, operand2)

		stack_item := new(StackItem)
		stack_item.InitExpr(binary_sub_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterCiOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.STOP {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		ci_op_code_expr := new(expr.Expr)
		ci_op_code_expr.InitCiOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(ci_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterDdciOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.MOVD || token_type == lexer.SWAPD {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		ddci_op_code_expr := new(expr.Expr)
		ddci_op_code_expr.InitDdciOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(ddci_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterDmaRriOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.LDMA || token_type == lexer.LDMAI || token_type == lexer.SDMA {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		dma_rri_op_code_expr := new(expr.Expr)
		dma_rri_op_code_expr.InitDmaRriOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(dma_rri_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterDrdiciOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.DIV_STEP || token_type == lexer.MUL_STEP {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		drdici_op_code_expr := new(expr.Expr)
		drdici_op_code_expr.InitDrdiciOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(drdici_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterConditionExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.TRUE ||
					token_type == lexer.FALSE ||
					token_type == lexer.Z ||
					token_type == lexer.NZ ||
					token_type == lexer.E ||
					token_type == lexer.O ||
					token_type == lexer.PL ||
					token_type == lexer.MI ||
					token_type == lexer.OV ||
					token_type == lexer.NOV ||
					token_type == lexer.C ||
					token_type == lexer.NC ||
					token_type == lexer.SZ ||
					token_type == lexer.SNZ ||
					token_type == lexer.SPL ||
					token_type == lexer.SMI ||
					token_type == lexer.SO ||
					token_type == lexer.SE ||
					token_type == lexer.NC5 ||
					token_type == lexer.NC5 ||
					token_type == lexer.NC6 ||
					token_type == lexer.NC7 ||
					token_type == lexer.NC8 ||
					token_type == lexer.NC9 ||
					token_type == lexer.NC10 ||
					token_type == lexer.NC11 ||
					token_type == lexer.NC12 ||
					token_type == lexer.NC13 ||
					token_type == lexer.NC14 ||
					token_type == lexer.MAX ||
					token_type == lexer.NMAX ||
					token_type == lexer.SH32 ||
					token_type == lexer.NSH32 ||
					token_type == lexer.EQ ||
					token_type == lexer.NEQ ||
					token_type == lexer.LTU ||
					token_type == lexer.LEU ||
					token_type == lexer.GTU ||
					token_type == lexer.GEU ||
					token_type == lexer.LTS ||
					token_type == lexer.LES ||
					token_type == lexer.GTS ||
					token_type == lexer.GES ||
					token_type == lexer.XZ ||
					token_type == lexer.XNZ ||
					token_type == lexer.XLEU ||
					token_type == lexer.XGTU ||
					token_type == lexer.XLES ||
					token_type == lexer.XGTS ||
					token_type == lexer.SMALL ||
					token_type == lexer.LARGE {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		condition_expr := new(expr.Expr)
		condition_expr.InitConditionExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(condition_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterEndianExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.LITTLE || token_type == lexer.BIG {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		endian_expr := new(expr.Expr)
		endian_expr.InitEndianExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(endian_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterIOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.FAULT || token_type == lexer.BKP {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		i_op_code_expr := new(expr.Expr)
		i_op_code_expr.InitIOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(i_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterJumpOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.JEQ ||
					token_type == lexer.JNEQ ||
					token_type == lexer.JZ ||
					token_type == lexer.JNZ ||
					token_type == lexer.JLTU ||
					token_type == lexer.JGTU ||
					token_type == lexer.JLEU ||
					token_type == lexer.JGEU ||
					token_type == lexer.JLTS ||
					token_type == lexer.JGTS ||
					token_type == lexer.JLES ||
					token_type == lexer.JGES ||
					token_type == lexer.JUMP {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		jump_op_code_expr := new(expr.Expr)
		jump_op_code_expr.InitJumpOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(jump_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterLoadOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.LBS ||
					token_type == lexer.LBU ||
					token_type == lexer.LD ||
					token_type == lexer.LHS ||
					token_type == lexer.LHU ||
					token_type == lexer.LW {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		load_op_code_expr := new(expr.Expr)
		load_op_code_expr.InitLoadOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(load_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterNegativeNumberExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.MINUS &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.POSITIVIE_NUMBER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[1].Token()

		negativer_number_expr := new(expr.Expr)
		negativer_number_expr.InitNegativeNumberExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(negativer_number_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterPrimaryExpr() {
	precedence := map[lexer.TokenType]bool{lexer.IDENTIFIER: true}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.POSITIVIE_NUMBER || token_type == lexer.HEX_NUMBER || token_type == lexer.IDENTIFIER {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		primary_expr := new(expr.Expr)
		primary_expr.InitPrimaryExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(primary_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterProgramCounterExpr() {
	precedence := map[lexer.TokenType]bool{lexer.PLUS: true, lexer.MINUS: true}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR {
				expr_type := stack_items[0].Expr().ExprType()

				if expr_type == expr.PRIMARY ||
					expr_type == expr.NEGATIVE_NUMBER ||
					expr_type == expr.BINARY_ADD ||
					expr_type == expr.BINARY_SUB {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[0].Expr()

		program_counter_expr := new(expr.Expr)
		program_counter_expr.InitProgramCounterExpr(expr_)

		stack_item := new(StackItem)
		stack_item.InitExpr(program_counter_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterROpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.TIME || token_type == lexer.NOP {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		r_op_code_expr := new(expr.Expr)
		r_op_code_expr.InitROpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(r_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterRiciOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.ACQUIRE ||
					token_type == lexer.RELEASE ||
					token_type == lexer.BOOT ||
					token_type == lexer.RESUME {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		rici_op_code_expr := new(expr.Expr)
		rici_op_code_expr.InitRiciOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(rici_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterRrOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.CAO ||
					token_type == lexer.CLO ||
					token_type == lexer.CLS ||
					token_type == lexer.CLZ ||
					token_type == lexer.EXTSB ||
					token_type == lexer.EXTSH ||
					token_type == lexer.EXTUB ||
					token_type == lexer.EXTUH ||
					token_type == lexer.SATS ||
					token_type == lexer.TIME_CFG ||
					token_type == lexer.MOVE ||
					token_type == lexer.NEG ||
					token_type == lexer.NOT {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		rr_op_code_expr := new(expr.Expr)
		rr_op_code_expr.InitRrOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(rr_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterRriOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.ADD ||
					token_type == lexer.ADDC ||
					token_type == lexer.AND ||
					token_type == lexer.ANDN ||
					token_type == lexer.ASR ||
					token_type == lexer.CMPB4 ||
					token_type == lexer.LSL ||
					token_type == lexer.LSL1 ||
					token_type == lexer.LSL1X ||
					token_type == lexer.LSLX ||
					token_type == lexer.LSR ||
					token_type == lexer.LSR1 ||
					token_type == lexer.LSR1X ||
					token_type == lexer.LSRX ||
					token_type == lexer.MUL_SH_SH ||
					token_type == lexer.MUL_SH_SL ||
					token_type == lexer.MUL_SH_UH ||
					token_type == lexer.MUL_SH_UL ||
					token_type == lexer.MUL_SL_SH ||
					token_type == lexer.MUL_SL_SL ||
					token_type == lexer.MUL_SL_UH ||
					token_type == lexer.MUL_SL_UL ||
					token_type == lexer.MUL_UH_UH ||
					token_type == lexer.MUL_UH_UL ||
					token_type == lexer.MUL_UL_UH ||
					token_type == lexer.MUL_UL_UL ||
					token_type == lexer.NAND ||
					token_type == lexer.NOR ||
					token_type == lexer.NXOR ||
					token_type == lexer.OR ||
					token_type == lexer.ORN ||
					token_type == lexer.ROL ||
					token_type == lexer.ROR ||
					token_type == lexer.RSUB ||
					token_type == lexer.RSUBC ||
					token_type == lexer.SUB ||
					token_type == lexer.SUBC ||
					token_type == lexer.XOR ||
					token_type == lexer.CALL {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		rri_op_code_expr := new(expr.Expr)
		rri_op_code_expr.InitRriOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(rri_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterRrriOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.LSL_ADD ||
					token_type == lexer.LSL_SUB ||
					token_type == lexer.LSR_ADD ||
					token_type == lexer.ROL_ADD {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		rrri_op_code_expr := new(expr.Expr)
		rrri_op_code_expr.InitRrriOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(rrri_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterStoreOpCodeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.SB ||
					token_type == lexer.SB_ID ||
					token_type == lexer.SD ||
					token_type == lexer.SD_ID ||
					token_type == lexer.SH ||
					token_type == lexer.SH_ID ||
					token_type == lexer.SW ||
					token_type == lexer.SW_ID {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		store_op_code_expr := new(expr.Expr)
		store_op_code_expr.InitStoreOpCodeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(store_op_code_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterSuffixExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.S || token_type == lexer.U {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		suffix_expr := new(expr.Expr)
		suffix_expr.InitSuffixExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(suffix_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterSectionNameExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.ATOMIC ||
					token_type == lexer.BSS ||
					token_type == lexer.DATA ||
					token_type == lexer.DEBUG_ABBREV ||
					token_type == lexer.DEBUG_FRAME ||
					token_type == lexer.DEBUG_INFO ||
					token_type == lexer.DEBUG_LINE ||
					token_type == lexer.DEBUG_LOC ||
					token_type == lexer.DEBUG_RANGES ||
					token_type == lexer.DEBUG_STR ||
					token_type == lexer.DPU_HOST ||
					token_type == lexer.MRAM ||
					token_type == lexer.RODATA ||
					token_type == lexer.STACK_SIZES ||
					token_type == lexer.TEXT {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		section_name_expr := new(expr.Expr)
		section_name_expr.InitSectionNameExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(section_name_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterSectionTypeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.PROGBITS || token_type == lexer.NOBITS {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		section_type_expr := new(expr.Expr)
		section_type_expr.InitSectionTypeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(section_type_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterSrcRegExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.GP_REG ||
					token_type == lexer.ZERO_REG ||
					token_type == lexer.ONE ||
					token_type == lexer.ID ||
					token_type == lexer.ID2 ||
					token_type == lexer.ID4 ||
					token_type == lexer.ID8 ||
					token_type == lexer.LNEG ||
					token_type == lexer.MNEG {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		src_reg_expr := new(expr.Expr)
		src_reg_expr.InitSrcRegExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(src_reg_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterSymbolTypeExpr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN {
				token_type := stack_items[0].Token().TokenType()

				if token_type == lexer.FUNCTION || token_type == lexer.OBJECT {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		symbol_type_expr := new(expr.Expr)
		symbol_type_expr.InitSymbolTypeExpr(token)

		stack_item := new(StackItem)
		stack_item.InitExpr(symbol_type_expr)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddExprRule(rule)
}

func (this *Parser) RegisterAddrsigStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.ADDRSIG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		addrsig_stmt := new(stmt.Stmt)
		addrsig_stmt.InitAddrsigStmt()

		stack_item := new(StackItem)
		stack_item.InitStmt(addrsig_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterAddrsigSymStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.ADDRSIG_SYM &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		addrsig_sym_stmt := new(stmt.Stmt)
		addrsig_sym_stmt.InitAddrsigSymStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(addrsig_sym_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterAsciiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.ASCII &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.STRING {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[1].Token()

		ascii_stmt := new(stmt.Stmt)
		ascii_stmt.InitAsciiStmt(token)

		stack_item := new(StackItem)
		stack_item.InitStmt(ascii_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterAscizStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.ASCIZ &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.STRING {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[1].Token()

		asciz_stmt := new(stmt.Stmt)
		asciz_stmt.InitAscizStmt(token)

		stack_item := new(StackItem)
		stack_item.InitStmt(asciz_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterByteStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.BYTE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		byte_stmt := new(stmt.Stmt)
		byte_stmt.InitByteStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(byte_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCfiDefCfaOffsetStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.CFI_DEF_CFA_OFFSET &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		cfi_def_cfa_offset_stmt := new(stmt.Stmt)
		cfi_def_cfa_offset_stmt.InitCfiDefCfaOffsetStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(cfi_def_cfa_offset_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCfiEndprocStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.CFI_ENDPROC {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		cfi_endproc_stmt := new(stmt.Stmt)
		cfi_endproc_stmt.InitCfiEndprocStmt()

		stack_item := new(StackItem)
		stack_item.InitStmt(cfi_endproc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCfiOffsetStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.CFI_OFFSET &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[3].Expr()

		cfi_offset_stmt := new(stmt.Stmt)
		cfi_offset_stmt.InitCfiOffsetStmt(expr1, expr2)

		stack_item := new(StackItem)
		stack_item.InitStmt(cfi_offset_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCfiSectionsStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.CFI_SECTIONS &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SECTION_NAME {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		cfi_sections_stmt := new(stmt.Stmt)
		cfi_sections_stmt.InitCfiSectionsStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(cfi_sections_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCfiStartprocStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.CFI_STARTPROC {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		cfi_startproc_stmt := new(stmt.Stmt)
		cfi_startproc_stmt.InitCfiStartprocStmt()

		stack_item := new(StackItem)
		stack_item.InitStmt(cfi_startproc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterFileNumberStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.FILE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.STRING &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.STRING {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()
		token1 := stack_items[2].Token()
		token2 := stack_items[3].Token()

		file_number_stmt := new(stmt.Stmt)
		file_number_stmt.InitFileNumberStmt(expr_, token1, token2)

		stack_item := new(StackItem)
		stack_item.InitStmt(file_number_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterFileStringStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.FILE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.STRING {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[1].Token()

		file_string_stmt := new(stmt.Stmt)
		file_string_stmt.InitFileStringStmt(token)

		stack_item := new(StackItem)
		stack_item.InitStmt(file_string_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterGlobalStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.GLOBL &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		global_stmt := new(stmt.Stmt)
		global_stmt.InitGlobalStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(global_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLocIsStmtStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.LOC &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.IS_STMT &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[2].Expr()
		expr3 := stack_items[3].Expr()
		expr4 := stack_items[5].Expr()

		loc_is_stmt_stmt := new(stmt.Stmt)
		loc_is_stmt_stmt.InitLocIsStmtStmt(expr1, expr2, expr3, expr4)

		stack_item := new(StackItem)
		stack_item.InitStmt(loc_is_stmt_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLocNumberStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.LOC &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[2].Expr()
		expr3 := stack_items[3].Expr()

		loc_number_stmt := new(stmt.Stmt)
		loc_number_stmt.InitLocNumberStmt(expr1, expr2, expr3)

		stack_item := new(StackItem)
		stack_item.InitStmt(loc_number_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLocPrologueEndStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 5 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.LOC &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.PROLOGUE_END {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[2].Expr()
		expr3 := stack_items[3].Expr()

		loc_prologue_end_stmt := new(stmt.Stmt)
		loc_prologue_end_stmt.InitLocPrologueEndStmt(expr1, expr2, expr3)

		stack_item := new(StackItem)
		stack_item.InitStmt(loc_prologue_end_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLongProgramCounterStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.LONG &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		long_program_counter_stmt := new(stmt.Stmt)
		long_program_counter_stmt.InitLongProgramCounterStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(long_program_counter_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLongSectionNameStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.LONG &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SECTION_NAME {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		long_section_name_stmt := new(stmt.Stmt)
		long_section_name_stmt.InitLongSectionNameStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(long_section_name_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterP2AlignStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.P2ALIGN &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		p2align_stmt := new(stmt.Stmt)
		p2align_stmt.InitP2AlignStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(p2align_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterQuadStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.QUAD &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		quad_stmt := new(stmt.Stmt)
		quad_stmt.InitQuadStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(quad_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSectionIdentifierNumberStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SECTION &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SECTION_NAME &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.STRING &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SECTION_TYPE &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[2].Expr()
		token := stack_items[4].Token()
		expr3 := stack_items[6].Expr()
		expr4 := stack_items[8].Expr()

		section_identifier_number_stmt := new(stmt.Stmt)
		section_identifier_number_stmt.InitSectionIdentifierNumberStmt(
			expr1,
			expr2,
			token,
			expr3,
			expr4,
		)

		stack_item := new(StackItem)
		stack_item.InitStmt(section_identifier_number_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSectionIdentifierStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 7 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SECTION &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SECTION_NAME &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.STRING &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SECTION_TYPE {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[2].Expr()
		token := stack_items[4].Token()
		expr3 := stack_items[6].Expr()

		section_identifier_stmt := new(stmt.Stmt)
		section_identifier_stmt.InitSectionIdentifierStmt(expr1, expr2, token, expr3)

		stack_item := new(StackItem)
		stack_item.InitStmt(section_identifier_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSectionStackSizesStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SECTION &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SECTION_NAME &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.STRING &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SECTION_TYPE &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.SECTION_NAME &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[3].Token()
		expr1 := stack_items[5].Expr()
		expr2 := stack_items[7].Expr()
		expr3 := stack_items[8].Expr()

		section_stack_sizes_stmt := new(stmt.Stmt)
		section_stack_sizes_stmt.InitSectionStackSizesStmt(token, expr1, expr2, expr3)

		stack_item := new(StackItem)
		stack_item.InitStmt(section_stack_sizes_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSectionStringNumberStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SECTION &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SECTION_NAME &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.STRING &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SECTION_TYPE &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		token := stack_items[3].Token()
		expr2 := stack_items[5].Expr()
		expr3 := stack_items[7].Expr()

		section_string_number_stmt := new(stmt.Stmt)
		section_string_number_stmt.InitSectionStringNumberStmt(expr1, token, expr2, expr3)

		stack_item := new(StackItem)
		stack_item.InitStmt(section_string_number_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSectionStringStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SECTION &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SECTION_NAME &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.STRING &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SECTION_TYPE {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		token := stack_items[3].Token()
		expr2 := stack_items[5].Expr()

		section_string_stmt := new(stmt.Stmt)
		section_string_stmt.InitSectionStringStmt(expr1, token, expr2)

		stack_item := new(StackItem)
		stack_item.InitStmt(section_string_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSetStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SET &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[3].Expr()

		set_stmt := new(stmt.Stmt)
		set_stmt.InitSetStmt(expr1, expr2)

		stack_item := new(StackItem)
		stack_item.InitStmt(set_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterShortStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SHORT &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		short_stmt := new(stmt.Stmt)
		short_stmt.InitShortStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(short_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSizeStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SIZE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[3].Expr()

		size_stmt := new(stmt.Stmt)
		size_stmt.InitSizeStmt(expr1, expr2)

		stack_item := new(StackItem)
		stack_item.InitStmt(size_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterTextStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.SECTION_NAME {
				expr_ := stack_items[0].Expr()
				if expr_.SectionNameExpr().Token().TokenType() == lexer.TEXT {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		text_stmt := new(stmt.Stmt)
		text_stmt.InitTextStmt()

		stack_item := new(StackItem)
		stack_item.InitStmt(text_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterTypeStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.TYPE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SYMBOL_TYPE {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[3].Expr()

		type_stmt := new(stmt.Stmt)
		type_stmt.InitTypeStmt(expr1, expr2)

		stack_item := new(StackItem)
		stack_item.InitStmt(type_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterWeakStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.WEAK &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		weak_stmt := new(stmt.Stmt)
		weak_stmt.InitWeakStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(weak_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterZeroDoubleNumberStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.ZERO_DIRECTIVE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr1 := stack_items[1].Expr()
		expr2 := stack_items[3].Expr()

		zero_double_number_stmt := new(stmt.Stmt)
		zero_double_number_stmt.InitZeroDoubleNumberStmt(expr1, expr2)

		stack_item := new(StackItem)
		stack_item.InitStmt(zero_double_number_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterZeroSingleNumberStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.ZERO_DIRECTIVE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[1].Expr()

		zero_single_number_stmt := new(stmt.Stmt)
		zero_single_number_stmt.InitZeroSingleNumberStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(zero_single_number_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.CI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.CONDITION &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		condition := stack_items[1].Expr()
		pc := stack_items[3].Expr()

		ci_stmt := new(stmt.Stmt)
		ci_stmt.InitCiStmt(op_code, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(ci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterDdciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.DDCI_OP_CODE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.CONDITION &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		dc := stack_items[1].Token()
		db := stack_items[3].Token()
		condition := stack_items[5].Expr()
		pc := stack_items[7].Expr()

		ddci_stmt := new(stmt.Stmt)
		ddci_stmt.InitDdciStmt(op_code, dc, db, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(ddci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterDmaRriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.DMA_RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		rb := stack_items[3].Expr()
		pc := stack_items[5].Expr()

		dma_rri_stmt := new(stmt.Stmt)
		dma_rri_stmt.InitDmaRriStmt(op_code, ra, rb, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(dma_rri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterDrdiciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 12 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.DRDICI_OP_CODE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[8].StackItemType() == TOKEN &&
				stack_items[8].Token().TokenType() == lexer.COMMA &&
				stack_items[9].StackItemType() == EXPR &&
				stack_items[9].Expr().ExprType() == expr.CONDITION &&
				stack_items[10].StackItemType() == TOKEN &&
				stack_items[10].Token().TokenType() == lexer.COMMA &&
				stack_items[11].StackItemType() == EXPR &&
				stack_items[11].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		dc := stack_items[1].Token()
		ra := stack_items[3].Expr()
		db := stack_items[5].Token()
		imm := stack_items[7].Expr()
		condition := stack_items[9].Expr()
		pc := stack_items[11].Expr()

		drdici_stmt := new(stmt.Stmt)
		drdici_stmt.InitDrdiciStmt(op_code, dc, ra, db, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(drdici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterEdriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.LOAD_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.ENDIAN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		endian := stack_items[1].Expr()
		dc := stack_items[3].Token()
		ra := stack_items[5].Expr()
		off := stack_items[7].Expr()

		edri_stmt := new(stmt.Stmt)
		edri_stmt.InitEdriStmt(op_code, endian, dc, ra, off)

		stack_item := new(StackItem)
		stack_item.InitStmt(edri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterEridStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.STORE_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.ENDIAN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.PAIR_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		endian := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		off := stack_items[5].Expr()
		db := stack_items[7].Token()

		edri_stmt := new(stmt.Stmt)
		edri_stmt.InitEridStmt(op_code, endian, ra, off, db)

		stack_item := new(StackItem)
		stack_item.InitStmt(edri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterEriiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.STORE_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.ENDIAN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		endian := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		off := stack_items[5].Expr()
		imm := stack_items[7].Expr()

		erii_stmt := new(stmt.Stmt)
		erii_stmt.InitEriiStmt(op_code, endian, ra, off, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(erii_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterErirStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.STORE_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.ENDIAN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		endian := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		off := stack_items[5].Expr()
		rb := stack_items[7].Expr()

		erir_stmt := new(stmt.Stmt)
		erir_stmt.InitErirStmt(op_code, endian, ra, off, rb)

		stack_item := new(StackItem)
		stack_item.InitStmt(erir_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterErriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.LOAD_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.ENDIAN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		endian := stack_items[1].Expr()
		rc := stack_items[3].Expr()
		ra := stack_items[5].Expr()
		off := stack_items[7].Expr()

		erri_stmt := new(stmt.Stmt)
		erri_stmt.InitErriStmt(op_code, endian, rc, ra, off)

		stack_item := new(StackItem)
		stack_item.InitStmt(erri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterIStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.I_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		imm := stack_items[1].Expr()

		i_stmt := new(stmt.Stmt)
		i_stmt.InitIStmt(op_code, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(i_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.R_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.CONDITION &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		condition := stack_items[3].Expr()
		pc := stack_items[5].Expr()

		rci_stmt := new(stmt.Stmt)
		rci_stmt.InitRciStmt(op_code, rc, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(rci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRiciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RICI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.CONDITION &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		imm := stack_items[3].Expr()
		condition := stack_items[5].Expr()
		pc := stack_items[7].Expr()

		rici_stmt := new(stmt.Stmt)
		rici_stmt.InitRiciStmt(op_code, ra, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(rici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRirciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 10 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.CONDITION &&
				stack_items[8].StackItemType() == TOKEN &&
				stack_items[8].Token().TokenType() == lexer.COMMA &&
				stack_items[9].StackItemType() == EXPR &&
				stack_items[9].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		imm := stack_items[3].Expr()
		ra := stack_items[5].Expr()
		condition := stack_items[7].Expr()
		pc := stack_items[9].Expr()

		rirci_stmt := new(stmt.Stmt)
		rirci_stmt.InitRirciStmt(op_code, rc, imm, ra, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(rirci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRircStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		imm := stack_items[3].Expr()
		ra := stack_items[5].Expr()
		condition := stack_items[7].Expr()

		rirc_stmt := new(stmt.Stmt)
		rirc_stmt.InitRircStmt(op_code, rc, imm, ra, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(rirc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRirStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		imm := stack_items[3].Expr()
		ra := stack_items[5].Expr()

		rir_stmt := new(stmt.Stmt)
		rir_stmt.InitRirStmt(op_code, rc, imm, ra)

		stack_item := new(StackItem)
		stack_item.InitStmt(rir_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.CONDITION &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		condition := stack_items[5].Expr()
		pc := stack_items[7].Expr()

		rrci_stmt := new(stmt.Stmt)
		rrci_stmt.InitRrciStmt(op_code, rc, ra, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrcStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		condition := stack_items[5].Expr()

		rrc_stmt := new(stmt.Stmt)
		rrc_stmt.InitRrcStmt(op_code, rc, ra, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRriciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 10 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.CONDITION &&
				stack_items[8].StackItemType() == TOKEN &&
				stack_items[8].Token().TokenType() == lexer.COMMA &&
				stack_items[9].StackItemType() == EXPR &&
				stack_items[9].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		imm := stack_items[5].Expr()
		condition := stack_items[7].Expr()
		pc := stack_items[9].Expr()

		rrici_stmt := new(stmt.Stmt)
		rrici_stmt.InitRriciStmt(op_code, rc, ra, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRricStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		imm := stack_items[5].Expr()
		condition := stack_items[7].Expr()

		rric_stmt := new(stmt.Stmt)
		rric_stmt.InitRricStmt(op_code, rc, ra, imm, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(rric_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		imm := stack_items[5].Expr()

		rri_stmt := new(stmt.Stmt)
		rri_stmt.InitRriStmt(op_code, rc, ra, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(rri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrrciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 10 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.CONDITION &&
				stack_items[8].StackItemType() == TOKEN &&
				stack_items[8].Token().TokenType() == lexer.COMMA &&
				stack_items[9].StackItemType() == EXPR &&
				stack_items[9].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		rb := stack_items[5].Expr()
		condition := stack_items[7].Expr()
		pc := stack_items[9].Expr()

		rrrci_stmt := new(stmt.Stmt)
		rrrci_stmt.InitRrrciStmt(op_code, rc, ra, rb, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrrci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrrcStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		rb := stack_items[5].Expr()
		condition := stack_items[7].Expr()

		rrrc_stmt := new(stmt.Stmt)
		rrrc_stmt.InitRrrcStmt(op_code, rc, ra, rb, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrrc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrriciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 12 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[8].StackItemType() == TOKEN &&
				stack_items[8].Token().TokenType() == lexer.COMMA &&
				stack_items[9].StackItemType() == EXPR &&
				stack_items[9].Expr().ExprType() == expr.CONDITION &&
				stack_items[10].StackItemType() == TOKEN &&
				stack_items[10].Token().TokenType() == lexer.COMMA &&
				stack_items[11].StackItemType() == EXPR &&
				stack_items[11].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		rb := stack_items[5].Expr()
		imm := stack_items[7].Expr()
		condition := stack_items[9].Expr()
		pc := stack_items[11].Expr()

		rrrici_stmt := new(stmt.Stmt)
		rrrici_stmt.InitRrriciStmt(op_code, rc, ra, rb, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrrici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		rb := stack_items[5].Expr()
		imm := stack_items[7].Expr()

		rrri_stmt := new(stmt.Stmt)
		rrri_stmt.InitRrriStmt(op_code, rc, ra, rb, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrrStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		rb := stack_items[5].Expr()

		rrr_stmt := new(stmt.Stmt)
		rrr_stmt.InitRrrStmt(op_code, rc, ra, rb)

		stack_item := new(StackItem)
		stack_item.InitStmt(rrr_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRrStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()

		rr_stmt := new(stmt.Stmt)
		rr_stmt.InitRrStmt(op_code, rc, ra)

		stack_item := new(StackItem)
		stack_item.InitStmt(rr_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterRStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.R_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()

		r_stmt := new(stmt.Stmt)
		r_stmt.InitRStmt(op_code, rc)

		stack_item := new(StackItem)
		stack_item.InitStmt(r_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSErriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.LOAD_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[2].Expr().ExprType() == expr.ENDIAN &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		endian := stack_items[2].Expr()
		dc := stack_items[4].Token()
		ra := stack_items[6].Expr()
		off := stack_items[8].Expr()

		s_erri_stmt := new(stmt.Stmt)
		s_erri_stmt.InitSErriStmt(op_code, suffix, endian, dc, ra, off)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_erri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 7 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.R_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.CONDITION &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		condition := stack_items[4].Expr()
		pc := stack_items[6].Expr()

		s_rci_stmt := new(stmt.Stmt)
		s_rci_stmt.InitSRciStmt(op_code, suffix, dc, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRirciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 11 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.CONDITION &&
				stack_items[9].StackItemType() == TOKEN &&
				stack_items[9].Token().TokenType() == lexer.COMMA &&
				stack_items[10].StackItemType() == EXPR &&
				stack_items[10].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		imm := stack_items[4].Expr()
		ra := stack_items[6].Expr()
		condition := stack_items[8].Expr()
		pc := stack_items[10].Expr()

		s_rirci_stmt := new(stmt.Stmt)
		s_rirci_stmt.InitSRirciStmt(op_code, suffix, dc, imm, ra, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rirci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRircStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		imm := stack_items[4].Expr()
		ra := stack_items[6].Expr()
		condition := stack_items[8].Expr()

		s_rirc_stmt := new(stmt.Stmt)
		s_rirc_stmt.InitSRircStmt(op_code, suffix, dc, imm, ra, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rirc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.CONDITION &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		condition := stack_items[6].Expr()
		pc := stack_items[8].Expr()

		s_rrci_stmt := new(stmt.Stmt)
		s_rrci_stmt.InitSRrciStmt(op_code, suffix, dc, ra, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrcStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 7 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		condition := stack_items[6].Expr()

		s_rrc_stmt := new(stmt.Stmt)
		s_rrc_stmt.InitSRrcStmt(op_code, suffix, dc, ra, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRriciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 11 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.CONDITION &&
				stack_items[9].StackItemType() == TOKEN &&
				stack_items[9].Token().TokenType() == lexer.COMMA &&
				stack_items[10].StackItemType() == EXPR &&
				stack_items[10].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		imm := stack_items[6].Expr()
		condition := stack_items[8].Expr()
		pc := stack_items[10].Expr()

		s_rrici_stmt := new(stmt.Stmt)
		s_rrici_stmt.InitSRriciStmt(op_code, suffix, dc, ra, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRricStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		imm := stack_items[6].Expr()
		condition := stack_items[8].Expr()

		s_rric_stmt := new(stmt.Stmt)
		s_rric_stmt.InitSRricStmt(op_code, suffix, dc, ra, imm, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rric_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 7 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		imm := stack_items[6].Expr()

		s_rri_stmt := new(stmt.Stmt)
		s_rri_stmt.InitSRriStmt(op_code, suffix, dc, ra, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrrciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 11 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.CONDITION &&
				stack_items[9].StackItemType() == TOKEN &&
				stack_items[9].Token().TokenType() == lexer.COMMA &&
				stack_items[10].StackItemType() == EXPR &&
				stack_items[10].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		rb := stack_items[6].Expr()
		condition := stack_items[7].Expr()
		pc := stack_items[10].Expr()

		s_rrrci_stmt := new(stmt.Stmt)
		s_rrrci_stmt.InitSRrrciStmt(op_code, suffix, dc, ra, rb, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrrci_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrrcStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.CONDITION {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		rb := stack_items[6].Expr()
		condition := stack_items[7].Expr()

		s_rrrc_stmt := new(stmt.Stmt)
		s_rrrc_stmt.InitSRrrcStmt(op_code, suffix, dc, ra, rb, condition)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrrc_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrriciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 13 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[9].StackItemType() == TOKEN &&
				stack_items[9].Token().TokenType() == lexer.COMMA &&
				stack_items[10].StackItemType() == EXPR &&
				stack_items[10].Expr().ExprType() == expr.CONDITION &&
				stack_items[11].StackItemType() == TOKEN &&
				stack_items[11].Token().TokenType() == lexer.COMMA &&
				stack_items[12].StackItemType() == EXPR &&
				stack_items[12].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		rb := stack_items[6].Expr()
		imm := stack_items[8].Expr()
		condition := stack_items[10].Expr()
		pc := stack_items[12].Expr()

		s_rrrici_stmt := new(stmt.Stmt)
		s_rrrici_stmt.InitSRrriciStmt(op_code, suffix, dc, ra, rb, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrrici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		rb := stack_items[6].Expr()
		imm := stack_items[8].Expr()

		s_rrri_stmt := new(stmt.Stmt)
		s_rrri_stmt.InitSRrriStmt(op_code, suffix, dc, ra, rb, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrrStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 7 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		rb := stack_items[6].Expr()

		s_rrr_stmt := new(stmt.Stmt)
		s_rrr_stmt.InitSRrrStmt(op_code, suffix, dc, ra, rb)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rrr_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRrStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 5 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()

		s_rr_stmt := new(stmt.Stmt)
		s_rr_stmt.InitSRrStmt(op_code, suffix, dc, ra)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_rr_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSRStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.R_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()

		s_r_stmt := new(stmt.Stmt)
		s_r_stmt.InitSRStmt(op_code, suffix, dc)

		stack_item := new(StackItem)
		stack_item.InitStmt(s_r_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterNopStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.R_OP_CODE {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()

		if op_code.ROpCodeExpr().Token().TokenType() != lexer.NOP {
			err := errors.New("op code is not NOP")
			panic(err)
		}

		nop_stmt := new(stmt.Stmt)
		nop_stmt.InitNopStmt(op_code)

		stack_item := new(StackItem)
		stack_item.InitStmt(nop_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterBkpStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.I_OP_CODE {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()

		if op_code.IOpCodeExpr().Token().TokenType() != lexer.BKP {
			err := errors.New("op code is not BKP")
			panic(err)
		}

		bkp_stmt := new(stmt.Stmt)
		bkp_stmt.InitBkpStmt()

		stack_item := new(StackItem)
		stack_item.InitStmt(bkp_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterBootRiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RICI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		imm := stack_items[3].Expr()

		token_type := op_code.RiciOpCodeExpr().Token().TokenType()
		if token_type != lexer.BOOT && token_type != lexer.RESUME {
			err := errors.New("op code is not BOOT nor RESUME")
			panic(err)
		}

		boot_ri_stmt := new(stmt.Stmt)
		boot_ri_stmt.InitBootRiStmt(op_code, ra, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(boot_ri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCallRiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		imm := stack_items[3].Expr()

		token_type := op_code.RriOpCodeExpr().Token().TokenType()
		if token_type != lexer.CALL {
			err := errors.New("op code is not CALL")
			panic(err)
		}

		call_ri_stmt := new(stmt.Stmt)
		call_ri_stmt.InitCallRiStmt(rc, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(call_ri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterCallRrStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RRI_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()

		token_type := op_code.RriOpCodeExpr().Token().TokenType()
		if token_type != lexer.CALL {
			err := errors.New("op code is not CALL")
			panic(err)
		}

		call_rr_stmt := new(stmt.Stmt)
		call_rr_stmt.InitCallRrStmt(rc, ra)

		stack_item := new(StackItem)
		stack_item.InitStmt(call_rr_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterDivStepDrdiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.DRDICI_OP_CODE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		dc := stack_items[1].Token()
		ra := stack_items[3].Expr()
		db := stack_items[5].Token()
		pc := stack_items[7].Expr()

		div_step_drdi_stmt := new(stmt.Stmt)
		div_step_drdi_stmt.InitDivStepDrdiStmt(op_code, dc, ra, db, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(div_step_drdi_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterJeqRiiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.JUMP_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		imm := stack_items[3].Expr()
		pc := stack_items[5].Expr()

		jeq_rii_stmt := new(stmt.Stmt)
		jeq_rii_stmt.InitJeqRiiStmt(op_code, ra, imm, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(jeq_rii_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterJeqRriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.JUMP_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		rb := stack_items[3].Expr()
		pc := stack_items[5].Expr()

		jeq_rri_stmt := new(stmt.Stmt)
		jeq_rri_stmt.InitJeqRriStmt(op_code, ra, rb, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(jeq_rri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterJnzRiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.JUMP_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		pc := stack_items[3].Expr()

		jnz_ri_stmt := new(stmt.Stmt)
		jnz_ri_stmt.InitJnzRiStmt(op_code, ra, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(jnz_ri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterJumpIStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.JUMP_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		pc := stack_items[1].Expr()

		token_type := op_code.JumpOpCodeExpr().Token().TokenType()
		if token_type != lexer.JUMP {
			err := errors.New("op code is not JUMP")
			panic(err)
		}

		jump_i_stmt := new(stmt.Stmt)
		jump_i_stmt.InitJumpIStmt(pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(jump_i_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterJumpRStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.JUMP_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()

		token_type := op_code.JumpOpCodeExpr().Token().TokenType()
		if token_type != lexer.JUMP {
			err := errors.New("op code is not JUMP")
			panic(err)
		}

		jump_r_stmt := new(stmt.Stmt)
		jump_r_stmt.InitJumpRStmt(ra)

		stack_item := new(StackItem)
		stack_item.InitStmt(jump_r_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLbsRriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.LOAD_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		ra := stack_items[3].Expr()
		off := stack_items[5].Expr()

		lbs_rri_stmt := new(stmt.Stmt)
		lbs_rri_stmt.InitLbsRriStmt(op_code, rc, ra, off)

		stack_item := new(StackItem)
		stack_item.InitStmt(lbs_rri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLbsSRriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 7 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.LOAD_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.SRC_REG &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		ra := stack_items[4].Expr()
		off := stack_items[6].Expr()

		lbs_s_rri_stmt := new(stmt.Stmt)
		lbs_s_rri_stmt.InitLbsSRriStmt(op_code, suffix, dc, ra, off)

		stack_item := new(StackItem)
		stack_item.InitStmt(lbs_s_rri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLdDriStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.LOAD_OP_CODE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.SRC_REG &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		dc := stack_items[1].Token()
		ra := stack_items[3].Expr()
		off := stack_items[5].Expr()

		ld_dri_stmt := new(stmt.Stmt)
		ld_dri_stmt.InitLdDriStmt(op_code, dc, ra, off)

		stack_item := new(StackItem)
		stack_item.InitStmt(ld_dri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterMovdDdStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.DDCI_OP_CODE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.PAIR_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		dc := stack_items[1].Token()
		db := stack_items[3].Token()

		movd_dd_stmt := new(stmt.Stmt)
		movd_dd_stmt.InitMovdDdStmt(op_code, dc, db)

		stack_item := new(StackItem)
		stack_item.InitStmt(movd_dd_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterMoveRiciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 8 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.CONDITION &&
				stack_items[6].StackItemType() == TOKEN &&
				stack_items[6].Token().TokenType() == lexer.COMMA &&
				stack_items[7].StackItemType() == EXPR &&
				stack_items[7].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		imm := stack_items[3].Expr()
		condition := stack_items[5].Expr()
		pc := stack_items[7].Expr()

		token_type := op_code.RrOpCodeExpr().Token().TokenType()
		if token_type != lexer.MOVE {
			err := errors.New("op code is not MOVE")
			panic(err)
		}

		move_rici_stmt := new(stmt.Stmt)
		move_rici_stmt.InitMoveRiciStmt(rc, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(move_rici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterMoveRiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		rc := stack_items[1].Expr()
		imm := stack_items[3].Expr()

		token_type := op_code.RrOpCodeExpr().Token().TokenType()
		if token_type != lexer.MOVE {
			err := errors.New("op code is not MOVE")
			panic(err)
		}

		move_ri_stmt := new(stmt.Stmt)
		move_ri_stmt.InitMoveRiStmt(rc, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(move_ri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterMoveSRiciStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 9 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.COMMA &&
				stack_items[6].StackItemType() == EXPR &&
				stack_items[6].Expr().ExprType() == expr.CONDITION &&
				stack_items[7].StackItemType() == TOKEN &&
				stack_items[7].Token().TokenType() == lexer.COMMA &&
				stack_items[8].StackItemType() == EXPR &&
				stack_items[8].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		imm := stack_items[4].Expr()
		condition := stack_items[6].Expr()
		pc := stack_items[8].Expr()

		token_type := op_code.RrOpCodeExpr().Token().TokenType()
		if token_type != lexer.MOVE {
			err := errors.New("op code is not MOVE")
			panic(err)
		}

		move_s_rici_stmt := new(stmt.Stmt)
		move_s_rici_stmt.InitMoveSRiciStmt(suffix, dc, imm, condition, pc)

		stack_item := new(StackItem)
		stack_item.InitStmt(move_s_rici_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterMoveSRiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 5 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SUFFIX &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.PAIR_REG &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COMMA &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[4].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		suffix := stack_items[1].Expr()
		dc := stack_items[2].Token()
		imm := stack_items[4].Expr()

		token_type := op_code.RrOpCodeExpr().Token().TokenType()
		if token_type != lexer.MOVE {
			err := errors.New("op code is not MOVE")
			panic(err)
		}

		move_s_ri_stmt := new(stmt.Stmt)
		move_s_ri_stmt.InitMoveSRiStmt(suffix, dc, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(move_s_ri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSbIdRiiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.STORE_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		off := stack_items[3].Expr()
		imm := stack_items[5].Expr()

		sb_id_rii_stmt := new(stmt.Stmt)
		sb_id_rii_stmt.InitSbIdRiiStmt(op_code, ra, off, imm)

		stack_item := new(StackItem)
		stack_item.InitStmt(sb_id_rii_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSbIdRiStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.STORE_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		off := stack_items[3].Expr()

		sb_id_ri_stmt := new(stmt.Stmt)
		sb_id_ri_stmt.InitSbIdRiStmt(op_code, ra, off)

		stack_item := new(StackItem)
		stack_item.InitStmt(sb_id_ri_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSbRirStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.STORE_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == EXPR &&
				stack_items[5].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		off := stack_items[3].Expr()
		rb := stack_items[5].Expr()

		sb_rir_stmt := new(stmt.Stmt)
		sb_rir_stmt.InitSbRirStmt(op_code, ra, off, rb)

		stack_item := new(StackItem)
		stack_item.InitStmt(sb_rir_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterSdRidStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 6 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.STORE_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.COMMA &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.COMMA &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.PAIR_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()
		off := stack_items[3].Expr()
		db := stack_items[5].Token()

		sd_rid_stmt := new(stmt.Stmt)
		sd_rid_stmt.InitSdRidStmt(op_code, ra, off, db)

		stack_item := new(StackItem)
		stack_item.InitStmt(sd_rid_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterStopStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.CI_OP_CODE {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()

		token_type := op_code.CiOpCodeExpr().Token().TokenType()
		if token_type != lexer.STOP {
			err := errors.New("op code is not STOP")
			panic(err)
		}

		stop_stmt := new(stmt.Stmt)
		stop_stmt.InitStopStmt()

		stack_item := new(StackItem)
		stack_item.InitStmt(stop_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterTimeCfgRStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.RR_OP_CODE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.SRC_REG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		op_code := stack_items[0].Expr()
		ra := stack_items[1].Expr()

		token_type := op_code.RrOpCodeExpr().Token().TokenType()
		if token_type != lexer.TIME_CFG {
			err := errors.New("op code is not TIME_CFG")
			panic(err)
		}

		time_cfg_r_stmt := new(stmt.Stmt)
		time_cfg_r_stmt.InitTimeCfgRStmt(ra)

		stack_item := new(StackItem)
		stack_item.InitStmt(time_cfg_r_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}

func (this *Parser) RegisterLabelStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.PROGRAM_COUNTER &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.COLON {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[0].Expr()

		label_stmt := new(stmt.Stmt)
		label_stmt.InitLabelStmt(expr_)

		stack_item := new(StackItem)
		stack_item.InitStmt(label_stmt)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddStmtRule(rule)
}
