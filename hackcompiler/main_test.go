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

/*
func TestSingleProgram(t *testing.T) {
	cmd := exec.Command("hackcompiler", example_input_filepath, "--xml")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	s_result := strings.TrimSpace(string(stdoutStderr))
	r_result := strings.NewReader(s_result)

	b_correct, err := ioutil.ReadFile(correct_token_files[0])
	if err != nil {
		t.Error(err)
	}
	s_correct := strings.TrimSpace(string(b_correct))
	r_correct := strings.NewReader(s_correct)

	scanner1 := bufio.NewScanner(r_result)
	scanner2 := bufio.NewScanner(r_correct)

	c1 := make(chan string)
	c2 := make(chan string)
	quit := make(chan bool)

	// begin the first scanner.
	go func() {
		for scanner1.Scan() {
			if scanner1.Text() == "" {
				continue
			}
			c1 <- scanner1.Text()
		}
		if err := scanner1.Err(); err != nil {
			t.Error(err)
		}
		close(c1)
		quit <- true
	}()

	// begin the second scanner.
	go func() {
		for scanner2.Scan() {
			if scanner2.Text() == "" {
				continue
			}
			c2 <- scanner2.Text()
		}
		if err := scanner2.Err(); err != nil {
			t.Error(err)
		}
		close(c2)
		quit <- true
	}()

	var s1, s2 string
	line_num := 1
	for {
		select {
		case <-quit:
			return
		default:
			s1 = <-c1
			s2 = <-c2
			if s1 != s2 {
				t.Errorf("Line %d does not match:\n\t%s\n\t%s\n", line_num, s1, s2)
			}
			t.Log(line_num, s1, s2)
			line_num++
		}
	}
}
*/

func TestTokenizer(t *testing.T) {
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

func TestParser(t *testing.T) {
	// Iterate through each of the source files
	for file_num, source_file := range input_source_files {

		// Call the hackcompiler parser on the current file.
		cmd := exec.Command("hackcompiler", source_file, "--parse")

		// Place output into a string, trim and split.
		stdoutStderr, err := cmd.CombinedOutput()
		errCheck(t, err)
		result := strings.Split(strings.TrimSpace(string(stdoutStderr)), "\n")

		// Load the comparison file into a string: trim and split.
		// The file has \r\n endings, so remove all of the \r bytes.
		bytes_correct, err := ioutil.ReadFile(correct_token_files[file_num])
		errCheck(t, err)
		string_temp := strings.TrimSpace(string(bytes_correct))
		string_temp = strings.Replace(string_temp, "\r", "", -1)
		correct := strings.Split(string_temp, "\n")

		// compare the length of the arrays to eliminate any obvious fails.
		if len(correct) != len(result) {
			t.Errorf("\n%s\n\tFile lengths do not match. got:(%d), expected:(%d)", source_file, len(result), len(correct))
			t.FailNow()
		}
	}
}

func errCheck(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
