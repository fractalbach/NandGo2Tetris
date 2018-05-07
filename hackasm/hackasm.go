// package hackasm is the Hack assembler.
//
// This assembler takes a functional approach to parsing,
// instead of an object-oriented one.  There are no declared
// data structures or types.  There are many functions,
// most of which do operations on strings or byte arrays.
//
// An assembly file is input into this program. First, all of
// its whitespace and comments are removed.  Next, it is
// split into an array of strings, line-by-line.  Each
// element in this array is treated as a separate "command".
//
//
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	_A_COMMAND       = 1
	_C_COMMAND       = 2
	_L_COMMAND       = 3
	_INVALID_COMMAND = 4
	BYTE_SPACE       = 32
	BYTE_SLASH       = 47
)

var predefined_computations = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"M":   "1110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"!M":  "1110001",
	"-D":  "0001101",
	"-A":  "0110011",
	"-M":  "1110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"M+1": "1110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"M-1": "1110010",
	"D+A": "0000010",
	"D+M": "1000010",
	"D-A": "0010011",
	"D-M": "1010011",
	"A-D": "0000111",
	"M-D": "1000111",
	"D&A": "0000000",
	"D&M": "1000000",
	"D|A": "0010101",
	"D|M": "1010101",
}

// Jumps are used in C-commands.  They are often used to create
// "If x, goto y" statements.  Jumps in the hack assembly
// language are also required for ending the program.
// All of the jump commands need a "destination" variable,
// which is compared to 0.
//
// 		GT (greater than)
//		EQ (equal to)
//		GE (greater or equal)
// 		LT (less than)
// 		NE (not equal to)
// 		LE (less than or equal to)
// 	    JMP (unconditional jump).
//
var predefined_jumps = map[string]string{
	"null": "000",
	"JGT":  "001",
	"JEQ":  "010",
	"JGE":  "011",
	"JLT":  "100",
	"JNE":  "101",
	"JLE":  "110",
	"JMP":  "111",
}

// SymbolTable is used for A-commands.  They correspond
// to locations of registers in memory.
//
// There are a handful of default symbols that are
// predefined.  For example, The first 16 registers begin
// with the letter "R", and the screen and keyboard I/O
// have their first register's predefined.
//
// When a custom variable is used in the assembly language,
// it will be added to this table by the assembler during
// its initial scan.  During a second scan, the variables
// will be resolved to their numeric equivalents.
//
// Custom variables start at index 16, which is the 17th
// register in the memory.
//
var SymbolTable = map[string]int{
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"SCREEN": 16384,
	"KBD":    24576,
}

var output_filename string
var verbose bool
var output_table_only bool

func main() {

	// Handle command line flags such as input/output
	// filenames, and verbosity.
	flag.StringVar(&output_filename,
		"o", "", "Output File Location")
	flag.BoolVar(&verbose, "v", false, "Verbose mode.")
	flag.BoolVar(&output_table_only, "t", false,
		"Print Symbol Table Only.")
	flag.Parse()
	input_filename := flag.Arg(0)

	// Loads the Data and cleans it up.
	data := EasyLoad(input_filename)
	RemoveCommentsFromArray(data)
	data = RemoveEmptyLines(data)

	// First Pass: Add Unknown Symbols to the Table.
	AddAllSymbols(data)

	// If the t - flag was passed to the command line,
	// then only output the symbol table that was just created.
	// and exit the program.
	if output_table_only {
		for i, v := range SymbolTable {
			fmt.Printf("%6v %v\n", v, i)
		}
		return
	}

	// Use a special output format for verbose debugging.
	// Does not create a file.
	// Instead, output is printed to the stdout.
	if verbose {
		FancyDisplayData(data)
		fmt.Println("Done!")
		return
	}

	// If no output file path has been given,
	// Resolve paths and create a suitable path.
	if output_filename == "" {
		output_filename = ResolveOutputPath(input_filename)
		fmt.Println(output_filename)
	}

	// Send the data to the parser, create a file,
	// and write the machine code into the file!
	parsed := ParseData(data)
	err := ioutil.WriteFile(output_filename, []byte(parsed), 0600)
	if err != nil {
		log.Fatal(err)
	}
}

// ResolveOutputPath accepts the input filename and returns
// a suitable path for the output file.  Using this path
// resolver ensures that the .hack file is in the same directory
// as the .asm file.
func ResolveOutputPath(s string) string {

	// Resolve the absolute path of the input file.
	absolute, err := filepath.Abs(s)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the directory of the input file.
	dir := filepath.Dir(absolute)

	// Create a name for the output file.
	name := filepath.Base(s)
	name = strings.Replace(name, ".asm", ".hack", 1)

	// Combines the directory path and the .hack name
	// and return the filepath string.
	return filepath.Join(dir, name)
}

