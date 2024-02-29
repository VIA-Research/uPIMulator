package cc

import (
	"errors"
)

type TrueFalseCc struct {
	condition Condition
}

func (this *TrueFalseCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		TRUE:  true,
		FALSE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *TrueFalseCc) Condition() Condition {
	return this.condition
}
