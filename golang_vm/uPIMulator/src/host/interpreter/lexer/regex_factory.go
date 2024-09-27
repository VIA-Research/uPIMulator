package lexer

import (
	"errors"
)

type RegexFactory struct {
	regexes []*Regex
}

func (this *RegexFactory) Init() {
	this.regexes = make([]*Regex, 0)
}

func (this *RegexFactory) HasRegex(expr string) bool {
	for _, regex := range this.regexes {
		if regex.Expr() == expr {
			return true
		}
	}
	return false
}

func (this *RegexFactory) AddRegex(expr string, token_type TokenType) {
	if this.HasRegex(expr) {
		err := errors.New("regex already exists")
		panic(err)
	}

	regex := new(Regex)
	regex.Init(expr, token_type)

	this.regexes = append(this.regexes, regex)
}

func (this *RegexFactory) IsTokenizable(word string) bool {
	for _, regex := range this.regexes {
		if regex.IsTokenizable(word) {
			return true
		}
	}
	return false
}

func (this *RegexFactory) Tokenize(word string) *Token {
	if !this.IsTokenizable(word) {
		err := errors.New("word is not tokenizable")
		panic(err)
	}

	for _, regex := range this.regexes {
		if regex.IsTokenizable(word) {
			token := new(Token)
			token.Init(regex.TokenType(), word)
			return token
		}
	}

	return nil
}
