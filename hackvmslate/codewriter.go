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
	return_counter   = 0
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

func (cmd *Command) WritePushPop() (string, error) {

	// Constant push commands don't need to go through the
	// whole process
	if cmd.Kind == C_PUSH && cmd.Arg1 == "constant" {
		return constant_push(cmd.Arg2), nil
	}

	// retrieve the assembly equivalent of the given segment.
	// if it can't be found, then you've been given a bad segment.
	segment, ok := segment_map[cmd.Arg1]
	if !ok {
		return "", fmt.Errorf("Push/Pop: Unknown argument: (%v).", cmd.Arg1)
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
func push(s string, n int) (string, error) {
	switch s {
	case "TMP":
		return temp.Push(n), nil
	case "pointer":
		return pointer.Push(n), nil
	case "static":
		return static.Push(current_filename, n), nil
	case "LCL", "ARG", "THIS", "THAT":
		return pointer.PushThrough(s, n), nil
	}
	return "", fmt.Errorf("Push: Unknown first argument: (%v)", s)
}

// pop accepts the segment and index from a pop command.
// returns string containing assembly instructions
// that will pop a value from the stack, and place it
// somewhere in memory based on the given segment and index.
func pop(s string, n int) (string, error) {
	switch s {
	case "TMP":
		return temp.Pop(n), nil

	case "pointer":
		return pointer.Pop(n), nil

	case "static":
		return static.Pop(current_filename, n), nil

	case "LCL", "ARG", "THIS", "THAT":
		return pointer.PopThrough(s, n), nil
	}
	return "", fmt.Errorf("Pop: Unknown first argument: (%v)", s)
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
		if cmd.Arg2 < 0 {
			return "", fmt.Errorf("Can't have a function with %d local variables! That doesn't make sense!", cmd.Arg2)
		}
		return WriteFunction(cmd.Arg1, cmd.Arg2), nil

	case C_RETURN:
		return s_return, nil

	case C_CALL:
		if cmd.Arg2 < 0 {
			return "", fmt.Errorf("Can't call a function with %d arguments! That doesn't make sense!", cmd.Arg2)
		}
		return WriteCall(cmd.Arg1, cmd.Arg2), nil
	}
	panic("This command should not be writing a program control.")
}

func errNotImplemented(s string) error {
	return fmt.Errorf("%s hasn't been implemented yet.", s)
}

func (cmd *Command) WriteFunction() (string, error) {
	return "", nil
}

// at the beginning, create next unique id, and save it.
// that will also be used to create the return label.
// Format Args
// - name
// - nArgs
// - return label
// - nArgs
// - name
// - return label
const s_call = ` // call %s %d 
@RETURN.%s
D=A
` + stack.PUSHD + `
@LCL
D=M
` + stack.PUSHD + `
@ARG
D=M
` + stack.PUSHD + `
@THIS
D=M
` + stack.PUSHD + `
@%d
D=A
@SP
D=D-A
@ARG
M=D
@SP
D=M
@LCL
M=D
@FUNCTION.%s
0; JMP
(RETURN.%s)
`

func WriteCall(name string, nArgs int) string {
	id := next(&return_counter)
	ret := name + "." + string(id)
	return fmt.Sprintf(s_call, name, nArgs, ret, nArgs, name, ret)
}

func WriteFunction(name string, nLocal int) string {
	many_push := ""
	if nLocal < 0 {
		return ""
	}
	for i := 0; i < nLocal; i++ {
		many_push += s_push_0
	}
	return fmt.Sprintf("(FUNCTION.%s)\n", name) + many_push
}

const s_push_0 = `@SP
M=M+1
A=M-1
M=0
`

const s_return = `// return
` + stack.POPD + `
@ARG 	// places return value in the right spot.
A=M
M=D
@ARG	// restores stack pointer.  SP <- ARG + 1
D=M
@SP
M=D+1
@LCL 	// init frame pointer
D=M
@FRAME
M=D
` + s_POP_FRAME + `
@THAT 	// restore that.
M=D
` + s_POP_FRAME + `
@THIS	// restore this.
M=D
` + s_POP_FRAME + `
@ARG	// restore arg.
M=D
` + s_POP_FRAME + `
@LCL 	// restore local.
M=D
` + s_POP_FRAME + `
A=D 	// jumps to return address.
0;JMP
`

const s_POP_FRAME = `@FRAME
AM=M-1
D=M`
