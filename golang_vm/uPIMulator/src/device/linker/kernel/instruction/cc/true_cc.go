package cc

import (
	"errors"
)

type TrueCc struct {
	condition Condition
}

func (this *TrueCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		TRUE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *TrueCc) Condition() Condition {
	return this.condition
}
