package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SRrrStmt struct {
	op_code *expr.Expr
	suffix  *expr.Expr
	dc      *lexer.Token
	ra      *expr.Expr
	rb      *expr.Expr
}

func (this *SRrrStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	rb *expr.Expr,
) {
	if op_code.ExprType() != expr.RRI_OP_CODE {
		err := errors.New("op code is not an RRI op code")
		panic(err)
	}

	if suffix.ExprType() != expr.SUFFIX {
		err := errors.New("suffix is not a suffix")
		panic(err)
	}

	if dc.TokenType() != lexer.PAIR_REG {
		err := errors.New("dc is not a pair reg")
		panic(err)
	}

	if ra.ExprType() != expr.SRC_REG {
		err := errors.New("ra is not a src reg")
		panic(err)
	}

	if rb.ExprType() != expr.SRC_REG {
		err := errors.New("rb is not a src reg")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.rb = rb
}

func (this *SRrrStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRrrStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRrrStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRrrStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SRrrStmt) Rb() *expr.Expr {
	return this.rb
}
