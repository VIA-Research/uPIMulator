package cc

import (
	"errors"
)

type AcquireCc struct {
	condition Condition
}

func (this *AcquireCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:    true,
		NZ:   true,
		TRUE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *AcquireCc) Condition() Condition {
	return this.condition
}
