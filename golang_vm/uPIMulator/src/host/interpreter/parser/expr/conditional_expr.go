package expr

import (
	"errors"
)

type ConditionalExpr struct {
	conditional_expr *Expr
	true_expr        *Expr
	false_expr       *Expr
}

func (this *ConditionalExpr) Init(conditional_expr *Expr, true_expr *Expr, false_expr *Expr) {
	if conditional_expr.ExprType() != PRIMARY && conditional_expr.ExprType() != POSTFIX &&
		conditional_expr.ExprType() != UNARY &&
		conditional_expr.ExprType() != MULTIPLICATIVE &&
		conditional_expr.ExprType() != ADDITIVE &&
		conditional_expr.ExprType() != SHIFT &&
		conditional_expr.ExprType() != RELATIONAL &&
		conditional_expr.ExprType() != EQUALITY &&
		conditional_expr.ExprType() != BITWISE_AND &&
		conditional_expr.ExprType() != BITWISE_XOR &&
		conditional_expr.ExprType() != BITWISE_OR &&
		conditional_expr.ExprType() != LOGICAL_AND &&
		conditional_expr.ExprType() != LOGICAL_OR {
		err := errors.New("conditional_expr expr type is wrong")
		panic(err)
	} else if true_expr.ExprType() != PRIMARY && true_expr.ExprType() != POSTFIX &&
		true_expr.ExprType() != UNARY &&
		true_expr.ExprType() != MULTIPLICATIVE &&
		true_expr.ExprType() != ADDITIVE &&
		true_expr.ExprType() != SHIFT &&
		true_expr.ExprType() != RELATIONAL &&
		true_expr.ExprType() != EQUALITY &&
		true_expr.ExprType() != BITWISE_AND &&
		true_expr.ExprType() != BITWISE_XOR &&
		true_expr.ExprType() != BITWISE_OR &&
		true_expr.ExprType() != LOGICAL_AND &&
		true_expr.ExprType() != LOGICAL_OR &&
		true_expr.ExprType() != CONDITIONAL {
		err := errors.New("true_expr expr type is wrong")
		panic(err)
	} else if false_expr.ExprType() != PRIMARY && false_expr.ExprType() != POSTFIX &&
		false_expr.ExprType() != UNARY &&
		false_expr.ExprType() != MULTIPLICATIVE &&
		false_expr.ExprType() != ADDITIVE &&
		false_expr.ExprType() != SHIFT &&
		false_expr.ExprType() != RELATIONAL &&
		false_expr.ExprType() != EQUALITY &&
		false_expr.ExprType() != BITWISE_AND &&
		false_expr.ExprType() != BITWISE_XOR &&
		false_expr.ExprType() != BITWISE_OR &&
		false_expr.ExprType() != LOGICAL_AND &&
		false_expr.ExprType() != LOGICAL_OR &&
		false_expr.ExprType() != CONDITIONAL {
		err := errors.New("true_expr expr type is wrong")
		panic(err)
	}

	this.conditional_expr = conditional_expr
	this.true_expr = true_expr
	this.false_expr = false_expr
}

func (this *ConditionalExpr) ConditionalExpr() *Expr {
	return this.conditional_expr
}

func (this *ConditionalExpr) TrueExpr() *Expr {
	return this.true_expr
}

func (this *ConditionalExpr) FalseExpr() *Expr {
	return this.false_expr
}
