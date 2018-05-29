// CompilationEngine parses source code for the Jack programming language
package CompilationEngine

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/CompilationEngine/ParseTree"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackGrammar"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackTokenizer"
	// "github.com/fractalbach/nandGo2tetris/hackcompiler/Token"
	"io"
)

type engine struct {
	i JackTokenizer.TokenIterator
	w io.Writer
}

func Run(w io.Writer, tokenizer JackTokenizer.TokenIterator) {
	e := engine{
		i: tokenizer,
		w: w,
	}
	tree := e.CompileClass()
	fmt.Fprintln(w, tree)
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
	e.println(e.i.Current())
}

func (e *engine) Write(p []byte) (int, error) {
	return e.w.Write(p)
}

// Class:  'class' className '{' classVarDec* subroutineDec* '}'
//
/*
	Structure of a Class

	class
		className
			String (terminal)
		symbol (terminal)
		classVarDec*
			...
		subroutineDec*
			...
		symbol (terminal)

*/
func (e *engine) compileClass() {

	e.tag("class")
	// advance to class name.
	e.i.Advance()

	e.tag("className")
	e.printCurrentToken()
	e.endtag("className")
	// advance to '{'
	e.i.Advance()
	e.printCurrentToken()
	// the next advance could end up at any of the following:
	// classVarDec* | subroutineDec* | '}',
	// it depends on what is given as input.
	e.i.Advance()
	// classVarDec*
	for JackGrammar.IsClassVarDec(e.i.Current().Content()) {
		e.CompileClassVarDec()
	}
	// subroutineDec*
	for JackGrammar.IsSubroutineDec(e.i.Current().Content()) {
		e.CompileSubroutineDec()
	}
	// '}'
	e.printCurrentToken()
	e.endtag("class")
}

// CompileClass is the first function that is called by Run(), so it
// creates and returns the Parse Tree.  All other compile methods are
// called from CompileClass().
func (e *engine) CompileClass() ParseTree.ParseTree {
	t := ParseTree.NewParseTree("class")
	t.Leaf(e.i.Current()) // keyword 'class'
	e.i.Advance()
	t.Leaf(e.i.Current()) // identifier className
	e.i.Advance()
	t.Leaf(e.i.Current()) // symbol {
	e.i.Advance()
	// closure:  (classVarDec)*
	for JackGrammar.IsClassVarDec(e.i.Current().Content()) {
		e.CompileClassVarDec()
	}
	// closure: (subroutineDec)*
	for JackGrammar.IsSubroutineDec(e.i.Current().Content()) {
		e.CompileSubroutineDec()
	}
	t.Leaf(e.i.Current()) // symbol }
	return t
}

// ClassVarDec = ('static' | 'field')
func (e *engine) CompileClassVarDec() {
	e.tag("ClassVarDec")
	defer e.endtag("ClassVarDec")

}

func (e *engine) CompileSubroutineDec() {

}

func tag(s string) string {
	return fmt.Sprintf("<%s>", s)
}

func endtag(s string) string {
	return fmt.Sprintf("</%s>", s)
}
