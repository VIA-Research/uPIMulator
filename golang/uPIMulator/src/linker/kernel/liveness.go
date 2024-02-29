package kernel

type Liveness struct {
	defs           map[string]bool
	uses           map[string]bool
	global_symbols map[string]bool
}

func (this *Liveness) Init() {
	this.defs = make(map[string]bool, 0)
	this.uses = make(map[string]bool, 0)
	this.global_symbols = make(map[string]bool, 0)
}

func (this *Liveness) Defs() map[string]bool {
	return this.defs
}

func (this *Liveness) AddDef(def string) {
	this.defs[def] = true
}

func (this *Liveness) Uses() map[string]bool {
	return this.uses
}

func (this *Liveness) AddUse(use string) {
	this.uses[use] = true
}

func (this *Liveness) GlobalSymbols() map[string]bool {
	return this.global_symbols
}

func (this *Liveness) AddGlobalSymbol(global_symbol string) {
	this.global_symbols[global_symbol] = true
}

func (this *Liveness) LocalSymbols() map[string]bool {
	local_symbols := make(map[string]bool, 0)
	for def, _ := range this.defs {
		if _, found := this.global_symbols[def]; !found {
			local_symbols[def] = true
		}
	}
	return local_symbols
}

func (this *Liveness) UnresolvedSymbols() map[string]bool {
	unresolved_symbols := make(map[string]bool, 0)
	for use, _ := range this.uses {
		if _, found := this.defs[use]; !found {
			unresolved_symbols[use] = true
		}
	}
	return unresolved_symbols
}
