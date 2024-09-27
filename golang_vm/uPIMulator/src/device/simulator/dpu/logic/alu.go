package logic

import (
	"errors"
	"uPIMulator/src/device/abi"
	"uPIMulator/src/misc"
)

type Alu struct {
}

func (this *Alu) Init() {
}

func (this *Alu) Fini() {
}

func (this *Alu) AtomicAddressHash(operand1 int64, operand2 int64) int64 {
	if operand1+operand2 >= 256 {
		err := errors.New("operand1 + operand2 >= 256")
		panic(err)
	}

	result, _, _ := this.Add(operand1, operand2)
	return result
}

func (this *Alu) Add(operand1 int64, operand2 int64) (int64, bool, bool) {
	return this.Addc(operand1, operand2, false)
}

func (this *Alu) Addc(operand1 int64, operand2 int64, carry_flag bool) (int64, bool, bool) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	var result int64
	if carry_flag {
		result = word1.Value(abi.UNSIGNED) + word2.Value(abi.UNSIGNED) + 1
	} else {
		result = word1.Value(abi.UNSIGNED) + word2.Value(abi.UNSIGNED)
	}

	max_unsigned_value := this.Pow2(mram_data_width) - 1

	var carry bool
	if result > max_unsigned_value {
		result %= this.Pow2(mram_data_width)
		carry = true
	} else {
		carry = false
	}

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	result_word.SetValue(result)

	var overflow bool
	if word1.SignBit() && word2.SignBit() && !result_word.SignBit() {
		overflow = true
	} else if !word1.SignBit() && !word2.SignBit() && result_word.SignBit() {
		overflow = true
	} else {
		overflow = false
	}

	return result_word.Value(abi.SIGNED), carry, overflow
}

func (this *Alu) Sub(operand1 int64, operand2 int64) (int64, bool, bool) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	var result int64
	var carry bool
	if word1.Value(abi.UNSIGNED) >= word2.Value(abi.UNSIGNED) {
		result = word1.Value(abi.UNSIGNED) - word2.Value(abi.UNSIGNED)
		carry = false
	} else {
		result = this.Pow2(mram_data_width) + word1.Value(abi.UNSIGNED) - word2.Value(abi.UNSIGNED)
		carry = true
	}

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	result_word.SetValue(result)

	var overflow bool
	if word1.SignBit() && !word2.SignBit() && result_word.SignBit() {
		overflow = true
	} else if !word1.SignBit() && word2.SignBit() && !result_word.SignBit() {
		overflow = true
	} else {
		overflow = false
	}

	return result, carry, overflow
}

func (this *Alu) Subc(operand1 int64, operand2 int64, carry_flag bool) (int64, bool, bool) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	var result int64
	var carry bool
	if carry_flag {
		if word1.Value(abi.UNSIGNED)+1 >= word2.Value(abi.UNSIGNED) {
			result = word1.Value(abi.UNSIGNED) - word2.Value(abi.UNSIGNED) - 1
			carry = false
		} else {
			result = this.Pow2(mram_data_width) + word1.Value(abi.UNSIGNED) - word2.Value(abi.UNSIGNED) - 1
			carry = true
		}
	} else {
		if word1.Value(abi.UNSIGNED) >= word2.Value(abi.UNSIGNED) {
			result = word1.Value(abi.UNSIGNED) - word2.Value(abi.UNSIGNED)
			carry = false
		} else {
			result = this.Pow2(mram_data_width) + word1.Value(abi.UNSIGNED) - word2.Value(abi.UNSIGNED)
			carry = true
		}
	}

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	result_word.SetValue(result)

	var overflow bool
	if word1.SignBit() && !word2.SignBit() && result_word.SignBit() {
		overflow = true
	} else if !word1.SignBit() && word2.SignBit() && !result_word.SignBit() {
		overflow = true
	} else {
		overflow = false
	}

	return result, carry, overflow
}

