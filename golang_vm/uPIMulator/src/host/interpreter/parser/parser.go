package parser

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/decl"
	"uPIMulator/src/host/interpreter/parser/directive"
	"uPIMulator/src/host/interpreter/parser/expr"
	"uPIMulator/src/host/interpreter/parser/param_list"
	"uPIMulator/src/host/interpreter/parser/stmt"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
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

	this.RegisterTypeSpecifierVoid()
	this.RegisterTypeSpecifierChar()
	this.RegisterTypeSpecifierShort()
	this.RegisterTypeSpecifierInt()
	this.RegisterTypeSpecifierLong()
	this.RegisterTypeSpecifierStruct()
	this.RegisterTypeSpecifierStar()

	this.RegisterParamListBegin()
	this.RegisterParamListAppend()

	this.RegisterArgListBegin()
	this.RegisterArgListAppend()

	this.RegisterPrimaryExprIdentifier()
	this.RegisterPrimaryExprNumber()
	this.RegisterPrimaryExprString()
	this.RegisterPrimaryExprNullptr()
	this.RegisterPrimaryExprParen()

	this.RegisterPostfixExprBracket()
	this.RegisterPostfixExprCallEmpty()
	this.RegisterPostfixExprCallSingle()
	this.RegisterPostfixExprCallMultiple()
	this.RegisterPostfixExprDot()
	this.RegisterPostfixExprArrow()
	this.RegisterPostfixExprPlusPlus()
	this.RegisterPostfixExprMinusMinus()

	this.RegisterUnaryExprPlusPlus()
	this.RegisterUnaryExprMinusMinus()
	this.RegisterUnaryExprAnd()
	this.RegisterUnaryExprStar()
	this.RegisterUnaryExprPlus()
	this.RegisterUnaryExprMinus()
	this.RegisterUnaryExprTilde()
	this.RegisterUnaryExprNot()
	this.RegisterUnaryExprSizeof()

	this.RegisterMultiplicativeExprMul()
	this.RegisterMultiplicativeExprDiv()
	this.RegisterMultiplicativeExprMod()

	this.RegisterAdditiveExprAdd()
	this.RegisterAdditiveExprSub()

	this.RegisterShiftExprLshift()
	this.RegisterShiftExprRshift()

	this.RegisterRelationalExprLess()
	this.RegisterRelationalExprLessEq()
	this.RegisterRelationalExprGreater()
	this.RegisterRelationalExprGreaterEq()

	this.RegisterEqualityExprEq()
	this.RegisterEqualityExprNotEq()

	this.RegisterBitwiseAndExpr()
	this.RegisterBitwiseXorExpr()
	this.RegisterBitwiseOrExpr()

	this.RegisterLogicalAndExpr()
	this.RegisterLogicalOrExpr()

	this.RegisterConditionalExpr()

	this.RegisterAssignmentExprAssign()
	this.RegisterAssignmentExprStarAssign()
	this.RegisterAssignmentExprDivAssign()
	this.RegisterAssignmentExprModAssign()
	this.RegisterAssignmentExprPlusAssign()
	this.RegisterAssignmentExprMinusAssign()
	this.RegisterAssignmentExprLshiftAssign()
	this.RegisterAssignmentExprRshiftAssign()
	this.RegisterAssignmentExprAndAssign()
	this.RegisterAssignmentExprCaretAssign()
	this.RegisterAssignmentExprOrAssign()

	this.RegisterConcatExpr()

	this.RegisterEmptyStmt()
	this.RegisterVarDeclStmt()
	this.RegisterVarDeclInitStmt()
	this.RegisterForStmt()
	this.RegisterDpuForeachStmt()
	this.RegisterWhileStmt()
	this.RegisterContinueStmt()
	this.RegisterBreakStmt()

	this.RegisterIfStmt()
	this.RegisterElseIfStmt()
	this.RegisterElseStmt()

	this.RegisterReturnStmtWithoutValue()
	this.RegisterReturnStmtWithValue()

	this.RegisterExprStmt()

	this.RegisterBlockStmt()

	this.RegisterStructDef()

	this.RegisterFuncDeclEmpty()
	this.RegisterFuncDeclNonEmpty()

	this.RegisterFuncDefEmpty()
	this.RegisterFuncDefNonEmpty()

	this.RegisterIncludeDirective()
	this.RegisterDefineDirective()
}

