// Arithmetic Stack Operations
//
package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (

	// Represents any Arithemtic command.
	C_ARITHMETIC = iota

	// Basic Arithmetic
	C_ADD = iota // Addition
	C_SUB = iota // Subtraction

	// Arithmetic Comparisons
	C_LT = iota // Less than
	C_GT = iota // Greater than
	C_EQ = iota // Equal to

	// Bit-wise Operators
	C_NOT = iota
	C_AND = iota
	C_OR  = iota

	// Memory Commands
	C_PUSH = iota
	C_POP  = iota
	C_THIS = iota

	//Program Control commmands
	C_LABEL    = iota
	C_IF       = iota
	C_GOTO     = iota
	C_FUNCTION = iota
	C_RETURN   = iota
	C_CALL     = iota
)

type Command struct {
	Kind, Arg2 int
	Arg1       string
}

func ParseLine(line string) (string, error) {

	// Parse the line received from the file scanner
	x, err := getCommandFromLine(line)
	if len(x) < 1 {
		return "", err
	}
	return x, err
}

func getCommandFromLine(line string) (string, error) {

	// Remove Comments
	// Looks for the first occurence of the comment: //
	// If it's found, it removes all characters after it.
	comment_index := strings.Index(line, "//")
	if comment_index != -1 {
		line = line[:comment_index]
	}

	// converts all letters to lowercase
	line = strings.ToLower(line)

	// Fields
	// This creates the array of strings, split by whitespace.
	fields := strings.Fields(line)

	// If the fields are empty, then it is most likely a comment
	// or a blank line.  This is normal, so it's not an error.
	if len(fields) < 1 {
		return "", nil
	}

	// If there are more than 3 fields, then this is not a
	// valid command for the hack VM.  Report the error and exit.
	if len(fields) > 3 {
		return "", fmt.Errorf("Too Many Arguments: %v", fields)
	}

	// Create an command using the fields.
	cmd, err := getCommandFromFields(fields)
	if err != nil {
		return "", err
	}

	// Converts commands to assembly code based on their Kind.
	// Invokes Codewriter.go
	switch cmd.Kind {

	case C_ARITHMETIC:
		return WriteArithmetic(cmd.Arg1), nil

	case C_POP, C_PUSH:
		return cmd.WritePushPop(), nil

	case C_LABEL, C_IF, C_GOTO, C_FUNCTION, C_RETURN, C_CALL:
		return cmd.WriteProgramControl()
	}

	// Convert into a string.
	s := fmt.Sprintln(*cmd)

	return s, nil
}

// getCommandFromFields takes an array of strings that has already
// been split and trimmed of whitespace.  It creates and returns
// a new command object, containing all of the information required
// by the codewriter to make assembly code.
func getCommandFromFields(fields []string) (*Command, error) {
	cmd := Command{}
	var err error
	switch fields[0] {

	// Arithemtic commands are simple 1 word commands.
	// That word is placed in the arg1.
	// arg2 is not used.
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		cmd.Kind = C_ARITHMETIC
		cmd.Arg1 = fields[0]

	case "label":
		cmd.Kind = C_LABEL
		cmd.Arg1 = fields[1]

	case "goto":
		cmd.Kind = C_GOTO
		cmd.Arg1 = fields[1]

	case "if-goto":
		cmd.Kind = C_IF
		if len(fields) != 1 {

		}
		cmd.Arg1 = fields[1]

	case "function":
		if len(fields) != 3 {
			err = errWrongArguments("function", 3, len(fields)-1)
			break
		}
		cmd.Kind = C_FUNCTION
		cmd.Arg1 = fields[1]
		cmd.Arg2, err = strconv.Atoi(fields[2])

	case "return":
		cmd.Kind = C_RETURN

	case "call":
		if len(fields) != 3 {
			err = errWrongArguments("call", 3, len(fields)-1)
			break
		}
		cmd.Kind = C_CALL
		cmd.Arg1 = fields[1]
		cmd.Arg2, err = strconv.Atoi(fields[2])

	// Push and Pop are memory access commands.
	// Values can be pushed/popped from/to different places.
	// Arg1 is a string defining the location in memory.
	// Arg2 is a integer defining offset from location.
	case "push":
		cmd.Kind = C_PUSH
		err = cmd.addPushPopArgs(fields)

	case "pop":
		cmd.Kind = C_POP
		err = cmd.addPushPopArgs(fields)

	// If the command doesn't fit any of the cases,
	// then it is not a command the parser understands, which
	// means we can't compile it. Return an error.
	default:
		return &cmd, fmt.Errorf("Invalid command kind: %v", fields)
	}
	return &cmd, err
}

// addPushPopArgs takes the fields of a source command, and
// places in the arguments into the command object accordingly
// Includes error checking for badly formed commands.
func (cmd *Command) addPushPopArgs(fields []string) error {

	// Confirm there are the correct number of arguments.
	if len(fields) != 3 {
		return errWrongArguments("push/pop", 2, len(fields)-1)
	}

	// The first argument is a string.  Pass it directly.
	cmd.Arg1 = fields[1]

	// The second argument is an integer, so it must be
	// converted from ASCII before it can saved.
	// Return any errors in the process of doing this.
	i, err := strconv.Atoi(fields[2])
	if err != nil {
		return err
	}
	cmd.Arg2 = i

	// No errors!  Arguments have succesfully been added!
	return nil
}

func errWrongArguments(name string, expects, got int) error {
	return fmt.Errorf("Invalid number of arguments for %s command. Expects:(%d), Got:(%d).", name, expects, got)
}
