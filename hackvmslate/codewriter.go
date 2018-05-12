package main

import (
	"fmt"
	// "github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/stack"
)

var (
	location_counter = 0
	static_counter   = 0
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
func WriteArithmetic(command string) string {
	switch command {

	// Basic arithemtic and bit-wise commands are straight forward,
	// and always have the same representation in assembly.

	case "add":
		return S_ADD
	case "sub":
		return S_SUB
	case "not":
		return S_NOT
	case "and":
		return S_AND
	case "or":
		return S_OR
	case "neg":
		return S_NEG

	// When comparing inequalities, we need to include jumps.
	// in order to avoid conflicts of jumps, increment
	// the location counter by 1 prior to return the string.

	case "eq":
		return S_JEQ(next(&location_counter))
	case "gt":
		return S_JGT(next(&location_counter))
	case "lt":
		return S_JLT(next(&location_counter))
	}

	// If an invalid command has been given, then it is due to
	// an error in the compiler itself, and not the user.
	// Panic to inform creator of the compiler (that's you!)
	// that something is wrong.
	panic("ERROR: INVALID ARITHMETIC COMMAND GIVEN.")
}

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

// Bit-wise NOT on the top element on the stack.
// Only affects one element.
// Does not pop the stack.
const S_NOT = `// bitwise NOT
@SP
A=M-1
M=!M
`

// Negation of the top element on the stack.
// Only affects one element.
// Does not pop the stack.
const S_NEG = `// negation
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
const S_INEQ = `%s
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
	return fmt.Sprintf(S_INEQ, comment, POPD, counter, assembly, counter)
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

var segment_map = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
	"temp":     "TMP",
	"pointer":  "pointer",
}

func (cmd *Command) WritePushPop() string {

	// Constant push commands don't need to go through the
	// whole process
	if cmd.Kind == C_PUSH && cmd.Arg1 == "constant" {
		return constant_push(cmd.Arg2)
	}

	// retrieve the assembly equivalent of the given segment.
	// if it can't be found, then you've been given a bad segment.
	segment, ok := segment_map[cmd.Arg1]
	if !ok {
		return ""
	}

	switch cmd.Kind {
	case C_POP:
		return pop(segment, cmd.Arg2)
	case C_PUSH:
		return push(segment, cmd.Arg2)
	default:
		panic("Needs to be a Push or Pop command.")
	}

	panic("Invalid push/pop command.")
}

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

func push(s string, n int) string {
	if s == "TMP" {
		return pushTemp(n)
	}
	if s == "pointer" {
		return pushPointer(n)
	}
	return fmt.Sprintf(push_thru_pointer, s, n, n, s, PUSHD)
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

func pop(s string, n int) string {
	if s == "TMP" {
		return popTemp(n)
	}
	if s == "pointer" {
		return popPointer(n)
	}
	return fmt.Sprintf(pop_thru_pointer, s, n, n, s, POPD)
}

const s_constant_push = `// push constant %d
@%d
D=A
%v
`

func constant_push(n int) string {
	return fmt.Sprintf(s_constant_push, n, n, PUSHD)
}

const s_end_program = ` // End of Program.
(END)
@END
0; JMP
`

func pushTemp(n int) string {
	return fmt.Sprintf(s_push_temp, n, n, PUSHD)
}

const s_push_temp = `// push temp %d
@%d
D=A
@5
A=A+D
D=M
%v
`

func popTemp(n int) string {
	return fmt.Sprintf(s_pop_temp, n, n, POPD)
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

func popPointer(n int) string {
	var x string
	switch n {
	case 0:
		x = s_pop_pointer_0
	case 1:
		x = s_pop_pointer_1
	}
	return fmt.Sprintf(x, POPD)
}

func pushPointer(n int) string {
	var x string
	switch n {
	case 0:
		x = s_push_pointer_0
	case 1:
		x = s_push_pointer_1
	}
	return fmt.Sprintf(x, PUSHD)
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

// usage: (n, filename, counter, PUSHD)
const s_push_static = `// push static %d
@%s.%s
D=M
%v
`

// usage: (n, POPD, filename, counter)
const s_pop_static = `// pop static %d
%v
@%s.%s
M=D
`

func popStatic(filename string, n int) {

}
