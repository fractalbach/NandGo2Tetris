// Arithmetic Stack Operations
//
package main

const (

	// Arithmetic Stack Operations

	// Less than
	C_LT = iota

	// Greater than
	C_GT = iota

	// Equal to
	C_EQ = iota

	// Addition
	C_ADD = iota

	// Subtraction
	C_SUB = iota

	// Multiplication
	C_MUL = iota

	// Memory Commands

	C_PUSH = iota
	C_POP  = iota
	C_THIS = iota
)

type Command struct {
	Kind int
}
