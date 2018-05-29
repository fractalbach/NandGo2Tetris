package JackTokenizer

import (
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackGrammar"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/Token"
	"strconv"
	"unicode"
)

type tokenSplitter struct {
	arr         []Token.Token
	buf         string
	string_mode bool
}

func tokenize(s string) []Token.Token {

	symbols := JackGrammar.LIST_OF_SYMBOLS
	t := new(tokenSplitter)
	t.string_mode = false
	skip := false
	comment_mode := false
	extended_comment_mode := false
	exiting_extended_comment_mode := false

	for i, r := range s {

		if exiting_extended_comment_mode {
			if r == '/' {
				exiting_extended_comment_mode = false
				extended_comment_mode = false
			}
			continue
		}

		// Check for comments.
		// 1. Check if current character is /
		// 2. If it is, make sure we aren't at the end of file.
		// 3. Look ahead 1 character to see if there is another /
		// If we confirm that there is a // present, then enter comment mode.
		if r == '/' && i+1 < len(s) {
			// regular comment mode.
			if s[i+1] == '/' {
				t.buf = ""
				comment_mode = true
			}
			// extended comment mode.
			if s[i+1] == '*' {
				t.buf = ""
				extended_comment_mode = true
			}
		}

		// If we are in comment mode, ignore all characters until endline.
		if comment_mode {
			if r == '\n' {
				comment_mode = false
			}
			continue
		}
		if extended_comment_mode {
			if r == '*' && i+1 < len(s) && s[i+1] == '/' {
				exiting_extended_comment_mode = true
			}
			continue
		}

		// Toggle string mode when encounting a " character.
		// Push the current buffer.
		if r == '"' {
			// If we were already in string mode,
			// then we are currently exiting string mode, so the buffer should
			// contain a constant string.
			if t.string_mode {
				t.push(JackGrammar.STRING_CONST)
			}
			// If we were NOT already in string mode, then we are entering string right now,
			// and we don't know what is in the buffer, so push an unknown token.
			if !t.string_mode {
				t.push("?")
			}
			// Now that buffers have been taken care of... toggle string mode to the new value.
			t.string_mode = !t.string_mode
			continue
		}

		// During string mode, always add the character to buffer.
		// Add the contents of the string to the buffer.
		if t.string_mode {
			t.add2buf(r)
			continue
		}

		// Split by spaces.  Any empty elements will be removed later.
		if r == ' ' {
			t.push(JackGrammar.UNKNOWN)
			continue
		}

		// Check the symbol chart for a match.
		// Split by symbols, but include the symbol as it's own token.
		skip = false
		for i, _ := range symbols {
			if r == symbols[i] {
				t.push(JackGrammar.UNKNOWN) // could be anything.
				t.add2buf(r)                // add symbol to buffer
				t.push(JackGrammar.SYMBOL)  // definitely a symbol.
				skip = true
				break
			}
		}

		// default action is to add the character to the buffer.
		if !skip {
			t.add2buf(r)
		}
	}
	return t.arr
}

// push appends the current buffered string to the array.
// If the buffer is empty, then do not push it.
func (t *tokenSplitter) push(kind string) {
	// t.buf = strings.Replace(t.buf, " ", "", -1)
	if t.buf == "" {
		return
	}
	if kind == JackGrammar.UNKNOWN {
		kind = resolveUnknownToken(t.buf)
	}
	tok := Token.NewToken(kind, t.buf)
	t.arr = append(t.arr, tok)
	t.buf = ""
}

// add2buff adds a character to the buffer.
// input will be treated as ASCII.
// Control characters will be ignored (including newlines).
// Spaces are ignored unless we are currently tokenizing a constant string.
func (t *tokenSplitter) add2buf(r rune) {
	if r < 32 {
		return
	}
	if r == 32 && !t.string_mode {
		return
	}
	t.buf += string(r)
}

// ResolveUnknownToken accepts the content of a single token,
// and Checks (in this order) if the given token is a...
// 		- Symbol
// 		- Keyword
//		- Constant integer.
// 		- Valid identifier
// Otherwise, there is something wrong with the content, and it is
// becomes an INVALID token.
//
// We assume that constant strings have been taken care of,
// because they can't be determined by just looking at the content,
// since the tokenizer process removes the quotation marks.
//
func resolveUnknownToken(content string) string {

	// preliminary check for empty string.
	if len(content) <= 0 {
		return JackGrammar.INVALID
	}

	// Symbol?
	if len(content) == 1 {
		for _, token := range JackGrammar.LIST_OF_SYMBOLS {
			if content == string(token) {
				return JackGrammar.SYMBOL
			}
		}
	}

	// Keyword?
	for _, keyword := range JackGrammar.LIST_OF_KEYWORDS {
		if content == keyword {
			return JackGrammar.KEYWORD
		}
	}

	// Constant Integer?
	_, validInteger := strconv.Atoi(content)
	if validInteger == nil {
		return JackGrammar.INT_CONST
	}

	// Valid Identifier?
	// check for a leading digit, which would be invalid.
	if !unicode.IsLetter(rune(content[0])) {
		return JackGrammar.INVALID
	}
	// check each rune to se if each is either a letter, digit, or underscore.
	for _, r := range content {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			return JackGrammar.INVALID
		}
	}
	return JackGrammar.IDENTIFIER
}
