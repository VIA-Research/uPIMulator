package expr

import (
	"errors"
)

type AssignmentExprType int

const (
	ASSIGN AssignmentExprType = iota
	STAR_ASSIGN
	DIV_ASSIGN
	MOD_ASSIGN
	PLUS_ASSIGN
	MINUS_ASSIGN
	LSHIFT_ASSIGN
	RSHIFT_ASSIGN
	AND_ASSIGN
	CARET_ASSIGN
	OR_ASSIGN
)

type AssignmentExpr struct {
	assignment_expr_type AssignmentExprType
	lvalue               *Expr
	rvalue               *Expr
}

func (this *AssignmentExpr) Init(
	assignment_expr_type AssignmentExprType,
	lvalue *Expr,
	rvalue *Expr,
) {
	if lvalue.ExprType() != PRIMARY && lvalue.ExprType() != POSTFIX && lvalue.ExprType() != UNARY {
		err := errors.New("lvalue expr type is wrong")
		panic(err)
	} else if rvalue.ExprType() != PRIMARY && rvalue.ExprType() != POSTFIX &&
		rvalue.ExprType() != UNARY &&
		rvalue.ExprType() != MULTIPLICATIVE &&
		rvalue.ExprType() != ADDITIVE &&
		rvalue.ExprType() != SHIFT &&
		rvalue.ExprType() != RELATIONAL &&
		rvalue.ExprType() != EQUALITY &&
		rvalue.ExprType() != BITWISE_AND &&
		rvalue.ExprType() != BITWISE_XOR &&
		rvalue.ExprType() != BITWISE_OR &&
		rvalue.ExprType() != LOGICAL_AND &&
		rvalue.ExprType() != LOGICAL_OR &&
		rvalue.ExprType() != CONDITIONAL {
		err := errors.New("rvalue expr type is wrong")
		panic(err)
	}

	this.assignment_expr_type = assignment_expr_type
	this.lvalue = lvalue
	this.rvalue = rvalue
}

func (this *AssignmentExpr) AssignmentExprType() AssignmentExprType {
	return this.assignment_expr_type
}

func (this *AssignmentExpr) Lvalue() *Expr {
	return this.lvalue
}

func (this *AssignmentExpr) Rvalue() *Expr {
	return this.rvalue
}
