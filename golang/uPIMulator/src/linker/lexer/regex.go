package lexer

import (
	"errors"
	"regexp"
)

type Regex struct {
	expr  string
	regex *regexp.Regexp

	token_type TokenType
}

func (this *Regex) Init(expr string, token_type TokenType) {
	this.expr = expr

	regex, err := regexp.Compile(expr)

	if err != nil {
		panic(err)
	}

	this.regex = regex

	this.token_type = token_type
}

func (this *Regex) Expr() string {
	return this.expr
}

func (this *Regex) TokenType() TokenType {
	return this.token_type
}

func (this *Regex) IsTokenizable(word string) bool {
	return this.regex.MatchString(word)
}

func (this *Regex) Tokenize(word string) *Token {
	if !this.IsTokenizable(word) {
		err := errors.New("word is not matched")
		panic(err)
	}

	token := new(Token)
	token.Init(this.token_type, word)
	return token
}
