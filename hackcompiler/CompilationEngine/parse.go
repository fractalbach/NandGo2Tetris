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
	nLocals             int
	expression_buf      []string
	className           string
	subroutineName      string
	symbol_table_output string
	label_counter       int = 0
	isConstructor       bool
	inside_method       bool
	// void_functions          = make(map[string]bool)
)

type OPTION int

const (
	OP_SYM_TBL OPTION = 1 << iota
	OP_XML
	OP_CODE
)

var kindToSeg = map[SymbolTable.Kind]vmWriter.Segment{
	SymbolTable.ARG:    vmWriter.ARG,
	SymbolTable.STATIC: vmWriter.STATIC,
	SymbolTable.VAR:    vmWriter.LOCAL,
	SymbolTable.FIELD:  vmWriter.THIS,
}

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
	symbol_table_output = "" // clean up global variable
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
		e.CompileSubroutineDec()
		vm.WriteReturn()
		symbol_table_output += ("Subroutine Table:" + className + "." + subroutineName + "\n")
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

// Parameter List = ((type varName) (',' type varName)*)?
func (e *engine) CompileParameterList() int {
	n := 0
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
		n++
		for e.o.Current().Content() == "," {
			e.CompileToken() // ','
			sType = e.o.Current().Content()
			e.CompileToken() // type
			sName = e.o.Current().Content()
			e.CompileToken() // varName
			st.Define(sName, sType, sKind)
			n++
		}
	}
	e.t = e.t.Up()
	return n
}

func (e *engine) CompileSubroutineDec() {
	e.t = e.t.Branch("subroutineDec")
	isConstructor = false
	switch e.o.Current().Content() {
	case "constructor":
		isConstructor = true
		inside_method = true
	case "method":
		st.Define("this", className, SymbolTable.ARG)
		inside_method = true
	}
	e.CompileToken() // ('constructor' | 'function' | 'method')
	e.CompileToken() // 'void' | type
	subroutineName = e.o.Current().Content()
	e.CompileToken() // subroutineName
	e.CompileToken() // '('
	e.CompileParameterList()
	e.CompileToken() // ')'
	e.CompileSubroutineBody()
	e.t = e.t.Up()
	isConstructor = false // clean up global variable.
	inside_method = false // no longer inside a method/constructor
}

// subroutineBody = '{' varDec* statements '}'
func (e *engine) CompileSubroutineBody() {
	e.t = e.t.Branch("subroutineBody")
	e.CompileToken() // '{'
	for e.o.Current().Content() == "var" {
		e.CompileVarDec()
	}
	// must call WriteFunction() after CompileVarDec() has been called,
	// because the VM command needs to know how many local variables
	// there are in the function.
	fullname := className + "." + subroutineName
	nLocals = st.VarCount(SymbolTable.VAR)
	vm.WriteFunction(fullname, nLocals)

	// check the current symbol table to see if we added a "this" to
	// the subroutine.  If we have, then we must make sure to write
	// "push argument 0".
	if st.Has("this") && st.KindOf("this") == SymbolTable.ARG {
		vm.WritePush(vmWriter.ARG, 0)
		vm.WritePop(vmWriter.POINTER, 0)
	}

	// constructors are special because they need to allocate memory.
	// determine the sizeOf() by counting the number of fields,
	// 1. push the size to the stack,
	// 2. call Memory.alloc(size), which leave a pointer on the stack,
	// 3. pop stack, changing THIS pointer to the new pointer.
	if subroutineName == "new" {
		size := st.VarCount(SymbolTable.FIELD)
		vm.WritePush(vmWriter.CONST, size)
		vm.WriteCall("Memory.alloc", 1)
		vm.WritePop(vmWriter.POINTER, 0)
	}
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
	sName := e.o.Current().Content()
	sKind := st.KindOf(sName)
	sIndex := st.IndexOf(sName)
	e.CompileToken() // varName
	for e.o.Current().Content() == "[" {
		e.CompileToken() // '['
		e.CompileExpression()
		e.CompileToken() // ']'
	}
	e.CompileToken() // '='
	e.CompileExpression()
	e.CompileToken() // ';'
	vm.WritePop(kindToSeg[sKind], sIndex)
	e.t = e.t.Up()
}

