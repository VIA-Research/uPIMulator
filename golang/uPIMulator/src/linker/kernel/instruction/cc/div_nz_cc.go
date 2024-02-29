package cc

import (
	"errors"
)

type DivNzCc struct {
	condition Condition
}

func (this *DivNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		SZ:   true,
		SNZ:  true,
		SPL:  true,
		SMI:  true,
		TRUE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *DivNzCc) Condition() Condition {
	return this.condition
}
