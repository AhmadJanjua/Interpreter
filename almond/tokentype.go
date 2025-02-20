package almond

type TokenType int

const (
	// Single characters
	L_PAREN TokenType = iota
	R_PAREN
	L_BRACE
	R_BRACE
	COMMA
	PERIOD
	MINUS
	PLUS
	STAR
	SLASH
	ASSIGNMENT
	BANG
	GREATER
	LESS
	SEMI_COLON
	HASH

	// Double characters
	EQUALS
	NOT_EQUALS
	GREATER_EQUAL
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	CLASS
	FN
	RETURN
	AUTO
	AND
	OR
	IF
	ELSE
	TRUE
	FALSE
	FOR
	WHILE
	PRINT
	SUPER
	THIS
	NULL
	EOF
)

// TokenType to string mapping
var tokenTypeNames = map[TokenType]string{
	// Single characters
	L_PAREN:    "L_PAREN",
	R_PAREN:    "R_PAREN",
	L_BRACE:    "L_BRACE",
	R_BRACE:    "R_BRACE",
	COMMA:      "COMMA",
	PERIOD:     "PERIOD",
	MINUS:      "MINUS",
	PLUS:       "PLUS",
	STAR:       "STAR",
	SLASH:      "SLASH",
	ASSIGNMENT: "ASSIGNMENT",
	BANG:       "BANG",
	GREATER:    "GREATER",
	LESS:       "LESS",
	SEMI_COLON: "SEMI_COLON",
	HASH:       "HASH",
	AND:        "AND",
	OR:         "OR",

	// Double characters
	EQUALS:        "EQUALS",
	NOT_EQUALS:    "NOT_EQUALS",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS_EQUAL:    "LESS_EQUAL",

	// Literals
	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	NUMBER:     "NUMBER",

	// Keywords
	CLASS:  "CLASS",
	FN:     "FUNCTION",
	RETURN: "RETURN",
	AUTO:   "VARIABLE",
	IF:     "IF",
	ELSE:   "ELSE",
	TRUE:   "TRUE",
	FALSE:  "FALSE",
	FOR:    "FOR",
	WHILE:  "WHILE",
	PRINT:  "PRINT",
	SUPER:  "SUPER",
	THIS:   "THIS",
	NULL:   "NULL",
	EOF:    "EOF",
}

func (t TokenType) String() string {
	if name, ok := tokenTypeNames[t]; ok {
		return name
	}
	return "<ERROR>"
}
