package cc

import (
	"errors"
)

type BootCc struct {
	condition Condition
}

func (this *BootCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:     true,
		NZ:    true,
		XZ:    true,
		XNZ:   true,
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

func (this *BootCc) Condition() Condition {
	return this.condition
}
