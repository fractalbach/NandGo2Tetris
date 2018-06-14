package main

import (
	"bufio"
	"fmt"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/CompilationEngine"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackTokenizer"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var example = `
// Example to test the CompileClass() command
class DerpClass {
	static int x;
	static int y;
	field string s;
	field bool b;
	function void derp() {}
}
`

const help_message = `
Compiles Jack code into Hack Programs for Nand2Tetris.

USAGE:         hackcompiler (<filename>|-wd) [option]

FILENAME FLAG:
-wd            Uses all .jack files from the working directory.

OPTIONS:
-t, --token    print tokens in debugging mode.
-p, --parse    print the parse tree as XML.
-x, --xml      print tokens as XML, split by line.
-s, --symbol   print symbol tables generated by parser.

HOW TO USE:
Use "-wd" in place of the filename argument to iterate
through each file in the working directory to use as input.
Each input file creates one output file of the same name.
`

type Mode int

const (
	mode_default Mode = 1 << iota
	mode_tokens_debug
	mode_tokens_xml
	mode_parse_debug
	mode_symbol_table
)

func DebugParse(w io.Writer, r io.Reader) {
	tokenizer := JackTokenizer.Create(r)
	CompilationEngine.Run(w, tokenizer, CompilationEngine.OP_XML)
}

func ParseSymbolTables(w io.Writer, r io.Reader) {
	tokenizer := JackTokenizer.Create(r)
	CompilationEngine.Run(w, tokenizer, CompilationEngine.OP_SYM_TBL)
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

func TokensXML(w io.Writer, r io.Reader) {
	t := JackTokenizer.Create(r)
	fmt.Fprintln(w, "<tokens>")
	for t.HasMoreTokens() {
		fmt.Fprintln(w, t.Current())
		t.Advance()
	}
	fmt.Fprintln(w, "</tokens>")
}

func HelpfulExit() {
	fmt.Fprint(os.Stderr, help_message)
	os.Exit(1)
}

func ReadFullFile(filename string) io.Reader {
	input_file_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		failrar(err)
	}
	return strings.NewReader(string(input_file_bytes))
}

func handle(filename string, mode Mode) {
	var r *bufio.Reader
	var w *bufio.Writer
	r = bufio.NewReader(ReadFullFile(filename))
	switch mode {

	case mode_tokens_debug:
		DebugTokens(r)

	case mode_tokens_xml:
		// w = MakeFile(filename, ".xml")
		w = bufio.NewWriter(os.Stdout)
		TokensXML(w, r)
		w.Flush()

	case mode_parse_debug:
		DebugParse(os.Stdout, r)

	case mode_symbol_table:
		ParseSymbolTables(os.Stdout, r)

	default:
		panic("Invalid mode")
	}
}

func GetJackFilesFromWorkingDir() []string {
	path, err := os.Getwd()
	if err != nil {
		failrar(err)
	}
	// create a list containing each of the files in the directory.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		failrar(err)
	}
	// check each of the file names for the extention ".jack",
	// if a .jack is found, then add it to the list of filenames,
	// which will be returned at the end of the function.
	var filename_list []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".jack" {
			filename_list = append(filename_list, file.Name())
		}
	}
	if len(filename_list) <= 0 {
		failrar("There are no .jack files in the working directory.")
	}
	return filename_list
}

func failrar(a ...interface{}) {
	fmt.Fprint(os.Stderr, "[EPIC FAIL]: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func MakeFile(input_filename, output_suffix string) *bufio.Writer {
	output_filename := strings.TrimSuffix(input_filename, ".jack")
	output_filename += output_suffix
	output_file, err := os.Create(output_filename)
	if err != nil {
		failrar(err)
	}
	return bufio.NewWriter(output_file)
}

func MultiFile(mode Mode) {
	file_list := GetJackFilesFromWorkingDir()
	for _, filename := range file_list {
		fmt.Fprintln(os.Stderr, filename)
		handle(filename, mode)
	}
}

func main() {
	if len(os.Args) < 2 {
		HelpfulExit()
	}
	switch os.Args[1] {
	case "-h", "--help":
		HelpfulExit()
	}
	if len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "Too many arguments.")
		HelpfulExit()
	}

	option := ""
	mode := mode_default
	if len(os.Args) == 3 {
		option = os.Args[2]
	}

	switch option {
	case "-t", "--token":
		mode = mode_tokens_debug
	case "-x", "--xml":
		mode = mode_tokens_xml
	case "-p", "--parse":
		mode = mode_parse_debug
	case "-s", "--symbol":
		mode = mode_symbol_table
	case "":
		fmt.Fprintln(os.Stderr, "[TODO]: No option args given.")
		os.Exit(0)
		mode = mode_default
	default:
		fmt.Fprintln(os.Stderr, "Unknown option argument given.")
		HelpfulExit()
	}

	switch os.Args[1] {
	case "-wd":
		MultiFile(mode)
	default:
		filename := os.Args[1]
		handle(filename, mode)
	}
}
