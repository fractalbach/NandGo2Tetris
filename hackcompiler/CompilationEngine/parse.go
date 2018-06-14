// CompilationEngine parses source code for the Jack programming language
package CompilationEngine

import (
	"bytes"
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/CompilationEngine/ParseTree"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackGrammar"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackTokenizer"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/SymbolTable"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/vmWriter"
	"io"
	"strconv"
)

var (
	st                  SymbolTable.SymbolTable
	vm                  vmWriter.VMWriter
	counter             int
	nLocals             int
	nArgs               int
	expression_buf      []string
	className           string
	subroutineName      string
	symbol_table_output string
)

type OPTION int

const (
	OP_SYM_TBL OPTION = 1 << iota
	OP_XML
	OP_CODE
)

type engine struct {
	o JackTokenizer.TokenIterator
	w io.Writer
	t ParseTree.ParseTree
}

func Run(w io.Writer, tokenizer JackTokenizer.TokenIterator, opt OPTION) {
	e := engine{
		o: tokenizer,
		w: w,
		t: ParseTree.NewParseTree("class"),
	}
	buf := new(bytes.Buffer)
	vm = vmWriter.NewVMWriter(buf)
	st = SymbolTable.NewSymbolTable()
	e.CompileClass()
	switch opt {
	case OP_SYM_TBL:
		fmt.Fprint(w, symbol_table_output)
	case OP_XML:
		fmt.Fprintln(w, e.t.Root())
	case OP_CODE:
		buf.WriteTo(w)
	}
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
	className = e.o.Current().Content()
	e.CompileToken() // identifier className
	e.CompileToken() // symbol {
	// closure:  (classVarDec)*
	if e.hasClassVarDec() {
		for e.hasClassVarDec() == true {
			e.CompileClassVarDec()
		}
		symbol_table_output += fmt.Sprintln("Class Table:", className)
		symbol_table_output += st.PrintClassTable()
	}
	// closure: (subroutineDec)*
	for e.hasSubroutineDec() {
		st.StartSubroutine()
		st.Define("this", className, SymbolTable.ARG)
		e.CompileSubroutineDec()
		vm.WriteReturn()
		symbol_table_output += ("Subroutine Table:" + className + "." + subroutineName)
		symbol_table_output += st.PrintSubroutineTable()
	}
	e.CurrentToLeaf() // symbol }
}

// ClassVarDec = ('static' | 'field') type varName (',' varName)* ';'
func (e *engine) CompileClassVarDec() {
	e.t = e.t.Branch("classVarDec")
	varKind := SymbolTable.StringToKind(e.o.Current().Content())
	e.CompileToken() // ('static' | 'field')
	varType := e.o.Current().Content()
	e.CompileToken() // type
	varName := e.o.Current().Content()
	e.CompileToken() // varName
	st.Define(varName, varType, varKind)
	for e.o.Current().Content() == "," {
		e.CompileToken()                     // ','
		varName = e.o.Current().Content()    // saves the varName
		e.CompileToken()                     // varName
		st.Define(varName, varType, varKind) // adds to symbol table.
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
	sName := ""
	sType := ""
	sKind := SymbolTable.ARG
	e.t = e.t.Branch("parameterList")
	for e.o.Current().Content() != ")" {
		sType = e.o.Current().Content()
		e.CompileToken() // type
		sName = e.o.Current().Content()
		e.CompileToken() // varName
		st.Define(sName, sType, sKind)
		for e.o.Current().Content() == "," {
			e.CompileToken() // ','
			sType = e.o.Current().Content()
			e.CompileToken() // type
			sName = e.o.Current().Content()
			e.CompileToken() // varName
			st.Define(sName, sType, sKind)
		}
	}
	e.t = e.t.Up()
}

func (e *engine) CompileSubroutineDec() {
	e.t = e.t.Branch("subroutineDec")
	e.CompileToken() // ('constructor' | 'function' | 'method')
	e.CompileToken() // 'void' | type
	subroutineName = e.o.Current().Content()
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
	// writes the vm code for declaring a function at this point,
	// since all of the local variable declarations have finished:
	// The subroutine symbol table is complete, and nLocals can be counted.
	fullname := className + "." + subroutineName
	nLocals = st.VarCount(SymbolTable.VAR)
	vm.WriteFunction(fullname, nLocals)
	// continue on to compile the proccesses within the function.
	e.CompileStatements()
	e.CompileToken() // '}'
	e.t = e.t.Up()
}

// Compile Variable Declarations.  Adds Variables to the Symbol Table,
// at the scope of Subroutine.
func (e *engine) CompileVarDec() {
	sName := ""
	sType := ""
	sKind := SymbolTable.VAR
	e.t = e.t.Branch("varDec")
	e.CompileToken() // 'var'
	sType = e.o.Current().Content()
	e.CompileToken() // type
	sName = e.o.Current().Content()
	e.CompileToken() // varName
	st.Define(sName, sType, sKind)
	for e.o.Current().Content() == "," {
		e.CompileToken() // ','
		sName = e.o.Current().Content()
		e.CompileToken() // varName
		st.Define(sName, sType, sKind)
		nLocals++
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
			e.t = e.t.Up()
			return
		}
	}
	e.t = e.t.Up()
}

