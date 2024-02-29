package directive

import (
	"uPIMulator/src/abi/encoding"
)

type AscizDirective struct {
	characters string
}

func (this *AscizDirective) Init(characters string) {
	this.characters = characters
}

func (this *AscizDirective) Characters() string {
	return this.characters
}

func (this *AscizDirective) Size() int64 {
	return int64(len(this.characters)) + 1
}

func (this *AscizDirective) Encode() *encoding.ByteStream {
	ascii_encoder := new(encoding.AsciiEncoder)
	ascii_encoder.Init()

	return ascii_encoder.Encode(this.characters + ascii_encoder.Unknown())
}
