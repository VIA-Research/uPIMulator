package cc

import (
	"errors"
)

type NoCc struct {
	condition Condition
}

func (this *NoCc) Init(condition Condition) {
	conditions := map[Condition]bool{}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *NoCc) Condition() Condition {
	return this.condition
}
