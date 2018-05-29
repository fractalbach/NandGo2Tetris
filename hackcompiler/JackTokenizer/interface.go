package JackTokenizer

import (
	"github.com/fractalbach/nandGo2tetris/hackcompiler/Token"
	"io"
	"io/ioutil"
)

// TokenIterator advances through each of the tokens in the source file,
// after it has been initialized.  Should use the command "NewTokenIterator"
// to initialize the token iterator.
type TokenIterator interface {
	HasMoreTokens() bool
	Advance()
	Current() Token.Token
}

// NewTokenIterator takes a reader (which should contain the source code), and returns
// an interface that allows the user to advance through the tokens in the source code.
func Create(r io.Reader) TokenIterator {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return createTokenIteratorFromString(string(b))
}
