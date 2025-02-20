package almond

import (
	"bufio"
	"fmt"
	"os"
)

// Process input string
func run(inter *Interpreter, line string) {
	tokenizer := NewTokenizer(line)
	allToks := tokenizer.Tokenize()
	parser := NewParser(allToks)
	statements := parser.Parse()

	if HadFault {
		return
	}

	inter.Interpret(statements)
}

// Run the code from a file
func RunFile(filename string) error {
	// Create a new interpreter
	inter := NewInterpreter()

	// read file bytes and get error
	data, e := os.ReadFile(filename)

	// propogate error
	if e != nil {
		return e
	}

	// convert byte to string and run code
	run(inter, string(data))

	// exit if there is an error in the code
	if HadFault {
		os.Exit(65)
	}
	if HadRuntimeFault {
		os.Exit(70)
	}
	return e
}

// Run interactive console
func RunPrompt() error {
	// Create a new interpreter
	inter := NewInterpreter()

	// Scan console inputs
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("To exit press Ctrl+C...\n> ")

	// Readline and process text
	for scanner.Scan() {
		text := scanner.Text()

		// Run and store line input
		run(inter, text)

		// Dont kill session if there is an error
		HadFault = false

		fmt.Print("> ")
	}

	return scanner.Err()
}
