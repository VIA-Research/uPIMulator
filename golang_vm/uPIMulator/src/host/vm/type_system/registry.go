package type_system

import (
	"errors"
	"fmt"
)

type Registry struct {
	skeletons map[string]*Skeleton
}

func (this *Registry) Init() {
	this.skeletons = make(map[string]*Skeleton)
}

func (this *Registry) HasSkeleton(skeleton_name string) bool {
	_, found := this.skeletons[skeleton_name]
	return found
}

func (this *Registry) Skeleton(skeleton_name string) *Skeleton {
	if !this.HasSkeleton(skeleton_name) {
		err_msg := fmt.Sprintf("skeleton (%s) is not found", skeleton_name)
		err := errors.New(err_msg)
		panic(err)
	}

	return this.skeletons[skeleton_name]
}

func (this *Registry) AddSkeleton(skeleton *Skeleton) {
	if this.HasSkeleton(skeleton.Name()) {
		err_msg := fmt.Sprintf("skeleton (%s) already exists", skeleton.Name())
		err := errors.New(err_msg)
		panic(err)
	}

	this.skeletons[skeleton.Name()] = skeleton
}

func (this *Registry) SkeletonSize(skeleton_name string) int64 {
	if !this.HasSkeleton(skeleton_name) {
		err_msg := fmt.Sprintf("skeleton (%s) is not found", skeleton_name)
		err := errors.New(err_msg)
		panic(err)
	}

	skeleton := this.Skeleton(skeleton_name)

	offset := int64(0)
	for i := 0; i < skeleton.Length(); i++ {
		field := skeleton.Get(i)

		if field.TypeVariable().NumStars() > 0 {
			offset += 4
		} else {
			if field.TypeVariable().TypeVariableType() == VOID {
				err := errors.New("type variable type is void")
				panic(err)
			} else if field.TypeVariable().TypeVariableType() == CHAR {
				offset += 1
			} else if field.TypeVariable().TypeVariableType() == SHORT {
				offset += 2
			} else if field.TypeVariable().TypeVariableType() == INT {
				offset += 4
			} else if field.TypeVariable().TypeVariableType() == LONG {
				offset += 8
			} else if field.TypeVariable().TypeVariableType() == STRUCT {
				offset += this.SkeletonSize(field.TypeVariable().StructName())
			}
		}
	}
	return offset
}

func (this *Registry) FieldOffset(skeleton_name string, field_name string) int64 {
	if !this.HasSkeleton(skeleton_name) {
		err_msg := fmt.Sprintf("skeleton (%s) is not found", skeleton_name)
		err := errors.New(err_msg)
		panic(err)
	}

	skeleton := this.Skeleton(skeleton_name)

	offset := int64(0)
	for i := 0; i < skeleton.Length(); i++ {
		field := skeleton.Get(i)

		if field.Name() == field_name {
			break
		}

		if field.TypeVariable().NumStars() > 0 {
			offset += 4
		} else {
			if field.TypeVariable().TypeVariableType() == VOID {
				err := errors.New("type variable type is void")
				panic(err)
			} else if field.TypeVariable().TypeVariableType() == CHAR {
				offset += 1
			} else if field.TypeVariable().TypeVariableType() == SHORT {
				offset += 2
			} else if field.TypeVariable().TypeVariableType() == INT {
				offset += 4
			} else if field.TypeVariable().TypeVariableType() == LONG {
				offset += 8
			} else if field.TypeVariable().TypeVariableType() == STRUCT {
				offset += this.SkeletonSize(field.TypeVariable().StructName())
			} else {
				err := errors.New("type variable type is not valid")
				panic(err)
			}
		}
	}
	return offset
}

func (this *Registry) FieldSize(skeleton_name string, field_name string) int64 {
	if !this.HasSkeleton(skeleton_name) {
		err_msg := fmt.Sprintf("skeleton (%s) is not found", skeleton_name)
		err := errors.New(err_msg)
		panic(err)
	}

	skeleton := this.Skeleton(skeleton_name)
	field := skeleton.Field(field_name)

	if field.TypeVariable().NumStars() > 0 {
		return 4
	} else {
		if field.TypeVariable().TypeVariableType() == VOID {
			err := errors.New("type variable type is void")
			panic(err)
		} else if field.TypeVariable().TypeVariableType() == CHAR {
			return 1
		} else if field.TypeVariable().TypeVariableType() == SHORT {
			return 2
		} else if field.TypeVariable().TypeVariableType() == INT {
			return 4
		} else if field.TypeVariable().TypeVariableType() == LONG {
			return 8
		} else if field.TypeVariable().TypeVariableType() == STRUCT {
			return this.SkeletonSize(field.TypeVariable().StructName())
		} else {
			err := errors.New("type variable type is not valid")
			panic(err)
		}
	}
}
