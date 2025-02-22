package almond

import (
	"fmt"
)

var HadFault bool = false
var HadRuntimeFault bool = false

// Helper to format error
func report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	HadFault = true
}

// Report errors using tokens
func TokenError(token Token, message string) {
	if token.GetType() == EOF {
		report(token.GetLine(), " at end", message)
	} else {
		report(token.GetLine(), " at '"+token.GetLexeme()+"'", message)
	}
}

// Submit an error
func Error(line int, message string) {
	report(line, "", message)
}

// Runtime Error
func RuntimeError(message string, tok Token) {
	fmt.Printf("%s\n[line %d] ", message, tok.GetLine())
	HadRuntimeFault = true
}
