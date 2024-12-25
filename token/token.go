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

// Get the token type
func (t Token) GetType() tokentype.TokenType {
	return t.tok_type
}

// Get the literal value
func (t Token) GetLiteral() string {
	return t.literal
}

// Get the lexeme value
func (t Token) GetLexeme() string {
	return t.lexeme
}

// Get the line value
func (t Token) GetLine() int {
	return t.line
}

// Convert token content to string
func (t *Token) String() string {
	return "Type:" + t.tok_type.String() + " Lexeme:" + t.lexeme + " Literal:" + t.literal
}
