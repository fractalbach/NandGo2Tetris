package Token

import (
	"fmt"
	"strings"
)

var xml_symbol_replacements_table = map[string]string{
	">":  "&gt;",
	"<":  "&lt;",
	"\"": "&quot;",
}

// Token has a kind and content, the different kinds are
// keyword, symbol, identifier, int_const, string_const"
type Token interface {
	Kind() string
	Content() string
}

type token struct {
	kind, content string
}

func NewToken(kind, content string) Token {
	return &token{
		kind:    kind,
		content: content,
	}
}

func (t *token) Kind() string {
	return t.kind
}

func (t *token) Content() string {
	return t.content
}

// Default string representation for the Token is XML.
// Any symbols literals will be replaced with their XML equivalent.
func (t *token) String() string {
	s := t.content
	// replace ampersands first so that it doesn't overwrite any of the
	// other changes.  Then go through the table.
	s = strings.Replace(s, "&", "&amp;", -1)
	for old, new := range xml_symbol_replacements_table {
		s = strings.Replace(s, old, new, -1)
	}
	return fmt.Sprintf("<%s> %s </%s>", t.kind, s, t.kind)
	// return ("<" + t.kind + "> " + s + " </" + t.kind + ">")
}

// The default Go String representation is most likely for debugging.
// This will be formatted and include both the token's kind and content.
func (t *token) GoString() string {
	return fmt.Sprintf("%-11s %s", t.kind, t.content)
}