func (this *Alu) And(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if word1.Bit(i) && word2.Bit(i) {
			result_word.SetBit(i)
		} else {
			result_word.ClearBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Nand(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if word1.Bit(i) && word2.Bit(i) {
			result_word.ClearBit(i)
		} else {
			result_word.SetBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Andn(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if !word1.Bit(i) && word2.Bit(i) {
			result_word.SetBit(i)
		} else {
			result_word.ClearBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Or(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if word1.Bit(i) || word2.Bit(i) {
			result_word.SetBit(i)
		} else {
			result_word.ClearBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Nor(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if word1.Bit(i) || word2.Bit(i) {
			result_word.ClearBit(i)
		} else {
			result_word.SetBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Orn(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if !word1.Bit(i) || word2.Bit(i) {
			result_word.SetBit(i)
		} else {
			result_word.ClearBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Xor(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if !word1.Bit(i) && word2.Bit(i) {
			result_word.SetBit(i)
		} else if word1.Bit(i) && !word2.Bit(i) {
			result_word.SetBit(i)
		} else {
			result_word.ClearBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Nxor(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if !word1.Bit(i) && word2.Bit(i) {
			result_word.ClearBit(i)
		} else if word1.Bit(i) && word2.Bit(i) {
			result_word.ClearBit(i)
		} else {
			result_word.SetBit(i)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Asr(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)
	msb := word_.SignBit()

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if i+int(shift_value) >= mram_data_width {
			if msb {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		} else {
			if word_.Bit(i + int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Lsl(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if i < int(shift_value) {
			result_word.ClearBit(i)
		} else {
			if word_.Bit(i - int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) LslAdd(operand1 int64, operand2 int64, shift int64) (int64, bool, bool) {
	return this.Add(operand1, this.Lsl(operand2, shift))
}

func (this *Alu) LslSub(operand1 int64, operand2 int64, shift int64) (int64, bool, bool) {
	return this.Sub(operand1, this.Lsl(operand2, shift))
}

func (this *Alu) Lsl1(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if i < int(shift_value) {
			result_word.SetBit(i)
		} else {
			if word_.Bit(i - int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Lsl1x(operand int64, shift int64) int64 {
	err := errors.New("lsl1x is not yet implemented")
	panic(err)
}

func (this *Alu) Lslx(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	if shift_value == 0 {
		return 0
	} else {
		return this.Lsr(operand, int64(mram_data_width)-shift_value)
	}
}

func (this *Alu) Lsr(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if i+int(shift_value) >= mram_data_width {
			result_word.ClearBit(i)
		} else {
			if word_.Bit(i + int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) LsrAdd(operand1 int64, operand2 int64, shift int64) (int64, bool, bool) {
	return this.Add(operand1, this.Lsr(operand2, shift))
}

func (this *Alu) Lsr1(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if i+int(shift_value) >= mram_data_width {
			result_word.SetBit(i)
		} else {
			if word_.Bit(i + int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Lsr1x(operand int64, shift int64) int64 {
	err := errors.New("lsr1x is not yet implemented")
	panic(err)
}

func (this *Alu) Lsrx(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	if shift_value == 0 {
		return 0
	} else {
		return this.Lsl(operand, int64(mram_data_width)-shift_value)
	}
}

func (this *Alu) Rol(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if i < int(shift_value) {
			if word_.Bit(i + mram_data_width - int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
			result_word.SetBit(i)
		} else {
			if word_.Bit(i - int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) RolAdd(operand1 int64, operand2 int64, shift int64) (int64, bool, bool) {
	return this.Add(operand1, this.Rol(operand2, shift))
}

func (this *Alu) Ror(operand int64, shift int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	shift_word := new(abi.Word)
	shift_word.Init(mram_data_width)
	shift_word.SetValue(shift)
	shift_value := shift_word.BitSlice(abi.UNSIGNED, 0, 5)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < mram_data_width; i++ {
		if i+int(shift_value) >= mram_data_width {
			if word_.Bit((i + int(shift_value)) % mram_data_width) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
			result_word.SetBit(i)
		} else {
			if word_.Bit(i + int(shift_value)) {
				result_word.SetBit(i)
			} else {
				result_word.ClearBit(i)
			}
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Cao(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	ones := int64(0)
	for i := 0; i < mram_data_width; i++ {
		if word_.Bit(i) {
			ones++
		}
	}
	return ones
}

func (this *Alu) Clo(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	leading_ones := int64(0)
	for i := 0; i < mram_data_width; i++ {
		if word_.Bit(mram_data_width - 1 - i) {
			leading_ones++
		} else {
			break
		}
	}
	return leading_ones
}

func (this *Alu) Cls(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)
	msb := word_.SignBit()

	leading_sign_bits := int64(0)
	for i := 0; i < mram_data_width; i++ {
		if word_.Bit(mram_data_width-1-i) == msb {
			leading_sign_bits++
		} else {
			break
		}
	}
	return leading_sign_bits
}

func (this *Alu) Clz(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	leading_zeros := int64(0)
	for i := 0; i < mram_data_width; i++ {
		if !word_.Bit(mram_data_width - 1 - i) {
			leading_zeros++
		} else {
			break
		}
	}
	return leading_zeros
}

func (this *Alu) Cmpb4(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)
	for i := 0; i < 4; i++ {
		begin := 8 * i
		end := 8 * (i + 1)

		byte1 := word1.BitSlice(abi.UNSIGNED, begin, end)
		byte2 := word2.BitSlice(abi.UNSIGNED, begin, end)

		if byte1 == byte2 {
			result_word.SetBitSlice(begin, end, 1)
		} else {
			result_word.SetBitSlice(begin, end, 0)
		}
	}

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Extsb(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	return word_.BitSlice(abi.SIGNED, 0, 8)
}

func (this *Alu) Extsh(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	return word_.BitSlice(abi.SIGNED, 0, 16)
}

func (this *Alu) Extub(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	return word_.BitSlice(abi.UNSIGNED, 0, 8)
}

func (this *Alu) Extuh(operand int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	return word_.BitSlice(abi.UNSIGNED, 0, 16)
}

func (this *Alu) MulShSh(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 8, 16) * word2.BitSlice(abi.SIGNED, 8, 16)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulShSl(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 8, 16) * word2.BitSlice(abi.SIGNED, 0, 8)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulShUh(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 8, 16) * word2.BitSlice(abi.UNSIGNED, 8, 16)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulShUl(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 8, 16) * word2.BitSlice(abi.UNSIGNED, 0, 8)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulSlSh(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 0, 8) * word2.BitSlice(abi.SIGNED, 8, 16)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulSlSl(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 0, 8) * word2.BitSlice(abi.SIGNED, 0, 8)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulSlUh(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 0, 8) * word2.BitSlice(abi.UNSIGNED, 8, 16)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulSlUl(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.SIGNED, 0, 8) * word2.BitSlice(abi.UNSIGNED, 0, 8)
	result_word.SetValue(result)

	return result_word.Value(abi.SIGNED)
}

func (this *Alu) MulUhUh(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.UNSIGNED, 8, 16) * word2.BitSlice(abi.UNSIGNED, 8, 16)
	result_word.SetValue(result)

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) MulUhUl(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.UNSIGNED, 8, 16) * word2.BitSlice(abi.UNSIGNED, 0, 8)
	result_word.SetValue(result)

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) MulUlUh(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.UNSIGNED, 0, 8) * word2.BitSlice(abi.UNSIGNED, 8, 16)
	result_word.SetValue(result)

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) MulUlUl(operand1 int64, operand2 int64) int64 {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	if mram_data_width != 4*8 {
		err := errors.New("MRAM data width != 4 * 8")
		panic(err)
	}

	word1 := new(abi.Word)
	word1.Init(mram_data_width)
	word1.SetValue(operand1)

	word2 := new(abi.Word)
	word2.Init(mram_data_width)
	word2.SetValue(operand2)

	result_word := new(abi.Word)
	result_word.Init(mram_data_width)

	result := word1.BitSlice(abi.UNSIGNED, 0, 8) * word2.BitSlice(abi.UNSIGNED, 0, 8)
	result_word.SetValue(result)

	return result_word.Value(abi.UNSIGNED)
}

func (this *Alu) Sats(operand int64) int64 {
	err := errors.New("sats is not yet implemented")
	panic(err)
}

func (this *Alu) Hash(operand1 int64, operand2 int64) int64 {
	err := errors.New("hash is not yet implemented")
	panic(err)
}

func (this *Alu) SignedExtension(operand int64) (int64, int64) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	even_word := new(abi.Word)
	even_word.Init(mram_data_width)
	if word_.SignBit() {
		even_word.SetValue(-1)
	}

	odd_word := new(abi.Word)
	odd_word.Init(mram_data_width)
	odd_word.SetValue(word_.Value(abi.UNSIGNED))

	return even_word.Value(abi.UNSIGNED), odd_word.Value(abi.UNSIGNED)
}

func (this *Alu) UnsignedExtension(operand int64) (int64, int64) {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	mram_data_width := config_loader.MramDataWidth()

	word_ := new(abi.Word)
	word_.Init(mram_data_width)
	word_.SetValue(operand)

	even_word := new(abi.Word)
	even_word.Init(mram_data_width)
	even_word.SetValue(0)

	odd_word := new(abi.Word)
	odd_word.Init(mram_data_width)
	odd_word.SetValue(word_.Value(abi.UNSIGNED))

	return even_word.Value(abi.UNSIGNED), odd_word.Value(abi.UNSIGNED)
}

func (this *Alu) Pow2(exponent int) int64 {
	if exponent < 0 {
		err := errors.New("exponent < 0")
		panic(err)
	}

	value := int64(1)
	for i := 0; i < exponent; i++ {
		value *= 2
	}
	return value
}
