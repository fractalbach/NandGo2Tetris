/*
package SymbolTable helps the Jack programming language compiler associate identifiers in the source code with their properties,
to be used in code generation.

The symbol table keeps track of identifiers as they are found in
the code.  There are 2 different scopes: Class and Subroutine,
and a running index to keep track of various levels of nested scope.
Each identifier has a "type", and "Kind", and a running index.

These properties are used by the VMwriter during the code generation
phase of the compiler.
*/
package SymbolTable

type Kind int

const (
	STATIC Kind = 1 << iota
	FIELD
	ARG
	VAR
	NONE
)

type SymbolTable interface {
	StartSubroutine()
	Define(string, string, Kind)
	VarCount(Kind) int
	KindOf(string) Kind
	TypeOf(string) string
	IndexOf(string) int
}

func NewSymbolTable() SymbolTable {
	return &symbolTable{}
}

type symbol struct {
	kind  Kind
	type_ string
	index int
}

type symbolTable struct {
	class_table      map[string]symbol
	subroutine_table map[string]symbol
	nStatic          int
	nField           int
	nArg             int
	nVar             int
}

// StartSubroutine creates a new subroutine scope by resetting the
// subroutine symbol table.  Resets running indicies for Args and Vars.
func (st *symbolTable) StartSubroutine() {
	st.subroutine_table = map[string]symbol{}
	st.nArg = 0
	st.nVar = 0
}

// Creates a new Symbol and adds it to the symbol table.  The symbol "kind"
// will determine if it is subroutine scope or class scope.
func (st *symbolTable) Define(name, type_ string, kind Kind) {
	index := -1
	switch kind {
	case STATIC:
		index = st.nStatic
		st.nStatic++
	case FIELD:
		index = st.nField
		st.nField++
	case ARG:
		index = st.nArg
		st.nArg++
	case VAR:
		index = st.nVar
		st.nVar++
	default:
		panic("Unknown Symbol Kind")
	}
	switch kind {
	case STATIC, FIELD:
		st.class_table[name] = symbol{
			kind:  kind,
			type_: type_,
			index: index,
		}
	case ARG, VAR:
		st.subroutine_table[name] = symbol{
			kind:  kind,
			type_: type_,
			index: index,
		}
	}
}

func (st *symbolTable) VarCount(kind Kind) int {
	switch kind {
	case STATIC:
		return st.nStatic
	case FIELD:
		return st.nField
	case ARG:
		return st.nArg
	case VAR:
		return st.nVar
	default:
		panic("Unknown Symbol Kind")
	}
}

// KindOf looks up the symbol and returns it's Kind. Panics if not found.
func (st *symbolTable) KindOf(identifier string) Kind {
	return st.lookupSymbol(identifier).kind
}

// TypeOf looks up the symbol and returns it's Type. Panics if not found.
func (st *symbolTable) TypeOf(identifier string) string {
	return st.lookupSymbol(identifier).type_
}

// IndexOf looks up the symbol and returns it's Index. Panics if not found.
func (st *symbolTable) IndexOf(identifier string) int {
	return st.lookupSymbol(identifier).index
}

func (st *symbolTable) lookupSymbol(identifier string) symbol {
	symbol, ok := st.subroutine_table[identifier]
	if ok {
		return symbol
	}
	symbol, ok = st.class_table[identifier]
	if ok {
		return symbol
	}
	panic("Cannot find symbol.")
}
