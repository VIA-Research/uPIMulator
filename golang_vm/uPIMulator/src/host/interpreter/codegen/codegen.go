package codegen

import (
	"errors"
	"strconv"
	"uPIMulator/src/host/abi"
	"uPIMulator/src/host/interpreter/codegen/type_system"
	"uPIMulator/src/host/interpreter/parser"
	"uPIMulator/src/host/interpreter/parser/decl"
	"uPIMulator/src/host/interpreter/parser/directive"
	"uPIMulator/src/host/interpreter/parser/expr"
	"uPIMulator/src/host/interpreter/parser/stmt"
	"uPIMulator/src/host/interpreter/parser/type_specifier"
)

type Codegen struct {
	benchmark        string
	num_dpus         int
	num_tasklets     int
	data_prep_params int

	dpu_mram_heap_pointer_name int64

	type_system *type_system.TypeSystem

	relocatable *abi.Relocatable
}

func (this *Codegen) Init(
	benchmark string,
	num_dpus int,
	num_tasklets int,
	data_prep_params int,
	dpu_mram_heap_pointer_name int64,
) {
	this.benchmark = benchmark
	this.num_dpus = num_dpus
	this.num_tasklets = num_tasklets
	this.data_prep_params = data_prep_params
	this.dpu_mram_heap_pointer_name = dpu_mram_heap_pointer_name

	this.type_system = new(type_system.TypeSystem)
	this.type_system.Init()

	this.relocatable = new(abi.Relocatable)
	this.relocatable.Init()
}

func (this *Codegen) Codegen(ast *parser.Ast) *abi.Relocatable {
	this.CodegenInitBootstrap()

	for i := 0; i < ast.Length(); i++ {
		stack_item := ast.Get(i)

		if stack_item.StackItemType() == parser.DIRECTIVE {
			this.CodegenDirective(stack_item.Directive())
		} else if stack_item.StackItemType() == parser.DECL {
			this.CodegenDecl(stack_item.Decl())
		} else {
			err := errors.New("stack item is not directive nor decl")
			panic(err)
		}
	}

	this.CodegenFiniBootstrap()

	return this.relocatable
}

func (this *Codegen) CodegenInitBootstrap() {
	this.relocatable.NewFunc("__bootstrap")
	this.relocatable.SwitchLabel("__bootstrap")

	this.relocatable.NewBytecode(abi.NEW_SCOPE, []int64{}, []string{})

	this.relocatable.NewBytecode(abi.BEGIN_STRUCT, []int64{}, []string{"dpu_set_t"})
	this.relocatable.NewBytecode(abi.APPEND_INT, []int64{0}, []string{"foo"})
	this.relocatable.NewBytecode(abi.END_STRUCT, []int64{}, []string{})
}

func (this *Codegen) CodegenFiniBootstrap() {
	this.relocatable.SwitchLabel("__bootstrap")

	this.relocatable.NewBytecode(abi.CALL, []int64{}, []string{"main"})
	this.relocatable.NewBytecode(abi.DELETE_SCOPE, []int64{}, []string{})
}

func (this *Codegen) CodegenDirective(directive_ *directive.Directive) {
	if directive_.DirectiveType() == directive.INCLUDE {
		this.CodegenIncludeDirective(directive_.IncludeDirective())
	} else if directive_.DirectiveType() == directive.DEFINE {
		this.CodegenDefineDirective(directive_.DefineDirective())
	} else {
		err := errors.New("directive type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenIncludeDirective(include_directive *directive.IncludeDirective) {
}

func (this *Codegen) CodegenDefineDirective(define_directive *directive.DefineDirective) {
	this.relocatable.SwitchLabel("__bootstrap")

	lvalue := define_directive.Lvalue()
	rvalue := define_directive.Rvalue()

	if rvalue.ExprType() == expr.PRIMARY {
		if rvalue.PrimaryExpr().PrimaryExprType() == expr.IDENTIFIER {
			err := errors.New("rvalue primary expr type is identifier")
			panic(err)
		} else if rvalue.PrimaryExpr().PrimaryExprType() == expr.NUMBER {
			if lvalue.Attribute() == "NUM_DPUS" {
				this.relocatable.NewBytecode(abi.NEW_GLOBAL_INT, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.PUSH_INT, []int64{int64(this.num_dpus)}, []string{})
				this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})
			} else if lvalue.Attribute() == "NUM_TASKLETS" {
				this.relocatable.NewBytecode(abi.NEW_GLOBAL_INT, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.PUSH_INT, []int64{int64(this.num_tasklets)}, []string{})
				this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})
			} else if lvalue.Attribute() == "DATA_PREP_PARAMS" {
				this.relocatable.NewBytecode(abi.NEW_GLOBAL_INT, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.PUSH_INT, []int64{int64(this.data_prep_params)}, []string{})
				this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})
			} else {
				rvalue_value, err := strconv.ParseInt(rvalue.PrimaryExpr().Token().Attribute(), 10, 64)
				if err != nil {
					panic(err)
				}

				this.relocatable.NewBytecode(abi.NEW_GLOBAL_INT, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{0}, []string{lvalue.Attribute()})
				this.relocatable.NewBytecode(abi.PUSH_INT, []int64{rvalue_value}, []string{})
				this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})
			}
		} else if rvalue.PrimaryExpr().PrimaryExprType() == expr.STRING {
		} else if rvalue.PrimaryExpr().PrimaryExprType() == expr.NULLPTR {
			err := errors.New("rvalue primary expr type is nullptr")
			panic(err)
		} else if rvalue.PrimaryExpr().PrimaryExprType() == expr.PAREN {
			err := errors.New("rvalue primary expr type is paren")
			panic(err)
		} else {
			err := errors.New("primary expr type is valid")
			panic(err)
		}
	} else {
		err := errors.New("rvalue expr type is primary")
		panic(err)
	}
}

