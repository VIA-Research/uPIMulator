package parser

type Ast struct {
	stack_items []*StackItem
}

func (this *Ast) Init(stack_items []*StackItem) {
	this.stack_items = stack_items
}

func (this *Ast) Length() int {
	return len(this.stack_items)
}

func (this *Ast) Get(pos int) *StackItem {
	return this.stack_items[pos]
}
