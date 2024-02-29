package cc

import (
	"errors"
)

type ShiftNzCc struct {
	condition Condition
}

func (this *ShiftNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:     true,
		NZ:    true,
		XZ:    true,
		XNZ:   true,
		E:     true,
		O:     true,
		PL:    true,
		MI:    true,
		SZ:    true,
		SNZ:   true,
		SE:    true,
		SO:    true,
		SPL:   true,
		SMI:   true,
		SH32:  true,
		NSH32: true,
		TRUE:  true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *ShiftNzCc) Condition() Condition {
	return this.condition
}
