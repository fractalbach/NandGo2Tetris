package JackGrammar

const (
	KEYWORD      = "keyword"
	SYMBOL       = "symbol"
	IDENTIFIER   = "identifier"
	INT_CONST    = "integer"
	STRING_CONST = "string"
	UNKNOWN      = "?"
	INVALID      = "!INVALID!"
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
