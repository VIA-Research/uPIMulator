package expr

import (
	"errors"
)

type ShiftExprType int

const (
	LSHIFT ShiftExprType = iota
	RSHIFT
)

type ShiftExpr struct {
	shift_expr_type ShiftExprType
	loperand        *Expr
	roperand        *Expr
}

func (this *ShiftExpr) Init(
	shift_expr_type ShiftExprType,
	loperand *Expr,
	roperand *Expr,
) {
	if loperand.ExprType() != PRIMARY && loperand.ExprType() != POSTFIX &&
		loperand.ExprType() != UNARY &&
		loperand.ExprType() != MULTIPLICATIVE &&
		loperand.ExprType() != ADDITIVE &&
		loperand.ExprType() != SHIFT {
		err := errors.New("loperand expr type is wrong")
		panic(err)
	} else if roperand.ExprType() != PRIMARY && roperand.ExprType() != POSTFIX &&
		roperand.ExprType() != UNARY &&
		roperand.ExprType() != MULTIPLICATIVE &&
		roperand.ExprType() != ADDITIVE &&
		roperand.ExprType() != SHIFT {
		err := errors.New("roperand expr type is wrong")
		panic(err)
	}

	this.shift_expr_type = shift_expr_type
	this.loperand = loperand
	this.roperand = roperand
}

func (this *ShiftExpr) ShiftExprType() ShiftExprType {
	return this.shift_expr_type
}

func (this *ShiftExpr) Loperand() *Expr {
	return this.loperand
}

func (this *ShiftExpr) Roperand() *Expr {
	return this.roperand
}
