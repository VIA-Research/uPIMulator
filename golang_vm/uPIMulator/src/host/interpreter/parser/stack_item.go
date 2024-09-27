package parser

import (
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser/decl"
	"uPIMulator/src/host/interpreter/parser/directive"
	"uPIMulator/src/host/interpreter/parser/expr"
	"uPIMulator/src/host/interpreter/parser/param_list"
	"uPIMulator/src/host/interpreter/parser/stmt"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type StackItemType int

const (
	TOKEN StackItemType = iota
	TYPE_SPECIFIER
	PARAM_LIST
	ARG_LIST
	EXPR
	STMT
	DECL
	DIRECTIVE
)

type StackItem struct {
	stack_item_type StackItemType

	token          *lexer.Token
	type_specifier *type_specifier.TypeSpecifier
	param_list     *param_list.ParamList
	arg_list       *expr.ArgList
	expr           *expr.Expr
	stmt           *stmt.Stmt
	decl           *decl.Decl
	directive      *directive.Directive
}

func (this *StackItem) InitToken(token *lexer.Token) {
	this.stack_item_type = TOKEN

	this.token = token
}

func (this *StackItem) InitTypeSpecifier(type_specifier *type_specifier.TypeSpecifier) {
	this.stack_item_type = TYPE_SPECIFIER

	this.type_specifier = type_specifier
}

func (this *StackItem) InitParamList(param_list *param_list.ParamList) {
	this.stack_item_type = PARAM_LIST

	this.param_list = param_list
}

func (this *StackItem) InitArgList(arg_list *expr.ArgList) {
	this.stack_item_type = ARG_LIST

	this.arg_list = arg_list
}

func (this *StackItem) InitExpr(expr *expr.Expr) {
	this.stack_item_type = EXPR

	this.expr = expr
}

func (this *StackItem) InitStmt(stmt *stmt.Stmt) {
	this.stack_item_type = STMT

	this.stmt = stmt
}

func (this *StackItem) InitDecl(decl *decl.Decl) {
	this.stack_item_type = DECL

	this.decl = decl
}

func (this *StackItem) InitDirective(directive *directive.Directive) {
	this.stack_item_type = DIRECTIVE

	this.directive = directive
}

func (this *StackItem) StackItemType() StackItemType {
	return this.stack_item_type
}

func (this *StackItem) Token() *lexer.Token {
	return this.token
}

func (this *StackItem) TypeSpecifier() *type_specifier.TypeSpecifier {
	return this.type_specifier
}

func (this *StackItem) ParamList() *param_list.ParamList {
	return this.param_list
}

func (this *StackItem) ArgList() *expr.ArgList {
	return this.arg_list
}

func (this *StackItem) Expr() *expr.Expr {
	return this.expr
}

func (this *StackItem) Stmt() *stmt.Stmt {
	return this.stmt
}

func (this *StackItem) Decl() *decl.Decl {
	return this.decl
}

func (this *StackItem) Directive() *directive.Directive {
	return this.directive
}
