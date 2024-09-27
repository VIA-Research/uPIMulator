package linker

import (
	"fmt"
	"uPIMulator/src/device/linker/kernel"
	"uPIMulator/src/device/linker/lexer"
)

type LexJob struct {
	relocatable *kernel.Relocatable
}

func (this *LexJob) Init(relocatable *kernel.Relocatable) {
	this.relocatable = relocatable
}

func (this *LexJob) Execute() {
	fmt.Printf("Lexing %s...\n", this.relocatable.Path())

	lexer_ := new(lexer.Lexer)
	lexer_.Init()

	token_stream := lexer_.Lex(this.relocatable.Path())
	this.relocatable.SetTokenStream(token_stream)
}
