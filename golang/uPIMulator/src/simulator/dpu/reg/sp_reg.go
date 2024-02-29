package reg

import (
	"errors"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/linker/kernel/instruction/reg_descriptor"
	"uPIMulator/src/misc"
)

type SpReg struct {
	zero *word.Word
	one  *word.Word
	lneg *word.Word
	mneg *word.Word
	id   *word.Word
	id2  *word.Word
	id4  *word.Word
	id8  *word.Word
}

func (this *SpReg) Init(thread_id int) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.zero = new(word.Word)
	this.zero.Init(config_loader.MramDataWidth())
	this.zero.SetValue(0)

	this.one = new(word.Word)
	this.one.Init(config_loader.MramDataWidth())
	this.one.SetValue(1)

	this.lneg = new(word.Word)
	this.lneg.Init(config_loader.MramDataWidth())
	this.lneg.SetValue(-1)

	this.mneg = new(word.Word)
	this.mneg.Init(config_loader.MramDataWidth())
	this.mneg.SetValue(int64(this.mneg.Width()) - 1)

	this.id = new(word.Word)
	this.id.Init(config_loader.MramDataWidth())
	this.id.SetValue(int64(thread_id))

	this.id2 = new(word.Word)
	this.id2.Init(config_loader.MramDataWidth())
	this.id2.SetValue(int64(2 * thread_id))

	this.id4 = new(word.Word)
	this.id4.Init(config_loader.MramDataWidth())
	this.id4.SetValue(int64(4 * thread_id))

	this.id8 = new(word.Word)
	this.id8.Init(config_loader.MramDataWidth())
	this.id8.SetValue(int64(8 * thread_id))
}

func (this *SpReg) Fini() {
}

func (this *SpReg) Read(
	sp_reg_descriptor *reg_descriptor.SpRegDescriptor,
	representation word.Representation,
) int64 {
	if *sp_reg_descriptor == reg_descriptor.ZERO {
		return this.zero.Value(representation)
	} else if *sp_reg_descriptor == reg_descriptor.ONE {
		return this.one.Value(representation)
	} else if *sp_reg_descriptor == reg_descriptor.LNEG {
		return this.lneg.Value(representation)
	} else if *sp_reg_descriptor == reg_descriptor.MNEG {
		return this.mneg.Value(representation)
	} else if *sp_reg_descriptor == reg_descriptor.ID {
		return this.id.Value(representation)
	} else if *sp_reg_descriptor == reg_descriptor.ID2 {
		return this.id2.Value(representation)
	} else if *sp_reg_descriptor == reg_descriptor.ID4 {
		return this.id4.Value(representation)
	} else if *sp_reg_descriptor == reg_descriptor.ID8 {
		return this.id8.Value(representation)
	} else {
		err := errors.New("sp reg descriptor is not valid")
		panic(err)
	}
}
