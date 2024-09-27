package cc

import (
	"errors"
)

type FalseCc struct {
	condition Condition
}

func (this *FalseCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		FALSE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *FalseCc) Condition() Condition {
	return this.condition
}
