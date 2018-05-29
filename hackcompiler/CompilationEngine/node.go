package CompilationEngine

import (
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackGrammar"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/Token"
)

// SyntaxTree is the interface for the parsing engine.
// The only way it is used from the outside is by adding a new token to it.
// Iternal logic handles what is done with that knowledge.
// If all is well, error will return nil, which means it is ready for a new token,
// otherwise, there is a syntax error of some kind.
type SyntaxTree interface {
	FeedToken(Token.Token) error
	String() string
}

// NewSyntaxTree creates and initializes a SyntaxTree interface
// with a root node, and is ready to be fed tokens.
func NewSyntaxTree() SyntaxTree {
	root := new(node)
	return root
}

// Node data structure is unexported.
// It represents a syntatic element in a Syntax Tree.
// Terminal nodes contain tokens.
// Non-terminal nodes do not.
// The way they are printed depends upon their terminal status.
type node struct {
	kind   string
	leaves []*node
	token  Token.Token
}

// FeedToken accepts the token kind and content as strings,
// instead of the actual token object, and returns an error.
// If there is no syntax error, then it is okay to feed
// the tree more tokens.
func (n *node) FeedToken(t Token.Token) error {
	return nil
}

// Recursively prints the contents of the root node,
// and all leaf nodes attached to it.
func (n *node) String() string {
	s := tag(n.kind)
	for _, v := range n.leaves {
		s += v.String()
	}
	s += endtag(n.kind)
	return s
}

func (n *node) isTerminal() bool {
	if n.kind == JackGrammar.TERMINAL {
		return true
	}
	return false
}

func (n *node) addLeaf(new_leaf *node) {
	n.leaves = append(n.leaves, new_leaf)
}

func (n *node) addTerminalLeaf(token Token.Token) {
	n.addLeaf(&node{
		kind:  JackGrammar.TERMINAL,
		token: token,
	})
}

func (n *node) addNonTerminalLeaf(kind string) {
	if kind == JackGrammar.TERMINAL {
		panic("Don't pretend to be a terminal leaf!")
	}
	n.addLeaf(&node{
		kind: kind,
	})
}
