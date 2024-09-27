package cc

import (
	"errors"
)

type ConstCcZero struct {
	condition Condition
}

func (this *ConstCcZero) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *ConstCcZero) Condition() Condition {
	return this.condition
}
