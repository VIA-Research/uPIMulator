package linker

import (
	"fmt"
	"uPIMulator/src/device/linker/kernel"
	"uPIMulator/src/device/linker/parser"
)

type ParseJob struct {
	relocatable *kernel.Relocatable
}

func (this *ParseJob) Init(relocatable *kernel.Relocatable) {
	this.relocatable = relocatable
}

func (this *ParseJob) Execute() {
	fmt.Printf("Parsing %s...\n", this.relocatable.Path())

	parser_ := new(parser.Parser)
	parser_.Init()

	ast := parser_.Parse(this.relocatable.TokenStream())
	this.relocatable.SetAst(ast)
}
