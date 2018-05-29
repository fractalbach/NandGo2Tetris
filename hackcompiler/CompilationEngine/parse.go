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
	jt JackTokenizer.TokenIterator
	w  io.Writer
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

func Run(w io.Writer, tokenizer JackTokenizer.TokenIterator) {
	e := engine{
		jt: tokenizer,
		w:  w,
	}
	e.CompileClass()
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
	// the next advance could end up at any of the following:
	// classVarDec* | subroutineDec* | '}',
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

// ClassVarDec = ('static' | 'field')
func (e *engine) CompileClassVarDec() {
	e.tag("ClassVarDec")
	defer e.endtag("ClassVarDec")

}

func (e *engine) CompileSubroutineDec() {

}
