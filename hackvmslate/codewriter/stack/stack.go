// package stack contains assembly instructions for stack based VM commands.
package stack

// Pushes the value in the D register onto the stack.
// Common building block for other commands.
const PUSHD = `@SP // push d
M=M+1
A=M-1
M=D`

// Pops the top element from the stack and places it in register D.
// Common building block for other commands.
const POPD = `@SP // pop d
AM=M-1
D=M
M=0`
