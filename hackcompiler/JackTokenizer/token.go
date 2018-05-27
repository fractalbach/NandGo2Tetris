package JackTokenizer

import "fmt"

type token struct {
	kind, content string
}

func (t *token) Kind() string {
	return t.kind
}

func (t *token) Content() string {
	return t.content
}

func (t *token) String() string {
	return ("<" + t.kind + ">" + t.content + "</" + t.kind + ">")
}

func (t *token) GoString() string {
	return fmt.Sprintf("%-11s %s", t.kind, t.content)
}
