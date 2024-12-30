package token

import (
	"Interpreter/object"
	"Interpreter/tokentype"
	"fmt"
	"os"
	"strconv"
)

type Token struct {
	obj    object.Object
	lexeme string
	line   int
}

// Constructor
func NewToken(kind tokentype.TokenType, lexeme string, literal string, line int) *Token {
	// Based on the token type create the literal
	var obj object.Object

	switch kind {
	case tokentype.STRING:
		obj = *object.NewObject(kind, literal)
	case tokentype.NUMBER:
		// Convert to number
		s, err := strconv.ParseFloat(literal, 64)

		// Error in conversion means the tokenizer has issues
		if err != nil {
			fmt.Println("Implementation Error: tokenizer incorrectly parsed number")
			os.Exit(3)
		}
		obj = *object.NewObject(kind, s)
	default:
		obj = *object.NewObject(kind, nil)
	}

	this_token := Token{obj, lexeme, line}
	return &this_token
}

// Get the token type
func (t *Token) GetType() tokentype.TokenType {
	return t.obj.GetKind()
}

// Get the literal value
func (t *Token) GetLiteral() any {
	return t.obj.GetLiteral()
}

func (t *Token) GetLiteralStr() string {
	switch t.obj.GetKind() {
	case tokentype.NUMBER:
		// Assert type
		s, ok := t.obj.GetLiteral().(float64)
		if ok {
			return fmt.Sprintf("%f", s)
		}

		// If type assertion fails, error in implementation
		fmt.Println("Implementation Error: error in number tokenizer.")

		os.Exit(4)
	case tokentype.STRING:
		// Assert type
		s, ok := t.obj.GetLiteral().(string)
		if ok {
			return s
		}

		// If type assertion fails, error in implementation
		fmt.Println("Implementation Error: error in string tokenizer.")
		os.Exit(5)
	}

	// all other instances the literal is blank
	return ""
}

// Get the lexeme value
func (t *Token) GetLexeme() string {
	return t.lexeme
}

// Get the line value
func (t *Token) GetLine() int {
	return t.line
}

// Convert token content to string
func (t *Token) String() string {
	return "Type:" + t.obj.GetKindStr() + " Lexeme:" + t.lexeme + " Literal:" + t.GetLiteralStr()
}
