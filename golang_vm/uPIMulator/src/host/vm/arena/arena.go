package arena

import (
	"uPIMulator/src/encoding"
	"uPIMulator/src/host/vm/base"
	"uPIMulator/src/host/vm/type_system"
)

type Arena struct {
	pool *Pool
}

func (this *Arena) Init() {
	this.pool = new(Pool)
	this.pool.Init()
}

func (this *Arena) Pool() *Pool {
	return this.pool
}

func (this *Arena) NewChar(value int64) *base.Object {
	object := this.pool.Alloc(base.TEMPORARY, 1)

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	byte_stream.Append(uint8(value))

	this.pool.Memory().Write(object.Address(), 1, byte_stream)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.CHAR, 0)

	object.SetTypeVariable(type_variable)

	return object
}

func (this *Arena) NewShort(value int64) *base.Object {
	object := this.pool.Alloc(base.TEMPORARY, 2)

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	byte_stream.Append(uint8(value & 0xFF))
	byte_stream.Append(uint8((value >> 8) & 0xFF))
	this.pool.Memory().Write(object.Address(), 2, byte_stream)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.SHORT, 0)

	object.SetTypeVariable(type_variable)

	return object
}

func (this *Arena) NewInt(value int64) *base.Object {
	object := this.pool.Alloc(base.TEMPORARY, 4)

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	byte_stream.Append(uint8(value & 0xFF))
	byte_stream.Append(uint8((value >> 8) & 0xFF))
	byte_stream.Append(uint8((value >> 16) & 0xFF))
	byte_stream.Append(uint8((value >> 24) & 0xFF))
	this.pool.Memory().Write(object.Address(), 4, byte_stream)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.INT, 0)

	object.SetTypeVariable(type_variable)

	return object
}

func (this *Arena) NewLong(value int64) *base.Object {
	object := this.pool.Alloc(base.TEMPORARY, 8)

	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	byte_stream.Append(uint8(value & 0xFF))
	byte_stream.Append(uint8((value >> 8) & 0xFF))
	byte_stream.Append(uint8((value >> 16) & 0xFF))
	byte_stream.Append(uint8((value >> 24) & 0xFF))
	byte_stream.Append(uint8((value >> 32) & 0xFF))
	byte_stream.Append(uint8((value >> 40) & 0xFF))
	byte_stream.Append(uint8((value >> 48) & 0xFF))
	byte_stream.Append(uint8((value >> 56) & 0xFF))
	this.pool.Memory().Write(object.Address(), 8, byte_stream)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.LONG, 0)

	object.SetTypeVariable(type_variable)

	return object
}

func (this *Arena) NewString(value string) *base.Object {
	object := this.pool.Alloc(base.TEMPORARY, int64(len(value)))

	ascii_encoder := new(encoding.AsciiEncoder)
	ascii_encoder.Init()

	encoded_byte_stream := ascii_encoder.Encode(value)

	this.pool.Memory().Write(object.Address(), encoded_byte_stream.Size(), encoded_byte_stream)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitPrimitive(type_system.STRING, 0)

	object.SetTypeVariable(type_variable)

	return object
}

func (this *Arena) NewStruct(struct_name string, size int64) *base.Object {
	object := this.pool.Alloc(base.TEMPORARY, size)

	type_variable := new(type_system.TypeVariable)
	type_variable.InitStruct(type_system.STRUCT, struct_name, 0)

	object.SetTypeVariable(type_variable)

	return object
}

func (this *Arena) NewPointer(size int64) *base.Object {
	object := this.pool.Alloc(base.UNTEMPORARY, size)

	return object
}

func (this *Arena) Free(address int64) {
	this.pool.Free(address)
}
