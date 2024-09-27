package instruction

import (
	"errors"
	"uPIMulator/src/device/linker/lexer"
	"uPIMulator/src/device/linker/parser/expr"
)

type SRStmt struct {
	op_code *expr.Expr
	suffix  *expr.Expr
	dc      *lexer.Token
}

func (this *SRStmt) Init(op_code *expr.Expr, suffix *expr.Expr, dc *lexer.Token) {
	if op_code.ExprType() != expr.R_OP_CODE {
		err := errors.New("op code is not an R op code")
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

	this.op_code = op_code
	this.suffix = suffix
	this.dc = dc
}

func (this *SRStmt) OpCode() *expr.Expr {
	return this.op_code
}

func (this *SRStmt) Suffix() *expr.Expr {
	return this.suffix
}

func (this *SRStmt) Dc() *lexer.Token {
	return this.dc
}