func (e *engine) CompileExpression() {
	e.t = e.t.Branch("expression")
	e.CompileTerm()
	for e.isOperator() {
		op := e.o.Current().Content()
		e.CompileToken() // op
		e.CompileTerm()
		cmd := vmWriter.OpToCmd(op)
		vm.WriteArithmetic(cmd)
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
		e.CompileToken() // ']'
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
	for e.o.Current().Content() == "else" {
		e.CompileToken() // 'else'
		e.CompileToken() // '{'
		e.CompileStatements()
		e.CompileToken() // '}'
	}
	e.t = e.t.Up()
}

func (e *engine) CompileWhile() {
	e.t = e.t.Branch("whileStatement")
	e.CompileToken() // 'while'
	e.CompileToken() // '('
	e.CompileExpression()
	e.CompileToken() // ')'
	e.CompileToken() // '{'
	e.CompileStatements()
	e.CompileToken() // '}'
	e.t = e.t.Up()
}

func (e *engine) CompileDo() {
	e.t = e.t.Branch("doStatement")
	e.CompileToken() // 'do'

	// for subroutine calls, you need to do 1 char look-ahead.
	current_token := e.o.Current()
	e.o.Advance()
	next_token := e.o.Current()

	// do statements create function calls, so we want to save the name.
	// nArgs will be automatically handled in the compileExpressionList()
	sName := current_token.Content()

	// decide parsing method based on that next token.
	switch next_token.Content() {
	case "(": // subroutineCall
		e.t.Leaf(current_token) // subroutineName
		e.CompileToken()        // '('
		e.CompileExpressionList()
		e.CompileToken() // ')'

	case ".": //subroutineCall
		e.t.Leaf(current_token) // className | varName
		e.CompileToken()        // '.'
		sName += ("." + e.o.Current().Content())
		e.CompileToken() // subroutineName
		e.CompileToken() // '('
		e.CompileExpressionList()
		e.CompileToken() // ')'

	default:
		panic("wtf your do statement sucks!")
	}
	e.CompileToken() // ';'
	e.t = e.t.Up()
	vm.WriteCall(sName, nArgs)
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
	case JackGrammar.STRING_CONST:
		e.t.Leaf(current_token)
		return

	case JackGrammar.INT_CONST:
		e.t.Leaf(current_token)
		val, _ := strconv.Atoi(current_token.Content())
		vm.WritePush(vmWriter.CONST, val)
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

	case JackGrammar.KEYWORD:
		switch current_token.Content() {
		case "true", "false", "null", "this":
			e.t.Leaf(current_token) // keyword const
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

		default:
			e.t.Leaf(current_token) // varName
			return
		}
	}
}

func (e *engine) CompileExpressionList() {
	nArgs = 0
	e.t = e.t.Branch("expressionList")
	if e.o.Current().Content() == ")" {
		e.t = e.t.Up()
		return
	}
	e.CompileExpression()
	nArgs++
	for e.o.Current().Content() == "," {
		e.CompileToken() // ','
		e.CompileExpression()
		nArgs++
	}
	e.t = e.t.Up()
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