func (this *Parser) Parse(token_stream *lexer.TokenStream) *Ast {
	for i := 0; i < token_stream.Length(); i++ {
		token := token_stream.Get(i)

		for this.table.IsReducible(token) {
			this.table.Reduce(token)
		}

		stack_item := new(StackItem)
		stack_item.InitToken(token)
		this.stack.Push(stack_item)
	}

	return this.stack.Accept()
}

func (this *Parser) RegisterTypeSpecifierVoid() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.VOID {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := new(type_specifier.TypeSpecifier)
		type_specifier_.InitPrimitive(type_specifier.VOID)

		stack_item := new(StackItem)
		stack_item.InitTypeSpecifier(type_specifier_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterTypeSpecifierChar() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.CHAR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := new(type_specifier.TypeSpecifier)
		type_specifier_.InitPrimitive(type_specifier.CHAR)

		stack_item := new(StackItem)
		stack_item.InitTypeSpecifier(type_specifier_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterTypeSpecifierShort() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SHORT {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := new(type_specifier.TypeSpecifier)
		type_specifier_.InitPrimitive(type_specifier.SHORT)

		stack_item := new(StackItem)
		stack_item.InitTypeSpecifier(type_specifier_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterTypeSpecifierInt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.INT {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := new(type_specifier.TypeSpecifier)
		type_specifier_.InitPrimitive(type_specifier.INT)

		stack_item := new(StackItem)
		stack_item.InitTypeSpecifier(type_specifier_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterTypeSpecifierLong() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.LONG {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := new(type_specifier.TypeSpecifier)
		type_specifier_.InitPrimitive(type_specifier.LONG)

		stack_item := new(StackItem)
		stack_item.InitTypeSpecifier(type_specifier_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterTypeSpecifierStruct() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.STRUCT &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].Expr().PrimaryExpr().PrimaryExprType() == expr.IDENTIFIER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		struct_identifier := stack_items[1].Expr().PrimaryExpr().Token()

		type_specifier_ := new(type_specifier.TypeSpecifier)
		type_specifier_.InitStruct(type_specifier.STRUCT, struct_identifier)

		stack_item := new(StackItem)
		stack_item.InitTypeSpecifier(type_specifier_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterTypeSpecifierStar() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TYPE_SPECIFIER &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.STAR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := stack_items[0].TypeSpecifier()
		type_specifier_.AddStar()

		stack_item := new(StackItem)
		stack_item.InitTypeSpecifier(type_specifier_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterParamListBegin() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TYPE_SPECIFIER &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].Expr().PrimaryExpr().PrimaryExprType() == expr.IDENTIFIER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := stack_items[0].TypeSpecifier()
		identifier := stack_items[1].Expr().PrimaryExpr().Token()

		param := new(param_list.Param)
		param.Init(type_specifier_, identifier)

		param_list_ := new(param_list.ParamList)
		param_list_.Init()
		param_list_.Append(param)

		stack_item := new(StackItem)
		stack_item.InitParamList(param_list_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterParamListAppend() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == PARAM_LIST &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.COMMA &&
				stack_items[2].StackItemType() == TYPE_SPECIFIER &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PRIMARY &&
				stack_items[3].Expr().PrimaryExpr().PrimaryExprType() == expr.IDENTIFIER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		param_list_ := stack_items[0].ParamList()
		type_specifier_ := stack_items[2].TypeSpecifier()
		identifier := stack_items[3].Expr().PrimaryExpr().Token()

		param := new(param_list.Param)
		param.Init(type_specifier_, identifier)

		param_list_.Append(param)

		stack_item := new(StackItem)
		stack_item.InitParamList(param_list_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterArgListBegin() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:      true,
		lexer.LPAREN:        true,
		lexer.DOT:           true,
		lexer.ARROW:         true,
		lexer.STAR:          true,
		lexer.DIV:           true,
		lexer.MOD:           true,
		lexer.PLUS:          true,
		lexer.MINUS:         true,
		lexer.LSHIFT:        true,
		lexer.RSHIFT:        true,
		lexer.LESS:          true,
		lexer.LESS_EQ:       true,
		lexer.GREATER:       true,
		lexer.GREATER_EQ:    true,
		lexer.EQ:            true,
		lexer.NOT_EQ:        true,
		lexer.AND:           true,
		lexer.CARET:         true,
		lexer.OR:            true,
		lexer.AND_AND:       true,
		lexer.OR_OR:         true,
		lexer.QUESTION:      true,
		lexer.ASSIGN:        true,
		lexer.STAR_ASSIGN:   true,
		lexer.DIV_ASSIGN:    true,
		lexer.MOD_ASSIGN:    true,
		lexer.PLUS_ASSIGN:   true,
		lexer.MINUS_ASSIGN:  true,
		lexer.LSHIFT_ASSIGN: true,
		lexer.RSHIFT_ASSIGN: true,
		lexer.AND_ASSIGN:    true,
		lexer.CARET_ASSIGN:  true,
		lexer.OR_ASSIGN:     true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.COMMA &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		arg1 := stack_items[0].Expr()
		arg2 := stack_items[2].Expr()

		arg_list := new(expr.ArgList)
		arg_list.Init()
		arg_list.Append(arg1)
		arg_list.Append(arg2)

		stack_item := new(StackItem)
		stack_item.InitArgList(arg_list)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterArgListAppend() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:      true,
		lexer.LPAREN:        true,
		lexer.DOT:           true,
		lexer.ARROW:         true,
		lexer.STAR:          true,
		lexer.DIV:           true,
		lexer.MOD:           true,
		lexer.PLUS:          true,
		lexer.MINUS:         true,
		lexer.LSHIFT:        true,
		lexer.RSHIFT:        true,
		lexer.LESS:          true,
		lexer.LESS_EQ:       true,
		lexer.GREATER:       true,
		lexer.GREATER_EQ:    true,
		lexer.EQ:            true,
		lexer.NOT_EQ:        true,
		lexer.AND:           true,
		lexer.CARET:         true,
		lexer.OR:            true,
		lexer.AND_AND:       true,
		lexer.OR_OR:         true,
		lexer.QUESTION:      true,
		lexer.ASSIGN:        true,
		lexer.STAR_ASSIGN:   true,
		lexer.DIV_ASSIGN:    true,
		lexer.MOD_ASSIGN:    true,
		lexer.PLUS_ASSIGN:   true,
		lexer.MINUS_ASSIGN:  true,
		lexer.LSHIFT_ASSIGN: true,
		lexer.RSHIFT_ASSIGN: true,
		lexer.AND_ASSIGN:    true,
		lexer.CARET_ASSIGN:  true,
		lexer.OR_ASSIGN:     true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == ARG_LIST &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.COMMA &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		arg_list := stack_items[0].ArgList()
		arg := stack_items[2].Expr()

		arg_list.Append(arg)

		stack_item := new(StackItem)
		stack_item.InitArgList(arg_list)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPrimaryExprIdentifier() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.IDENTIFIER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		primary_expr := new(expr.PrimaryExpr)
		primary_expr.InitIdentifier(token)

		expr_ := new(expr.Expr)
		expr_.InitPrimaryExpr(primary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPrimaryExprNumber() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.NUMBER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		primary_expr := new(expr.PrimaryExpr)
		primary_expr.InitNumber(token)

		expr_ := new(expr.Expr)
		expr_.InitPrimaryExpr(primary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPrimaryExprString() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.STRING {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		primary_expr := new(expr.PrimaryExpr)
		primary_expr.InitString(token)

		expr_ := new(expr.Expr)
		expr_.InitPrimaryExpr(primary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPrimaryExprNullptr() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.NULL {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		token := stack_items[0].Token()

		primary_expr := new(expr.PrimaryExpr)
		primary_expr.InitNullptr(token)

		expr_ := new(expr.Expr)
		expr_.InitPrimaryExpr(primary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPrimaryExprParen() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.LPAREN &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.RPAREN {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		e := stack_items[1].Expr()

		primary_expr := new(expr.PrimaryExpr)
		primary_expr.InitParen(e)

		expr_ := new(expr.Expr)
		expr_.InitPrimaryExpr(primary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprBracket() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LBRACKET &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.RBRACKET {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()
		offset_expr := stack_items[2].Expr()

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitBracket(base, offset_expr)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprCallEmpty() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.RPAREN {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()

		arg_list := new(expr.ArgList)
		arg_list.Init()

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitCall(base, arg_list)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprCallSingle() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].Expr().PrimaryExpr().PrimaryExprType() == expr.PAREN {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()

		arg_list := new(expr.ArgList)
		arg_list.Init()
		arg_list.Append(stack_items[1].Expr().PrimaryExpr().Expr())

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitCall(base, arg_list)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprCallMultiple() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == ARG_LIST &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.RPAREN {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()
		arg_list := stack_items[2].ArgList()

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitCall(base, arg_list)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprDot() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.DOT &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.IDENTIFIER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()
		offset_token := stack_items[2].Token()

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitDot(base, offset_token)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprArrow() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.ARROW &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.IDENTIFIER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()
		offset_token := stack_items[2].Token()

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitArrow(base, offset_token)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprPlusPlus() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PLUS_PLUS {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitPlusPlus(base)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterPostfixExprMinusMinus() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.MINUS_MINUS {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[0].Expr()

		postfix_expr := new(expr.PostfixExpr)
		postfix_expr.InitMinusMinus(base)

		expr_ := new(expr.Expr)
		expr_.InitPostfixExpr(postfix_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprPlusPlus() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.PLUS_PLUS &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitPlusPlus(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprMinusMinus() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.MINUS_MINUS &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitMinusMinus(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprAnd() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.AND &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitAnd(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprStar() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.STAR &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitStar(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprPlus() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.PLUS &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitPlus(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprMinus() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.MINUS &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitMinus(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprTilde() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.TILDE &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitTilde(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprNot() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.NOT &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		base := stack_items[1].Expr()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitNot(base)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterUnaryExprSizeof() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SIZEOF &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == TYPE_SPECIFIER &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.RPAREN {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := stack_items[2].TypeSpecifier()

		unary_expr := new(expr.UnaryExpr)
		unary_expr.InitSizeof(type_specifier_)

		expr_ := new(expr.Expr)
		expr_.InitUnaryExpr(unary_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterMultiplicativeExprMul() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.STAR &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		multiplicative_expr := new(expr.MultiplicativeExpr)
		multiplicative_expr.Init(expr.MUL, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitMultiplicativeExpr(multiplicative_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterMultiplicativeExprDiv() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.DIV &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		multiplicative_expr := new(expr.MultiplicativeExpr)
		multiplicative_expr.Init(expr.DIV, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitMultiplicativeExpr(multiplicative_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterMultiplicativeExprMod() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.MOD &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		multiplicative_expr := new(expr.MultiplicativeExpr)
		multiplicative_expr.Init(expr.MOD, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitMultiplicativeExpr(multiplicative_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAdditiveExprAdd() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PLUS &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		additive_expr := new(expr.AdditiveExpr)
		additive_expr.Init(expr.ADD, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitAdditiveExpr(additive_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAdditiveExprSub() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.MINUS &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		additive_expr := new(expr.AdditiveExpr)
		additive_expr.Init(expr.SUB, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitAdditiveExpr(additive_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterShiftExprLshift() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LSHIFT &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		shift_expr := new(expr.ShiftExpr)
		shift_expr.Init(expr.LSHIFT, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitShiftExpr(shift_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterShiftExprRshift() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.RSHIFT &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		shift_expr := new(expr.ShiftExpr)
		shift_expr.Init(expr.RSHIFT, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitShiftExpr(shift_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterRelationalExprLess() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LESS &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		relational_expr := new(expr.RelationalExpr)
		relational_expr.Init(expr.LESS, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitRelationalExpr(relational_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterRelationalExprLessEq() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LESS_EQ &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		relational_expr := new(expr.RelationalExpr)
		relational_expr.Init(expr.LESS_EQ, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitRelationalExpr(relational_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterRelationalExprGreater() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.GREATER &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		relational_expr := new(expr.RelationalExpr)
		relational_expr.Init(expr.GREATER, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitRelationalExpr(relational_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterRelationalExprGreaterEq() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.GREATER_EQ &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		relational_expr := new(expr.RelationalExpr)
		relational_expr.Init(expr.GREATER_EQ, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitRelationalExpr(relational_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterEqualityExprEq() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.EQ &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		equality_expr := new(expr.EqualityExpr)
		equality_expr.Init(expr.EQ, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitEqualityExpr(equality_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterEqualityExprNotEq() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.NOT_EQ &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		equality_expr := new(expr.EqualityExpr)
		equality_expr.Init(expr.NOT_EQ, loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitEqualityExpr(equality_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterBitwiseAndExpr() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.AND &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		bitwise_and_expr := new(expr.BitwiseAndExpr)
		bitwise_and_expr.Init(loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitBitwiseAndExpr(bitwise_and_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterBitwiseXorExpr() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.CARET &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		bitwise_xor_expr := new(expr.BitwiseXorExpr)
		bitwise_xor_expr.Init(loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitBitwiseXorExpr(bitwise_xor_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterBitwiseOrExpr() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.OR &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		bitwise_or_expr := new(expr.BitwiseOrExpr)
		bitwise_or_expr.Init(loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitBitwiseOrExpr(bitwise_or_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterLogicalAndExpr() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.AND_AND &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		logical_and_expr := new(expr.LogicalAndExpr)
		logical_and_expr.Init(loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitLogicalAndExpr(logical_and_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterLogicalOrExpr() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.OR_OR &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		loperand := stack_items[0].Expr()
		roperand := stack_items[2].Expr()

		logical_or_expr := new(expr.LogicalOrExpr)
		logical_or_expr.Init(loperand, roperand)

		expr_ := new(expr.Expr)
		expr_.InitLogicalOrExpr(logical_or_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterConditionalExpr() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 5 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.QUESTION &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.COLON &&
				stack_items[4].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		condition_expr := stack_items[0].Expr()
		true_expr := stack_items[2].Expr()
		false_expr := stack_items[4].Expr()

		conditional_expr := new(expr.ConditionalExpr)
		conditional_expr.Init(condition_expr, true_expr, false_expr)

		expr_ := new(expr.Expr)
		expr_.InitConditionalExpr(conditional_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprStarAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.STAR_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.STAR_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprDivAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.DIV_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.DIV_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprModAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.MOD_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.MOD_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprPlusAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.PLUS_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.PLUS_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprMinusAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.MINUS_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.MINUS_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprLshiftAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LSHIFT_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.LSHIFT_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprRshiftAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.RSHIFT_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.RSHIFT_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprAndAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.AND_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.AND_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprCaretAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.CARET_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.CARET_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterAssignmentExprOrAssign() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
		lexer.PLUS:        true,
		lexer.MINUS:       true,
		lexer.LSHIFT:      true,
		lexer.RSHIFT:      true,
		lexer.LESS:        true,
		lexer.LESS_EQ:     true,
		lexer.GREATER:     true,
		lexer.GREATER_EQ:  true,
		lexer.EQ:          true,
		lexer.NOT_EQ:      true,
		lexer.AND:         true,
		lexer.CARET:       true,
		lexer.OR:          true,
		lexer.AND_AND:     true,
		lexer.OR_OR:       true,
		lexer.QUESTION:    true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.OR_ASSIGN &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[2].Expr()

		assignment_expr := new(expr.AssignmentExpr)
		assignment_expr.Init(expr.OR_ASSIGN, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAssignmentExpr(assignment_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterConcatExpr() {
	precedence := map[lexer.TokenType]bool{
		lexer.LBRACKET:    true,
		lexer.LPAREN:      true,
		lexer.DOT:         true,
		lexer.ARROW:       true,
		lexer.PLUS_PLUS:   true,
		lexer.MINUS_MINUS: true,
		lexer.STAR:        true,
		lexer.DIV:         true,
		lexer.MOD:         true,
	}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[0].Expr()
		rvalue := stack_items[1].Expr()

		additive_expr := new(expr.AdditiveExpr)
		additive_expr.Init(expr.ADD, lvalue, rvalue)

		expr_ := new(expr.Expr)
		expr_.InitAdditiveExpr(additive_expr)

		stack_item := new(StackItem)
		stack_item.InitExpr(expr_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterEmptyStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 1 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		empty_stmt := new(stmt.EmptyStmt)
		empty_stmt.Init()

		stmt_ := new(stmt.Stmt)
		stmt_.InitEmptyStmt(empty_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterVarDeclStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == PARAM_LIST &&
				stack_items[0].ParamList().Length() == 1 &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		param := stack_items[0].ParamList().Get(0)

		type_specifier_ := param.TypeSpecifier()
		identifier := param.Identifier()

		var_decl_stmt := new(stmt.VarDeclStmt)
		var_decl_stmt.Init(type_specifier_, identifier)

		stmt_ := new(stmt.Stmt)
		stmt_.InitVarDeclStmt(var_decl_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterVarDeclInitStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == PARAM_LIST &&
				stack_items[0].ParamList().Length() == 1 &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.ASSIGN &&
				stack_items[2].StackItemType() == EXPR &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		param := stack_items[0].ParamList().Get(0)
		rvalue := stack_items[2].Expr()

		type_specifier_ := param.TypeSpecifier()
		identifier := param.Identifier()

		var_decl_init_stmt := new(stmt.VarDeclInitStmt)
		var_decl_init_stmt.Init(type_specifier_, identifier, rvalue)

		stmt_ := new(stmt.Stmt)
		stmt_.InitVarDeclInitStmt(var_decl_init_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterForStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 7 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.FOR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == STMT &&
				stack_items[3].StackItemType() == STMT &&
				stack_items[4].StackItemType() == EXPR &&
				stack_items[5].StackItemType() == TOKEN &&
				stack_items[5].Token().TokenType() == lexer.RPAREN &&
				stack_items[6].StackItemType() == STMT &&
				stack_items[6].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		initialization := stack_items[2].Stmt()
		condition_stmt := stack_items[3].Stmt()
		update_expr := stack_items[4].Expr()
		body := stack_items[6].Stmt()

		if condition_stmt.StmtType() != stmt.EXPR {
			err := errors.New("condition stmt type is not expr")
			panic(err)
		}

		condition_expr := condition_stmt.ExprStmt().Expr()

		u := new(stmt.ExprStmt)
		u.Init(update_expr)

		update_stmt := new(stmt.Stmt)
		update_stmt.InitExprStmt(u)

		for_stmt := new(stmt.ForStmt)
		for_stmt.Init(initialization, condition_expr, update_stmt, body)

		stmt_ := new(stmt.Stmt)
		stmt_.InitForStmt(for_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterDpuForeachStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[0].Expr().ExprType() == expr.POSTFIX &&
				stack_items[0].Expr().PostfixExpr().PostfixExprType() == expr.CALL &&
				stack_items[0].Expr().PostfixExpr().Base().ExprType() == expr.PRIMARY &&
				stack_items[0].Expr().PostfixExpr().Base().PrimaryExpr().PrimaryExprType() == expr.IDENTIFIER &&
				stack_items[0].Expr().PostfixExpr().Base().PrimaryExpr().Token().Attribute() == "DPU_FOREACH" &&
				stack_items[0].Expr().PostfixExpr().ArgList().Length() == 3 &&
				stack_items[1].StackItemType() == STMT &&
				stack_items[1].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		foreach := stack_items[0].Expr().PostfixExpr().ArgList()
		body := stack_items[1].Stmt()

		dpu_foreach_stmt := new(stmt.DpuForeachStmt)
		dpu_foreach_stmt.Init(foreach, body)

		stmt_ := new(stmt.Stmt)
		stmt_.InitDpuForeachStmt(dpu_foreach_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterWhileStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.WHILE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].Expr().PrimaryExpr().PrimaryExprType() == expr.PAREN &&
				stack_items[2].StackItemType() == STMT &&
				stack_items[2].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		condition := stack_items[1].Expr().PrimaryExpr().Expr()
		body := stack_items[2].Stmt()

		while_stmt := new(stmt.WhileStmt)
		while_stmt.Init(condition, body)

		stmt_ := new(stmt.Stmt)
		stmt_.InitWhileStmt(while_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterContinueStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.CONTINUE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		continue_stmt := new(stmt.ContinueStmt)
		continue_stmt.Init()

		stmt_ := new(stmt.Stmt)
		stmt_.InitContinueStmt(continue_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterBreakStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.BREAK &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		break_stmt := new(stmt.BreakStmt)
		break_stmt.Init()

		stmt_ := new(stmt.Stmt)
		stmt_.InitBreakStmt(break_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterIfStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.IF &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].Expr().PrimaryExpr().PrimaryExprType() == expr.PAREN &&
				stack_items[2].StackItemType() == STMT &&
				stack_items[2].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		condition := stack_items[1].Expr().PrimaryExpr().Expr()
		body := stack_items[2].Stmt()

		if_stmt := new(stmt.IfStmt)
		if_stmt.Init(condition, body)

		stmt_ := new(stmt.Stmt)
		stmt_.InitIfStmt(if_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterElseIfStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 5 {
			return false
		} else {
			if stack_items[0].StackItemType() == STMT &&
				stack_items[0].Stmt().StmtType() == stmt.IF &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.ELSE &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.IF &&
				stack_items[3].StackItemType() == EXPR &&
				stack_items[3].Expr().ExprType() == expr.PRIMARY &&
				stack_items[3].Expr().PrimaryExpr().PrimaryExprType() == expr.PAREN &&
				stack_items[4].StackItemType() == STMT &&
				stack_items[4].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		if_stmt := stack_items[0].Stmt().IfStmt()
		condition := stack_items[3].Expr().PrimaryExpr().Expr()
		body := stack_items[4].Stmt()

		if_stmt.AppendElseIf(condition, body)

		stmt_ := new(stmt.Stmt)
		stmt_.InitIfStmt(if_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterElseStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == STMT &&
				stack_items[0].Stmt().StmtType() == stmt.IF &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.ELSE &&
				stack_items[2].StackItemType() == STMT &&
				stack_items[2].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		if_stmt := stack_items[0].Stmt().IfStmt()
		body := stack_items[2].Stmt()

		if_stmt.SetElseBody(body)

		stmt_ := new(stmt.Stmt)
		stmt_.InitIfStmt(if_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterReturnStmtWithoutValue() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.RETURN &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		return_stmt := new(stmt.ReturnStmt)
		return_stmt.InitWithoutValue()

		stmt_ := new(stmt.Stmt)
		stmt_.InitReturnStmt(return_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterReturnStmtWithValue() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.RETURN &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		value := stack_items[1].Expr()

		return_stmt := new(stmt.ReturnStmt)
		return_stmt.InitWithValue(value)

		stmt_ := new(stmt.Stmt)
		stmt_.InitReturnStmt(return_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterExprStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == EXPR &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		expr_ := stack_items[0].Expr()

		expr_stmt := new(stmt.ExprStmt)
		expr_stmt.Init(expr_)

		stmt_ := new(stmt.Stmt)
		stmt_.InitExprStmt(expr_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterBlockStmt() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if stack_items[0].StackItemType() == TOKEN &&
			stack_items[0].Token().TokenType() == lexer.LBRACE &&
			stack_items[len(stack_items)-1].StackItemType() == TOKEN &&
			stack_items[len(stack_items)-1].Token().TokenType() == lexer.RBRACE {
			for i := 1; i < len(stack_items)-1; i++ {
				if stack_items[i].StackItemType() != STMT {
					return false
				}
			}
			return true
		} else {
			return false
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		block_stmt := new(stmt.BlockStmt)
		block_stmt.Init()

		for i := 1; i < len(stack_items)-1; i++ {
			block_stmt.Append(stack_items[i].Stmt())
		}

		stmt_ := new(stmt.Stmt)
		stmt_.InitBlockStmt(block_stmt)

		stack_item := new(StackItem)
		stack_item.InitStmt(stmt_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterStructDef() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == TYPE_SPECIFIER &&
				stack_items[0].TypeSpecifier().TypeSpecifierType() == type_specifier.STRUCT &&
				stack_items[0].TypeSpecifier().NumStars() == 0 &&
				stack_items[1].StackItemType() == STMT &&
				stack_items[1].Stmt().StmtType() == stmt.BLOCK &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.SEMI {
				body := stack_items[1].Stmt().BlockStmt()
				for i := 0; i < body.Length(); i++ {
					stmt_ := body.Get(i)

					if stmt_.StmtType() != stmt.VAR_DECL {
						return false
					}
				}
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		type_specifier_ := stack_items[0].TypeSpecifier()
		body := stack_items[1].Stmt()

		identifier := type_specifier_.StructIdentifier()

		struct_def := new(decl.StructDef)
		struct_def.Init(identifier, body)

		decl_ := new(decl.Decl)
		decl_.InitStructDef(struct_def)

		stack_item := new(StackItem)
		stack_item.InitDecl(decl_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterFuncDeclEmpty() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == PARAM_LIST &&
				stack_items[0].ParamList().Length() == 1 &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.RPAREN &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		param := stack_items[0].ParamList().Get(0)

		type_specifier_ := param.TypeSpecifier()
		identifier := param.Identifier()

		param_list_ := new(param_list.ParamList)
		param_list_.Init()

		func_decl := new(decl.FuncDecl)
		func_decl.Init(type_specifier_, identifier, param_list_)

		decl_ := new(decl.Decl)
		decl_.InitFuncDecl(func_decl)

		stack_item := new(StackItem)
		stack_item.InitDecl(decl_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterFuncDeclNonEmpty() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 5 {
			return false
		} else {
			if stack_items[0].StackItemType() == PARAM_LIST &&
				stack_items[0].ParamList().Length() == 1 &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == PARAM_LIST &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.RPAREN &&
				stack_items[4].StackItemType() == TOKEN &&
				stack_items[4].Token().TokenType() == lexer.SEMI {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		param := stack_items[0].ParamList().Get(0)
		param_list_ := stack_items[2].ParamList()

		type_specifier_ := param.TypeSpecifier()
		identifier := param.Identifier()

		func_decl := new(decl.FuncDecl)
		func_decl.Init(type_specifier_, identifier, param_list_)

		decl_ := new(decl.Decl)
		decl_.InitFuncDecl(func_decl)

		stack_item := new(StackItem)
		stack_item.InitDecl(decl_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterFuncDefEmpty() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 4 {
			return false
		} else {
			if stack_items[0].StackItemType() == PARAM_LIST &&
				stack_items[0].ParamList().Length() == 1 &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == TOKEN &&
				stack_items[2].Token().TokenType() == lexer.RPAREN &&
				stack_items[3].StackItemType() == STMT &&
				stack_items[3].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		param := stack_items[0].ParamList().Get(0)
		body := stack_items[3].Stmt()

		type_specifier_ := param.TypeSpecifier()
		identifier := param.Identifier()

		param_list_ := new(param_list.ParamList)
		param_list_.Init()

		func_def := new(decl.FuncDef)
		func_def.Init(type_specifier_, identifier, param_list_, body)

		decl_ := new(decl.Decl)
		decl_.InitFuncDef(func_def)

		stack_item := new(StackItem)
		stack_item.InitDecl(decl_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterFuncDefNonEmpty() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 5 {
			return false
		} else {
			if stack_items[0].StackItemType() == PARAM_LIST &&
				stack_items[0].ParamList().Length() == 1 &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.LPAREN &&
				stack_items[2].StackItemType() == PARAM_LIST &&
				stack_items[3].StackItemType() == TOKEN &&
				stack_items[3].Token().TokenType() == lexer.RPAREN &&
				stack_items[4].StackItemType() == STMT &&
				stack_items[4].Stmt().StmtType() == stmt.BLOCK {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		param := stack_items[0].ParamList().Get(0)
		param_list_ := stack_items[2].ParamList()
		body := stack_items[4].Stmt()

		type_specifier_ := param.TypeSpecifier()
		identifier := param.Identifier()

		func_def := new(decl.FuncDef)
		func_def.Init(type_specifier_, identifier, param_list_, body)

		decl_ := new(decl.Decl)
		decl_.InitFuncDef(func_def)

		stack_item := new(StackItem)
		stack_item.InitDecl(decl_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterIncludeDirective() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 2 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.INCLUDE &&
				stack_items[1].StackItemType() == TOKEN &&
				stack_items[1].Token().TokenType() == lexer.HEADER {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		header := stack_items[1].Token()

		include_directive := new(directive.IncludeDirective)
		include_directive.Init(header)

		directive_ := new(directive.Directive)
		directive_.InitIncludeDirective(include_directive)

		stack_item := new(StackItem)
		stack_item.InitDirective(directive_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}

func (this *Parser) RegisterDefineDirective() {
	precedence := map[lexer.TokenType]bool{}

	reducible := func(stack_items []*StackItem) bool {
		if len(stack_items) != 3 {
			return false
		} else {
			if stack_items[0].StackItemType() == TOKEN &&
				stack_items[0].Token().TokenType() == lexer.DEFINE &&
				stack_items[1].StackItemType() == EXPR &&
				stack_items[1].Expr().ExprType() == expr.PRIMARY &&
				stack_items[1].Expr().PrimaryExpr().PrimaryExprType() == expr.IDENTIFIER &&
				stack_items[2].StackItemType() == EXPR {
				return true
			} else {
				return false
			}
		}
	}

	reduce := func(stack_items []*StackItem) *StackItem {
		lvalue := stack_items[1].Expr().PrimaryExpr().Token()
		rvalue := stack_items[2].Expr()

		define_directive := new(directive.DefineDirective)
		define_directive.Init(lvalue, rvalue)

		directive_ := new(directive.Directive)
		directive_.InitDefineDirective(define_directive)

		stack_item := new(StackItem)
		stack_item.InitDirective(directive_)

		return stack_item
	}

	rule := new(Rule)
	rule.Init(precedence, reducible, reduce)

	this.table.AddRule(rule)
}
