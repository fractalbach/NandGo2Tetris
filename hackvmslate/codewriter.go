package main

import (
	"strconv"
)

// Location of Stuff in the Stack
//
// 			| ... |
// 			|  x  |
// 			|  y  |
//          |     |  <- Stack Pointer (SP)
//
// Arithmetic is in the form of x _ y,
// where _ is an operator (like +, -, <, >, =)
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

// Bit-wise NOT on the top element on the stack.
// Only affects one element.
// Does not pop the stack.
const S_NOT = `
@SP
A=M-1
M=!M
`

// Negation of the top element on the stack.
// Only affects one element.
// Does not pop the stack.
const S_NEG = `
@SP
A=M-1
M=-M
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

const S_ = ``

// INEQUALITY requires similar parameters to the ALU_ARITHMETIC,
// but also requires a name for the "checkpoint" variable.  It does not
// determine this on it's own, because it may cause conflicts with other
// jump variables depending on the context.
//
// For the assembly parameter, write one of the following:
//
// 		JEQ   Jump if equal to
// 		JLT   Jump if less than
//      JGT   jump if greater than
//
func inequality(comment string, assembly string, location_count int) string {
	count := strconv.Itoa(location_count)
	return comment + `
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@INEQLOCATION` + count + "\nD;" + assembly + `
M=0
(INEQLOCATION` + count + `)
`
}

func S_JEQ(location int) string {
	return inequality("// true if x = y, else false", "JEQ", location)
}

func S_JLT(location int) string {
	return inequality("// true if x < y, else false", "JLT", location)
}

func S_JGT(location int) string {
	return inequality("// true if x > y, else false", "JGT", location)
}

// next accepts a pointer to an integer, increments it,
// and then returns its new value.  The incremented value
// will be saved in it's original variable.
func next(count *int) int {
	*count++
	return *count
}
