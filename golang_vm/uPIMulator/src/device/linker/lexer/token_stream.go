package lexer

type TokenStream struct {
	tokens []*Token
}

func (this *TokenStream) Init() {
	this.tokens = make([]*Token, 0)
}

func (this *TokenStream) Length() int {
	return len(this.tokens)
}

func (this *TokenStream) Get(pos int) *Token {
	return this.tokens[pos]
}

func (this *TokenStream) Append(token *Token) {
	this.tokens = append(this.tokens, token)
}

func (this *TokenStream) Merge(token_stream *TokenStream) {
	for i := 0; i < token_stream.Length(); i++ {
		this.Append(token_stream.Get(i))
	}
}
