// CompilationEngine parses source code for the Jack programming language
package CompilationEngine

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/CompilationEngine/ParseTree"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackGrammar"
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
	fmt.Fprintln(w, e.t.Root())
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
		for e.hasClassVarDec() == true {
			e.CompileClassVarDec()
		}
	}
	// closure: (subroutineDec)*
	for e.hasSubroutineDec() {
		e.CompileSubroutine()
	}
	e.CurrentToLeaf() // symbol }
}

// ClassVarDec = ('static' | 'field') type varName (',' varName)* ';'
func (e *engine) CompileClassVarDec() {
	e.t = e.t.Branch("classVarDec")
	e.CompileToken() // ('static' | 'field')
	e.CompileToken() // type
	e.CompileToken() // varName
	for e.o.Current().Content() == "," {
		e.CompileToken() // ','
		e.CompileToken() // varName
	}
	e.CompileToken() // ';'
	e.t = e.t.Up()
}

// CompileSubroutine a single subroutine declaration.
func (e *engine) CompileSubroutine() {
	e.t = e.t.Branch("subroutine")
	e.CompileToken() // ('static' | 'field' | 'constructor')
	e.CompileToken() // ('void' | type)
	e.CompileToken() // subroutineName
	e.CompileToken() // '('
	e.CompileParameterList()
	e.CompileToken() // ')'
	e.CompileSubroutineBody()
	e.t = e.t.Up()
}

// Parameter List = ((type varName) (',' type varName)*)?
func (e *engine) CompileParameterList() {
	e.t = e.t.Branch("parameterList")
	if e.o.Current().Content() == ")" {
		e.t = e.t.Up()
		return
	}
	e.CompileToken() // type
	e.CompileToken() // varName
	for e.o.Current().Content() == "," {
		e.CompileToken() // type
		e.CompileToken() // varName
	}
	e.t = e.t.Up()
}

func (e *engine) CompileSubroutineDec() {
	e.t = e.t.Branch("subroutineDec")
	e.CompileToken() // ('constructor' | 'function' | 'method')
	e.CompileToken() // subroutineName
	e.CompileToken() // '('
	e.CompileParameterList()
	e.CompileToken() // ')'
	e.CompileSubroutineBody()
	e.t = e.t.Up()
}

// subroutineBody = '{' varDec* statements '}'
func (e *engine) CompileSubroutineBody() {
	e.t = e.t.Branch("subroutineBody")
	e.CompileToken() // '{'
	for e.o.Current().Content() == "var" {
		e.CompileVarDec()
	}
	e.CompileStatements()
	e.CompileToken() // '}'
	e.t = e.t.Up()
}

// Compile Variable Declarations.
func (e *engine) CompileVarDec() {
	e.t = e.t.Branch("varDec")
	e.CompileToken() // 'var'
	e.CompileToken() // type
	e.CompileToken() // varName
	for e.o.Current().Content() == "," {
		e.CompileToken() // ','
		e.CompileToken() // varName
	}
	e.CompileToken() // ';'
	e.t = e.t.Up()
}

func (e *engine) CompileStatements() {
	e.t = e.t.Branch("statements")
	for {
		switch e.o.Current().Content() {
		case "let":
			e.CompileLet()
		case "if":
			e.CompileIf()
		case "while":
			e.CompileWhile()
		case "do":
			e.CompileDo()
		case "return":
			e.CompileReturn()
		default:
			e.t.Up()
			return
		}
	}
	e.t.Up()
}

func (e *engine) CompileExpression() {
	e.t = e.t.Branch("expression")
	e.CompileTerm()
	for e.isOperator() {
		e.CompileToken() // op
		e.CompileTerm()
	}
	e.t = e.t.Up()
}

func (e *engine) CompileLet() {
	e.t = e.t.Branch("letStatement")
	e.CompileToken() // 'let'
	e.CompileToken() // varName
	for e.o.Current().Content() == "[" {
		e.CompileToken() // '['
		e.CompileExpression()
	}
	e.CompileToken() // '='
	e.CompileExpression()
	e.CompileToken() // ';'
	e.t = e.t.Up()
}

