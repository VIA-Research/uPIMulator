package lexer

type TokenType int

const (
	END_OF_FILE TokenType = iota

	IDENTIFIER
	NUMBER
	STRING
	HEADER

	INCLUDE
	DEFINE
	IFNDEF
	IFDEF
	BEGINIF
	ENDIF

	BREAK
	CHAR
	CONTINUE
	ELSE
	FOR
	IF
	INT
	LONG
	NULL
	RETURN
	SHORT
	SIZEOF
	STRUCT
	VOID
	WHILE

	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE

	LESS
	LESS_EQ
	GREATER
	GREATER_EQ
	EQ
	NOT_EQ

	PLUS
	PLUS_PLUS
	MINUS
	MINUS_MINUS
	STAR
	DIV
	MOD

	LSHIFT
	RSHIFT

	AND
	AND_AND
	OR
	OR_OR
	CARET
	NOT
	TILDE

	QUESTION
	SEMI
	COLON
	COMMA

	ASSIGN
	PLUS_ASSIGN
	MINUS_ASSIGN
	STAR_ASSIGN
	DIV_ASSIGN
	MOD_ASSIGN
	LSHIFT_ASSIGN
	RSHIFT_ASSIGN
	AND_ASSIGN
	OR_ASSIGN
	CARET_ASSIGN

	ARROW
	DOT
)

type Token struct {
	token_type TokenType
	attribute  string
}

func (this *Token) Init(token_type TokenType, attribute string) {
	this.token_type = token_type
	this.attribute = attribute
}

func (this *Token) TokenType() TokenType {
	return this.token_type
}

func (this *Token) Attribute() string {
	return this.attribute
}
