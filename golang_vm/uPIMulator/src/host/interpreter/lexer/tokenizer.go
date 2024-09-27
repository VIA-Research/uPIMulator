package lexer

import (
	"errors"
)

type Tokenizer struct {
	keyword_factory *KeywordFactory
	regex_factory   *RegexFactory
}

func (this *Tokenizer) Init() {
	this.InitKeywordFactory()
	this.InitRegexFactory()
}

func (this *Tokenizer) InitKeywordFactory() {
	this.keyword_factory = new(KeywordFactory)
	this.keyword_factory.Init()

	this.keyword_factory.AddKeyword("#include", INCLUDE)
	this.keyword_factory.AddKeyword("#define", DEFINE)
	this.keyword_factory.AddKeyword("#ifndef", IFNDEF)
	this.keyword_factory.AddKeyword("#ifdef", IFDEF)
	this.keyword_factory.AddKeyword("#if", BEGINIF)
	this.keyword_factory.AddKeyword("#endif", ENDIF)

	this.keyword_factory.AddKeyword("break", BREAK)
	this.keyword_factory.AddKeyword("char", CHAR)
	this.keyword_factory.AddKeyword("continue", CONTINUE)
	this.keyword_factory.AddKeyword("else", ELSE)
	this.keyword_factory.AddKeyword("for", FOR)
	this.keyword_factory.AddKeyword("if", IF)
	this.keyword_factory.AddKeyword("int", INT)
	this.keyword_factory.AddKeyword("long", LONG)
	this.keyword_factory.AddKeyword("NULL", NULL)
	this.keyword_factory.AddKeyword("return", RETURN)
	this.keyword_factory.AddKeyword("short", SHORT)
	this.keyword_factory.AddKeyword("sizeof", SIZEOF)
	this.keyword_factory.AddKeyword("struct", STRUCT)
	this.keyword_factory.AddKeyword("void", VOID)
	this.keyword_factory.AddKeyword("while", WHILE)

	this.keyword_factory.AddKeyword("(", LPAREN)
	this.keyword_factory.AddKeyword(")", RPAREN)
	this.keyword_factory.AddKeyword("[", LBRACKET)
	this.keyword_factory.AddKeyword("]", RBRACKET)
	this.keyword_factory.AddKeyword("{", LBRACE)
	this.keyword_factory.AddKeyword("}", RBRACE)

	this.keyword_factory.AddKeyword("<", LESS)
	this.keyword_factory.AddKeyword("<=", LESS_EQ)
	this.keyword_factory.AddKeyword(">", GREATER)
	this.keyword_factory.AddKeyword(">=", GREATER_EQ)
	this.keyword_factory.AddKeyword("==", EQ)
	this.keyword_factory.AddKeyword("!=", NOT_EQ)

	this.keyword_factory.AddKeyword("+", PLUS)
	this.keyword_factory.AddKeyword("++", PLUS_PLUS)
	this.keyword_factory.AddKeyword("-", MINUS)
	this.keyword_factory.AddKeyword("--", MINUS_MINUS)
	this.keyword_factory.AddKeyword("*", STAR)
	this.keyword_factory.AddKeyword("/", DIV)
	this.keyword_factory.AddKeyword("%", MOD)

	this.keyword_factory.AddKeyword("<<", LSHIFT)
	this.keyword_factory.AddKeyword(">>", RSHIFT)

	this.keyword_factory.AddKeyword("&", AND)
	this.keyword_factory.AddKeyword("&&", AND_AND)
	this.keyword_factory.AddKeyword("|", OR)
	this.keyword_factory.AddKeyword("||", OR_OR)
	this.keyword_factory.AddKeyword("^", CARET)
	this.keyword_factory.AddKeyword("!", NOT)
	this.keyword_factory.AddKeyword("~", TILDE)

	this.keyword_factory.AddKeyword("?", QUESTION)
	this.keyword_factory.AddKeyword(";", SEMI)
	this.keyword_factory.AddKeyword(":", COLON)
	this.keyword_factory.AddKeyword(",", COMMA)

	this.keyword_factory.AddKeyword("=", ASSIGN)
	this.keyword_factory.AddKeyword("+=", PLUS_ASSIGN)
	this.keyword_factory.AddKeyword("-=", MINUS_ASSIGN)
	this.keyword_factory.AddKeyword("*=", STAR_ASSIGN)
	this.keyword_factory.AddKeyword("/=", DIV_ASSIGN)
	this.keyword_factory.AddKeyword("%=", MOD_ASSIGN)
	this.keyword_factory.AddKeyword("<<=", LSHIFT_ASSIGN)
	this.keyword_factory.AddKeyword(">>=", RSHIFT_ASSIGN)
	this.keyword_factory.AddKeyword("&=", AND_ASSIGN)
	this.keyword_factory.AddKeyword("|=", OR_ASSIGN)
	this.keyword_factory.AddKeyword("^=", CARET_ASSIGN)

	this.keyword_factory.AddKeyword("->", ARROW)
	this.keyword_factory.AddKeyword(".", DOT)
}

func (this *Tokenizer) InitRegexFactory() {
	this.regex_factory = new(RegexFactory)
	this.regex_factory.Init()

	this.regex_factory.AddRegex("^([A-Za-z_])([A-Za-z0-9_]*)$", IDENTIFIER)
	this.regex_factory.AddRegex("^([0-9]+)$", NUMBER)
	this.regex_factory.AddRegex("^\"([A-Za-z0-9/_.])*\"$", STRING)
	this.regex_factory.AddRegex("^<([A-Za-z0-9/_.]*)>$", HEADER)
}

func (this *Tokenizer) IsTokenizable(word string) bool {
	return this.keyword_factory.IsTokenizable(word) || this.regex_factory.IsTokenizable(word)
}

func (this *Tokenizer) Tokenize(word string) *Token {
	if this.keyword_factory.IsTokenizable(word) {
		return this.keyword_factory.Tokenize(word)
	} else if this.regex_factory.IsTokenizable(word) {
		return this.regex_factory.Tokenize(word)
	} else {
		err := errors.New("word is not tokenizable")
		panic(err)
	}
}
