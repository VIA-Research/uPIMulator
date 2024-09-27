package arena

import (
	"errors"
	"fmt"
	"slices"
	"uPIMulator/src/encoding"
	"uPIMulator/src/host/vm/base"
	"uPIMulator/src/misc"
)

type Pool struct {
	offset  int64
	objects []*base.Object
	memory  *Memory
}

func (this *Pool) Init() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.offset = config_loader.VmBankOffset()

	this.objects = make([]*base.Object, 0)

	this.memory = new(Memory)
	this.memory.Init(config_loader.VmMemorySize())
}

func (this *Pool) Memory() *Memory {
	return this.memory
}

func (this *Pool) Alloc(object_type base.ObjectType, size int64) *base.Object {
	cur_offset := this.offset
	if len(this.objects) == 0 {
		object := new(base.Object)
		object.Init(object_type, cur_offset, size)
		this.objects = append(this.objects, object)

		this.DoZeros(object)

		return object
	} else {
		for _, obj := range this.objects {
			if cur_offset+size <= obj.Address() {
				object := new(base.Object)
				object.Init(object_type, cur_offset, size)
				this.objects = append(this.objects, object)

				sort_fn := func(obj1 *base.Object, obj2 *base.Object) int {
					if obj1.Address() < obj2.Address() {
						return -1
					} else if obj1.Address() == obj2.Address() {
						return 0
					} else {
						return 1
					}
				}

				slices.SortFunc(this.objects, sort_fn)

				this.DoZeros(object)

				return object
			} else {
				cur_offset = obj.Address() + obj.Size()
			}
		}

		object := new(base.Object)
		object.Init(object_type, cur_offset, size)
		this.objects = append(this.objects, object)

		this.DoZeros(object)

		return object
	}
}

func (this *Pool) Free(address int64) {
	if !this.HasObject(address) {
		err_msg := fmt.Sprintf("object with address (%d) is not found", address)
		err := errors.New(err_msg)
		panic(err)
	}

	for i, object := range this.objects {
		if object.Address() == address {
			this.objects = append(this.objects[:i], this.objects[i+1:]...)
			break
		}
	}
}

func (this *Pool) HasObject(address int64) bool {
	for _, object := range this.objects {
		if object.Address() == address {
			return true
		}
	}
	return false
}

func (this *Pool) Object(address int64) *base.Object {
	if !this.HasObject(address) {
		err_msg := fmt.Sprintf("object with address (%d) is not found", address)
		err := errors.New(err_msg)
		panic(err)
	}

	for _, object := range this.objects {
		if object.Address() == address {
			return object
		}
	}

	return nil
}

func (this *Pool) HasLastObject(address int64) bool {
	for i := len(this.objects) - 1; i >= 0; i-- {
		object := this.objects[i]

		if object.Address() == address {
			return true
		}
	}

	return false
}

func (this *Pool) LastObject(address int64) *base.Object {
	if !this.HasLastObject(address) {
		err_msg := fmt.Sprintf("object with address (%d) is not found", address)
		err := errors.New(err_msg)
		panic(err)
	}

	for i := len(this.objects) - 1; i >= 0; i-- {
		object := this.objects[i]

		if object.Address() == address {
			return object
		}
	}

	return nil
}

func (this *Pool) Objects() []*base.Object {
	return this.objects
}

func (this *Pool) DoZeros(object *base.Object) {
	byte_stream := new(encoding.ByteStream)
	byte_stream.Init()

	for i := int64(0); i < object.Size(); i++ {
		byte_stream.Append(0)
	}

	this.memory.Write(object.Address(), object.Size(), byte_stream)
}
