package main

import (
	// "io/ioutil"
	// "bufio"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	cmd := exec.Command("go", "install", "-v")
	cmd.Run()
}

var (
	input_source_files = []string{
		"tests/ExpressionLessSquare/Main.jack",
		"tests/ExpressionLessSquare/Square.jack",
		"tests/ExpressionLessSquare/SquareGame.jack",
		"tests/ArrayTest/Main.jack",
		"tests/Square/Main.jack",
		"tests/Square/Square.jack",
		"tests/Square/SquareGame.jack",
	}
	// test files of the correct Tokenizer output.
	correct_token_files = []string{
		"tests/ExpressionLessSquare/expected/MainT.xml",
		"tests/ExpressionLessSquare/expected/SquareT.xml",
		"tests/ExpressionLessSquare/expected/SquareGameT.xml",
		"tests/ArrayTest/expected/MainT.xml",
		"tests/Square/expected/MainT.xml",
		"tests/Square/expected/SquareT.xml",
		"tests/Square/expected/SquareGameT.xml",
	}
	// test files of the correct Parser output.
	correct_parse_files = []string{
		"tests/ExpressionLessSquare/expected/Main.xml",
		"tests/ExpressionLessSquare/expected/Square.xml",
		"tests/ExpressionLessSquare/expected/SquareGame.xml",
		"tests/ArrayTest/expected/Main.xml",
		"tests/Square/expected/Main.xml",
		"tests/Square/expected/Square.xml",
		"tests/Square/expected/SquareGame.xml",
	}
)

// creates absolute paths to the testing files.
func TestFilepathResolution(t *testing.T) {
	var err error
	for i, path_suffix := range correct_token_files {
		correct_token_files[i], err = filepath.Abs(path_suffix)
		if err != nil {
			t.Error(err)
		}
		// t.Log(correct_token_files[i])
	}
	for i, path_suffix := range correct_parse_files {
		correct_parse_files[i], err = filepath.Abs(path_suffix)
		if err != nil {
			t.Error(err)
		}
		// t.Log(correct_token_files[i])
	}
	for i, path_suffix := range input_source_files {
		input_source_files[i], err = filepath.Abs(path_suffix)
		if err != nil {
			t.Error(err)
		}
	}
	t.Log("Absolute filepaths successfully resolved.")
}

func __TestTokenizer(t *testing.T) {
	for file_num, source_file := range input_source_files {
		cmd := exec.Command("hackcompiler", source_file, "--xml")
		// Place tokenized output into a string, trim and split.
		stdoutStderr, err := cmd.CombinedOutput()
		errCheck(t, err)
		s_result := strings.Split(strings.TrimSpace(string(stdoutStderr)), "\n")

		// The comparison file: trim space and split by lines
		b_correct, err := ioutil.ReadFile(correct_token_files[file_num])
		errCheck(t, err)
		string_correct := strings.TrimSpace(string(b_correct))
		string_correct = strings.Replace(string_correct, "\r", "", -1)
		s_correct := strings.Split(string_correct, "\n")

		if len(s_result) != len(s_correct) {
			t.Errorf("\n%s\n\tFile lengths do not match. got:(%d), expected:(%d)", source_file, len(s_result), len(s_correct))
			t.FailNow()
		}

		for i := range s_result {
			if s_result[i] != s_correct[i] {
				t.Errorf("\n%s:%d:\n\tLine:(%d) does not match. got:(%s), expected:(%s)", source_file, i, i, s_result[i], s_correct[i])
				t.FailNow()
			}
		}
		t.Log(source_file)
	}
}

func __TestParser(t *testing.T) {
	// Iterate through each of the source files
	for file_num, source_file := range input_source_files {

		// Print name of current file so we can reference bugs to it later.
		t.Log(source_file)

		// Call the hackcompiler parser on the current file.
		cmd := exec.Command("hackcompiler", source_file, "--parse")

		// Place output into a string, trim and split.
		stdoutStderr, err := cmd.CombinedOutput()
		errCheck(t, err)
		result := strings.Split(strings.TrimSpace(string(stdoutStderr)), "\n")

		// Load the comparison file into a string: trim and split.
		// The file has \r\n endings, so remove all of the \r bytes.
		bytes_correct, err := ioutil.ReadFile(correct_parse_files[file_num])
		errCheck(t, err)
		string_temp := strings.TrimSpace(string(bytes_correct))
		string_temp = strings.Replace(string_temp, "\r", "", -1)
		correct := strings.Split(string_temp, "\n")

		// compare the length of the arrays to eliminate any obvious fails.
		// save the minimum length, and use it to compare lines.
		minlen := len(result)
		if len(correct) != len(result) {
			t.Errorf("File lengths do not match. got:(%d), expected:(%d)", len(result), len(correct))
			// use length of the shortest file when matching lines,
			// in order to prevent out-of-bounds errors when checking each line.
			if len(correct) < len(result) {
				minlen = len(correct)
			}
		}
		if minlen <= 0 {
			t.Fatalf("Out of bounds. Cannot have a file length:(%d) <= 0", minlen)
			return
		}
		// report information about mismatched lines.
		// compare lines even if mismatched length, to provide more helpful info.

		for i := 0; i < minlen; i++ {
			if result[i] != correct[i] {
				t.Fatalf("Line:(%d) does not match. \n\tgot:(%s) \n\texpected:(%s)", i+1, result[i], correct[i])
			}
		}
	}
}

func errCheck(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
