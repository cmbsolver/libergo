package main

import (
	"flag"
	"fmt"
	"os"
	"runer"
	"titler"
)

// main reads input text, encodes it, and writes the result to an output file or stdout.
func main() {
	titler.PrintTitle("Gematria Encoder")

	// Define flags
	textFlag := flag.String("text", "", "Text to be encoded")
	fileFlag := flag.String("file", "", "File containing text to be encoded (overrides text flag)")
	encodingType := flag.String("type", "latin-to-rune", "Type of encoding: 'latin-to-rune' or 'rune-to-latin'")
	outputFile := flag.String("output", "", "Output file to write the encoded text")
	helpFlag := flag.Bool("help", false, "Display help")

	// Parse flags
	flag.Parse()

	// Check if no flags are provided
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Display help if requested
	if *helpFlag {
		flag.Usage()
		return
	}

	// Read input text
	var inputText string
	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
		inputText = string(data)
	} else {
		inputText = *textFlag
	}

	// Perform encoding
	var encodedText string
	switch *encodingType {
	case "latin-to-rune":
		prepped := runer.PrepLatinToRune(inputText)
		encodedText = runer.TransposeLatinToRune(prepped)
	case "rune-to-latin":
		encodedText = runer.TransposeRuneToLatin(inputText)
	default:
		fmt.Println("Invalid encoding type:", *encodingType)
		os.Exit(1)
	}

	// Output result
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(encodedText), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(encodedText)
	}
}
