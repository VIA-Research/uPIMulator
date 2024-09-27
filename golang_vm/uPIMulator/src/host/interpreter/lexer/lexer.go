package lexer

import (
	"errors"
	"strings"
	"uPIMulator/src/misc"
)

type Lexer struct {
	tokenizer *Tokenizer
}

func (this *Lexer) Init() {
	this.tokenizer = new(Tokenizer)
	this.tokenizer.Init()
}

func (this *Lexer) Lex(path string) *TokenStream {
	file_scanner := new(misc.FileScanner)
	file_scanner.Init(path)

	token_stream := new(TokenStream)
	token_stream.Init()

	for _, line := range file_scanner.ReadLines() {
		line = strings.Replace(line, "\n", "", -1)
		token_stream.Merge(this.Tokenize(line))
	}

	end_of_file := new(Token)
	end_of_file.Init(END_OF_FILE, "")

	token_stream.Append(end_of_file)

	return token_stream
}

func (this *Lexer) Tokenize(line string) *TokenStream {
	token_stream := new(TokenStream)
	token_stream.Init()

	prev_pos := 0
	for prev_pos < len(line) {
		if token_stream.Length() > 0 {
			last_token := token_stream.Get(token_stream.Length() - 1)

			if last_token.TokenType() == INCLUDE {
				token := this.FindHeader(line, prev_pos)
				token_stream.Append(token)

				return token_stream
			}
		}

		token, length := this.FindTokenWithMaxLength(line, prev_pos)

		if token != nil {
			token_stream.Append(token)
		}

		prev_pos += length
	}

	return token_stream
}

func (this *Lexer) FindHeader(line string, prev_pos int) *Token {
	for this.IsWhiteSpace(string(line[prev_pos])) {
		prev_pos++
	}

	word := line[prev_pos:]
	token := this.tokenizer.Tokenize(word)
	return token
}

func (this *Lexer) FindTokenWithMaxLength(line string, prev_pos int) (*Token, int) {
	if prev_pos < 0 {
		err := errors.New("prev pos < 0")
		panic(err)
	}

	if this.IsWhiteSpace(string(line[prev_pos])) {
		return nil, 1
	}

	if this.IsQuote(string(line[prev_pos])) {
		next_quote_pos := this.FindNextQuote(line, prev_pos+1)

		word := line[prev_pos : next_quote_pos+1]

		token := this.tokenizer.Tokenize(word)

		return token, next_quote_pos - prev_pos + 1
	}

	for i := prev_pos + 1; i <= len(line); i++ {
		word := line[prev_pos:i]

		if i+1 <= len(line) {
			next_word := line[prev_pos : i+1]

			if this.tokenizer.IsTokenizable(word) && !this.tokenizer.IsTokenizable(next_word) {
				token := this.tokenizer.Tokenize(word)
				return token, i - prev_pos
			}
		} else {
			token := this.tokenizer.Tokenize(word)
			return token, i - prev_pos
		}
	}

	err := errors.New("line is not further tokenizable")
	panic(err)
}

func (this *Lexer) IsWhiteSpace(word string) bool {
	if len(word) != 1 {
		err := errors.New("word size != 1")
		panic(err)
	}

	return word == " " || word == "\t" || word == "\n"
}

func (this *Lexer) IsQuote(word string) bool {
	if len(word) != 1 {
		err := errors.New("word size != 1")
		panic(err)
	}

	return word == "\""
}

func (this *Lexer) FindNextQuote(line string, pos int) int {
	for i := pos; i < len(line); i++ {
		if this.IsQuote(string(line[i])) {
			return i
		}
	}

	err := errors.New("line does not have the next quote")
	panic(err)
}
