package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const help_message = `
INFO:
	Hack Virtual Machine Translator (hackvmslate)

	Translates virtual machine code into assembly code for the
	"hack virtual machine" from the Nand2Tetris course.

USAGE:
	hackvmslate [options]

	When no options are used, the Current Working Directory
	will be examined for .vm files.  If there exist .vm files,
	then each one will be parsed, and a single file .asm will
	be output.

	The name for that output file will be the name of the 
	current working directory, followed by ".asm".  It will 
	overwrite an existing file of the same name.

OPTIONS:

-h | --help
	Displays this help message to stderr, then exits.

-i
	Enters Interactive Mode. This uses stdin and stdout
	for reading and writing, instead of creating a file
	to store the output.  This can be useful if you want 
	to only process a single file, or if you want to 
	experiment by typing directly.
`

var (
	working_directory = ""
	interactive_mode  = false
	current_filename  = ""
)

func main() {

	// Checks to see if there are any arguments given.
	// If not, then exit with help message.
	// if len(os.Args) < 2 {
	// 	fmt.Fprintln(os.Stderr, help_message)
	// 	os.Exit(1)
	// }

	// checks for the help message flag "-h" or "--help"
	for _, v := range os.Args {
		if v == "-h" || v == "--help" {
			fmt.Fprintln(os.Stderr, help_message)
			os.Exit(1)
		}
	}

	// initialize variables that will hold pointers to the
	// source reader and target writer.
	var r *bufio.Reader
	var w *bufio.Writer

	// Check arguments for "-i", and enter Interactive mode.
	// Set input and output to stdin and stdout.
	// Enter parsing mode, and then exit the program.
	if len(os.Args) >= 2 && os.Args[1] == "-i" {
		interactive_mode = true
		current_filename = "stdin"
		r = bufio.NewReader(os.Stdin)
		w = bufio.NewWriter(os.Stdout)
		parseFile(r, w)
		os.Exit(0)
	}

	// Creates an output file if we aren't in interactive mode.
	output_file, err := os.Create(GetOutputFilename())
	if err != nil {
		failrar(err)
	}
	defer output_file.Close()
	defer fmt.Println("Created File: ", output_file.Name())
	w = bufio.NewWriter(output_file)

	// Get the working directory, and save it in the global variable for later.
	working_directory, err := os.Getwd()
	if err != nil {
		failrar(err)
	}
	file_list := GetVmFilesFromDir(working_directory)
	if len(file_list) <= 0 {
		fmt.Fprintln(os.Stderr, "There are no .vm files to translate in this directory.")
		os.Exit(1)
	}

	fmt.Println("List of .vm Files in this directory:", file_list)

	// Before writing each of the files, write Sys.Init
	fmt.Fprintln(w, WriteInit())

	// Go through each file, parsing it's code, and writing to the .asm file.
	for _, filename := range file_list {
		fmt.Fprintln(w, "// ~~~~ "+filename+" ~~~~")
		HandleFile(r, w, filename)
	}

	// Write the final line of assembly code to end program.
	fmt.Fprintln(w, s_end_program)
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
// func resolveOutputPath(s string) string {

// 	// Resolve the absolute path of the input file.
// 	// absolute, err := filepath.Abs(s)
// 	// if err != nil {
// 	// 	failrar(err)
// 	// }

// 	// Retrieve the directory of the input file.
// 	// dir := filepath.Dir(s)

// 	// Create a name for the output file based on
// 	// the directory name.
// 	// name := filepath.Base(s) + ".asm"

// 	// Combines the directory path and the .hack name
// 	// and return the filepath string.
// 	// return
// }

// parseFile accepts a buffered reader and buffered writer,
// the reader should contain VM code,
// the writer is the destination for the assembly code.
// Does NOT close the reader at the end.
func parseFile(r *bufio.Reader, w *bufio.Writer) {

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
			failrar("Line", source_line_count, ":", err, "\n[TEXT]:", scanner.Text())
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

	// Flush out the buffer, writing all of the stored data to file.
	w.Flush()
}

// GetVmFilesFromDir returns a an array of .vm files in the given directory.
// If there are none, returned array will contain no elements.
func GetVmFilesFromDir(path string) []string {

	// save a variable that will contain the array of file names.
	var filename_list []string

	// create a list containing each of the files in the directory.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		failrar(err)
	}

	// check each of the file names for the extention ".vm",
	// if a .vm is found, then add it to the list of filenames,
	// which will be returned at the end of the function.
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".vm" {
			filename_list = append(filename_list, file.Name())
		}
	}
	return filename_list
}

// GetOutputFilename returns the absolute path of the .asm output file.
func GetOutputFilename() string {
	working_directory, err := os.Getwd()
	if err != nil {
		failrar(err)
	}
	name := filepath.Base(working_directory)
	name = strings.TrimSuffix(name, "/")
	name = strings.TrimSuffix(name, "\\")
	name += ".asm"
	return name
}

func HandleFile(r *bufio.Reader, w *bufio.Writer, filename string) {
	// Attempt to open the file for reading.
	input_file, err := os.Open(filename)
	if err != nil {
		failrar(err)
	}
	defer input_file.Close()

	// Sets the global variable containing current filename
	// this is used by the translator when encounting
	// static varibles
	current_filename = filepath.Base(filename)

	// Creates buffered i/o for the files.
	r = bufio.NewReader(input_file)
	parseFile(r, w)
	defer w.Flush()
}
