package expr

import (
	"errors"
)

type RelationalExprType int

const (
	LESS RelationalExprType = iota
	LESS_EQ
	GREATER
	GREATER_EQ
)

type RelationalExpr struct {
	relational_expr_type RelationalExprType
	loperand             *Expr
	roperand             *Expr
}

func (this *RelationalExpr) Init(
	relational_expr_type RelationalExprType,
	loperand *Expr,
	roperand *Expr,
) {
	if loperand.ExprType() != PRIMARY && loperand.ExprType() != POSTFIX &&
		loperand.ExprType() != UNARY &&
		loperand.ExprType() != MULTIPLICATIVE &&
		loperand.ExprType() != ADDITIVE &&
		loperand.ExprType() != SHIFT &&
		loperand.ExprType() != RELATIONAL {
		err := errors.New("loperand expr type is wrong")
		panic(err)
	} else if roperand.ExprType() != PRIMARY && roperand.ExprType() != POSTFIX &&
		roperand.ExprType() != UNARY &&
		roperand.ExprType() != MULTIPLICATIVE &&
		roperand.ExprType() != ADDITIVE &&
		roperand.ExprType() != SHIFT &&
		roperand.ExprType() != RELATIONAL {
		err := errors.New("roperand expr type is wrong")
		panic(err)
	}

	this.relational_expr_type = relational_expr_type
	this.loperand = loperand
	this.roperand = roperand
}

func (this *RelationalExpr) RelationalExprType() RelationalExprType {
	return this.relational_expr_type
}

func (this *RelationalExpr) Loperand() *Expr {
	return this.loperand
}

func (this *RelationalExpr) Roperand() *Expr {
	return this.roperand
}
