package control

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/stack"
)

// s_label args:
// 1. label name
const s_label = `(LABEL.%s)`

// s_if args:
// 1. label name
// 2. pop stack
// 3. label name
const s_if = `// if-goto %s
%v
@LABEL.%s
D; JNE
`

// s_goto args:
// 1. label name
// 2. label name
const s_goto = `// goto %s
@LABEL.%s
0; JMP
`

func WriteLabel(label string) string {
	return fmt.Sprintf(s_label, label)
}

func WriteIf(label string) string {
	return fmt.Sprintf(s_if, label, stack.POPD, label)
}

func WriteGoto(label string) string {
	return fmt.Sprintf(s_goto, label, label)
}
