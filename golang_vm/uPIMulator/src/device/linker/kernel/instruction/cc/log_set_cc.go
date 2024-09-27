package cc

import (
	"errors"
)

type LogSetCc struct {
	condition Condition
}

func (this *LogSetCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		Z:   true,
		NZ:  true,
		XZ:  true,
		XNZ: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *LogSetCc) Condition() Condition {
	return this.condition
}
