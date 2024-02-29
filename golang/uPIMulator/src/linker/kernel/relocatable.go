package kernel

import (
	"errors"
	"strings"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser"
	"uPIMulator/src/misc"
)

type Relocatable struct {
	name string
	path string

	token_stream *lexer.TokenStream
	ast          *parser.Ast
	liveness     *Liveness

	renames map[string]string
}

func (this *Relocatable) Init(name string) {
	this.name = name

	this.renames = make(map[string]string, 0)
}

func (this *Relocatable) Name() string {
	return this.name
}

func (this *Relocatable) Path() string {
	return this.path
}

func (this *Relocatable) SetPath(path string) {
	this.path = path
}

func (this *Relocatable) TokenStream() *lexer.TokenStream {
	return this.token_stream
}

func (this *Relocatable) SetTokenStream(token_stream *lexer.TokenStream) {
	this.token_stream = token_stream
}

func (this *Relocatable) Ast() *parser.Ast {
	return this.ast
}

func (this *Relocatable) SetAst(ast *parser.Ast) {
	this.ast = ast
}

func (this *Relocatable) Liveness() *Liveness {
	return this.liveness
}

func (this *Relocatable) SetLiveness(liveness *Liveness) {
	this.liveness = liveness
}

func (this *Relocatable) Lines() []string {
	file_scanner := new(misc.FileScanner)
	file_scanner.Init(this.path)

	lines := file_scanner.ReadLines()
	for i, line := range lines {
		lines[i] = this.RenameLine(line)
	}
	return lines
}

func (this *Relocatable) RenameLocalSymbol(old_name string, new_name string) {
	if _, found := this.liveness.LocalSymbols()[old_name]; !found {
		err := errors.New("local symbol is not found")
		panic(err)
	}

	if rename, found := this.renames[old_name]; found {
		if rename != new_name {
			err := errors.New("rename is already set")
			panic(err)
		}
	}

	this.renames[old_name] = new_name
}

func (this *Relocatable) RenameLine(line string) string {
	for old_name, new_name := range this.renames {
		line = strings.ReplaceAll(line, old_name+",", new_name+",")
		line = strings.ReplaceAll(line, old_name+" ", new_name+" ")
		line = strings.ReplaceAll(line, old_name+"\t", new_name+"\t")
		line = strings.ReplaceAll(line, old_name+":", new_name+":")
		line = strings.ReplaceAll(line, old_name+"+", new_name+"+")
		line = strings.ReplaceAll(line, old_name+"-", new_name+"-")

		if len(line) > len(old_name) {
			pos := len(line) - len(old_name)

			if line[pos:] == old_name {
				line = line[:pos] + new_name
			}
		}
	}
	return line
}
