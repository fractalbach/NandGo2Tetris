// package pointer has assembly translations for pointers.
//
// The "pointer" segment refers to THIS and THAT, and
// will directly change those values.  For those, use
// pointer.Pop(index) and pointer.Push(index) directly.
//
// Use "PopThrough()" and "PushThrough()" for changing
// values that are referenced by pointers.
//
package pointer

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/stack"
)

func Pop(n int) string {
	var x string
	switch n {
	case 0:
		x = s_pop_pointer_0
	case 1:
		x = s_pop_pointer_1
	}
	return fmt.Sprintf(x, stack.POPD)
}

func Push(n int) string {
	var x string
	switch n {
	case 0:
		x = s_push_pointer_0
	case 1:
		x = s_push_pointer_1
	}
	return fmt.Sprintf(x, stack.PUSHD)
}

const s_push_pointer_0 = `//push pointer 1
@THIS
D=M
%v
`

const s_push_pointer_1 = `//push pointer 0
@THAT
D=M
%v
`

const s_pop_pointer_0 = `// pop pointer 0
%v
@THIS
M=D
`
const s_pop_pointer_1 = `// pop pointer 1
%v
@THAT
M=D
`

// PopThrough accepts a segment string and index int.
// Returns assembly.
func PopThrough(s string, n int) string {
	return fmt.Sprintf(pop_thru_pointer, s, n, n, s, stack.POPD)
}

// PushThrough accepts a segment string and index int.
// Returns assembly.
func PushThrough(s string, n int) string {
	return fmt.Sprintf(push_thru_pointer, s, n, n, s, stack.PUSHD)
}

// Special Memory Access Pop Command.
// 		Variable Order:
// 		(string, int, int, string)
const pop_thru_pointer = `// pop %s %d
@%d
D=A
@%s
D=D+M
@R13
M=D
%v
@R13
A=M
M=D
`

// Special Push Example.
// 		Variable Order:
// 		(string, int, int, string)
const push_thru_pointer = `// push %s %d
@%d
D=A
@%s
A=M+D
D=M
%v
`
