package cc

import (
	"errors"
)

type ConstCcGeu struct {
	condition Condition
}

func (this *ConstCcGeu) Init(condition Condition) {
	conditions := map[Condition]bool{
		GEU: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *ConstCcGeu) Condition() Condition {
	return this.condition
}
