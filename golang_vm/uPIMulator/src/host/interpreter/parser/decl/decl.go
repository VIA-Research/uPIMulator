package decl

type DeclType int

const (
	STRUCT_DEF DeclType = iota
	FUNC_DECL
	FUNC_DEF
)

type Decl struct {
	decl_type DeclType

	struct_def *StructDef
	func_decl  *FuncDecl
	func_def   *FuncDef
}

func (this *Decl) InitStructDef(struct_def *StructDef) {
	this.decl_type = STRUCT_DEF

	this.struct_def = struct_def
}

func (this *Decl) InitFuncDecl(func_decl *FuncDecl) {
	this.decl_type = FUNC_DECL

	this.func_decl = func_decl
}

func (this *Decl) InitFuncDef(func_def *FuncDef) {
	this.decl_type = FUNC_DEF

	this.func_def = func_def
}

func (this *Decl) DeclType() DeclType {
	return this.decl_type
}

func (this *Decl) StructDef() *StructDef {
	return this.struct_def
}

func (this *Decl) FuncDecl() *FuncDecl {
	return this.func_decl
}

func (this *Decl) FuncDef() *FuncDef {
	return this.func_def
}
