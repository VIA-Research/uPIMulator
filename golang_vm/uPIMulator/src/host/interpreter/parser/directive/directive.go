package directive

type DirectiveType int

const (
	INCLUDE DirectiveType = iota
	DEFINE
)

type Directive struct {
	directive_type DirectiveType

	include_directive *IncludeDirective
	define_directive  *DefineDirective
}

func (this *Directive) InitIncludeDirective(include_directive *IncludeDirective) {
	this.directive_type = INCLUDE

	this.include_directive = include_directive
}

func (this *Directive) InitDefineDirective(define_directive *DefineDirective) {
	this.directive_type = DEFINE

	this.define_directive = define_directive
}

func (this *Directive) DirectiveType() DirectiveType {
	return this.directive_type
}

func (this *Directive) IncludeDirective() *IncludeDirective {
	return this.include_directive
}

func (this *Directive) DefineDirective() *DefineDirective {
	return this.define_directive
}
