// Arithmetic Stack Operations
//
package main

const (

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
)

type Command struct {
	Kind int
}
