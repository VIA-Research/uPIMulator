package cc

import (
	"errors"
)

type ReleaseCc struct {
	condition Condition
}

func (this *ReleaseCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		NZ: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *ReleaseCc) Condition() Condition {
	return this.condition
}
