package expr

type ArgList struct {
	args []*Expr
}

func (this *ArgList) Init() {
	this.args = make([]*Expr, 0)
}

func (this *ArgList) Length() int {
	return len(this.args)
}

func (this *ArgList) Get(pos int) *Expr {
	return this.args[pos]
}

func (this *ArgList) Append(arg *Expr) {
	this.args = append(this.args, arg)
}
