/*
package static implements static variables for the hackVM translator.

Static Segment

Static Segments are shared among all functions within a .vm file.
Most segments, like argument and local, have their own instance
for each function.

In the VM code, a static segment is created using a command like
"pop static 5".  In assembly this will be rewritten as a symbol
in the form of "filename.index", where index=5 from our example.

The assembler handles the convertion of symbols into memory,
so the VM language need not be concerned with actual registers.
*/
package static

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/stack"
)

// usage: (n, filename, counter, PUSHD)
const s_push_static = `// push static %d
@%s.%d
D=M
%v
`

// usage: (n, POPD, filename, counter)
const s_pop_static = `// pop static %d
%v
@%s.%d
M=D
`

// pushes the value at static[index] to the stack.
// the symbol will be named @filename.index
func push(filename string, index int) string {
	return fmt.Sprintf(s_push_static, index, filename, index, stack.PUSHD)
}

// pops the stack, and places the value into static[index].
// the symbol will be named @filename.index
func pop(filename string, index int) string {
	return fmt.Sprintf(s_pop_static, index, stack.POPD, filename, index)
}