func (e *engine) CompileIf() {
	e.t = e.t.Branch("ifStatement")
	e.CompileToken() // 'if'
	e.CompileToken() // '('
	e.CompileExpression()
	e.CompileToken() // ')'
	e.CompileToken() // '{'
	e.CompileStatements()
	e.CompileToken() // '}'
	e.t = e.t.Up()
}

func (e *engine) CompileWhile() {
	e.t = e.t.Branch("whileStatement")
	for e.o.Current().Content() == "else" {
		e.CompileToken() // 'while'
		e.CompileToken() // '('
		e.CompileExpression()
		e.CompileToken() // ')'
		e.CompileToken() // '{'
		e.CompileStatements()
		e.CompileToken() // '}'
	}
	e.t = e.t.Up()
}

func (e *engine) CompileDo() {
	e.t = e.t.Branch("doStatement")
	e.CompileToken() // 'do'
	// begin subroutineCall
	e.CompileToken() // subroutineName
	e.CompileToken() // '('
	e.CompileExpressionList()
	e.CompileToken() // ')'
	// end subroutineCall
	e.CompileToken() // ';'
	e.t = e.t.Up()
}

func (e *engine) CompileReturn() {
	e.t = e.t.Branch("returnStatement")
	e.CompileToken() // 'return'
	for e.o.Current().Content() != ";" {
		e.CompileExpression()
	}
	e.CompileToken() // ';'
	e.t = e.t.Up()
}

func (e *engine) CompileTerm() {
	e.t = e.t.Branch("term")
	defer func() {
		e.t = e.t.Up()
	}()
	// save the current token and advance. We need to do a single look-ahead.
	current_token := e.o.Current()
	e.o.Advance()
	next_token := e.o.Current()

	// switch based on the first token's kind.
	switch current_token.Kind() {
	case JackGrammar.STRING_CONST, JackGrammar.INT_CONST:
		e.t.Leaf(current_token)
		return

	case JackGrammar.SYMBOL:
		switch current_token.Content() {
		case "(":
			e.t.Leaf(current_token) // '('
			e.CompileExpression()
			e.CompileToken() // ')'
			return
		case "-", "~":
			e.t.Leaf(current_token) // unaryOp
			e.CompileTerm()
			return
		}

	case JackGrammar.IDENTIFIER:
		switch next_token.Content() {
		case "[":
			e.t.Leaf(current_token) // varName
			e.CompileToken()        // '['
			e.CompileExpression()   // exp
			e.CompileToken()        // ']'
			return
		case "(": // subroutineCall
			e.t.Leaf(current_token) // subroutineName
			e.CompileToken()        // '('
			e.CompileExpressionList()
			e.CompileToken() // ')'
			return
		case ".": //subroutineCall
			e.t.Leaf(current_token) // className | varName
			e.CompileToken()        // '.'
			e.CompileToken()        // subroutineName
			e.CompileToken()        // '('
			e.CompileExpressionList()
			e.CompileToken() // ')'
			return

		case "true", "false", "null", "this":
			e.t.Leaf(current_token) // keyword const
			return

		default:
			e.t.Leaf(current_token) // varName
		}
	}
}

func (e *engine) CompileExpressionList() {
	if e.o.Current().Content() == ")" {
		return
	}
	e.CompileExpression()
	for e.o.Current().Content() == "," {
		e.CompileToken() // ','
		e.CompileExpression()
	}
}

// -----------------------------------------------------

// Template
func (e *engine) CompileThing() {
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
	switch e.o.Current().Content() {
	case "function", "method", "constructor":
		return true
	}
	return false
}

func (e *engine) isOperator() bool {
	switch e.o.Current().Content() {
	case "+", "-", "*", "/", "&", "|", "<", ">", "=":
		return true
	}
	return false
}
