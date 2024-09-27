package expr

import (
	"errors"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type UnaryExprType int

const (
	UNARY_PLUS_PLUS UnaryExprType = iota
	UNARY_MINUS_MINUS
	AND
	STAR
	PLUS
	MINUS
	TILDE
	NOT
	SIZEOF
)

type UnaryExpr struct {
	unary_expr_type UnaryExprType
	base            *Expr
	type_specifier  *type_specifier.TypeSpecifier
}

func (this *UnaryExpr) InitPlusPlus(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = UNARY_PLUS_PLUS
	this.base = base
}

func (this *UnaryExpr) InitMinusMinus(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = UNARY_MINUS_MINUS
	this.base = base
}

func (this *UnaryExpr) InitAnd(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = AND
	this.base = base
}

func (this *UnaryExpr) InitStar(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = STAR
	this.base = base
}

func (this *UnaryExpr) InitPlus(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = PLUS
	this.base = base
}

func (this *UnaryExpr) InitMinus(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = MINUS
	this.base = base
}

func (this *UnaryExpr) InitTilde(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = TILDE
	this.base = base
}

func (this *UnaryExpr) InitNot(base *Expr) {
	if base.ExprType() != PRIMARY && base.ExprType() != POSTFIX && base.ExprType() != UNARY {
		err := errors.New("base expr type is wrong")
		panic(err)
	}

	this.unary_expr_type = NOT
	this.base = base
}

func (this *UnaryExpr) InitSizeof(type_specifier *type_specifier.TypeSpecifier) {
	this.unary_expr_type = SIZEOF
	this.type_specifier = type_specifier
}

func (this *UnaryExpr) UnaryExprType() UnaryExprType {
	return this.unary_expr_type
}

func (this *UnaryExpr) Base() *Expr {
	return this.base
}

func (this *UnaryExpr) TypeSpecifier() *type_specifier.TypeSpecifier {
	return this.type_specifier
}
