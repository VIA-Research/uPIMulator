package cc

import (
	"errors"
)

type CountNzCc struct {
	condition Condition
}

func (this *CountNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:    true,
		NZ:   true,
		XZ:   true,
		XNZ:  true,
		SZ:   true,
		SNZ:  true,
		SPL:  true,
		SMI:  true,
		MAX:  true,
		NMAX: true,
		TRUE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *CountNzCc) Condition() Condition {
	return this.condition
}
