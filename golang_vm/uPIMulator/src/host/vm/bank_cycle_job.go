package vm

import (
	"uPIMulator/src/host/vm/dram/bank"
)

type BankCycleJob struct {
	bank *bank.Bank
}

func (this *BankCycleJob) Init(bank_ *bank.Bank) {
	this.bank = bank_
}

func (this *BankCycleJob) Execute() {
	this.bank.Cycle()
}
