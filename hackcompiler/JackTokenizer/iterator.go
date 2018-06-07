package JackTokenizer

import "github.com/fractalbach/nandGo2tetris/hackcompiler/Token"

// import "fmt"

type tokenIterator struct {
	token_list          []Token.Token
	current_token_index int
}

func createTokenIteratorFromString(s string) *tokenIterator {
	ti := new(tokenIterator)
	ti.token_list = tokenize(s)
	// fmt.Println("Total number of tokens:", len(ti.token_list))
	return ti
}

func (ti *tokenIterator) Advance() {
	if ti.HasMoreTokens() {
		ti.current_token_index++
		return
	}
	// panic("Tried to Advance to the next token, but the Token Iterator has already reached it's last token.")
}

func (ti *tokenIterator) Current() Token.Token {
	return ti.token_list[ti.current_token_index]
}

func (ti *tokenIterator) HasMoreTokens() bool {
	if len(ti.token_list) <= 0 {
		return false
	}
	if ti.current_token_index < len(ti.token_list) {
		return true
	}
	return false
}
