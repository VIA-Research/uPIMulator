package cc

import (
	"errors"
)

type SubSetCc struct {
	condition Condition
}

func (this *SubSetCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:   true,
		NZ:  true,
		XZ:  true,
		XNZ: true,
		EQ:  true,
		NEQ: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *SubSetCc) Condition() Condition {
	return this.condition
}
