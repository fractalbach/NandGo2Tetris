package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	help_message = "USAGE: hackvmslate INPUT_FILENAME"
)

var (
	interactive_mode = false
)

func main() {

	// Checks to see if there are any arguments given.
	// If not, then exit with help message.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, help_message)
		os.Exit(1)
	}

	// initialize variables that will hold pointers to the
	// source reader and target writer.
	var r *bufio.Reader
	var w *bufio.Writer

	// Check command line arguments for interactive mode,
	// if the -i flag wasn't specified, then assume that it
	// will be a file path.
	if os.Args[1] == "-i" {

		// Open interactive mode by using stdin instead of file.
		interactive_mode = true
		r = bufio.NewReader(os.Stdin)
		w = bufio.NewWriter(os.Stdout)

	} else {

		// Attempt to open the file for reading.
		input_file, err := os.Open(os.Args[1])
		if err != nil {
			failrar(err)
		}
		defer input_file.Close()

		// Creates and Opens a file for writing.
		// The filepath will be resolved automatically,
		// and will be placed in the same directory as the
		// source code file.
		output_file, err := os.Create(resolveOutputPath(os.Args[1]))
		if err != nil {
			failrar(err)
		}
		defer output_file.Close()
		defer fmt.Println("Created File: ", output_file.Name())

		// Creates buffered i/o for the files.
		r = bufio.NewReader(input_file)
		w = bufio.NewWriter(output_file)
		defer w.Flush()
	}

	// Creates a scanner for the string, which will allow us
	// to scan through it line-by-line.
	scanner := bufio.NewScanner(r)

	// Keep track of the numbers of lines we have scanned.
	source_line_count := 0

	// Parses each line found by the scanner.
	for scanner.Scan() {

		// Increment the source code counter by 1,
		// to reflect the position of the scanner.
		source_line_count++

		// Parse the current line of source code,
		// the result will be a string of assembly code.
		// Check for any errors and loudly report them.
		s, err := ParseLine(scanner.Text())
		if err != nil {
			failrar("Line:", source_line_count, err)
		}

		// Skip past empty strings.
		if len(s) < 1 {
			continue
		}

		// Print the parsed string to the buffered output.
		fmt.Fprintln(w, s)

		// Flush the buffer ONLY if in standard output mode.
		// Otherwise, hold the output in a buffer, and wait
		// to write it to a file until the compile finishes.
		if interactive_mode {
			w.Flush()
		}
	}

	// Report any errors that the scanner itself encounters.
	if err := scanner.Err(); err != nil {
		failrar(err)
	}

	// Write the final line of assembly code to end program.
	fmt.Fprintln(w, s_end_program)
	w.Flush()
}

// failrar prints to stderr and exits the program.
// this is just a helper function and makes the code a bit
// cleaner and easier to read and write.
func failrar(a ...interface{}) {
	fmt.Fprint(os.Stderr, "[EPIC FAIL]: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

// ResolveOutputPath accepts the input filename and returns
// a suitable path for the output file.  Using this path
// resolver ensures that the .asm file is in the same directory
// as the .vm file.
func resolveOutputPath(s string) string {

	// Resolve the absolute path of the input file.
	absolute, err := filepath.Abs(s)
	if err != nil {
		failrar(err)
	}

	// Retrieve the directory of the input file.
	dir := filepath.Dir(absolute)

	// Create a name for the output file.
	// Removes old file name extenstion (if needed).
	// Appends the new file name extension.
	name := filepath.Base(s)
	if strings.HasSuffix(name, ".vm") {
		name = strings.TrimSuffix(name, ".vm")
	}
	name += ".asm"

	// Combines the directory path and the .hack name
	// and return the filepath string.
	return filepath.Join(dir, name)
}
