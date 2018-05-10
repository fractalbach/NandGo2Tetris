package main

import (
	"strconv"
)

var (
	location_counter = 0
)

// WriteArithemtic accepts an arithemtic command and returns the
// assembly instructions as a string.  These commands come from
// the parser, which has already interprted the source code.
// The string generated from this command can be written into the
// output .asm file.
//
// Location of Stuff in the Stack
//
// 			| ... |
// 			|  x  |
// 			|  y  |
// 			|     |  <- Stack Pointer (SP)
//
// Arithmetic is in the form of x _ y,
// where _ is an operator (like +, -, <, >, =)
//
func WriteArithmetic(command int) string {
	switch command {

	// Basic arithemtic and bit-wise commands are straight forward,
	// and always have the same representation in assembly.

	case C_ADD:
		return S_ADD
	case C_SUB:
		return S_SUB
	case C_NOT:
		return S_NOT
	case C_AND:
		return S_AND
	case C_OR:
		return S_OR

	// When comparing inequalities, we need to include jumps.
	// in order to avoid conflicts of jumps, increment
	// the location counter by 1 prior to return the string.

	case C_EQ:
		return S_JEQ(next(&location_counter))
	case C_GT:
		return S_JGT(next(&location_counter))
	case C_LT:
		return S_JLT(next(&location_counter))
	}

	// If an invalid command has been given, then it is due to
	// an error in the compiler itself, and not the user.
	// Panic to inform creator of the compiler (that's you!)
	// that something is wrong.
	panic("ERROR: INVALID ARITHMETIC COMMAND GIVEN.")
}

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
// 		5. M=M_A	Evaluate arithemtic, replace "_" with an operator.
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
func alu_arithmetic(comment string, assembly string) string {
	return S_POP + comment + "\nA=A-1\n" + assembly + "\n"
}

var (
	S_ADD = alu_arithmetic("// add", "M=M+D")
	S_SUB = alu_arithmetic("// subtract", "M=M-D")
	S_AND = alu_arithmetic("// AND", "M=M&D")
	S_OR  = alu_arithmetic("// OR", "M=M|D")
)

// INEQUALITY requires similar parameters to the alu_arithmetic,
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
