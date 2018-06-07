// package stack contains assembly instructions for stack based VM commands.
//
// Included with the stack package:
// 		- basic arithmetic
// 		- bitwise logic
//		- inequality logic
//
package stack

import (
	"fmt"
)

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
D=M`

// Bit-wise NOT on the top element on the stack.
// Only affects one element.
// Does not pop the stack.
const NOT = `// bitwise NOT
@SP
A=M-1
M=!M
`

// Negation of the top element on the stack.
// Only affects one element.
// Does not pop the stack.
const NEG = `// negation
@SP
A=M-1
M=-M
`

var (
	ADD = alu_arithmetic("// add", "M=M+D")
	SUB = alu_arithmetic("// subtract", "M=M-D")
	AND = alu_arithmetic("// AND", "M=M&D")
	OR  = alu_arithmetic("// OR", "M=M|D")
)

func JEQ(location int) string {
	return inequality("// true if x = y, else false", "JEQ", location)
}

func JLT(location int) string {
	return inequality("// true if x < y, else false", "JLT", location)
}

func JGT(location int) string {
	return inequality("// true if x > y, else false", "JGT", location)
}

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
const alu_arith_string = `%s
%s
A=A-1
%s
`

// arithmetic string that supports add, subtract, AND, OR.
// Usage: (POPD, comment, assembly)
func alu_arithmetic(comment string, assembly string) string {
	return fmt.Sprintf(alu_arith_string, comment, POPD, assembly)
}

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
const s_ineq = `%s
%v 
A=A-1
D=M-D
M=-1
@LOCATION%d
D; %s
@SP
A=M-1
M=0
(LOCATION%d)
`

func inequality(comment string, assembly string, counter int) string {
	return fmt.Sprintf(s_ineq, comment, POPD, counter, assembly, counter)
}
