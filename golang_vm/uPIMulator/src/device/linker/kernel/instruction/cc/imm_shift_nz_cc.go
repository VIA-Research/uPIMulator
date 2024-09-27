package cc

import (
	"errors"
)

type ImmShiftNzCc struct {
	condition Condition
}

func (this *ImmShiftNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:    true,
		NZ:   true,
		XZ:   true,
		XNZ:  true,
		E:    true,
		O:    true,
		PL:   true,
		MI:   true,
		SZ:   true,
		SNZ:  true,
		SPL:  true,
		SMI:  true,
		SE:   true,
		SO:   true,
		TRUE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *ImmShiftNzCc) Condition() Condition {
	return this.condition
}
