package lexer

import (
	"errors"
)

type KeywordFactory struct {
	keywords map[string]TokenType
}

func (this *KeywordFactory) Init() {
	this.keywords = make(map[string]TokenType, 0)
}

func (this *KeywordFactory) AddKeyword(keyword string, token_type TokenType) {
	this.keywords[keyword] = token_type
}

func (this *KeywordFactory) IsTokenizable(word string) bool {
	_, found := this.keywords[word]
	return found
}

func (this *KeywordFactory) Tokenize(word string) *Token {
	if _, found := this.keywords[word]; !found {
		err := errors.New("word is not tokenizable")
		panic(err)
	}

	token_type := this.keywords[word]

	token := new(Token)
	token.Init(token_type, "")
	return token
}
