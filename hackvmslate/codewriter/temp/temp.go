package temp

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/stack"
)

func Push(n int) string {
	return fmt.Sprintf(s_push_temp, n, n, stack.PUSHD)
}

const s_push_temp = `// push temp %d
@%d
D=A
@5
A=A+D
D=M
%v
`

func Pop(n int) string {
	return fmt.Sprintf(s_pop_temp, n, n, stack.POPD)
}

// Usage:
// (n, n, POP d, )
const s_pop_temp = `// pop temp %d
@%d
D=A
@5
D=A+D
@R13
M=D
%v
@R13
A=M
M=D
`
