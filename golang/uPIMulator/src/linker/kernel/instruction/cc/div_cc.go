package cc

import (
	"errors"
)

type DivCc struct {
	condition Condition
}

func (this *DivCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		SZ:    true,
		SNZ:   true,
		SPL:   true,
		SMI:   true,
		TRUE:  true,
		FALSE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *DivCc) Condition() Condition {
	return this.condition
}
