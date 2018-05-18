package main

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/control"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/pointer"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/stack"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/static"
	"github.com/fractalbach/nandGo2tetris/hackvmslate/codewriter/temp"
)

// Segment_map contains a list of acceptable segments,
// and maps certain segments to the symbol used in assembly.
var segment_map = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
	"temp":     "TMP",
	"pointer":  "pointer",
	"static":   "static",
}

// Counters is used internally for branching within the assembly code itself.
// They are used sparingly, and are always called when wrapped around
// a next(&counter) function.
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
		return stack.ADD
	case "sub":
		return stack.SUB
	case "not":
		return stack.NOT
	case "and":
		return stack.AND
	case "or":
		return stack.OR
	case "neg":
		return stack.NEG

	// When comparing inequalities, we need to include jumps.
	// in order to avoid conflicts of jumps, increment
	// the location counter by 1 prior to return the string.

	case "eq":
		return stack.JEQ(next(&location_counter))
	case "gt":
		return stack.JGT(next(&location_counter))
	case "lt":
		return stack.JLT(next(&location_counter))
	}

	// If an invalid command has been given, then it is due to
	// an error in the compiler itself, and not the user.
	// Panic to inform creator of the compiler (that's you!)
	// that something is wrong.
	panic("ERROR: INVALID ARITHMETIC COMMAND GIVEN.")
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

	// Decide between Push or Pop.
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

// push accepts a segment and index from a push command.
// returns the string containing assembly instructions,
// A value will be copied from memory based on the given
// segment and index.  That value is then pushed to the
// global stack.
func push(s string, n int) string {
	switch s {
	case "TMP":
		return temp.Push(n)
	case "pointer":
		return pointer.Push(n)
	case "static":
		return static.Push(current_filename, n)
	case "LCL", "ARG":
		return pointer.PushThrough(s, n)
	}
	return "// [ERROR] "
}

// pop accepts the segment and index from a pop command.
// returns string containing assembly instructions
// that will pop a value from the stack, and place it
// somewhere in memory based on the given segment and index.
func pop(s string, n int) string {
	switch s {
	case "TMP":
		return temp.Pop(n)

	case "pointer":
		return pointer.Pop(n)

	case "static":
		return static.Pop(current_filename, n)

	case "LCL", "ARG":
		return pointer.PopThrough(s, n)
	}
	return "// [ERROR]"
}

// The end of program is an infinite loop that will unconditionally
// jump back to itself.  This only happens once, and is only found
// at the very end of the .asm file.
const s_end_program = ` // End of Program.
(END)
@END
0; JMP
`

// The constant push is a simple integer value.
const s_constant_push = `// push constant %d
@%d
D=A
%v
`

// pushes a constant value to the stack.
func constant_push(n int) string {
	return fmt.Sprintf(s_constant_push, n, n, stack.PUSHD)
}

// next accepts a pointer to an integer, increments it,
// and then returns its new value.  The incremented value
// will be saved in it's original variable.
func next(count *int) int {
	*count++
	return *count
}

const s_init = `// bootstrap code
// ---------------
// set SP = 256
@256
D=A
@SP
M=D

// set LCL = 300
// MIGHT NOT NEED THIS PART
//
@300
D=A
@LCL
M=D

// Start executing sys.init
// call Sys.init
// ---------------
`

func WriteInit() string {
	return s_init
}

func (cmd *Command) WriteProgramControl() (string, error) {
	switch cmd.Kind {
	case C_LABEL:
		return control.WriteLabel(cmd.Arg1), nil
	case C_IF:
		return control.WriteIf(cmd.Arg1), nil
	case C_GOTO:
		return control.WriteGoto(cmd.Arg1), nil
	case C_FUNCTION:
		return "", errNotImplemented("function")
	case C_RETURN:
		return "", errNotImplemented("return")
	case C_CALL:
		return "", errNotImplemented("call")
	}
	panic("This command should not be writing a program control.")
}

func errNotImplemented(s string) error {
	return fmt.Errorf("%s hasn't been implemented yet.", s)
}
