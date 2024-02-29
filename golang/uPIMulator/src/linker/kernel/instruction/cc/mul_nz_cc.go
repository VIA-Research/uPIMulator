package cc

import (
	"errors"
)

type MulNzCc struct {
	condition Condition
}

func (this *MulNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:     true,
		NZ:    true,
		XZ:    true,
		XNZ:   true,
		SZ:    true,
		SNZ:   true,
		SPL:   true,
		SMI:   true,
		LARGE: true,
		SMALL: true,
		TRUE:  true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *MulNzCc) Condition() Condition {
	return this.condition
}
