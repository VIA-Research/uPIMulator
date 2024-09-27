package arena

import (
	"errors"
	"uPIMulator/src/host/vm/base"
	"uPIMulator/src/host/vm/frame"
	"uPIMulator/src/host/vm/symbol"
	"uPIMulator/src/host/vm/type_system"
	"uPIMulator/src/misc"
)

type GarbageCollector struct {
	cycle     int64
	threshold int64

	arena       *Arena
	frame_chain *frame.FrameChain
	registry    *type_system.Registry
}

func (this *GarbageCollector) Init() {
	config_loader := new(misc.ConfigLoader)
	config_loader.Init()

	this.cycle = 0
	this.threshold = config_loader.GarbageCollectionThreshold()
}

func (this *GarbageCollector) ConnectArena(arena *Arena) {
	this.arena = arena
}

func (this *GarbageCollector) ConnectFrameChain(frame_chain *frame.FrameChain) {
	this.frame_chain = frame_chain
}

func (this *GarbageCollector) ConnectRegistry(registry *type_system.Registry) {
	this.registry = registry
}

func (this *GarbageCollector) MarkAndSweep() {
	if this.cycle > this.threshold {
		objects := this.Mark()
		this.Sweep(objects)

		this.cycle = 0
	}

	this.cycle++
}

func (this *GarbageCollector) Mark() []*base.Object {
	objects := make([]*base.Object, 0)
	for _, symbol_ := range this.frame_chain.Symbols() {
		objects = append(objects, this.ChaseSymbol(symbol_)...)
	}
	return objects
}

func (this *GarbageCollector) Sweep(objects []*base.Object) {
	for _, obj := range this.arena.Pool().Objects() {
		if obj.ObjectType() == base.TEMPORARY && !this.IsMarked(obj, objects) &&
			!this.frame_chain.HasObject(obj.Address()) &&
			this.arena.Pool().HasObject(obj.Address()) {
			this.arena.Free(obj.Address())
		}
	}
}

func (this *GarbageCollector) ChaseSymbol(symbol_ *symbol.Symbol) []*base.Object {
	objects := make([]*base.Object, 0)

	type_variable := symbol_.TypeVariable()

	if type_variable.TypeVariableType() == type_system.VOID {
		if type_variable.NumStars() == 0 {
			err := errors.New("num stars == 0")
			panic(err)
		}

		objects = append(objects, this.ChaseObject(symbol_.Object())...)
	} else if type_variable.TypeVariableType() == type_system.CHAR {
		objects = append(objects, this.ChaseObject(symbol_.Object())...)
	} else if type_variable.TypeVariableType() == type_system.SHORT {
		objects = append(objects, this.ChaseObject(symbol_.Object())...)
	} else if type_variable.TypeVariableType() == type_system.INT {
		objects = append(objects, this.ChaseObject(symbol_.Object())...)
	} else if type_variable.TypeVariableType() == type_system.LONG {
		objects = append(objects, this.ChaseObject(symbol_.Object())...)
	} else if type_variable.TypeVariableType() == type_system.STRUCT {
		objects = append(objects, this.ChaseObject(symbol_.Object())...)
	} else {
		err := errors.New("type variable type is not valid")
		panic(err)
	}

	return objects
}

func (this *GarbageCollector) ChaseObject(object *base.Object) []*base.Object {
	objects := make([]*base.Object, 0)

	objects = append(objects, object)

	if !object.HasTypeVariable() {
		err := errors.New("object does not have a type variable")
		panic(err)
	}

	type_variable := object.TypeVariable()

	if type_variable.TypeVariableType() == type_system.VOID {
		if type_variable.NumStars() == 0 {
			err := errors.New("num stars == 0")
			panic(err)
		}

		if type_variable.NumStars() > 0 {
			value := this.arena.Pool().Memory().Read(object.Address(), 4).SignedValue()
			objects = append(objects, this.ChasePointer(value)...)
		}
	} else if type_variable.TypeVariableType() == type_system.CHAR {
		if type_variable.NumStars() > 0 {
			value := this.arena.Pool().Memory().Read(object.Address(), 4).SignedValue()
			objects = append(objects, this.ChasePointer(value)...)
		}
	} else if type_variable.TypeVariableType() == type_system.SHORT {
		if type_variable.NumStars() > 0 {
			value := this.arena.Pool().Memory().Read(object.Address(), 4).SignedValue()
			objects = append(objects, this.ChasePointer(value)...)
		}
	} else if type_variable.TypeVariableType() == type_system.INT {
		if type_variable.NumStars() > 0 {
			value := this.arena.Pool().Memory().Read(object.Address(), 4).SignedValue()
			objects = append(objects, this.ChasePointer(value)...)
		}
	} else if type_variable.TypeVariableType() == type_system.LONG {
		if type_variable.NumStars() > 0 {
			value := this.arena.Pool().Memory().Read(object.Address(), 4).SignedValue()
			objects = append(objects, this.ChasePointer(value)...)
		}
	} else if type_variable.TypeVariableType() == type_system.STRUCT {
		if type_variable.NumStars() > 0 {
			value := this.arena.Pool().Memory().Read(object.Address(), 4).SignedValue()
			objects = append(objects, this.ChasePointer(value)...)
		} else {
			struct_name := type_variable.StructName()

			skeleton := this.registry.Skeleton(struct_name)

			for i, field := range skeleton.Fields() {
				if i == 0 {
					if field.TypeVariable().NumStars() > 0 {
						value := this.arena.Pool().Memory().Read(object.Address(), 4).SignedValue()
						objects = append(objects, this.ChasePointer(value)...)
					} else {
						objects = append(objects, this.arena.Pool().LastObject(object.Address()))
					}
				} else {
					field_offset := this.registry.FieldOffset(skeleton.Name(), field.Name())

					field_address := object.Address() + field_offset

					if this.arena.Pool().HasObject(field_address) {
						field_object := this.arena.Pool().Object(field_address)
						objects = append(objects, this.ChaseObject(field_object)...)
					}
				}
			}
		}
	} else {
		err := errors.New("type variable type is not valid")
		panic(err)
	}

	return objects
}

func (this *GarbageCollector) ChasePointer(address int64) []*base.Object {
	object := this.arena.Pool().Object(address)
	return this.ChaseObject(object)
}

func (this *GarbageCollector) IsMarked(object *base.Object, objects []*base.Object) bool {
	for _, obj := range objects {
		if obj == object {
			return true
		}
	}
	return false
}
