package fault

import (
	"Interpreter/token"
	"Interpreter/tokentype"
	"fmt"
)

var Had_fault bool = false

// Helper to format error
func report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	Had_fault = true
}

// report errors using tokens
func TokenError(token token.Token, message string) {
	if token.GetType() == tokentype.EOF {
		report(token.GetLine(), " at end", message)
	} else {
		report(token.GetLine(), " at '"+token.GetLexeme()+"'", message)
	}
}

// Submit an error
func Error(line int, message string) {
	report(line, "", message)
}
