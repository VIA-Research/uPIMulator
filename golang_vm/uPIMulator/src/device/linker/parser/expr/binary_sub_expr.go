package expr

import (
	"errors"
)

type BinarySubExpr struct {
	operand1 *Expr
	operand2 *Expr
}

func (this *BinarySubExpr) Init(operand1 *Expr, operand2 *Expr) {
	if operand1.ExprType() != PRIMARY {
		err := errors.New("operand1 is not a primary expr")
		panic(err)
	}

	if operand2.ExprType() != PRIMARY {
		err := errors.New("operand2 is not a primary expr")
		panic(err)
	}

	this.operand1 = operand1
	this.operand2 = operand2
}

func (this *BinarySubExpr) Operand1() *Expr {
	return this.operand1
}

func (this *BinarySubExpr) Operand2() *Expr {
	return this.operand2
}
