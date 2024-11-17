package token

import (
	"Interpreter/tokentype"
)

type Token struct {
	tok_type tokentype.TokenType
	lexeme   string
	literal  string
	line     int
}

// Constructor
func NewToken(tok_type tokentype.TokenType, lexeme string, literal string, line int) *Token {
	this_token := Token{tok_type, lexeme, literal, line}
	return &this_token
}

// Convert token content to string
func (t *Token) String() string {
	return "Type:" + t.tok_type.String() + " Lexeme:" + t.lexeme + " Literal:" + t.literal
}
