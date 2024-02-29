package cc

import (
	"errors"
)

type ExtSubSetCc struct {
	condition Condition
}

func (this *ExtSubSetCc) Init(condition Condition) {
	conditions := map[Condition]bool{
		C:    true,
		NC:   true,
		Z:    true,
		NZ:   true,
		XZ:   true,
		XNZ:  true,
		OV:   true,
		NOV:  true,
		EQ:   true,
		NEQ:  true,
		PL:   true,
		MI:   true,
		SZ:   true,
		SNZ:  true,
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

func (this *ExtSubSetCc) Condition() Condition {
	return this.condition
}
