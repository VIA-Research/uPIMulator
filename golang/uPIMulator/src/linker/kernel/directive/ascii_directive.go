package directive

import (
	"uPIMulator/src/abi/encoding"
)

type AsciiDirective struct {
	characters string
}

func (this *AsciiDirective) Init(characters string) {
	this.characters = characters
}

func (this *AsciiDirective) Characters() string {
	return this.characters
}

func (this *AsciiDirective) Size() int64 {
	return int64(len(this.characters))
}

func (this *AsciiDirective) Encode() *encoding.ByteStream {
	ascii_encoder := new(encoding.AsciiEncoder)
	ascii_encoder.Init()

	return ascii_encoder.Encode(this.characters)
}
