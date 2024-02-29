package word

import (
	"errors"
	"uPIMulator/src/abi/encoding"
)

type Representation int

const (
	UNSIGNED Representation = iota
	SIGNED
)

type Word struct {
	bits []bool
}

func (this *Word) Init(width int) {
	if width <= 0 {
		err := errors.New("width <= 0")
		panic(err)
	}

	this.bits = make([]bool, width)
}

func (this *Word) Width() int {
	return len(this.bits)
}

func (this *Word) Size() int {
	if this.Width()%8 != 0 {
		err := errors.New("width is not a multiple of 8")
		panic(err)
	}

	return this.Width() / 8
}

func (this *Word) SignBit() bool {
	return this.bits[this.Width()-1]
}

func (this *Word) Bit(pos int) bool {
	return this.bits[pos]
}

func (this *Word) SetBit(pos int) {
	this.bits[pos] = true
}

func (this *Word) ClearBit(pos int) {
	this.bits[pos] = false
}

func (this *Word) BitSlice(representation Representation, begin int, end int) int64 {
	this.VerifySlice(begin, end)

	slice_width := end - begin
	value := int64(0)
	for i := 0; i < slice_width; i++ {
		if this.Bit(begin + i) {
			if representation == SIGNED && i == slice_width-1 {
				value -= this.Pow2(i)
			} else {
				value += this.Pow2(i)
			}
		}
	}
	return value
}

func (this *Word) SetBitSlice(begin int, end int, value int64) {
	this.VerifySlice(begin, end)

	if value >= 0 {
		this.SetPositiveBitSlice(begin, end, value)
	} else {
		this.SetNegativeBitSlice(begin, end, value)
	}
}

func (this *Word) Value(representation Representation) int64 {
	return this.BitSlice(representation, 0, this.Width())
}

func (this *Word) SetValue(value int64) {
	this.SetBitSlice(0, this.Width(), value)
}

func (this *Word) ToByteStream() *encoding.ByteStream {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := 0; i < this.Size(); i++ {
		begin := 8 * i
		end := begin + 8

		value := uint8(this.BitSlice(UNSIGNED, begin, end))
		byte_stream.Append(value)
	}

	return byte_stream
}

func (this *Word) FromByteStream(byte_stream *encoding.ByteStream) {
	for i := int64(0); i < byte_stream.Size(); i++ {
		begin := int(8 * i)
		end := int(begin + 8)

		value := int64(byte_stream.Get(int(i)))
		this.SetBitSlice(begin, end, value)
	}
}

func (this *Word) VerifySlice(begin int, end int) {
	if begin < 0 {
		err := errors.New("begin < 0")
		panic(err)
	}

	if begin >= end {
		err := errors.New("begin >= end")
		panic(err)
	}

	if end > this.Width() {
		err := errors.New("end > width")
		panic(err)
	}

	if end-begin > 64 {
		err := errors.New("end - begin >= 64")
		panic(err)
	}
}

func (this *Word) Pow2(exponent int) int64 {
	value := int64(1)
	for i := 0; i < exponent; i++ {
		value *= 2
	}
	return value
}

func (this *Word) SetPositiveBitSlice(begin int, end int, value int64) {
	this.VerifySlice(begin, end)

	if value < 0 {
		err := errors.New("value < 0")
		panic(err)
	}

	slice_width := end - begin
	for i := 0; i < slice_width; i++ {
		if value%2 == 1 {
			this.SetBit(begin + i)
		} else {
			this.ClearBit(begin + i)
		}

		value /= 2
	}

	if value != 0 {
		err := errors.New("value != 0")
		panic(err)
	}
}

func (this *Word) SetNegativeBitSlice(begin int, end int, value int64) {
	this.VerifySlice(begin, end)

	if value >= 0 {
		err := errors.New("value >= 0")
		panic(err)
	}

	this.SetBit(end - 1)

	if begin+1 < end {
		slice_width := end - begin
		value += this.Pow2(slice_width - 1)
		this.SetPositiveBitSlice(begin, end-1, value)
	}
}
