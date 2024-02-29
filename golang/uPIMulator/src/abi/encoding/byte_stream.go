package encoding

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

func (this *ByteStream) Merge(byte_stream *ByteStream) {
	for i := int64(0); i < byte_stream.Size(); i++ {
		value := byte_stream.Get(int(i))
		this.Append(value)
	}
}
