package instruction

import (
	"errors"
	"uPIMulator/src/linker/lexer"
	"uPIMulator/src/linker/parser/expr"
)

type SRrcStmt struct {
	op_code   *expr.Expr
	suffix    *expr.Expr
	dc        *lexer.Token
	ra        *expr.Expr
	condition *expr.Expr
}

func (this *SRrcStmt) Init(
	op_code *expr.Expr,
	suffix *expr.Expr,
	dc *lexer.Token,
	ra *expr.Expr,
	condition *expr.Expr,
) {
	if op_code.ExprType() != expr.RR_OP_CODE {
		err := errors.New("op code is not an RR op code")
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

	if condition.ExprType() != expr.CONDITION {
		err := errors.New("condition is not a condition")
		panic(err)
	}

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
	this.ra = ra
	this.condition = condition
}

func (this *SRrcStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRrcStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRrcStmt) Dc() *lexer.Token {
	return this.dc
}

func (this *SRrcStmt) Ra() *expr.Expr {
	return this.ra
}

func (this *SRrcStmt) Condition() *expr.Expr {
	return this.condition
}
