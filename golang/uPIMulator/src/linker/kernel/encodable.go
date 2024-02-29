package kernel

import (
	"uPIMulator/src/abi/encoding"
)

type Encodable interface {
	Encode() *encoding.ByteStream
}
