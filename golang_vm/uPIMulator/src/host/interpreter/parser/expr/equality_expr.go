package expr

import (
	"errors"
)

type EqualityExprType int

const (
	EQ EqualityExprType = iota
	NOT_EQ
)

type EqualityExpr struct {
	equality_expr_type EqualityExprType
	loperand           *Expr
	roperand           *Expr
}

func (this *EqualityExpr) Init(
	equality_expr_type EqualityExprType,
	loperand *Expr,
	roperand *Expr,
) {
	if loperand.ExprType() != PRIMARY && loperand.ExprType() != POSTFIX &&
		loperand.ExprType() != UNARY &&
		loperand.ExprType() != MULTIPLICATIVE &&
		loperand.ExprType() != ADDITIVE &&
		loperand.ExprType() != SHIFT &&
		loperand.ExprType() != RELATIONAL &&
		loperand.ExprType() != EQUALITY {
		err := errors.New("loperand expr type is wrong")
		panic(err)
	} else if roperand.ExprType() != PRIMARY && roperand.ExprType() != POSTFIX &&
		roperand.ExprType() != UNARY &&
		roperand.ExprType() != MULTIPLICATIVE &&
		roperand.ExprType() != ADDITIVE &&
		roperand.ExprType() != SHIFT &&
		roperand.ExprType() != RELATIONAL &&
		roperand.ExprType() != EQUALITY {
		err := errors.New("roperand expr type is wrong")
		panic(err)
	}

	this.equality_expr_type = equality_expr_type
	this.loperand = loperand
	this.roperand = roperand
}

func (this *EqualityExpr) EqualityExprType() EqualityExprType {
	return this.equality_expr_type
}

func (this *EqualityExpr) Loperand() *Expr {
	return this.loperand
}

func (this *EqualityExpr) Roperand() *Expr {
	return this.roperand
}
