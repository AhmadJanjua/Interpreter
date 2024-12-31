package tokenizer

import (
	"Interpreter/fault"
	"Interpreter/token"
	"Interpreter/tokentype"
	"unicode"
)

// Look-up table: string -> TokenType
var keyword_map = map[string]tokentype.TokenType{
	"class":  tokentype.CLASS,
	"fn":     tokentype.FN,
	"return": tokentype.RETURN,
	"var":    tokentype.AUTO,
	"if":     tokentype.IF,
	"else":   tokentype.ELSE,
	"true":   tokentype.TRUE,
	"false":  tokentype.FALSE,
	"for":    tokentype.FOR,
	"while":  tokentype.WHILE,
	"print":  tokentype.PRINT,
	"super":  tokentype.SUPER,
	"this":   tokentype.THIS,
	"null":   tokentype.NULL,
}

// Helper to check characters for identifer
func validIdentifier(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_'
}

type Tokenizer struct {
	start   int
	current int
	line    int
	source  string
	tokens  []token.Token
}

// Construct Tokenizer
func NewTokenizer(source string) *Tokenizer {
	tmp := Tokenizer{0, 0, 1, source, []token.Token{}}
	return &tmp
}

// Check if end is reached
func (s *Tokenizer) end() bool {
	return s.current >= len(s.source)
}

// Move to next character
func (s *Tokenizer) advance() rune {
	s.current++
	return rune(s.source[s.current-1])
}

// Get current character
func (s *Tokenizer) peek() rune {
	if s.end() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

// Get next character
func (s *Tokenizer) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return rune(0)
	} else {
		return rune(s.source[s.current+1])
	}
}

// Match the current character to a character
func (s *Tokenizer) match(expected rune) bool {
	if s.end() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

// Process keyword or get identifier
func (s *Tokenizer) processIdentifier() {
	for validIdentifier(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tok_type, ok := keyword_map[text]
	if ok {
		s.addToken(tok_type, "")
	} else {
		s.addToken(tokentype.IDENTIFIER, "")
	}
}

// Parse int or decimal number
func (s *Tokenizer) processNumber() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	s.addToken(tokentype.NUMBER, s.source[s.start:s.current])
}

// Parse string
func (s *Tokenizer) processString() {
	for s.peek() != '"' && !s.end() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.end() {
		fault.Error(s.line, "Unterminated string")
		return
	}
	s.advance()

	text := s.source[s.start+1 : s.current-1]
	s.addToken(tokentype.STRING, text)
}

// append token to list
func (s *Tokenizer) addToken(tok_type tokentype.TokenType, literal string) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *token.NewToken(tok_type, lexeme, literal, s.line))
}

// Search through string and get next token
func (s *Tokenizer) nextToken() {
	c := s.advance()

	switch c {
	// Single character
	case '(':
		s.addToken(tokentype.L_PAREN, "")
	case ')':
		s.addToken(tokentype.R_PAREN, "")
	case '{':
		s.addToken(tokentype.L_BRACE, "")
	case '}':
		s.addToken(tokentype.R_BRACE, "")
	case ',':
		s.addToken(tokentype.COMMA, "")
	case '.':
		s.addToken(tokentype.PERIOD, "")
	case '-':
		s.addToken(tokentype.MINUS, "")
	case '+':
		s.addToken(tokentype.PLUS, "")
	case '*':
		s.addToken(tokentype.STAR, "")
	case '/':
		s.addToken(tokentype.SLASH, "")
	case ';':
		s.addToken(tokentype.SEMI_COLON, "")
	case '|':
		s.addToken(tokentype.OR, "")
	case '&':
		s.addToken(tokentype.AND, "")

	// Single-Double characters
	case '=':
		if s.match('=') {
			s.addToken(tokentype.EQUALS, "")
		} else {
			s.addToken(tokentype.ASSIGNMENT, "")
		}
	case '!':
		if s.match('=') {
			s.addToken(tokentype.NOT_EQUALS, "")
		} else {
			s.addToken(tokentype.BANG, "")
		}
	case '>':
		if s.match('=') {
			s.addToken(tokentype.GREATER_EQUAL, "")
		} else {
			s.addToken(tokentype.GREATER, "")
		}
	case '<':
		if s.match('=') {
			s.addToken(tokentype.LESS_EQUAL, "")
		} else {
			s.addToken(tokentype.LESS, "")
		}

	// Deal with comments
	case '#':
		for s.peek() != '\n' && !s.end() {
			s.advance()
		}

	// New line
	case '\n':
		s.line++
	default:
		// Special cases
		if unicode.IsDigit(c) {
			s.processNumber()
		} else if c == '"' {
			s.processString()
		} else if unicode.IsLetter(c) {
			s.processIdentifier()
		} else if !unicode.IsSpace(c) {
			fault.Error(s.line, "Unexpected character.")
		}
	}
}

// Go through source and create a tokenized list
func (s *Tokenizer) Tokenize() []token.Token {
	for !s.end() {
		s.start = s.current
		s.nextToken()
	}
	s.tokens = append(s.tokens, *token.NewToken(tokentype.EOF, "", "", s.line))
	return s.tokens
}
