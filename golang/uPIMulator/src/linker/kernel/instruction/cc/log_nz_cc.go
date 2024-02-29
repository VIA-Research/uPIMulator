package cc

import (
	"errors"
)

type LogNzCc struct {
	condition Condition
}

func (this *LogNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:    true,
		NZ:   true,
		XZ:   true,
		XNZ:  true,
		PL:   true,
		MI:   true,
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

func (this *LogNzCc) Condition() Condition {
	return this.condition
}