// EasyLoad accepts the byte array from a file,
// and returns an array of command strings.
// All Whitespace and comments will be removed.
// Each string is split line-by-line.
func EasyLoad(filename string) []string {

	// Copy File into a Byte Array
	content_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Remove all of the lowest bytes, except for LF.
	// Replace them with spaces, because spaces will get
	// removed later.
	for i := 0; i < len(content_bytes); i++ {
		x := content_bytes[i]
		if x <= 31 && x != 10 {
			content_bytes[i] = 32
		}
	}

	// Convert the byte array into a string.
	s := string(content_bytes)

	// Remove ALL empty spaces.
	s = strings.Replace(s, " ", "", -1)

	// split the string into an array of strings,
	// separated by newline feeds.
	return strings.Split(s, "\n")
}

// ParseData accepts a full file as a string, handles the
// parsing of each command, and returns a string which contains
// all of the assembled machine code. This is the main function
// that calls all other functions in the program.
func ParseData(file_lines []string) string {
	out := ""
	for _, line := range file_lines {
		if getCommandType(line) == _L_COMMAND {
			continue
		}
		out += fmt.Sprintf("%v\n", ParseCommand(line))
	}
	return out
}

// DisplayData prints out each separate line, adding a number
// to the front: indicating which line it is.
// Mainly for debugging and informative usage.
func FancyDisplayData(file_lines []string) {
	counter := 0
	for _, line := range file_lines {
		t := getCommandType(line)
		if t == _A_COMMAND || t == _C_COMMAND {
			counter++
			fmt.Printf("%3v %1v %16v %v \n", counter,
				getCommandName(line), ParseCommand(line), line)
		}
		if t == _L_COMMAND {
			fmt.Printf("    %1v %16v %v \n",
				getCommandName(line), ParseCommand(line), line)
		}
	}
}

// RemoveComments returns a string without the trailing comment.
// In this parser, a comment is defined as all text after and
// including the string "//".
func RemoveCommentFromLine(s string) string {
	return strings.SplitN(s, "//", 2)[0]
}

// RemoveCommentsFromArray returns an array where each
// line has the comments removed.
func RemoveCommentsFromArray(s []string) []string {
	for i, _ := range s {
		if strings.Contains(s[i], "//") {
			s[i] = RemoveCommentFromLine(s[i])
		}
	}
	return s
}

// RemoveEmptyLines accepts an array, and returns an array
// without the empty lines.
// Use this AFTER you have removed comments and spaces,
// otherwise the line will not be considered "empty".
func RemoveEmptyLines(s []string) []string {
	counter := 0
	out := make([]string, len(s))

	// Iterate through the array, checking each for empty lines
	// If the line is NOT empty, add it to the new array.
	// Save the length of the new array, and use it to
	// return a slice, so that the output is the correct length.
	for i := 0; i < len(s); i++ {
		if s[i] == "" {
			continue
		}
		out[counter] = s[i]
		counter++
	}
	return out[:counter]
}

// ParseCommand accepts a single line in assembly and
// returns that line in binary machine code.
// Symbols should be resolved before invoking this function,
// but it can handle everything else beyond that.
func ParseCommand(s string) string {
	t := getCommandType(s)
	switch t {
	case _A_COMMAND:
		return parseCommandA(s)
	case _C_COMMAND:
		return parseCommandC(s)
	}
	return ""
}

// getCommandType takes a line, and returns the type of
// command that it is.  Returns an error if the command
// cannot be found.
//
// Assumes that you have already removed whitespace and comments,
// because it uses the first character to determine the command
// type.
func getCommandType(s string) int {
	if strings.ContainsRune(s, '@') {
		return _A_COMMAND
	}
	if strings.Contains(s, "=") {
		return _C_COMMAND
	}
	if strings.Contains(s, ";") {
		return _C_COMMAND
	}
	if strings.Contains(s, "(") {
		return _L_COMMAND
	}
	return _INVALID_COMMAND
}

func getCommandName(s string) string {
	switch getCommandType(s) {
	case _A_COMMAND:
		return "A"
	case _C_COMMAND:
		return "C"
	case _L_COMMAND:
		return "L"
	}
	return "?"
}

// parseCommandA accepts an A-command in assembly,
// and returns the command in machine code.
// Since it is an A-command, the first number is "0",
// and the following numbers represent a memory location.
func parseCommandA(s string) string {
	A := ""
	resolved := int64(-1)

	// Get rid of the @ symbol to get the memory reference.
	A = strings.Replace(s, "@", "", 1)

	// First, check to see if the location is an integer.
	// If it can be parsed into an integer, it's resolved.
	parsed, err := strconv.ParseInt(A, 10, 64)
	if err == nil {
		resolved = parsed
	}

	// Next, treat the location as a symbol, and check
	// the symbol table to see if it's there.
	savedInt, ok := SymbolTable[A]
	if ok {
		resolved = int64(savedInt)
	}

	// If neither method resolves the memory location,
	// then exit with an error.  This could be due to an
	// input error, or because the symbol was not correctly
	// added to the Symbol Table.
	if resolved == -1 {
		log.Fatal("Invalid A-command: Cannot resolve:", s)
	}

	// Return a formatted 16-bit binary representation of
	// the integer.  This is the A-instruction machine code.
	b := strconv.FormatInt(resolved, 2)
	return "0" + fmt.Sprintf("%015v", b)
}

