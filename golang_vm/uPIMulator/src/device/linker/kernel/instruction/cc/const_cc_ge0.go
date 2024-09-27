package cc

import (
	"errors"
)

type ConstCcGe0 struct {
	condition Condition
}

func (this *ConstCcGe0) Init(condition Condition) {
	conditions := map[Condition]bool{
		PL: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *ConstCcGe0) Condition() Condition {
	return this.condition
}