func (this *Codegen) CodegenDecl(decl_ *decl.Decl) {
	if decl_.DeclType() == decl.STRUCT_DEF {
		this.CodegenStructDef(decl_.StructDef())
	} else if decl_.DeclType() == decl.FUNC_DECL {
		this.CodegenFuncDecl(decl_.FuncDecl())
	} else if decl_.DeclType() == decl.FUNC_DEF {
		this.CodegenFuncDef(decl_.FuncDef())
	} else {
		err := errors.New("decl type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenStructDef(struct_def *decl.StructDef) {
	this.relocatable.SwitchLabel("__bootstrap")

	struct_name := struct_def.Identifier().Attribute()

	this.relocatable.NewBytecode(abi.BEGIN_STRUCT, []int64{}, []string{struct_name})

	for i := 0; i < struct_def.Body().BlockStmt().Length(); i++ {
		stmt_ := struct_def.Body().BlockStmt().Get(i)

		type_specifier_ := stmt_.VarDeclStmt().TypeSpecifier()
		field_name := stmt_.VarDeclStmt().Identifier().Attribute()

		if type_specifier_.TypeSpecifierType() == type_specifier.VOID {
			if type_specifier_.NumStars() == 0 {
				err := errors.New("num stars == 0")
				panic(err)
			}

			this.relocatable.NewBytecode(
				abi.APPEND_VOID,
				[]int64{int64(type_specifier_.NumStars())},
				[]string{field_name},
			)
		} else if type_specifier_.TypeSpecifierType() == type_specifier.CHAR {
			this.relocatable.NewBytecode(abi.APPEND_CHAR, []int64{int64(type_specifier_.NumStars())}, []string{field_name})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.SHORT {
			this.relocatable.NewBytecode(abi.APPEND_SHORT, []int64{int64(type_specifier_.NumStars())}, []string{field_name})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.INT {
			this.relocatable.NewBytecode(abi.APPEND_INT, []int64{int64(type_specifier_.NumStars())}, []string{field_name})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.LONG {
			this.relocatable.NewBytecode(abi.APPEND_LONG, []int64{int64(type_specifier_.NumStars())}, []string{field_name})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.STRUCT {
			this.relocatable.NewBytecode(abi.APPEND_STRUCT, []int64{int64(type_specifier_.NumStars())}, []string{type_specifier_.StructIdentifier().Attribute(), field_name})
		} else {
			err := errors.New("type specifier type is not valid")
			panic(err)
		}
	}

	this.relocatable.NewBytecode(abi.END_STRUCT, []int64{}, []string{})
}

func (this *Codegen) CodegenFuncDecl(func_decl *decl.FuncDecl) {
	type_specifier_ := func_decl.TypeSpecifier()

	method := new(type_system.Method)

	if type_specifier_.TypeSpecifierType() == type_specifier.VOID {
		symbol := new(type_system.Symbol)
		symbol.InitPrimitive(
			type_system.VOID,
			type_specifier_.NumStars(),
			func_decl.Identifier().Attribute(),
		)

		method.Init(symbol)
	} else if type_specifier_.TypeSpecifierType() == type_specifier.CHAR {
		symbol := new(type_system.Symbol)
		symbol.InitPrimitive(type_system.CHAR, type_specifier_.NumStars(), func_decl.Identifier().Attribute())

		method.Init(symbol)
	} else if type_specifier_.TypeSpecifierType() == type_specifier.SHORT {
		symbol := new(type_system.Symbol)
		symbol.InitPrimitive(type_system.SHORT, type_specifier_.NumStars(), func_decl.Identifier().Attribute())

		method.Init(symbol)
	} else if type_specifier_.TypeSpecifierType() == type_specifier.INT {
		symbol := new(type_system.Symbol)
		symbol.InitPrimitive(type_system.INT, type_specifier_.NumStars(), func_decl.Identifier().Attribute())

		method.Init(symbol)
	} else if type_specifier_.TypeSpecifierType() == type_specifier.LONG {
		symbol := new(type_system.Symbol)
		symbol.InitPrimitive(type_system.LONG, type_specifier_.NumStars(), func_decl.Identifier().Attribute())

		method.Init(symbol)
	} else if type_specifier_.TypeSpecifierType() == type_specifier.STRUCT {
		symbol := new(type_system.Symbol)
		symbol.InitStruct(type_system.STRUCT, type_specifier_.StructIdentifier().Attribute(), type_specifier_.NumStars(), func_decl.Identifier().Attribute())

		method.Init(symbol)
	} else {
		err := errors.New("type specifier type is not valid")
		panic(err)
	}

	this.type_system.AddMethod(method)

	for i := 0; i < func_decl.ParamList().Length(); i++ {
		param := func_decl.ParamList().Get(i)

		if param.TypeSpecifier().TypeSpecifierType() == type_specifier.VOID {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(
				type_system.VOID,
				param.TypeSpecifier().NumStars(),
				param.Identifier().Attribute(),
			)

			method.AppendParam(symbol)
		} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.CHAR {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.CHAR, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

			method.AppendParam(symbol)
		} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.SHORT {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.SHORT, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

			method.AppendParam(symbol)
		} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.INT {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.INT, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

			method.AppendParam(symbol)
		} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.LONG {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.LONG, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

			method.AppendParam(symbol)
		} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.STRUCT {
			symbol := new(type_system.Symbol)
			symbol.InitStruct(type_system.STRUCT, param.TypeSpecifier().StructIdentifier().Attribute(), param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

			method.AppendParam(symbol)
		} else {
			err := errors.New("type specifier type is not valid")
			panic(err)
		}
	}

	this.relocatable.NewFunc(func_decl.Identifier().Attribute())
}

func (this *Codegen) CodegenFuncDef(func_def *decl.FuncDef) {
	type_specifier_ := func_def.TypeSpecifier()

	func_name := func_def.Identifier().Attribute()

	if !this.type_system.HasMethod(func_name) {
		method := new(type_system.Method)

		if type_specifier_.TypeSpecifierType() == type_specifier.VOID {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(
				type_system.VOID,
				type_specifier_.NumStars(),
				func_def.Identifier().Attribute(),
			)

			method.Init(symbol)
		} else if type_specifier_.TypeSpecifierType() == type_specifier.CHAR {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.CHAR, type_specifier_.NumStars(), func_def.Identifier().Attribute())

			method.Init(symbol)
		} else if type_specifier_.TypeSpecifierType() == type_specifier.SHORT {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.SHORT, type_specifier_.NumStars(), func_def.Identifier().Attribute())

			method.Init(symbol)
		} else if type_specifier_.TypeSpecifierType() == type_specifier.INT {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.INT, type_specifier_.NumStars(), func_def.Identifier().Attribute())

			method.Init(symbol)
		} else if type_specifier_.TypeSpecifierType() == type_specifier.LONG {
			symbol := new(type_system.Symbol)
			symbol.InitPrimitive(type_system.LONG, type_specifier_.NumStars(), func_def.Identifier().Attribute())

			method.Init(symbol)
		} else if type_specifier_.TypeSpecifierType() == type_specifier.STRUCT {
			symbol := new(type_system.Symbol)
			symbol.InitStruct(type_system.CHAR, type_specifier_.StructIdentifier().Attribute(), type_specifier_.NumStars(), func_def.Identifier().Attribute())

			method.Init(symbol)
		} else {
			err := errors.New("type specifier type is not valid")
			panic(err)
		}

		this.type_system.AddMethod(method)

		for i := 0; i < func_def.ParamList().Length(); i++ {
			param := func_def.ParamList().Get(i)

			if param.TypeSpecifier().TypeSpecifierType() == type_specifier.VOID {
				symbol := new(type_system.Symbol)
				symbol.InitPrimitive(
					type_system.VOID,
					param.TypeSpecifier().NumStars(),
					param.Identifier().Attribute(),
				)

				method.AppendParam(symbol)
			} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.CHAR {
				symbol := new(type_system.Symbol)
				symbol.InitPrimitive(type_system.CHAR, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

				method.AppendParam(symbol)
			} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.SHORT {
				symbol := new(type_system.Symbol)
				symbol.InitPrimitive(type_system.SHORT, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

				method.AppendParam(symbol)
			} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.INT {
				symbol := new(type_system.Symbol)
				symbol.InitPrimitive(type_system.INT, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

				method.AppendParam(symbol)
			} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.LONG {
				symbol := new(type_system.Symbol)
				symbol.InitPrimitive(type_system.LONG, param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

				method.AppendParam(symbol)
			} else if param.TypeSpecifier().TypeSpecifierType() == type_specifier.STRUCT {
				symbol := new(type_system.Symbol)
				symbol.InitStruct(type_system.STRUCT, param.TypeSpecifier().StructIdentifier().Attribute(), param.TypeSpecifier().NumStars(), param.Identifier().Attribute())

				method.AppendParam(symbol)
			} else {
				err := errors.New("type specifier type is not valid")
				panic(err)
			}
		}

		this.relocatable.NewFunc(func_def.Identifier().Attribute())
	}

	this.relocatable.SwitchFunc(func_name)
	this.relocatable.SwitchLabel(func_name)

	this.CodegenStmt(func_def.Body())
}

func (this *Codegen) CodegenStmt(stmt_ *stmt.Stmt) {
	if stmt_.StmtType() == stmt.EMPTY {
		this.CodegenEmptyStmt(stmt_.EmptyStmt())
	} else if stmt_.StmtType() == stmt.VAR_DECL {
		this.CodegenVarDeclStmt(stmt_.VarDeclStmt())
	} else if stmt_.StmtType() == stmt.VAR_DECL_INIT {
		this.CodegenVarDeclInitStmt(stmt_.VarDeclInitStmt())
	} else if stmt_.StmtType() == stmt.FOR {
		this.CodegenForStmt(stmt_.ForStmt())
	} else if stmt_.StmtType() == stmt.DPU_FOREACH {
		this.CodegenDpuForeachStmt(stmt_.DpuForeachStmt())
	} else if stmt_.StmtType() == stmt.WHILE {
		this.CodegenWhileStmt(stmt_.WhileStmt())
	} else if stmt_.StmtType() == stmt.CONTINUE {
		this.CodegenContinueStmt(stmt_.ContinueStmt())
	} else if stmt_.StmtType() == stmt.BREAK {
		this.CodegenBreakStmt(stmt_.BreakStmt())
	} else if stmt_.StmtType() == stmt.IF {
		this.CodegenIfStmt(stmt_.IfStmt())
	} else if stmt_.StmtType() == stmt.RETURN {
		this.CodegenReturnStmt(stmt_.ReturnStmt())
	} else if stmt_.StmtType() == stmt.EXPR {
		this.CodegenExprStmt(stmt_.ExprStmt())
	} else if stmt_.StmtType() == stmt.BLOCK {
		this.CodegenBlockStmt(stmt_.BlockStmt())
	} else {
		err := errors.New("stmt type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenEmptyStmt(empty_stmt *stmt.EmptyStmt) {
	this.relocatable.NewBytecode(abi.NOP, []int64{}, []string{})
}

func (this *Codegen) CodegenVarDeclStmt(var_decl_stmt *stmt.VarDeclStmt) {
	type_specifier_ := var_decl_stmt.TypeSpecifier()
	symbol_name := var_decl_stmt.Identifier().Attribute()

	if type_specifier_.TypeSpecifierType() == type_specifier.VOID {
		if type_specifier_.NumStars() == 0 {
			err := errors.New("num stars == 0")
			panic(err)
		}

		this.relocatable.NewBytecode(
			abi.NEW_FAST_VOID,
			[]int64{int64(type_specifier_.NumStars())},
			[]string{symbol_name},
		)
	} else if type_specifier_.TypeSpecifierType() == type_specifier.CHAR {
		this.relocatable.NewBytecode(abi.NEW_FAST_CHAR, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.SHORT {
		this.relocatable.NewBytecode(abi.NEW_FAST_SHORT, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.INT {
		this.relocatable.NewBytecode(abi.NEW_FAST_INT, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.LONG {
		this.relocatable.NewBytecode(abi.NEW_FAST_LONG, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.STRUCT {
		this.relocatable.NewBytecode(abi.NEW_FAST_STRUCT, []int64{int64(type_specifier_.NumStars())}, []string{type_specifier_.StructIdentifier().Attribute(), symbol_name})
	} else {
		err := errors.New("type specifier type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenVarDeclInitStmt(var_decl_init_stmt *stmt.VarDeclInitStmt) {
	type_specifier_ := var_decl_init_stmt.TypeSpecifier()
	symbol_name := var_decl_init_stmt.Identifier().Attribute()

	if type_specifier_.TypeSpecifierType() == type_specifier.VOID {
		if type_specifier_.NumStars() == 0 {
			err := errors.New("num stars == 0")
			panic(err)
		}

		this.relocatable.NewBytecode(
			abi.NEW_FAST_VOID,
			[]int64{int64(type_specifier_.NumStars())},
			[]string{symbol_name},
		)
	} else if type_specifier_.TypeSpecifierType() == type_specifier.CHAR {
		this.relocatable.NewBytecode(abi.NEW_FAST_CHAR, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.SHORT {
		this.relocatable.NewBytecode(abi.NEW_FAST_SHORT, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.INT {
		this.relocatable.NewBytecode(abi.NEW_FAST_INT, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.LONG {
		this.relocatable.NewBytecode(abi.NEW_FAST_LONG, []int64{int64(type_specifier_.NumStars())}, []string{symbol_name})
	} else if type_specifier_.TypeSpecifierType() == type_specifier.STRUCT {
		this.relocatable.NewBytecode(abi.NEW_FAST_STRUCT, []int64{int64(type_specifier_.NumStars())}, []string{type_specifier_.StructIdentifier().Attribute(), symbol_name})
	} else {
		err := errors.New("type specifier type is not valid")
		panic(err)
	}

	this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{}, []string{symbol_name})
	this.CodegenExpr(var_decl_init_stmt.Expr())
	this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})
}

func (this *Codegen) CodegenForStmt(for_stmt *stmt.ForStmt) {
	condition_label, body_label, end_label := this.relocatable.NewLoop()

	this.relocatable.NewBytecode(abi.NEW_SCOPE, []int64{}, []string{})
	this.CodegenStmt(for_stmt.Initialization())
	this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{condition_label.Name()})

	this.relocatable.SwitchLabel(condition_label.Name())
	this.CodegenExpr(for_stmt.Condition())
	this.relocatable.NewBytecode(abi.JUMP_IF_NONZERO, []int64{}, []string{body_label.Name()})
	this.CodegenExpr(for_stmt.Condition())
	this.relocatable.NewBytecode(abi.JUMP_IF_ZERO, []int64{}, []string{end_label.Name()})

	this.relocatable.SwitchLabel(body_label.Name())
	this.CodegenStmt(for_stmt.Body())
	this.CodegenStmt(for_stmt.Update())
	this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{condition_label.Name()})

	this.relocatable.SwitchLabel(end_label.Name())
	this.relocatable.NewBytecode(abi.DELETE_SCOPE, []int64{}, []string{})
}

func (this *Codegen) CodegenDpuForeachStmt(dpu_foreach_stmt *stmt.DpuForeachStmt) {
	if dpu_foreach_stmt.Foreach().Get(1).ExprType() != expr.PRIMARY {
		err := errors.New("second argument's expr type is not primary")
		panic(err)
	} else if dpu_foreach_stmt.Foreach().Get(1).PrimaryExpr().PrimaryExprType() != expr.IDENTIFIER {
		err := errors.New("second argument's primary expr type is not identifier")
		panic(err)
	}

	dpu_symbol_name := dpu_foreach_stmt.Foreach().Get(1).PrimaryExpr().Token().Attribute()

	if dpu_foreach_stmt.Foreach().Get(2).ExprType() != expr.PRIMARY {
		err := errors.New("third argument's expr type is not primary")
		panic(err)
	} else if dpu_foreach_stmt.Foreach().Get(2).PrimaryExpr().PrimaryExprType() != expr.IDENTIFIER {
		err := errors.New("third argument's primary expr type is not identifier")
		panic(err)
	}

	i_symbol_name := dpu_foreach_stmt.Foreach().Get(2).PrimaryExpr().Token().Attribute()

	this.relocatable.NewBytecode(abi.NEW_SCOPE, []int64{}, []string{})
	this.relocatable.NewBytecode(abi.NEW_FAST_INT, []int64{0}, []string{dpu_symbol_name})
	this.relocatable.NewBytecode(abi.NEW_FAST_INT, []int64{0}, []string{i_symbol_name})

	for i := 0; i < this.num_dpus; i++ {
		this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{}, []string{dpu_symbol_name})
		this.relocatable.NewBytecode(abi.PUSH_INT, []int64{int64(i)}, []string{})
		this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})

		this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{}, []string{i_symbol_name})
		this.relocatable.NewBytecode(abi.PUSH_INT, []int64{int64(i)}, []string{})
		this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})

		this.CodegenStmt(dpu_foreach_stmt.Body())
	}

	this.relocatable.NewBytecode(abi.DELETE_SCOPE, []int64{}, []string{})
}

func (this *Codegen) CodegenWhileStmt(while_stmt *stmt.WhileStmt) {
	condition_label, body_label, end_label := this.relocatable.NewLoop()

	this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{condition_label.Name()})

	this.relocatable.SwitchLabel(condition_label.Name())
	this.CodegenExpr(while_stmt.Condition())
	this.relocatable.NewBytecode(abi.JUMP_IF_NONZERO, []int64{}, []string{body_label.Name()})
	this.CodegenExpr(while_stmt.Condition())
	this.relocatable.NewBytecode(abi.JUMP_IF_ZERO, []int64{}, []string{end_label.Name()})

	this.relocatable.SwitchLabel(body_label.Name())
	this.CodegenStmt(while_stmt.Body())
	this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{condition_label.Name()})

	this.relocatable.SwitchLabel(end_label.Name())
}

func (this *Codegen) CodegenContinueStmt(continue_stmt *stmt.ContinueStmt) {
	this.relocatable.NewBytecode(
		abi.JUMP,
		[]int64{},
		[]string{this.relocatable.CurLoopCondition().Name()},
	)
}

func (this *Codegen) CodegenBreakStmt(break_stmt *stmt.BreakStmt) {
	this.relocatable.NewBytecode(
		abi.JUMP,
		[]int64{},
		[]string{this.relocatable.CurLoopEnd().Name()},
	)
}

func (this *Codegen) CodegenIfStmt(if_stmt *stmt.IfStmt) {
	var num_branches int
	if if_stmt.HasElseBody() {
		num_branches = 2 + if_stmt.NumElseIfs()
	} else {
		num_branches = 1 + if_stmt.NumElseIfs()
	}

	branches := make([]*abi.Label, 0)
	for i := 0; i < num_branches; i++ {
		branch := this.relocatable.NewUnnamedLabel()
		branches = append(branches, branch)
	}

	branch_end := this.relocatable.NewUnnamedLabel()

	this.CodegenExpr(if_stmt.IfCondition())
	this.relocatable.NewBytecode(abi.JUMP_IF_NONZERO, []int64{}, []string{branches[0].Name()})

	for i := 0; i < if_stmt.NumElseIfs(); i++ {
		this.CodegenExpr(if_stmt.ElseIfCondition(i))
		this.relocatable.NewBytecode(abi.JUMP_IF_NONZERO, []int64{}, []string{branches[i+1].Name()})
	}

	if if_stmt.HasElseBody() {
		this.relocatable.NewBytecode(
			abi.JUMP,
			[]int64{},
			[]string{branches[len(branches)-1].Name()},
		)
	} else {
		this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{branch_end.Name()})
	}

	this.relocatable.SwitchLabel(branches[0].Name())
	this.CodegenStmt(if_stmt.IfBody())
	this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{branch_end.Name()})

	for i := 0; i < if_stmt.NumElseIfs(); i++ {
		this.relocatable.SwitchLabel(branches[i+1].Name())
		this.CodegenStmt(if_stmt.ElseIfBody(i))
		this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{branch_end.Name()})
	}

	if if_stmt.HasElseBody() {
		this.relocatable.SwitchLabel(branches[len(branches)-1].Name())
		this.CodegenStmt(if_stmt.ElseBody())
		this.relocatable.NewBytecode(abi.JUMP, []int64{}, []string{branch_end.Name()})
	}

	this.relocatable.SwitchLabel(branch_end.Name())
}

func (this *Codegen) CodegenReturnStmt(return_stmt *stmt.ReturnStmt) {
	if return_stmt.HasValue() {
		func_name := this.relocatable.CurFunc().Name()
		method := this.type_system.Method(func_name)

		if method.Symbol().SymbolType() == type_system.VOID {
			if method.Symbol().NumStars() == 0 {
				err := errors.New("num stars == 0")
				panic(err)
			}

			this.relocatable.NewBytecode(
				abi.NEW_RETURN_VOID,
				[]int64{int64(method.Symbol().NumStars())},
				[]string{},
			)
		} else if method.Symbol().SymbolType() == type_system.CHAR {
			this.relocatable.NewBytecode(abi.NEW_RETURN_CHAR, []int64{int64(method.Symbol().NumStars())}, []string{})
		} else if method.Symbol().SymbolType() == type_system.SHORT {
			this.relocatable.NewBytecode(abi.NEW_RETURN_SHORT, []int64{int64(method.Symbol().NumStars())}, []string{})
		} else if method.Symbol().SymbolType() == type_system.INT {
			this.relocatable.NewBytecode(abi.NEW_RETURN_INT, []int64{int64(method.Symbol().NumStars())}, []string{})
		} else if method.Symbol().SymbolType() == type_system.LONG {
			this.relocatable.NewBytecode(abi.NEW_RETURN_LONG, []int64{int64(method.Symbol().NumStars())}, []string{})
		} else if method.Symbol().SymbolType() == type_system.STRUCT {
			this.relocatable.NewBytecode(abi.NEW_RETURN_STRUCT, []int64{int64(method.Symbol().NumStars())}, []string{method.Symbol().StructName()})
		} else {
			err := errors.New("symbol type is not valid")
			panic(err)
		}

		this.CodegenExpr(return_stmt.Value())
		this.relocatable.NewBytecode(abi.ASSIGN_RETURN, []int64{}, []string{})
	}

	this.relocatable.NewBytecode(abi.RETURN, []int64{}, []string{})
}

func (this *Codegen) CodegenExprStmt(expr_stmt *stmt.ExprStmt) {
	this.CodegenExpr(expr_stmt.Expr())
}

func (this *Codegen) CodegenBlockStmt(block_stmt *stmt.BlockStmt) {
	this.relocatable.NewBytecode(abi.NEW_SCOPE, []int64{}, []string{})

	for i := 0; i < block_stmt.Length(); i++ {
		stmt_ := block_stmt.Get(i)
		this.CodegenStmt(stmt_)
	}

	this.relocatable.NewBytecode(abi.DELETE_SCOPE, []int64{}, []string{})
}

func (this *Codegen) CodegenExpr(expr_ *expr.Expr) {
	if expr_.ExprType() == expr.PRIMARY {
		this.CodegenPrimaryExpr(expr_.PrimaryExpr())
	} else if expr_.ExprType() == expr.POSTFIX {
		this.CodegenPostfixExpr(expr_.PostfixExpr())
	} else if expr_.ExprType() == expr.UNARY {
		this.CodegenUnaryExpr(expr_.UnaryExpr())
	} else if expr_.ExprType() == expr.MULTIPLICATIVE {
		this.CodegenMultiplicativeExpr(expr_.MultiplicativeExpr())
	} else if expr_.ExprType() == expr.ADDITIVE {
		this.CodegenAdditiveExpr(expr_.AdditiveExpr())
	} else if expr_.ExprType() == expr.SHIFT {
		this.CodegenShiftExpr(expr_.ShiftExpr())
	} else if expr_.ExprType() == expr.RELATIONAL {
		this.CodegenRelationalExpr(expr_.RelationalExpr())
	} else if expr_.ExprType() == expr.EQUALITY {
		this.CodegenEqualityExpr(expr_.EqualityExpr())
	} else if expr_.ExprType() == expr.BITWISE_AND {
		this.CodegenBitwiseAndExpr(expr_.BitwiseAndExpr())
	} else if expr_.ExprType() == expr.BITWISE_XOR {
		this.CodegenBitwiseXorExpr(expr_.BitwiseXorExpr())
	} else if expr_.ExprType() == expr.BITWISE_OR {
		this.CodegenBitwiseOrExpr(expr_.BitwiseOrExpr())
	} else if expr_.ExprType() == expr.LOGICAL_AND {
		this.CodegenLogicalAndExpr(expr_.LogicalAndExpr())
	} else if expr_.ExprType() == expr.LOGICAL_OR {
		this.CodegenLogicalOrExpr(expr_.LogicalOrExpr())
	} else if expr_.ExprType() == expr.CONDITIONAL {
		this.CodegenConditionalExpr(expr_.ConditionalExpr())
	} else if expr_.ExprType() == expr.ASSIGNMENT {
		this.CodegenAssignmentExpr(expr_.AssignmentExpr())
	} else {
		err := errors.New("expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenPrimaryExpr(primary_expr *expr.PrimaryExpr) {
	if primary_expr.PrimaryExprType() == expr.IDENTIFIER {
		symbol_name := primary_expr.Token().Attribute()

		if symbol_name == "DPU_XFER_TO_DPU" {
			this.relocatable.NewBytecode(abi.PUSH_INT, []int64{0}, []string{})
		} else if symbol_name == "DPU_XFER_FROM_DPU" {
			this.relocatable.NewBytecode(abi.PUSH_INT, []int64{1}, []string{})
		} else if symbol_name == "DPU_MRAM_HEAP_POINTER_NAME" {
			this.relocatable.NewBytecode(abi.PUSH_INT, []int64{this.dpu_mram_heap_pointer_name}, []string{})
		} else if symbol_name == "DPU_XFER_DEFAULT" {
			this.relocatable.NewBytecode(abi.PUSH_INT, []int64{0}, []string{})
		} else if symbol_name == "DPU_SYNCHRONOUS" {
			this.relocatable.NewBytecode(abi.PUSH_INT, []int64{0}, []string{})
		} else {
			this.relocatable.NewBytecode(abi.GET_IDENTIFIER, []int64{}, []string{symbol_name})
		}
	} else if primary_expr.PrimaryExprType() == expr.NUMBER {
		number, err := strconv.ParseInt(primary_expr.Token().Attribute(), 10, 64)
		if err != nil {
			panic(err)
		}

		this.relocatable.NewBytecode(abi.PUSH_INT, []int64{number}, []string{})
	} else if primary_expr.PrimaryExprType() == expr.STRING {
		str := primary_expr.Token().Attribute()

		this.relocatable.NewBytecode(abi.PUSH_STRING, []int64{}, []string{str})
	} else if primary_expr.PrimaryExprType() == expr.NULLPTR {
		this.relocatable.NewBytecode(abi.PUSH_INT, []int64{0}, []string{})
	} else if primary_expr.PrimaryExprType() == expr.PAREN {
		this.CodegenExpr(primary_expr.Expr())
	} else {
		err := errors.New("primary expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenPostfixExpr(postfix_expr *expr.PostfixExpr) {
	if postfix_expr.PostfixExprType() == expr.BRACKET {
		this.CodegenExpr(postfix_expr.Base())
		this.CodegenExpr(postfix_expr.OffsetExpr())

		this.relocatable.NewBytecode(abi.GET_SUBSCRIPT, []int64{}, []string{})
	} else if postfix_expr.PostfixExprType() == expr.CALL {
		if postfix_expr.Base().ExprType() != expr.PRIMARY {
			err := errors.New("base expr type is not primary")
			panic(err)
		} else if postfix_expr.Base().PrimaryExpr().PrimaryExprType() != expr.IDENTIFIER {
			err := errors.New("base primary expr type is not identifier")
			panic(err)
		}

		func_name := postfix_expr.Base().PrimaryExpr().Token().Attribute()

		if func_name == "malloc" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.ALLOC, []int64{}, []string{})
		} else if func_name == "free" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.FREE, []int64{}, []string{})
		} else if func_name == "assert" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.ASSERT, []int64{}, []string{})
		} else if func_name == "sqrt" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.SQRT, []int64{}, []string{})
		} else if func_name == "dpu_alloc" {
			for i := 0; i < this.num_dpus; i++ {
				this.relocatable.NewBytecode(abi.DPU_ALLOC, []int64{int64(i)}, []string{})
			}
		} else if func_name == "dpu_load" {
			for i := 0; i < this.num_dpus; i++ {
				this.relocatable.NewBytecode(abi.DPU_LOAD, []int64{int64(i)}, []string{this.benchmark})
			}
		} else if func_name == "dpu_prepare_xfer" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.DPU_PREPARE, []int64{}, []string{})
		} else if func_name == "dpu_push_xfer" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.DPU_TRANSFER, []int64{}, []string{})
		} else if func_name == "dpu_copy_to" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.DPU_COPY_TO, []int64{}, []string{})
		} else if func_name == "dpu_copy_from" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.DPU_COPY_FROM, []int64{}, []string{})
		} else if func_name == "dpu_launch" {
			for i := 0; i < postfix_expr.ArgList().Length(); i++ {
				this.CodegenExpr(postfix_expr.ArgList().Get(i))
			}

			this.relocatable.NewBytecode(abi.DPU_LAUNCH, []int64{}, []string{})
		} else if func_name == "dpu_free" {
			this.relocatable.NewBytecode(abi.DPU_FREE, []int64{}, []string{})
		} else {
			method := this.type_system.Method(func_name)
			params := method.Params()

			for i, param := range params {
				arg := postfix_expr.ArgList().Get(i)

				if param.SymbolType() == type_system.VOID {
					this.relocatable.NewBytecode(abi.NEW_ARG_VOID, []int64{int64(param.NumStars())}, []string{param.Name()})
				} else if param.SymbolType() == type_system.CHAR {
					this.relocatable.NewBytecode(abi.NEW_ARG_CHAR, []int64{int64(param.NumStars())}, []string{param.Name()})
				} else if param.SymbolType() == type_system.SHORT {
					this.relocatable.NewBytecode(abi.NEW_ARG_SHORT, []int64{int64(param.NumStars())}, []string{param.Name()})
				} else if param.SymbolType() == type_system.INT {
					this.relocatable.NewBytecode(abi.NEW_ARG_INT, []int64{int64(param.NumStars())}, []string{param.Name()})
				} else if param.SymbolType() == type_system.LONG {
					this.relocatable.NewBytecode(abi.NEW_ARG_LONG, []int64{int64(param.NumStars())}, []string{param.Name()})
				} else if param.SymbolType() == type_system.STRUCT {
					this.relocatable.NewBytecode(abi.NEW_ARG_STRUCT, []int64{int64(param.NumStars())}, []string{param.StructName(), param.Name()})
				} else {
					err := errors.New("symbol type is not valid")
					panic(err)
				}

				this.relocatable.NewBytecode(abi.GET_ARG_IDENTIFIER, []int64{0}, []string{param.Name()})
				this.CodegenExpr(arg)
				this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})
			}

			this.relocatable.NewBytecode(abi.CALL, []int64{}, []string{func_name})
		}
	} else if postfix_expr.PostfixExprType() == expr.DOT {
		field_name := postfix_expr.OffsetToken().Attribute()

		this.CodegenExpr(postfix_expr.Base())

		this.relocatable.NewBytecode(abi.GET_ACCESS, []int64{}, []string{field_name})
	} else if postfix_expr.PostfixExprType() == expr.ARROW {
		field_name := postfix_expr.OffsetToken().Attribute()

		this.CodegenExpr(postfix_expr.Base())

		this.relocatable.NewBytecode(abi.GET_REFERENCE, []int64{}, []string{field_name})
	} else if postfix_expr.PostfixExprType() == expr.POSTFIX_PLUS_PLUS {
		this.CodegenExpr(postfix_expr.Base())
		this.relocatable.NewBytecode(abi.ASSIGN_PLUS_PLUS, []int64{}, []string{})
	} else if postfix_expr.PostfixExprType() == expr.POSTFIX_MINUS_MINUS {
		this.CodegenExpr(postfix_expr.Base())
		this.relocatable.NewBytecode(abi.ASSIGN_MINUS_MINUS, []int64{}, []string{})
	} else {
		err := errors.New("postfix expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenUnaryExpr(unary_expr *expr.UnaryExpr) {
	if unary_expr.UnaryExprType() == expr.UNARY_PLUS_PLUS {
		this.CodegenExpr(unary_expr.Base())

		this.relocatable.NewBytecode(abi.ASSIGN_PLUS_PLUS, []int64{}, []string{})
	} else if unary_expr.UnaryExprType() == expr.UNARY_MINUS_MINUS {
		this.CodegenExpr(unary_expr.Base())

		this.relocatable.NewBytecode(abi.ASSIGN_MINUS_MINUS, []int64{}, []string{})
	} else if unary_expr.UnaryExprType() == expr.AND {
		this.CodegenExpr(unary_expr.Base())

		this.relocatable.NewBytecode(abi.GET_ADDRESS, []int64{}, []string{})
	} else if unary_expr.UnaryExprType() == expr.STAR {
		this.CodegenExpr(unary_expr.Base())

		this.relocatable.NewBytecode(abi.GET_VALUE, []int64{}, []string{})
	} else if unary_expr.UnaryExprType() == expr.PLUS {
		this.CodegenExpr(unary_expr.Base())
	} else if unary_expr.UnaryExprType() == expr.MINUS {
		this.CodegenExpr(unary_expr.Base())

		this.relocatable.NewBytecode(abi.NEGATE, []int64{}, []string{})
	} else if unary_expr.UnaryExprType() == expr.TILDE {
		this.CodegenExpr(unary_expr.Base())

		this.relocatable.NewBytecode(abi.TILDE, []int64{}, []string{})
	} else if unary_expr.UnaryExprType() == expr.NOT {
		this.CodegenExpr(unary_expr.Base())

		this.relocatable.NewBytecode(abi.LOGICAL_NOT, []int64{}, []string{})
	} else if unary_expr.UnaryExprType() == expr.SIZEOF {
		type_specifier_ := unary_expr.TypeSpecifier()

		if type_specifier_.TypeSpecifierType() == type_specifier.VOID {
			if type_specifier_.NumStars() == 0 {
				err := errors.New("num stars == 0")
				panic(err)
			}

			this.relocatable.NewBytecode(abi.SIZEOF_VOID, []int64{int64(type_specifier_.NumStars())}, []string{})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.CHAR {
			this.relocatable.NewBytecode(abi.SIZEOF_CHAR, []int64{int64(type_specifier_.NumStars())}, []string{})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.SHORT {
			this.relocatable.NewBytecode(abi.SIZEOF_SHORT, []int64{int64(type_specifier_.NumStars())}, []string{})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.INT {
			this.relocatable.NewBytecode(abi.SIZEOF_INT, []int64{int64(type_specifier_.NumStars())}, []string{})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.LONG {
			this.relocatable.NewBytecode(abi.SIZEOF_LONG, []int64{int64(type_specifier_.NumStars())}, []string{})
		} else if type_specifier_.TypeSpecifierType() == type_specifier.STRUCT {
			this.relocatable.NewBytecode(abi.SIZEOF_STRUCT, []int64{int64(type_specifier_.NumStars())}, []string{type_specifier_.StructIdentifier().Attribute()})
		} else {
			err := errors.New("type specifier type is not valid")
			panic(err)
		}
	} else {
		err := errors.New("unary expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenMultiplicativeExpr(multiplicative_expr *expr.MultiplicativeExpr) {
	this.CodegenExpr(multiplicative_expr.Loperand())
	this.CodegenExpr(multiplicative_expr.Roperand())

	if multiplicative_expr.MultiplicativeExprType() == expr.MUL {
		this.relocatable.NewBytecode(abi.MUL, []int64{}, []string{})
	} else if multiplicative_expr.MultiplicativeExprType() == expr.DIV {
		this.relocatable.NewBytecode(abi.DIV, []int64{}, []string{})
	} else if multiplicative_expr.MultiplicativeExprType() == expr.MOD {
		this.relocatable.NewBytecode(abi.MOD, []int64{}, []string{})
	} else {
		err := errors.New("multiplicative expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenAdditiveExpr(additive_expr *expr.AdditiveExpr) {
	this.CodegenExpr(additive_expr.Loperand())
	this.CodegenExpr(additive_expr.Roperand())

	if additive_expr.AdditiveExprType() == expr.ADD {
		this.relocatable.NewBytecode(abi.ADD, []int64{}, []string{})
	} else if additive_expr.AdditiveExprType() == expr.SUB {
		this.relocatable.NewBytecode(abi.SUB, []int64{}, []string{})
	} else {
		err := errors.New("additive expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenShiftExpr(shift_expr *expr.ShiftExpr) {
	this.CodegenExpr(shift_expr.Loperand())
	this.CodegenExpr(shift_expr.Roperand())

	if shift_expr.ShiftExprType() == expr.LSHIFT {
		this.relocatable.NewBytecode(abi.LSHIFT, []int64{}, []string{})
	} else if shift_expr.ShiftExprType() == expr.RSHIFT {
		this.relocatable.NewBytecode(abi.RSHIFT, []int64{}, []string{})
	} else {
		err := errors.New("shift expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenRelationalExpr(relational_expr *expr.RelationalExpr) {
	this.CodegenExpr(relational_expr.Loperand())
	this.CodegenExpr(relational_expr.Roperand())

	if relational_expr.RelationalExprType() == expr.LESS {
		this.relocatable.NewBytecode(abi.LESS, []int64{}, []string{})
	} else if relational_expr.RelationalExprType() == expr.LESS_EQ {
		this.relocatable.NewBytecode(abi.LESS_EQ, []int64{}, []string{})
	} else if relational_expr.RelationalExprType() == expr.GREATER {
		this.relocatable.NewBytecode(abi.GREATER, []int64{}, []string{})
	} else if relational_expr.RelationalExprType() == expr.GREATER_EQ {
		this.relocatable.NewBytecode(abi.GREATER_EQ, []int64{}, []string{})
	} else {
		err := errors.New("relational expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenEqualityExpr(equality_expr *expr.EqualityExpr) {
	this.CodegenExpr(equality_expr.Loperand())
	this.CodegenExpr(equality_expr.Roperand())

	if equality_expr.EqualityExprType() == expr.EQ {
		this.relocatable.NewBytecode(abi.EQ, []int64{}, []string{})
	} else if equality_expr.EqualityExprType() == expr.NOT_EQ {
		this.relocatable.NewBytecode(abi.NOT_EQ, []int64{}, []string{})
	} else {
		err := errors.New("equality expr type is not valid")
		panic(err)
	}
}

func (this *Codegen) CodegenBitwiseAndExpr(bitwise_and_expr *expr.BitwiseAndExpr) {
	this.CodegenExpr(bitwise_and_expr.Loperand())
	this.CodegenExpr(bitwise_and_expr.Roperand())

	this.relocatable.NewBytecode(abi.BITWISE_AND, []int64{}, []string{})
}

func (this *Codegen) CodegenBitwiseXorExpr(bitwise_xor_expr *expr.BitwiseXorExpr) {
	this.CodegenExpr(bitwise_xor_expr.Loperand())
	this.CodegenExpr(bitwise_xor_expr.Roperand())

	this.relocatable.NewBytecode(abi.BITWISE_XOR, []int64{}, []string{})
}

func (this *Codegen) CodegenBitwiseOrExpr(bitwise_or_expr *expr.BitwiseOrExpr) {
	this.CodegenExpr(bitwise_or_expr.Loperand())
	this.CodegenExpr(bitwise_or_expr.Roperand())

	this.relocatable.NewBytecode(abi.BITWISE_OR, []int64{}, []string{})
}

func (this *Codegen) CodegenLogicalAndExpr(logical_and_expr *expr.LogicalAndExpr) {
	this.CodegenExpr(logical_and_expr.Loperand())
	this.CodegenExpr(logical_and_expr.Roperand())

	this.relocatable.NewBytecode(abi.LOGICAL_AND, []int64{}, []string{})
}

func (this *Codegen) CodegenLogicalOrExpr(logical_or_expr *expr.LogicalOrExpr) {
	this.CodegenExpr(logical_or_expr.Loperand())
	this.CodegenExpr(logical_or_expr.Roperand())

	this.relocatable.NewBytecode(abi.LOGICAL_OR, []int64{}, []string{})
}

func (this *Codegen) CodegenConditionalExpr(conditional_expr *expr.ConditionalExpr) {
	this.CodegenExpr(conditional_expr.ConditionalExpr())
	this.CodegenExpr(conditional_expr.TrueExpr())
	this.CodegenExpr(conditional_expr.FalseExpr())

	this.relocatable.NewBytecode(abi.CONDITIONAL, []int64{}, []string{})
}

func (this *Codegen) CodegenAssignmentExpr(assignment_expr *expr.AssignmentExpr) {
	this.CodegenExpr(assignment_expr.Lvalue())
	this.CodegenExpr(assignment_expr.Rvalue())

	if assignment_expr.AssignmentExprType() == expr.ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.STAR_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_STAR, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.DIV_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_DIV, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.MOD_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_MOD, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.PLUS_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_ADD, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.MINUS_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_SUB, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.LSHIFT_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_LSHIFT, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.RSHIFT_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_RSHIFT, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.AND_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_BITWISE_AND, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.CARET_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_BITWISE_XOR, []int64{}, []string{})
	} else if assignment_expr.AssignmentExprType() == expr.OR_ASSIGN {
		this.relocatable.NewBytecode(abi.ASSIGN_BITWISE_OR, []int64{}, []string{})
	} else {
		err := errors.New("assignment expr type is not valid")
		panic(err)
	}
}
