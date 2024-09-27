package frame

import (
	"errors"
	"uPIMulator/src/host/abi"
	"uPIMulator/src/host/vm/symbol"
)

type FrameChain struct {
	global_scope *symbol.Scope
	frames       []*Frame
}

func (this *FrameChain) Init() {
	this.global_scope = new(symbol.Scope)
	this.global_scope.Init()

	this.frames = make([]*Frame, 0)
}

func (this *FrameChain) Bootstrap(label *abi.Label) {
	frame := new(Frame)
	frame.Init(label)

	this.frames = append(this.frames, frame)
}

func (this *FrameChain) GlobalScope() *symbol.Scope {
	return this.global_scope
}

func (this *FrameChain) Length() int {
	return len(this.frames)
}

func (this *FrameChain) Symbols() []*symbol.Symbol {
	symbols := make([]*symbol.Symbol, 0)

	symbols = append(symbols, this.global_scope.Symbols()...)

	for _, frame := range this.frames {
		symbols = append(symbols, frame.Symbols()...)
	}

	return symbols
}

func (this *FrameChain) HasObject(address int64) bool {
	for _, frame := range this.frames {
		if frame.Stack().HasObject(address) || frame.ReturnStack().HasObject(address) {
			return true
		}
	}
	return false
}

func (this *FrameChain) NewFrame(label *abi.Label) {
	frame := new(Frame)
	frame.Init(label)

	frame.FastScopeChain().NewScope()
	for _, symbol_ := range this.LastFrame().ArgScope().Symbols() {
		frame.FastScopeChain().AddSymbol(symbol_)
	}
	this.LastFrame().ArgScope().Clear()

	this.frames = append(this.frames, frame)
}

func (this *FrameChain) DeleteFrame() {
	last_frame := this.LastFrame()

	if len(this.frames) >= 2 {
		for i := 0; i < last_frame.ReturnStack().Length(); i++ {
			stack_item := last_frame.ReturnStack().Front(0)
			last_frame.ReturnStack().Pop()

			this.frames[len(this.frames)-2].Stack().Push(stack_item)
		}
	}

	this.frames = this.frames[:len(this.frames)-1]
}

func (this *FrameChain) Frame(pos int) *Frame {
	return this.frames[pos]
}

func (this *FrameChain) LastFrame() *Frame {
	return this.frames[len(this.frames)-1]
}

func (this *FrameChain) CanAdvance() bool {
	for this.Length() > 0 {
		last_frame := this.LastFrame()
		if last_frame.Pc().CanAdvance() {
			return true
		} else {
			this.DeleteFrame()
		}
	}

	return false
}

func (this *FrameChain) Advance() *abi.Bytecode {
	if !this.CanAdvance() {
		err := errors.New("frame chain cannot advance")
		panic(err)
	}

	return this.LastFrame().Pc().Advance()
}
