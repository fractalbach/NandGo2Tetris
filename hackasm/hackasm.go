package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	input_filename  string
	output_filename string
)

func main() {
	flag.StringVar(&output_filename, "o", input_filename+".hack", "Output File Location")
	flag.Parse()

	fmt.Println(flag.Args(), "  0th one=", flag.Arg(0))
	fmt.Printf("Assembling %s...", input_filename)
	fmt.Println("Done!")
}

// GimmeFile reads the file and returns a string.
func GimmeFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
