package interpreter

import (
	"path/filepath"
	"uPIMulator/src/host/abi"
	"uPIMulator/src/host/interpreter/codegen"
	"uPIMulator/src/host/interpreter/lexer"
	"uPIMulator/src/host/interpreter/parser"
	"uPIMulator/src/misc"
)

type Interpreter struct {
	root_dirpath     string
	bin_dirpath      string
	benchmark        string
	num_dpus         int
	num_tasklets     int
	data_prep_params int

	dpu_mram_heap_pointer_name int64

	binary *abi.Binary
}

func (this *Interpreter) Init(
	command_line_parser *misc.CommandLineParser,
	dpu_mram_heap_pointer_name int64,
) {
	this.root_dirpath = command_line_parser.StringParameter("root_dirpath")
	this.bin_dirpath = command_line_parser.StringParameter("bin_dirpath")
	this.benchmark = command_line_parser.StringParameter("benchmark")

	num_channels := int(command_line_parser.IntParameter("num_channels"))
	num_ranks_per_channel := int(command_line_parser.IntParameter("num_ranks_per_channel"))
	num_dpus_per_rank := int(command_line_parser.IntParameter("num_dpus_per_rank"))
	this.num_dpus = num_channels * num_ranks_per_channel * num_dpus_per_rank

	this.num_tasklets = int(command_line_parser.IntParameter("num_tasklets"))

	this.data_prep_params = int(command_line_parser.IntParameter("data_prep_params"))

	this.dpu_mram_heap_pointer_name = dpu_mram_heap_pointer_name

	this.binary = new(abi.Binary)
	this.binary.Init(this.benchmark, this.num_dpus, this.num_tasklets)
}

func (this *Interpreter) Interpret() {
	this.Lex()
	this.Parse()
	this.Codegen()

	binary_path := filepath.Join(this.bin_dirpath, "bytecode.txt")
	this.binary.Dump(binary_path)
}

func (this *Interpreter) Lex() {
	source_code := this.FindSourceCode()

	lexer_ := new(lexer.Lexer)
	lexer_.Init()

	token_stream := lexer_.Lex(source_code)
	this.binary.SetTokenStream(token_stream)
}

func (this *Interpreter) Parse() {
	parser_ := new(parser.Parser)
	parser_.Init()

	ast := parser_.Parse(this.binary.TokenStream())
	this.binary.SetAst(ast)
}

func (this *Interpreter) Codegen() {
	codegen_ := new(codegen.Codegen)
	codegen_.Init(
		this.benchmark,
		this.num_dpus,
		this.num_tasklets,
		this.data_prep_params,
		this.dpu_mram_heap_pointer_name,
	)

	relocatable := codegen_.Codegen(this.binary.Ast())
	this.binary.SetRelocatable(relocatable)
}

func (this *Interpreter) FindSourceCode() string {
	benchmark_dirpath := filepath.Join(this.root_dirpath, "benchmark", this.benchmark)
	c_path := filepath.Join(benchmark_dirpath, "host", "app.c")
	return c_path
}
