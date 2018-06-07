// CompilationEngine parses source code for the Jack programming language
package CompilationEngine

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/CompilationEngine/ParseTree"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackTokenizer"
	"io"
)

type engine struct {
	o JackTokenizer.TokenIterator
	w io.Writer
	t ParseTree.ParseTree
}

func Run(w io.Writer, tokenizer JackTokenizer.TokenIterator) {
	e := engine{
		o: tokenizer,
		w: w,
		t: ParseTree.NewParseTree("class"),
	}
	e.CompileClass()
	fmt.Fprintln(w, e.t)
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
	e.println(e.o.Current())
}

func (e *engine) Write(p []byte) (int, error) {
	return e.w.Write(p)
}

// CompileClass is the first function that is called by Run(), so it
// creates and returns the Parse Tree.  All other compile methods are
// called from CompileClass().
//
// Class:  'class' className '{' classVarDec* subroutineDec* '}'
//
func (e *engine) CompileClass() {
	e.CompileToken() // keyword 'class'
	e.CompileToken() // identifier className
	e.CompileToken() // symbol {
	// closure:  (classVarDec)*
	if e.hasClassVarDec() {
		e.t = e.t.Branch("classVarDec")
		for e.hasClassVarDec() == true {
			e.CompileClassVarDec()
		}
		e.t = e.t.Up()
	}
	// closure: (subroutineDec)*
	for e.hasSubroutineDec() {
		e.t = e.t.Branch("subroutineDec")
		e.CompileSubroutine()
		e.t = e.t.Up()
	}
	e.CurrentToLeaf() // symbol }
}

// ClassVarDec = ('static' | 'field') type varName (',' varName)* ';'
func (e *engine) CompileClassVarDec() {
	e.CompileToken() // ('static' | 'field')
	e.CompileToken() // type
	e.CompileToken() // varName
	for e.o.Current().Content() == "," {
		e.CompileToken() // ','
		e.CompileToken() // varName
	}
	e.CompileToken() // ';'
	e.o.Advance()
}

// CompileSubroutine a single subroutine declaration.
func (e *engine) CompileSubroutine() {
	e.CompileToken() // ('static' | 'field' | 'constructor')
	e.CompileToken() // ('void' | type)
	e.CompileToken() // subroutineName
	e.CompileToken() // '('
	// parameterList
	e.t = e.t.Branch("parameterList")
	e.CompileParameterList()
	e.t = e.t.Up()

	e.CompileToken() // ')'

	// subroutineBody
	e.t = e.t.Branch("subroutineBody")
	// e.CompileSubroutineBody()
	e.t = e.t.Up()

}

// Parameter List = ((type varName) (',' type varName)*)?
func (e *engine) CompileParameterList() {
	if e.o.Current().Content() == ")" {
		return
	}
CompileParameter:
	e.CompileToken() // type
	e.CompileToken() // varName
	if e.o.Current().Content() == "," {
		goto CompileParameter
	}
}

// subroutineBody = '{' varDec* statements '}'
func (e *engine) CompileSubroutineBody() {
	// '{'
	e.CompileToken()
	// varDec*
	// if e.hasVarDec() {
	// 	//	 TODO
	// }
	// statements
	// '}'
	e.CompileToken()
}

func (e *engine) CurrentToLeaf() {
	e.t.Leaf(e.o.Current())
}

func (e *engine) CompileToken() {
	e.t.Leaf(e.o.Current())
	e.o.Advance()
}

func tag(s string) string {
	return fmt.Sprintf("<%s>", s)
}

func endtag(s string) string {
	return fmt.Sprintf("</%s>", s)
}

func (e *engine) hasClassVarDec() bool {
	c := e.o.Current().Content()
	return (c == "static" || c == "field")
}

func (e *engine) hasSubroutineDec() bool {
	c := e.o.Current().Content()
	return (c == "function" || c == "method" || c == "constructor")
}
