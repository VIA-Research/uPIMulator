package expr

import (
	"errors"
)

type AdditiveExprType int

const (
	ADD AdditiveExprType = iota
	SUB
)

type AdditiveExpr struct {
	additive_expr_type AdditiveExprType
	loperand           *Expr
	roperand           *Expr
}

func (this *AdditiveExpr) Init(
	additive_expr_type AdditiveExprType,
	loperand *Expr,
	roperand *Expr,
) {
	if loperand.ExprType() != PRIMARY && loperand.ExprType() != POSTFIX &&
		loperand.ExprType() != UNARY &&
		loperand.ExprType() != MULTIPLICATIVE &&
		loperand.ExprType() != ADDITIVE {
		err := errors.New("loperand expr type is wrong")
		panic(err)
	} else if roperand.ExprType() != PRIMARY && roperand.ExprType() != POSTFIX &&
		roperand.ExprType() != UNARY &&
		roperand.ExprType() != MULTIPLICATIVE &&
		roperand.ExprType() != ADDITIVE {
		err := errors.New("roperand expr type is wrong")
		panic(err)
	}

	this.additive_expr_type = additive_expr_type
	this.loperand = loperand
	this.roperand = roperand
}

func (this *AdditiveExpr) AdditiveExprType() AdditiveExprType {
	return this.additive_expr_type
}

func (this *AdditiveExpr) Loperand() *Expr {
	return this.loperand
}

func (this *AdditiveExpr) Roperand() *Expr {
	return this.roperand
}
