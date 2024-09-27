package reg_descriptor

type SpRegDescriptor int

const (
	ZERO SpRegDescriptor = iota
	ONE
	LNEG
	MNEG
	ID
	ID2
	ID4
	ID8
)