func (e *engine) CompileIf() {
	prefix := "L_"
	label1 := fmt.Sprint(prefix, "IF", label_counter)
	label2 := fmt.Sprint(prefix, "ENDIF", label_counter)
	// finished_first_else := false
	label_counter++
	e.t = e.t.Branch("ifStatement")
	e.CompileToken() // 'if'
	e.CompileToken() // '('
	e.CompileExpression()
	e.CompileToken() // ')'
	vm.WriteArithmetic(vmWriter.NOT)
	vm.WriteIf(label1)
	e.CompileToken() // '{'
	e.CompileStatements()
	e.CompileToken() // '}'
	vm.WriteGoto(label2)
	vm.WriteLabel(label1)
	for e.o.Current().Content() == "else" {
		e.CompileToken() // 'else'
		e.CompileToken() // '{'
		e.CompileStatements()
		e.CompileToken() // '}'
		// finished_first_else = true
	}
	vm.WriteLabel(label2)
	e.t = e.t.Up()
}

func (e *engine) CompileWhile() {
	prefix := "L_"
	label1 := fmt.Sprint(prefix, "WHILE", label_counter)
	label2 := fmt.Sprint(prefix, "ENDWHILE", label_counter)
	label_counter++
	vm.WriteLabel(label1)
	e.t = e.t.Branch("whileStatement")
	e.CompileToken() // 'while'
	e.CompileToken() // '('
	e.CompileExpression()
	e.CompileToken() // ')'
	vm.WriteArithmetic(vmWriter.NOT)
	vm.WriteIf(label2)
	e.CompileToken() // '{'
	e.CompileStatements()
	e.CompileToken() // '}'
	e.t = e.t.Up()
	vm.WriteGoto(label1)
	vm.WriteLabel(label2)
}

