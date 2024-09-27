package abi

import (
	"errors"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser"
	"uPIMulator/src/misc"
)

type Binary struct {
	benchmark    string
	num_dpus     int
	num_tasklets int

	token_stream *lexer.TokenStream
	ast          *parser.Ast
	relocatable  *Relocatable
}

func (this *Binary) Init(benchmark string, num_dpus int, num_tasklets int) {
	if num_dpus <= 0 {
		err := errors.New("num DPUs <= 0")
		panic(err)
	} else if num_tasklets <= 0 {
		err := errors.New("num tasklets <= 0")
		panic(err)
	}

	this.benchmark = benchmark
	this.num_dpus = num_dpus
	this.num_tasklets = num_tasklets

	this.token_stream = nil
	this.ast = nil
	this.relocatable = nil
}

func (this *Binary) Benchmark() string {
	return this.benchmark
}

func (this *Binary) NumDpus() int {
	return this.num_dpus
}

func (this *Binary) NumTasklets() int {
	return this.num_tasklets
}

func (this *Binary) TokenStream() *lexer.TokenStream {
	if this.token_stream == nil {
		err := errors.New("token stream == nil")
		panic(err)
	}

	return this.token_stream
}

func (this *Binary) SetTokenStream(token_stream *lexer.TokenStream) {
	if this.token_stream != nil {
		err := errors.New("token stream != nil")
		panic(err)
	}

	this.token_stream = token_stream
}

func (this *Binary) Ast() *parser.Ast {
	if this.ast == nil {
		err := errors.New("AST == nil")
		panic(err)
	}

	return this.ast
}

func (this *Binary) SetAst(ast *parser.Ast) {
	if this.ast != nil {
		err := errors.New("AST != nil")
		panic(err)
	}

	this.ast = ast
}

func (this *Binary) Relocatable() *Relocatable {
	if this.relocatable == nil {
		err := errors.New("relocatable == nil")
		panic(err)
	}

	return this.relocatable
}

func (this *Binary) SetRelocatable(relocatable *Relocatable) {
	if this.relocatable != nil {
		err := errors.New("relocatable != nil")
		panic(err)
	}

	this.relocatable = relocatable
}

func (this *Binary) Dump(path string) {
	lines := make([]string, 0)

	for _, label := range this.relocatable.Labels() {
		lines = append(lines, label.Stringify())
	}

	file_dumper := new(misc.FileDumper)
	file_dumper.Init(path)
	file_dumper.WriteLines(lines)
}
