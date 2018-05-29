// package ParseTree provides the interface for creating and printing a
// Parsing Tree for the Jack Programming Language.
package ParseTree

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/Token"
)

const terminal = "TERMINAL"

// ParseTree is the interface for the parsing engine.
// It should be Initialized by calling the function: NewParseTree(string).
// The ParseTree can be traversed by creating a variable like "CurrentNode",
// and re-assigning it whenever you call Up() or Branch(string),
// since those functions also return the ParseTree interface.
//
// When you have finished creating the tree, you can return the root of
// the tree by calling Root(), and saving it in a variable.
// Calling String(), or using the ParseTree as a string in fmt.Println,
// will traverse the tree automatically, calling string() at each node,
// and print the entire tree.
//
// One of the interesting and notable features of this ParseTree Interface
// is the Leaf() function, which returns the interface of the original node
// that called the Leaf() function.  See Documentation on Leaf() for an
// example.
type ParseTree interface {
	Root() ParseTree            // returns the root of the tree.
	Up() ParseTree              // traverse up the tree.
	Leaf(Token.Token) ParseTree // leaves are terminal nodes.
	Branch(string) ParseTree    // branches are non-terminal nodes.
	String() string
}

// NewParseTree creates an initial root node, and returns the ParseTree Interface.
// The root node is treated as non-terminal node. Most likely, the name of this
// root node will be "class".  This is the node that will be returned
// whenever you call Root().
func NewParseTree(name string) ParseTree {
	root := &node{
		kind:   name,
		parent: nil,
	}
	root.root = root // lol.
	return root
}

// Node data structure is unexported.
// Terminal nodes are identified by having kind = TERMINAL
// Non-terminal nodes do not.
// The way they are printed depends upon their terminal status.
type node struct {
	kind   string
	root   *node
	parent *node
	leaves []*node
	token  Token.Token
}

// Root returns the interface of the root node in the parse tree.
func (n *node) Root() ParseTree {
	return n.root
}

// Returns the ParseTree interface on this node's parent.
func (n *node) Up() ParseTree {
	return n.parent
}

// Leaf attachs a terminal node (aka leaf) to the given node.
// It returns an interface to the original node.  This allows you to call
// 		MyTree.Leaf(a).Leaf(b).Leaf(c)
// which will attach 3 leaf nodes (a, b, c) to MyTree.  Neat, huh?!
func (n *node) Leaf(token Token.Token) ParseTree {
	new_node := &node{
		kind:   terminal,
		root:   n.root,
		parent: n,
		token:  token,
	}
	n.leaves = append(n.leaves, new_node)
	return n
}

// Branch attaches a non-terminal node, and returns a pointer to the
// freshly created node.
func (n *node) Branch(name string) ParseTree {
	new_node := &node{
		kind:   name,
		root:   n.root,
		parent: n,
		token:  nil,
	}
	n.leaves = append(n.leaves, new_node)
	return new_node
}

/*
Recursively prints the contents of the root node,
and all leaf nodes attached to it.
	Example Structure of the String
	<kind>
		<terminal> ... </terminal>
		<non-terminal>
			...
		</non-terminal>
	</kind>
*/
func (n *node) String() string {
	if n.kind == terminal {
		return fmt.Sprint(n.token)
	}
	s := tag(n.kind) + "\n"
	for _, v := range n.leaves {
		s += v.String()
	}
	s += endtag(n.kind) + "\n"
	return s
}

func tag(s string) string {
	return fmt.Sprintf("<%s>", s)
}

func endtag(s string) string {
	return fmt.Sprintf("</%s>", s)
}
