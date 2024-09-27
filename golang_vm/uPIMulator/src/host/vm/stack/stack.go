package stack

type Stack struct {
	stack_items []*StackItem
}

func (this *Stack) Init() {
	this.stack_items = make([]*StackItem, 0)
}

func (this *Stack) Front(pos int) *StackItem {
	return this.stack_items[len(this.stack_items)-1-pos]
}

func (this *Stack) Push(stack_item *StackItem) {
	this.stack_items = append(this.stack_items, stack_item)
}

func (this *Stack) Pop() {
	this.stack_items = this.stack_items[:len(this.stack_items)-1]
}

func (this *Stack) HasObject(address int64) bool {
	for _, stack_item := range this.stack_items {
		if stack_item.Address() == address {
			return true
		}
	}
	return false
}

func (this *Stack) Length() int {
	return len(this.stack_items)
}
