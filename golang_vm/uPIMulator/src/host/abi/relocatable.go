package abi

import (
	"errors"
	"fmt"
)

type Relocatable struct {
	labels []*Label

	cur_func *Label

	cur_loop_condition *Label
	cur_loop_body      *Label
	cur_loop_end       *Label

	cur_label *Label

	label_brk int
}

func (this *Relocatable) Init() {
	this.labels = make([]*Label, 0)

	this.cur_func = nil

	this.cur_loop_condition = nil
	this.cur_loop_body = nil
	this.cur_loop_end = nil

	this.cur_label = nil

	this.label_brk = 0
}

func (this *Relocatable) HasFunc(func_name string) bool {
	return this.HasLabel(func_name)
}

func (this *Relocatable) CurFunc() *Label {
	if this.cur_func == nil {
		err := errors.New("cur func == nil")
		panic(err)
	}

	return this.cur_func
}

func (this *Relocatable) NewFunc(func_name string) *Label {
	label := this.NewNamedLabel(func_name)

	return label
}

func (this *Relocatable) SwitchFunc(func_name string) {
	this.cur_func = this.Label(func_name)
}

func (this *Relocatable) CurLoopCondition() *Label {
	if this.cur_loop_condition == nil {
		err := errors.New("cur loop condition == nil")
		panic(err)
	}

	return this.cur_loop_condition
}

func (this *Relocatable) CurLoopBody() *Label {
	if this.cur_loop_body == nil {
		err := errors.New("cur loop body == nil")
		panic(err)
	}

	return this.cur_loop_body
}

func (this *Relocatable) CurLoopEnd() *Label {
	if this.cur_loop_end == nil {
		err := errors.New("cur loop end == nil")
		panic(err)
	}

	return this.cur_loop_end
}

func (this *Relocatable) NewLoop() (*Label, *Label, *Label) {
	this.cur_loop_condition = this.NewUnnamedLabel()
	this.cur_loop_body = this.NewUnnamedLabel()
	this.cur_loop_end = this.NewUnnamedLabel()

	return this.cur_loop_condition, this.cur_loop_body, this.cur_loop_end
}

func (this *Relocatable) HasLabel(label_name string) bool {
	for _, label := range this.labels {
		if label.Name() == label_name {
			return true
		}
	}

	return false
}

func (this *Relocatable) Label(label_name string) *Label {
	if !this.HasLabel(label_name) {
		err_msg := fmt.Sprintf("label (%s) is not found", label_name)
		err := errors.New(err_msg)
		panic(err)
	}

	for _, label := range this.labels {
		if label.Name() == label_name {
			return label
		}
	}

	return nil
}

func (this *Relocatable) Labels() []*Label {
	return this.labels
}

func (this *Relocatable) NewNamedLabel(label_name string) *Label {
	if this.HasLabel(label_name) {
		err_msg := fmt.Sprintf("label (%s) already exists", label_name)
		err := errors.New(err_msg)
		panic(err)
	}

	label := new(Label)
	label.Init(label_name)

	this.labels = append(this.labels, label)

	return label
}

func (this *Relocatable) NewUnnamedLabel() *Label {
	label_name := fmt.Sprintf("L%d", this.label_brk)
	this.label_brk++

	label := new(Label)
	label.Init(label_name)

	this.labels = append(this.labels, label)

	return label
}

func (this *Relocatable) SwitchLabel(label_name string) {
	this.cur_label = this.Label(label_name)
}

func (this *Relocatable) NewBytecode(op_code OpCode, args []int64, strs []string) {
	bytecode := new(Bytecode)
	bytecode.Init(op_code, args, strs)

	this.cur_label.Append(bytecode)
}
