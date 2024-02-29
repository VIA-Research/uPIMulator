package logic

type LinkerConstant struct {
	name  string
	value int64
}

func (this *LinkerConstant) Init(name string) {
	this.name = name
	this.value = 0
}

func (this *LinkerConstant) Name() string {
	return this.name
}

func (this *LinkerConstant) Value() int64 {
	return this.value
}

func (this *LinkerConstant) SetValue(value int64) {
	this.value = value
}
