package main

import (
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/CompilationEngine"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackTokenizer"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var example = `
// Example to test the CompileClass() command
class DerpClass {
}

`

const help_message = `
USAGE:         JackAnalyzer <filename> [option]

FILENAME:      Use '-e' to use an example instead of a file.

OPTIONS:
-d, --debug    prints line number, token kind, and token content.
-p, --parse    Activate Parser: [CURRENTLY UNDER CONSTRUCTION].
-x, --xml      Print tokens as XML, split by line.
`

func DebugParse(r io.Reader) {
	w := os.Stdout
	tokenizer := JackTokenizer.Create(r)
	CompilationEngine.Run(w, tokenizer)
}

func DebugTokens(r io.Reader) {
	t := JackTokenizer.Create(r)
	i := 0
	for t.HasMoreTokens() {
		fmt.Printf("[%3d]: %#v\n", i, t.Current())
		t.Advance()
		i++
	}
}

func TokensXML(r io.Reader) {
	t := JackTokenizer.Create(r)
	fmt.Println("<tokens>")
	for t.HasMoreTokens() {
		fmt.Println(t.Current())
		t.Advance()
	}
	fmt.Println("</tokens>")
}

func HelpfulExit() {
	fmt.Fprint(os.Stderr, help_message)
	os.Exit(1)
}

func errexit(i interface{}) {
	fmt.Fprintln(os.Stderr, i)
	os.Exit(1)
}

func ReadFullFile(filename string) io.Reader {
	input_file_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errexit(err)
	}
	return strings.NewReader(string(input_file_bytes))
}

func handle(filename, option string) {
	var r io.Reader
	if filename == "-e" {
		r = strings.NewReader(example)
	} else {
		r = ReadFullFile(filename)
	}
	switch option {
	case "-d", "--debug":
		DebugTokens(r)
	case "-x", "--xml":
		TokensXML(r)
	case "-p", "--parse":
		DebugParse(r)
	case "":
		fmt.Fprintln(os.Stderr, "Default Behavior: No option args given. [TODO]")
	default:
		fmt.Fprintln(os.Stderr, "Unknown option argument given.")
		HelpfulExit()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Not enough arguments.")
		HelpfulExit()
	}
	if len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "Too many arguments.")
		HelpfulExit()
	}
	option := ""
	if len(os.Args) == 3 {
		option = os.Args[2]
	}
	filename := os.Args[1]
	handle(filename, option)
}
