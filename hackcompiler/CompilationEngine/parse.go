package CompilationEngine

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackGrammar"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackTokenizer"
	// "github.com/fractalbach/nandGo2tetris/hackcompiler/Token"
	"io"
)

type engine struct {
	jt JackTokenizer.TokenIterator
	w  io.Writer
}

func Run(w io.Writer, tokenizer JackTokenizer.TokenIterator) {
	e := engine{
		jt: tokenizer,
		w:  w,
	}
	e.CompileClass()
}

func (e *engine) CompileClassVarDec() {

}

func (e *engine) CompileSubroutineDec() {

}

// Class:  'class' className '{' classVarDec* subroutineDec* '}'
func (e *engine) CompileClass() {
	e.tag("class")

	// advance to class name.
	e.jt.Advance()
	e.tag("className")
	e.printCurrentToken()
	e.endtag("className")

	// advance to '{'
	e.jt.Advance()
	e.printCurrentToken()

	// advance to classVarDec* or subroutineDec* or '}',
	// it depends on what is given as input.
	e.jt.Advance()

	// classVarDec*
	for JackGrammar.IsClassVarDec(e.jt.Current().Content()) {
		e.CompileClassVarDec()
	}

	// subroutineDec*
	for JackGrammar.IsSubroutineDec(e.jt.Current().Content()) {
		e.CompileSubroutineDec()
	}

	// '}'
	e.printCurrentToken()

	e.endtag("class")
}

func ExampleRunMe() {
	root := &node{kind: "class"}
	leaf1 := &node{kind: "function"}
	leaf2 := &node{kind: "static"}
	leaf3 := &node{kind: "return"}

	// a := Token.NewToken("keyword", "class")
	// b := Token.NewToken("symbol", "{")
	// c := Token.NewToken("symbol", "}")

	root.addLeaf(leaf1)
	root.addLeaf(leaf2)
	root.addLeaf(leaf3)
	fmt.Println(root)
}

func tag(s string) string {
	return fmt.Sprintf("<%s>", s)
}

func endtag(s string) string {
	return fmt.Sprintf("</%s>", s)
}

func (e *engine) tag(s string) {
	fmt.Fprintln(e, tag(s))
}

func (e *engine) endtag(s string) {
	fmt.Fprintln(e, endtag(s))
}

func (e *engine) println(a ...interface{}) {
	fmt.Fprintln(e, a...)
}

func (e *engine) printCurrentToken() {
	e.println(e.jt.Current())
}

func (e *engine) Write(p []byte) (int, error) {
	return e.w.Write(p)
}
