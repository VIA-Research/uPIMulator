package expr

import (
	"errors"
)

type BitwiseXorExpr struct {
	loperand *Expr
	roperand *Expr
}

func (this *BitwiseXorExpr) Init(loperand *Expr, roperand *Expr) {
	if loperand.ExprType() != PRIMARY && loperand.ExprType() != POSTFIX &&
		loperand.ExprType() != UNARY &&
		loperand.ExprType() != MULTIPLICATIVE &&
		loperand.ExprType() != ADDITIVE &&
		loperand.ExprType() != SHIFT &&
		loperand.ExprType() != RELATIONAL &&
		loperand.ExprType() != EQUALITY &&
		loperand.ExprType() != BITWISE_AND &&
		loperand.ExprType() != BITWISE_XOR {
		err := errors.New("loperand expr type is wrong")
		panic(err)
	} else if roperand.ExprType() != PRIMARY && roperand.ExprType() != POSTFIX &&
		roperand.ExprType() != UNARY &&
		roperand.ExprType() != MULTIPLICATIVE &&
		roperand.ExprType() != ADDITIVE &&
		roperand.ExprType() != SHIFT &&
		roperand.ExprType() != RELATIONAL &&
		roperand.ExprType() != EQUALITY &&
		roperand.ExprType() != BITWISE_AND &&
		roperand.ExprType() != BITWISE_XOR {
		err := errors.New("roperand expr type is wrong")
		panic(err)
	}

	this.loperand = loperand
	this.roperand = roperand
}

func (this *BitwiseXorExpr) Loperand() *Expr {
	return this.loperand
}

func (this *BitwiseXorExpr) Roperand() *Expr {
	return this.roperand
}
