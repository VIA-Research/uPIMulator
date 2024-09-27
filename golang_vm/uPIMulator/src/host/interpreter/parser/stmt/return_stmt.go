package stmt

import (
	"errors"
	"uPIMulator/src/host/interpreter/parser/expr"
)

type ReturnStmt struct {
	value *expr.Expr
}

func (this *ReturnStmt) InitWithoutValue() {
	this.value = nil
}

func (this *ReturnStmt) InitWithValue(value *expr.Expr) {
	this.value = value
}

func (this *ReturnStmt) HasValue() bool {
	return this.value != nil
}

func (this *ReturnStmt) Value() *expr.Expr {
	if !this.HasValue() {
		err := errors.New("value does not exist")
		panic(err)
	}

	return this.value
}
