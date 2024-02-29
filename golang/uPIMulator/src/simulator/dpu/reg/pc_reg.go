package reg

import (
	"uPIMulator/src/abi/word"
	"uPIMulator/src/misc"
)

type PcReg struct {
	word *word.Word
}

func (this *PcReg) Init() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.word = new(word.Word)
	this.word.Init(config_loader.AddressWidth())
}

func (this *PcReg) Fini() {
}

func (this *PcReg) Read() int64 {
	return this.word.Value(word.UNSIGNED)
}

func (this *PcReg) Write(value int64) {
	this.word.SetValue(value)
}

func (this *PcReg) Increment() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	iram_data_size := int64(config_loader.IramDataWidth() / 8)

	this.Write(this.Read() + iram_data_size)
}
