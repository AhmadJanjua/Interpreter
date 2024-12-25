package main

import (
	"Interpreter/almond"
	"fmt"
	"os"
)

// method to interact with interpreter: shell || src file
func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: too many arguments -> try to pass path to source code or run with no arguments")
		os.Exit(64)

	} else if len(args) == 1 {
		// Run file
		almond.RunFile(args[0])

	} else {
		// Interative shell
		almond.RunPrompt()
	}
}