func (e *engine) CompileDo() {
	e.t = e.t.Branch("doStatement")
	e.CompileToken() // 'do'

	// for subroutine calls, you need to do 1 char look-ahead.
	current_token := e.o.Current()
	e.o.Advance()
	next_token := e.o.Current()

	// do statements create function calls, so we want to save the name.
	nArgs := 0
	sName := current_token.Content()

	// being inside a constructor or method, and having no explicit
	// receiver, implies that the current object is the receiver.
	// The pointer to current object is accesible via POINTER 0

	// decide parsing method based on that next token.
	switch next_token.Content() {
	case "(": // subroutineCall
		if inside_method {
			vm.WritePush(vmWriter.POINTER, 0)
			nArgs++
		}
		e.t.Leaf(current_token) // subroutineName
		e.CompileToken()        // '('
		nArgs += e.CompileExpressionList()
		e.CompileToken() // ')'
		// assume that the prefix is the current class name.
		sName = fmt.Sprint(className, ".", sName)

	case ".": //subroutineCall
		// lookup the name of the receiver in the symbol table.
		// If we can find it, then we know it's a variable which
		// contains a reference to an object.  (now called a receiver).
		// Push the receiver pointer onto the stack.  the pointer will
		// become an argument for calling a method.
		recName := current_token.Content()
		if st.Has(recName) {
			segment := kindToSeg[st.KindOf(recName)]
			index := st.IndexOf(recName)
			vm.WritePush(segment, index)
			sName = st.TypeOf(recName)
			nArgs++
		}
		// If we couldn't find the variable in the symbol table,
		// then assume it was a className, and there is no receiver.

		e.t.Leaf(current_token) // className | varName
		e.CompileToken()        // '.'
		sName += ("." + e.o.Current().Content())
		e.CompileToken() // subroutineName
		e.CompileToken() // '('
		nArgs += e.CompileExpressionList()
		e.CompileToken() // ')'

	default:
		panic("wtf your do statement sucks!")
	}
	e.CompileToken() // ';'
	e.t = e.t.Up()
	vm.WriteCall(sName, nArgs)
	vm.WritePop(vmWriter.TEMP, 0)
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
		// create a new string, and append characters to it.
		// Leave a pointer to the string on top of the stack.
		// An example (
		// 		https://play.golang.org/p/MxQsrUSU6dD
		// where a pointer to the array containing "hello world"
		// is pushed to the stack.
		str := current_token.Content()
		vm.WritePush(vmWriter.CONST, len(str))
		vm.WriteCall("String.new", 1)
		for _, c := range str {
			vm.WritePush(vmWriter.CONST, int(c))
			vm.WriteCall("String.appendChar", 2)
		}
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
		case "-":
			e.t.Leaf(current_token) // unaryOp - arithmetic negation
			e.CompileTerm()
			vm.WriteArithmetic(vmWriter.NEG)
			return
		case "~":
			e.t.Leaf(current_token) // unaryOp - boolean negation
			e.CompileTerm()
			vm.WriteArithmetic(vmWriter.NOT)
			return
		}

	case JackGrammar.KEYWORD:
		switch current_token.Content() {
		case "true":
			vm.WritePush(vmWriter.CONST, 1)
			vm.WriteArithmetic(vmWriter.NEG)
			e.t.Leaf(current_token) // keyword const
			return
		case "null", "false":
			vm.WritePush(vmWriter.CONST, 0)
		case "this":
			e.t.Leaf(current_token) // keyword const
			if st.Has("this") && st.KindOf("this") == SymbolTable.ARG {
				vm.WritePush(vmWriter.ARG, 0)
				return
			}
			vm.WritePush(vmWriter.POINTER, 0)
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
			nArgs := 0
			if inside_method {
				vm.WritePush(vmWriter.POINTER, 0)
				nArgs++
			}
			sName := current_token.Content()
			e.t.Leaf(current_token) // subroutineName
			e.CompileToken()        // '('
			nArgs += e.CompileExpressionList()
			e.CompileToken() // ')'
			vm.WriteCall((className + "." + sName), nArgs)
			return
		case ".": //subroutineCall
			nArgs := 0
			sName1 := current_token.Content()
			// If the name before a "." is a known variable, and is
			// in our symbol table, then it is being used as a reciever
			// for a method call.  If this is the case, push it to the
			// stack, because  it will become argument 0 for the calling
			// function.
			if st.Has(sName1) {
				segment := kindToSeg[st.KindOf(sName1)]
				index := st.IndexOf(sName1)
				vm.WritePush(segment, index)
				// additionally, we need to extract the name of the class
				// from the variable.  It's the TYPE.
				sName1 = st.TypeOf(sName1)
				nArgs++
			}
			e.t.Leaf(current_token) // className | varName
			e.CompileToken()        // '.'
			sName2 := e.o.Current().Content()
			e.CompileToken() // subroutineName
			e.CompileToken() // '('
			nArgs += e.CompileExpressionList()
			e.CompileToken() // ')'
			vm.WriteCall((sName1 + "." + sName2), nArgs)
			return

		default:
			varName := current_token.Content()
			segment := kindToSeg[st.KindOf(varName)]
			index := st.IndexOf(varName)
			vm.WritePush(segment, index)
			e.t.Leaf(current_token) // varName
			return
		}
	}
}

func (e *engine) CompileExpressionList() int {
	n := 0
	e.t = e.t.Branch("expressionList")
	if e.o.Current().Content() == ")" {
		e.t = e.t.Up()
		return n
	}
	e.CompileExpression()
	n++
	for e.o.Current().Content() == "," {
		e.CompileToken() // ','
		e.CompileExpression()
		n++
	}
	e.t = e.t.Up()
	return n
}

// -----------------------------------------------------

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
