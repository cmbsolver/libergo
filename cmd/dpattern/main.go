package main

import (
	"flag"
	"fmt"
	"lgstructs"
)

// main is the entry point of the application. It parses input flags, processes text, and outputs the corresponding rune pattern.
func main() {
	// Define flags
	textFlag := flag.String("text", "", "Text to get rune pattern")

	// Parse flags
	flag.Parse()

	// Check if the text flag is empty
	if *textFlag == "" {
		flag.Usage()
		return
	}

	// Create a DictionaryWord instance
	dw := lgstructs.DictionaryWord{
		RuneWordText: *textFlag,
	}

	// Get the rune pattern
	pattern := lgstructs.GetRunePattern(dw.RuneWordText)

	// Output the result
	fmt.Println("Pattern:", pattern)
}
