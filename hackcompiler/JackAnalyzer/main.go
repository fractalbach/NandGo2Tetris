package main

import (
	"fmt"
	// "io"
	"github.com/fractalbach/nandGo2tetris/hackcompiler/JackTokenizer"
	// "os"
	"strings"
)

var example = `
// COMMENT COMMENT
function thing {
	var static hello;
	var local there;
	var cool_story_bro;
	var string omgstring;

	let hello = 12;
	let there = 13;
	let omgstring = "well hello [ there ] [] ' ' ! ! ";

	if (hello > there) {   // COMMENT COMMENT 
		sys.println(cool_story_bro);
		return true;
	} else {
		return false;  // COMMENT COMMENT
	}
}
`

func main() {
	r := strings.NewReader(example)
	t := JackTokenizer.Create(r)
	i := 0

	for t.HasMoreTokens() {
		fmt.Printf("[%3d]: %#v\n", i, t.Current())
		t.Advance()
		i++
	}

}
