package almond

import (
	"unicode"
)

// Helper to check characters for identifer
func validIdentifier(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_'
}

type Tokenizer struct {
	start   int
	current int
	line    int
	source  string
	tokens  []Token
}

// Construct Tokenizer
func NewTokenizer(source string) *Tokenizer {
	tmp := Tokenizer{0, 0, 1, source, []Token{}}
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
	tokType := TokenTypeLUT(text)
	s.addToken(tokType, "")
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

	s.addToken(NUMBER, s.source[s.start:s.current])
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
		Error(s.line, "Unterminated string")
		return
	}
	s.advance()

	text := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, text)
}

// append token to list
func (s *Tokenizer) addToken(tokType TokenType, literal string) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *NewToken(tokType, lexeme, literal, s.line))
}

// Search through string and get next token
func (s *Tokenizer) nextToken() {
	c := s.advance()

	switch c {
	// Single character
	case '(':
		s.addToken(L_PAREN, "")
	case ')':
		s.addToken(R_PAREN, "")
	case '{':
		s.addToken(L_BRACE, "")
	case '}':
		s.addToken(R_BRACE, "")
	case ',':
		s.addToken(COMMA, "")
	case '.':
		s.addToken(PERIOD, "")
	case '-':
		s.addToken(MINUS, "")
	case '+':
		s.addToken(PLUS, "")
	case '*':
		s.addToken(STAR, "")
	case '/':
		s.addToken(SLASH, "")
	case ';':
		s.addToken(SEMI_COLON, "")
	case '|':
		s.addToken(OR, "")
	case '&':
		s.addToken(AND, "")

	// Single-Double characters
	case '=':
		if s.match('=') {
			s.addToken(EQUALS, "")
		} else {
			s.addToken(ASSIGNMENT, "")
		}
	case '!':
		if s.match('=') {
			s.addToken(NOT_EQUALS, "")
		} else {
			s.addToken(BANG, "")
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, "")
		} else {
			s.addToken(GREATER, "")
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, "")
		} else {
			s.addToken(LESS, "")
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
			Error(s.line, "Unexpected character.")
		}
	}
}

// Go through source and create a tokenized list
func (s *Tokenizer) Tokenize() []Token {
	for !s.end() {
		s.start = s.current
		s.nextToken()
	}
	s.tokens = append(s.tokens, *NewToken(EOF, "", "", s.line))
	return s.tokens
}
