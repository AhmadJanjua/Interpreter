package fault

import "fmt"

var Had_fault bool = false

// Helper to format error
func report(line int, where, message string) {
	fmt.Println("[line " + string(line) + "] Error" + where + ": " + message)
	Had_fault = true
}

// Submit an error
func Error(line int, message string) {
	report(line, "", message)
}
