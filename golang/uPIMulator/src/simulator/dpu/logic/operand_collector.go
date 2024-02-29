package logic

import (
	"errors"
	"uPIMulator/src/abi/encoding"
	"uPIMulator/src/abi/word"
	"uPIMulator/src/simulator/dpu/sram"
)

type OperandCollector struct {
	wram *sram.Wram
}

func (this *OperandCollector) Init() {
	this.wram = nil
}

func (this *OperandCollector) Fini() {
}

func (this *OperandCollector) ConnectWram(wram *sram.Wram) {
	if this.wram != nil {
		err := errors.New("wram is already set")
		panic(err)
	}

	this.wram = wram
}

func (this *OperandCollector) Lbs(address int64) int64 {
	byte_stream := this.wram.Read(address, 1)

	value := int64(byte_stream.Get(0))

	word_ := new(word.Word)
	word_.Init(8)
	word_.SetValue(value)

	return word_.Value(word.SIGNED)
}

func (this *OperandCollector) Lbu(address int64) int64 {
	byte_stream := this.wram.Read(address, 1)

	value := int64(byte_stream.Get(0))

	word_ := new(word.Word)
	word_.Init(8)
	word_.SetValue(value)

	return word_.Value(word.UNSIGNED)
}

func (this *OperandCollector) Lhs(address int64) int64 {
	word_ := new(word.Word)
	word_.Init(16)
	word_.SetBitSlice(0, 8, this.Lbs(address))
	word_.SetBitSlice(8, 16, this.Lbs(address+1))
	return word_.Value(word.SIGNED)
}

func (this *OperandCollector) Lhu(address int64) int64 {
	word_ := new(word.Word)
	word_.Init(16)
	word_.SetBitSlice(0, 8, this.Lbu(address))
	word_.SetBitSlice(8, 16, this.Lbu(address+1))
	return word_.Value(word.UNSIGNED)
}

func (this *OperandCollector) Lw(address int64) int64 {
	word_ := new(word.Word)
	word_.Init(32)
	word_.SetBitSlice(0, 8, this.Lbu(address))
	word_.SetBitSlice(8, 16, this.Lbu(address+1))
	word_.SetBitSlice(16, 24, this.Lbu(address+2))
	word_.SetBitSlice(24, 32, this.Lbu(address+3))
	return word_.Value(word.UNSIGNED)
}

func (this *OperandCollector) Ld(address int64) (int64, int64) {
	return this.Lw(address + 4), this.Lw(address)
}

func (this *OperandCollector) Sb(address int64, value int64) {
	word_ := new(word.Word)
	word_.Init(8)
	word_.SetValue(value)

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()
	byte_stream.Append(uint8(word_.Value(word.UNSIGNED)))

	this.wram.Write(address, 1, byte_stream)
}

func (this *OperandCollector) Sh(address int64, value int64) {
	word_ := new(word.Word)
	word_.Init(16)
	word_.SetValue(value)

	this.Sb(address, word_.BitSlice(word.UNSIGNED, 0, 8))
	this.Sb(address+1, word_.BitSlice(word.UNSIGNED, 8, 16))
}

func (this *OperandCollector) Sw(address int64, value int64) {
	word_ := new(word.Word)
	word_.Init(32)
	word_.SetValue(value)

	this.Sb(address, word_.BitSlice(word.UNSIGNED, 0, 8))
	this.Sb(address+1, word_.BitSlice(word.UNSIGNED, 8, 16))
	this.Sb(address+2, word_.BitSlice(word.UNSIGNED, 16, 24))
	this.Sb(address+3, word_.BitSlice(word.UNSIGNED, 24, 32))
}

func (this *OperandCollector) Sd(address int64, even int64, odd int64) {
	this.Sw(address+4, even)
	this.Sw(address, odd)
}
