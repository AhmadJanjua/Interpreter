package tokentype

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

func (t *TokenType) String() string {
	switch *t {
	// Single characters
	case L_PAREN:
		return "Left Parenthesis"
	case R_PAREN:
		return "Right Parenthesis"
	case L_BRACE:
		return "Left Brace"
	case R_BRACE:
		return "Right Brace"
	case COMMA:
		return "Comma"
	case PERIOD:
		return "Period"
	case MINUS:
		return "Minus"
	case PLUS:
		return "Plus"
	case STAR:
		return "Star"
	case SLASH:
		return "Slash"
	case ASSIGNMENT:
		return "Assignment"
	case BANG:
		return "Bang"
	case GREATER:
		return "Greater"
	case LESS:
		return "Less"
	case SEMI_COLON:
		return "Semi-Colon"
	case HASH:
		return "Hash"
	case AND:
		return "And"
	case OR:
		return "Or"

	// Double characters
	case EQUALS:
		return "Equals"
	case NOT_EQUALS:
		return "Not Equals"
	case GREATER_EQUAL:
		return "Greater than or Equals"
	case LESS_EQUAL:
		return "Less than or Equals"

	// Literals
	case IDENTIFIER:
		return "Identifier"
	case STRING:
		return "String"
	case NUMBER:
		return "Number"

	// Keywords
	case CLASS:
		return "Class"
	case FN:
		return "Function"
	case RETURN:
		return "Return"
	case AUTO:
		return "Variable"
	case IF:
		return "If"
	case ELSE:
		return "Else"
	case TRUE:
		return "True"
	case FALSE:
		return "False"
	case FOR:
		return "For"
	case WHILE:
		return "While"
	case PRINT:
		return "Print"
	case SUPER:
		return "Super"
	case THIS:
		return "This"
	case NULL:
		return "Null"
	case EOF:
		return "Eof"
	default:
		return ""
	}
}
