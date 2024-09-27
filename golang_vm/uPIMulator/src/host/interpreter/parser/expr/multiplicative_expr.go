package expr

import (
	"errors"
)

type MultiplicativeExprType int

const (
	MUL MultiplicativeExprType = iota
	DIV
	MOD
)

type MultiplicativeExpr struct {
	multiplicative_expr_type MultiplicativeExprType
	loperand                 *Expr
	roperand                 *Expr
}

func (this *MultiplicativeExpr) Init(
	multiplicative_expr_type MultiplicativeExprType,
	loperand *Expr,
	roperand *Expr,
) {
	if loperand.ExprType() != PRIMARY && loperand.ExprType() != POSTFIX &&
		loperand.ExprType() != UNARY &&
		loperand.ExprType() != MULTIPLICATIVE {
		err := errors.New("loperand expr type is wrong")
		panic(err)
	} else if roperand.ExprType() != PRIMARY && roperand.ExprType() != POSTFIX &&
		roperand.ExprType() != UNARY &&
		roperand.ExprType() != MULTIPLICATIVE {
		err := errors.New("roperand expr type is wrong")
		panic(err)
	}

	this.multiplicative_expr_type = multiplicative_expr_type
	this.loperand = loperand
	this.roperand = roperand
}

func (this *MultiplicativeExpr) MultiplicativeExprType() MultiplicativeExprType {
	return this.multiplicative_expr_type
}

func (this *MultiplicativeExpr) Loperand() *Expr {
	return this.loperand
}

func (this *MultiplicativeExpr) Roperand() *Expr {
	return this.roperand
}