// parseCommandC takes a c-instruction and converts it into
// binary machine code.  It can't have any whitespace or
// extra characters or it will exit with an error.
func parseCommandC(s string) string {

	// There are more combinations of the C-instruction,
	// and a better parsing algorithm would address those,
	// but this should do fine for all of the programs
	// we want to assemble.
	var a []string
	comp := "0000000"
	dest := "000"
	jump := "000"

	// handle jump commands:
	// Dest = Comp; Jump
	a = strings.Split(s, ";")
	if len(a) > 1 {
		comp, dest = convertCompAndDest(a[0])
		jump = convertJumps(a[1])
		return "111" + comp + dest + jump
	}

	// handle assignment commands for the ALU:
	// Dest = Comp.
	if strings.Contains(s, "=") {
		comp, dest = convertCompAndDest(s)
	} else {
		log.Fatal("Invalid C-command:")
	}
	return "111" + comp + dest + jump
}

// convertComputations inputs the fields "dest=comp",
// but also accepts the edge case where "comp" is standalone.
// It will return two strings: a computation of length 7,
// and a destination of length 3
func convertCompAndDest(s string) (string, string) {
	comp := "0000000"
	dest := "000"
	if strings.Contains(s, "=") {
		a := strings.Split(s, "=")
		dest = convertDestination(a[0])
		comp = convertComputation(a[1])
	} else {
		comp = convertComputation(s)
	}
	return comp, dest
}

// convertDestination returns a string of length 3,
// which is the binary representation of the destination
// in machine code.
func convertDestination(s string) string {
	if s == "null" {
		return "000"
	}
	if len(s) > 3 {
		log.Fatal("Syntax Error: invalid Destination.")
	}
	out := []string{"0", "0", "0"}
	for _, v := range s {
		switch v {
		case 'A':
			out[0] = "1"
		case 'D':
			out[1] = "1"
		case 'M':
			out[2] = "1"
		default:
			log.Fatal("Syntax Error: invalid destination")
		}
	}
	return strings.Join(out, "")
}

// convertAssignment returns a string of length 7,
// which is the binary representation of the assignment
// instructions that are given to the ALU.
func convertComputation(s string) string {
	bin, ok := predefined_computations[s]
	if !ok {
		log.Fatal("Syntax Error: Invalid Assignment.")
	}
	return bin
}

// convertJumps returns a string of length 3,
// which is the binary representation of the jump
// instructions for C-instructions.  This is the last
// 3 bits in the c-instruction.
func convertJumps(s string) string {
	bin, ok := predefined_jumps[s]
	if !ok {
		log.Fatal("Syntax Error: Invalid Jump Command.")
	}
	return bin
}

// AddAllSymbols passes through the array without modifying it.
// It checks for unknown symbols in L and A instructions,
// and adds them to the symbol table for later use.
func AddAllSymbols(s []string) {
	line_counter := 0
	var_counter := 16

	// L instructions take higher priority.
	// Looks for those first, and save them into the table.
	for i := 0; i < len(s); i++ {

		// Save the command type, it will help determine
		// how we count lines and variable locations.
		t := getCommandType(s[i])

		// If this is an L command, then we should add it's
		// following line number as the reference.
		if t == _L_COMMAND {
			_, ok := SymbolTable[s[i]]
			if ok {
				log.Fatal("Can't have same L command twice.")
			}

			// Remove the surrounding parenthesis, which
			// contain the symbol on the inside.
			l := s[i]
			l = strings.Replace(l, "(", "", 1)
			l = strings.Replace(l, ")", "", 1)

			SymbolTable[l] = line_counter
		}

		// Don't count L instructions in the line count,
		// because they will get removed from the actual
		// machine code.
		if t == _A_COMMAND || t == _C_COMMAND {
			line_counter++
		}
	}

	// Once we have scanned for all of our L instructions,
	// make a pass and look for other variables.
	for i := 0; i < len(s); i++ {

		// If we are looking at an A command, then we
		// want to check the memory reference to see if it
		// is a variable or not.

		if getCommandType(s[i]) == _A_COMMAND {
			// Remove the @ symbol from the A-instruction.
			A := strings.Replace(s[i], "@", "", 1)
			if isUnknownVariable(A) {
				SymbolTable[A] = var_counter
				var_counter++
			}
			continue
		}
	}
}

// isUnknownVariable checks an A instruction.
// If the memory address is a number, or if it is already
// part of the symbol table, return false.
// If it is a symbol we haven't seen before, return true.
func isUnknownVariable(A string) bool {

	// First, check to see if the location is an integer.
	// If it can be parsed into an integer, it's resolved.
	_, err := strconv.ParseInt(A, 10, 64)
	if err == nil {
		return false
	}

	// Next, treat the location as a symbol, and check
	// the symbol table to see if it's there.
	_, ok := SymbolTable[A]
	if ok {
		return false
	}

	// Otherwise, we haven't seen this symbol before.
	return true
}
