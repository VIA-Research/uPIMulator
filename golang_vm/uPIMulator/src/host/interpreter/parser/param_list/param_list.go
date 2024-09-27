package param_list

type ParamList struct {
	params []*Param
}

func (this *ParamList) Init() {
	this.params = make([]*Param, 0)
}

func (this *ParamList) Length() int {
	return len(this.params)
}

func (this *ParamList) Get(pos int) *Param {
	return this.params[pos]
}

func (this *ParamList) Append(param *Param) {
	this.params = append(this.params, param)
}
