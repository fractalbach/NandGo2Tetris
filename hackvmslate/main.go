// sup./
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	help_message = "USAGE: hackvmslate INPUT_FILENAME"
)

func main() {

	// Checks to see if there are any arguments given.
	// If not, then exit with help message.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, help_message)
		os.Exit(1)
	}

	// input file location as given by command argument.
	input := ""
	input = os.Args[1]
	fmt.Println(input)

	// Load the file into memory.
	data, err := ioutil.ReadFile(input)
	if err != nil {
		failrar(err)
	}

	// TESTING ZONE --------------------------------------
	count := 0

	x := S_JEQ(next(&count))
	y := S_JGT(next(&count))
	z := S_JLT(next(&count))

	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(z)

	fmt.Printf("%v\n%s", data, data)
	//  --------------------------------------------------

	// DEBUG:
	// Print the byte array to stdout

	// // Send the data to the parser, create a file,
	// // and write the machine code into the file!
	// parsed := ParseData(data)
	// err := ioutil.WriteFile(output_filename, []byte(parsed), 0600)
	// if err != nil {
	//     log.Fatal(err)
	// }
}

// failrar prints to stderr and exits the program.
func failrar(a ...interface{}) {
	fmt.Fprint(os.Stderr, "[EPIC FAIL]: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
