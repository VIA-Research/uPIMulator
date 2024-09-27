package encoding

import (
	"errors"
)

type ByteStream struct {
	bytes []uint8
}

func (this *ByteStream) Init() {
	this.bytes = make([]uint8, 0)
}

func (this *ByteStream) Size() int64 {
	return int64(len(this.bytes))
}

func (this *ByteStream) Get(pos int) uint8 {
	return this.bytes[pos]
}

func (this *ByteStream) Set(pos int, value uint8) {
	this.bytes[pos] = value
}

func (this *ByteStream) Append(value uint8) {
	this.bytes = append(this.bytes, value)
}

func (this *ByteStream) Remove(pos int) {
	this.bytes = append(this.bytes[:pos], this.bytes[pos+1:]...)
}

func (this *ByteStream) Merge(byte_stream *ByteStream) {
	for i := int64(0); i < byte_stream.Size(); i++ {
		value := byte_stream.Get(int(i))
		this.Append(value)
	}
}

func (this *ByteStream) Signbit() bool {
	last_byte := this.Get(int(this.Size() - 1))
	sign_bit := ((int(last_byte) & (1 << 7)) >> 7) == 1
	return sign_bit
}

func (this *ByteStream) SignedValue() int64 {
	if this.Size() > 8 {
		err := errors.New("byte stream cannot convert into an integer")
		panic(err)
	}

	value := int64(0)
	for i := int64(0); i < this.Size()-1; i++ {
		value += int64(this.Get(int(i))) * this.Pow2(int(8*i))
	}

	last_byte := this.Get(int(this.Size() - 1))

	for i := 0; i < 7; i++ {
		bit := ((int(last_byte) & (1 << i)) >> i) == 1

		if bit {
			value += this.Pow2(8*(len(this.bytes)-1) + i)
		}
	}

	if this.Signbit() {
		value -= this.Pow2(8*len(this.bytes) - 1)
	}

	return value
}

func (this *ByteStream) UnsignedValue() int64 {
	if this.Size() > 8 {
		err := errors.New("byte stream cannot convert into an integer")
		panic(err)
	}

	value := int64(0)
	for i := int64(0); i < this.Size(); i++ {
		value += int64(this.Get(int(i))) * this.Pow2(int(8*i))
	}
	return value
}

func (this *ByteStream) Pow2(exponent int) int64 {
	value := int64(1)
	for i := 0; i < exponent; i++ {
		value *= 2
	}
	return value
}
