/*
package JackGrammar contains definitions for the Jack programming language.

The JackTokenizer uses this grammar to identify the kinds of tokens and symbols
that exist.  The List of symbols is useful as possible delimiters when
splitting the source code by token.
Having a list of keywords is also useful for the tokenizer.
Prior to analysis, a string could be a constant string, an identifier, or a keyword.
Until checking the list of keywords, it could be anything.


Syntax

In addition to lists, this JackGrammar package also includes some helper
functions.  The purpose of these is to check whether or not a string
matches a certain syntatic element.  An example is the function:
	IsSubroutineDec(string) bool
which checks to see if the given string (which should be the content
of a token), is a Subroutine Declaration.  If this were a more advanced
and practical parser, it would also check additional tokens beyond
the first keyword.


Extra Constants

Additional constants that are not included in the Jack Grammar
have been added to this package for usage by the tokenizer and parser.
Specifically, those are UNKNOWN, INVALID, and TERMINAL.

Normally, all the constants here would be enumerated, but for
the sake of simplicity of this project, they are treated as constant strings.
The abstraction of naming them here, as variables, would allow for an easy
conversion to enumerated variables.
*/
package JackGrammar

const (
	KEYWORD      = "keyword"
	SYMBOL       = "symbol"
	IDENTIFIER   = "identifier"
	INT_CONST    = "integerConstant"
	STRING_CONST = "stringConstant"
	UNKNOWN      = "?"
	INVALID      = "!INVALID!"
	TERMINAL     = "terminal"
)

// The symbols are single character "runes".
// This list of symbols specifies all the possible symbols in
// Jack programming language grammar.
var LIST_OF_SYMBOLS = []rune{
	'{', '}', '(', ')', '[', ']',
	'+', '-', '*', '/', '=',
	'&', '|', '<', '>', '~',
	';', '.', ',',
}

// The list of keywords in the jack language grammar.
var LIST_OF_KEYWORDS = []string{
	"class",
	"method",
	"function",
	"constructor",
	"int",
	"boolean",
	"char",
	"void",
	"var",
	"static",
	"field",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
	"true",
	"false",
	"null",
	"this",
}

var classVarDec = map[string]bool{
	"constructor": true,
	"function":    true,
	"method":      true,
}

func IsClassVarDec(s string) bool {
	out := classVarDec[s]
	return out
}

func IsSubroutineDec(s string) bool {
	if s == "static" || s == "field" {
		return true
	}
	return false
}
