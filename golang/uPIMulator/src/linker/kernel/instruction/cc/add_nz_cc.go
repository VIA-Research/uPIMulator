package cc

import (
	"errors"
)

type AddNzCc struct {
	condition Condition
}

func (this *AddNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		C:    true,
		NC:   true,
		Z:    true,
		NZ:   true,
		XZ:   true,
		XNZ:  true,
		OV:   true,
		NOV:  true,
		PL:   true,
		MI:   true,
		SZ:   true,
		SNZ:  true,
		SPL:  true,
		SMI:  true,
		NC5:  true,
		NC6:  true,
		NC7:  true,
		NC8:  true,
		NC9:  true,
		NC10: true,
		NC11: true,
		NC12: true,
		NC13: true,
		NC14: true,
		TRUE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *AddNzCc) Condition() Condition {
	return this.condition
}
