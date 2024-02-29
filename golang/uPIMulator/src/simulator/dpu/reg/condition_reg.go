package reg

import (
	"errors"
	"uPIMulator/src/linker/kernel/instruction/cc"
)

type ConditionReg struct {
	conditions map[cc.Condition]bool
}

func (this *ConditionReg) Init() {
	this.conditions = make(map[cc.Condition]bool, 0)

	this.ClearConditions()
}

func (this *ConditionReg) Fini() {
}

func (this *ConditionReg) Condition(condition cc.Condition) bool {
	if condition == cc.TRUE {
		return true
	} else if condition == cc.FALSE {
		return false
	} else {
		return this.conditions[condition]
	}
}

func (this *ConditionReg) SetCondition(condition cc.Condition) {
	if condition == cc.TRUE || condition == cc.FALSE {
		err := errors.New("condition is true or false")
		panic(err)
	}

	this.conditions[condition] = true
}

func (this *ConditionReg) ClearCondition(condition cc.Condition) {
	if condition == cc.TRUE || condition == cc.FALSE {
		err := errors.New("condition is true or false")
		panic(err)
	}

	this.conditions[condition] = false
}

func (this *ConditionReg) ClearConditions() {
	for i := 0; i <= int(cc.LARGE); i++ {
		condition := cc.Condition(i)

		if condition == cc.TRUE || condition == cc.FALSE {
			continue
		} else {
			this.conditions[condition] = false
		}
	}
}
