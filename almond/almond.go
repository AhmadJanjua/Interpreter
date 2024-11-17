package almond

import (
	"Interpreter/fault"
	"Interpreter/tokenizer"
	"bufio"
	"fmt"
	"os"
)

// Process input string
func run(line string) {
	tokenizer := tokenizer.NewTokenizer(line)
	all_toks := tokenizer.Tokenize()

	for _, element := range all_toks {
		fmt.Println(element.String())
	}
}

// Run the code from a file
func RunFile(filename string) error {
	// read file bytes and get error
	data, e := os.ReadFile(filename)

	// propogate error
	if e != nil {
		return e
	}

	// convert byte to string and run code
	run(string(data))

	// exit if there is an error in the code
	if fault.Had_fault {
		os.Exit(65)
	}
	return e
}

// Run interactive console
func RunPrompt() error {
	// Scan console inputs
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("To exit press Ctrl+C...\n> ")

	// Readline and process text
	for scanner.Scan() {
		text := scanner.Text()

		// Run and store line input
		run(text)

		// Dont kill session if there is an error
		fault.Had_fault = false

		fmt.Print("> ")
	}

	return scanner.Err()
}
