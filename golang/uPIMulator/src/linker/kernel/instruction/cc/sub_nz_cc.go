package cc

import (
	"errors"
)

type SubNzCc struct {
	condition Condition
}

func (this *SubNzCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		C:    true,
		NC:   true,
		Z:    true,
		NZ:   true,
		XZ:   true,
		XNZ:  true,
		OV:   true,
		NOV:  true,
		MI:   true,
		PL:   true,
		EQ:   true,
		NEQ:  true,
		SPL:  true,
		SMI:  true,
		GES:  true,
		GEU:  true,
		GTS:  true,
		GTU:  true,
		LES:  true,
		LEU:  true,
		LTS:  true,
		LTU:  true,
		XGTS: true,
		XGTU: true,
		XLES: true,
		XLEU: true,
		TRUE: true,
	}

	if _, found := conditions[condition]; !found {
		err := errors.New("condition is not allowed")
		panic(err)
	}

	this.condition = condition
}

func (this *SubNzCc) Condition() Condition {
	return this.condition
}
