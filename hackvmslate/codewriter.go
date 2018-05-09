package main

// TODO LIST
//
// Multiply
// Divide?? (might not need to do this)
//
//

// Pushes the value in the D register onto the stack.
const S_PUSH = `// push d
@SP
A=M
M=D
@SP
M=M+1`

// Pops the top element from the stack and places it in register D.
const S_POP = `// pop d
@SP
AM=M-1
D=M`

// Negates the top element on the stack.  Doesn't change the stack in any other way.
const S_NOT = `
@SP
A=M
M=!M
`

// ALU_arithemtic returns the assembly instructions for stack arithemtic.
// Supports the basic ALU instructions: add, subtract, and, or.
// Each of these has a similar structure in assembly.
//
// 		1. @SP 		goto register holding the stack pointer
// 		2. AM=M-1   move top of stack down by 1, then goto top of stack.
// 		2. D=M 		Save the element at the top of the stack in register D.
// 		4. A=A-1    Move down 1 element in the stack.
// 		5. M=M?A	Evaluate arithemtic, replace "?" with an operator.
// 					and save the result in the current register.
//
// This is only meaningful if you have 2 values in the stack.
// Those 2 values are used together in the arithmetic to form a new value.
// This new value is also saved on the stack.
// So the net change in size of stack is -1.
//
// Given a different computer architecture, or a different chipset,
// there could be a wider range of commands that could be translated this way.
// However, the 16-bit ALU we use in the HACK computer can only support
// certain basic operations.
//
func ALU_arithmetic(comment string, assembly string) string {
	return S_POP + comment + "\nA=A-1\n" + assembly + "\n"
}

var (
	S_ADD = ALU_arithmetic("// add", "M=M+D")
	S_SUB = ALU_arithmetic("// subtract", "M=M-D")
	S_AND = ALU_arithmetic("// AND", "M=M&D")
	S_OR  = ALU_arithmetic("// OR", "M=M|D")
)

// empty string.
const S_ = ``

/*
``
@SP
A=M
M=D
@SP
M=M+1
*/
