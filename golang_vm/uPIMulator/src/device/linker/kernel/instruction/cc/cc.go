package cc

type Condition int

const (
	TRUE Condition = iota
	FALSE

	Z
	NZ

	E
	O

	PL
	MI

	OV
	NOV

	C
	NC

	SZ
	SNZ

	SPL
	SMI

	SO
	SE

	NC5
	NC6
	NC7
	NC8
	NC9
	NC10
	NC11
	NC12
	NC13
	NC14

	MAX
	NMAX

	SH32
	NSH32

	EQ
	NEQ

	LTU
	LEU
	GTU
	GEU

	LTS
	LES
	GTS
	GES

	XZ
	XNZ

	XLEU
	XGTU

	XLES
	XGTS

	SMALL
	LARGE
)

type Cc interface {
	Condition() Condition
}
