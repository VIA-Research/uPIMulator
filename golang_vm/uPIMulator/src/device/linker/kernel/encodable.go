package kernel

import (
	"uPIMulator/src/encoding"
)

type Encodable interface {
	Encode() *encoding.ByteStream
}
