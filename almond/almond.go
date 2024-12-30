package almond

import (
	"Interpreter/fault"
	"Interpreter/interpreter"
	"Interpreter/parser"
	"Interpreter/tokenizer"
	"bufio"
	"fmt"
	"os"
)

// Process input string
func run(inter *interpreter.Interpreter, line string) {
	tokenizer := tokenizer.NewTokenizer(line)
	all_toks := tokenizer.Tokenize()
	parser := parser.NewParser(all_toks)
	statements := parser.Parse()

	if fault.Had_fault {
		return
	}

	inter.Interpret(statements)
}

// Run the code from a file
func RunFile(filename string) error {
	// Create a new interpreter
	inter := interpreter.NewInterpreter()

	// read file bytes and get error
	data, e := os.ReadFile(filename)

	// propogate error
	if e != nil {
		return e
	}

	// convert byte to string and run code
	run(inter, string(data))

	// exit if there is an error in the code
	if fault.Had_fault {
		os.Exit(65)
	}
	if fault.Had_runtime_fault {
		os.Exit(70)
	}
	return e
}

// Run interactive console
func RunPrompt() error {
	// Create a new interpreter
	inter := interpreter.NewInterpreter()

	// Scan console inputs
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("To exit press Ctrl+C...\n> ")

	// Readline and process text
	for scanner.Scan() {
		text := scanner.Text()

		// Run and store line input
		run(inter, text)

		// Dont kill session if there is an error
		fault.Had_fault = false

		fmt.Print("> ")
	}

	return scanner.Err()
}
