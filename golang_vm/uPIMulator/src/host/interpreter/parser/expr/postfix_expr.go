package expr

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
)

type PostfixExprType int

const (
	BRACKET PostfixExprType = iota
	CALL
	DOT
	ARROW
	POSTFIX_PLUS_PLUS
	POSTFIX_MINUS_MINUS
)

type PostfixExpr struct {
	postfix_expr_type PostfixExprType
	base              *Expr
	offset_expr       *Expr
	arg_list          *ArgList
	offset_token      *lexer.Token
}

func (this *PostfixExpr) InitBracket(base *Expr, offset_expr *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.postfix_expr_type = BRACKET
	this.base = base
	this.offset_expr = offset_expr
}

func (this *PostfixExpr) InitCall(base *Expr, arg_list *ArgList) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.postfix_expr_type = CALL
	this.base = base
	this.arg_list = arg_list
}

func (this *PostfixExpr) InitDot(base *Expr, offset_token *lexer.Token) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX {
		err := errors.New("base expr type is wrong")
		panic(err)
	} else if offset_token.TokenType() != lexer.IDENTIFIER {
		err := errors.New("offset token's token type is not identifier")
		panic(err)
	}

	this.postfix_expr_type = DOT
	this.base = base
	this.offset_token = offset_token
}

func (this *PostfixExpr) InitArrow(base *Expr, offset_token *lexer.Token) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX {
		err := errors.New("base expr type is wrong")
		panic(err)
	} else if offset_token.TokenType() != lexer.IDENTIFIER {
		err := errors.New("offset token's token type is not identifier")
		panic(err)
	}

	this.postfix_expr_type = ARROW
	this.base = base
	this.offset_token = offset_token
}

func (this *PostfixExpr) InitPlusPlus(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.postfix_expr_type = POSTFIX_PLUS_PLUS
	this.base = base
}

func (this *PostfixExpr) InitMinusMinus(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.postfix_expr_type = POSTFIX_MINUS_MINUS
	this.base = base
}

func (this *PostfixExpr) PostfixExprType() PostfixExprType {
	return this.postfix_expr_type
}

func (this *PostfixExpr) Base() *Expr {
	return this.base
}

func (this *PostfixExpr) OffsetExpr() *Expr {
	return this.offset_expr
}

func (this *PostfixExpr) ArgList() *ArgList {
	return this.arg_list
}

func (this *PostfixExpr) OffsetToken() *lexer.Token {
	return this.offset_token
}
