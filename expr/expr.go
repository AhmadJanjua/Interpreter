package expr

import (
	"Interpreter/fault"
	"Interpreter/token"
	"Interpreter/tokentype"
	"errors"
	"fmt"
)

type Expr interface {
	accept() string
	evaluate() (Literal, error)
}

// Literal
type Literal struct {
	kind tokentype.TokenType
	str  string
	num  float64
}

func NewLiteral(kind tokentype.TokenType) *Literal {
	return &Literal{kind, "", 0}
}
func NewString(str string) *Literal {
	return &Literal{tokentype.STRING, str, 0}
}
func NewNumber(num float64) *Literal {
	return &Literal{tokentype.NUMBER, "", num}
}

func (l Literal) accept() string {
	// Return number
	if l.kind == tokentype.NUMBER {
		return l.kind.String() + "-" + fmt.Sprintf("%f", l.num)
	}

	// Return string
	if l.kind == tokentype.STRING {
		return l.kind.String() + "-" + l.str
	}

	// Return type
	return l.kind.String()
}
func (l Literal) evaluate() (Literal, error) {
	return l, nil
}

// Get the truth value of a literal
func (l Literal) isTruthy() bool {
	// 0, null, false
	switch l.kind {
	case tokentype.NULL:
		return false
	case tokentype.FALSE:
		return false
	case tokentype.NUMBER:
		if l.num == 0 {
			return false
		}
	}
	return true
}

func (left Literal) isEqual(right Literal) bool {
	if left.kind != right.kind {
		return false
	}

	switch right.kind {
	case tokentype.NUMBER:
		if right.num != left.num {
			return false
		}
	case tokentype.STRING:
		if right.str != left.str {
			return false
		}
	}

	// otherwise if the kind matches (true, false, null etc)
	return true

}

// Unary
type Unary struct {
	operator token.Token
	right    Expr
}

func NewUnary(operator token.Token, right Expr) *Unary {
	return &Unary{operator, right}
}
func (u Unary) accept() string {
	return "( " + u.operator.GetLexeme() + u.right.accept() + " )"
}
func (u Unary) evaluate() (Literal, error) {
	right, err := u.right.evaluate()

	if err != nil {
		return right, err
	}

	switch u.operator.GetType() {
	case tokentype.BANG:
		// if the value is true -> return false
		if right.isTruthy() {
			return *NewLiteral(tokentype.FALSE), nil
		} else {
			// true -> false
			return *NewLiteral(tokentype.TRUE), nil
		}
	case tokentype.MINUS:
		// negate number
		if right.kind == tokentype.NUMBER {
			return *NewNumber(-1 * right.num), nil
		}
	}

	// Report error
	fault.RuntimeError("Eval Error: illegal unary operator", u.operator)
	return Literal{}, errors.New("illegal unary operator")
}

// Binary
type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) *Binary {
	return &Binary{left, operator, right}
}
func (b Binary) accept() string {
	return "(" + b.operator.GetLexeme() + b.left.accept() + b.right.accept() + ")"
}
func (b Binary) evaluate() (Literal, error) {
	right, err := b.right.evaluate()

	// Cascade error
	if err != nil {
		return right, err
	}

	left, err := b.left.evaluate()

	// Cascade error
	if err != nil {
		return left, err
	}

	// check equality
	if b.operator.GetType() == tokentype.EQUALS {
		if left.isEqual(right) {
			return *NewLiteral(tokentype.TRUE), nil
		}
		return *NewLiteral(tokentype.FALSE), nil
	}
	if b.operator.GetType() == tokentype.NOT_EQUALS {
		if !left.isEqual(right) {
			return *NewLiteral(tokentype.TRUE), nil
		}
		return *NewLiteral(tokentype.FALSE), nil
	}

	// Report type missmatch for non equality tests
	if right.kind != left.kind {
		fault.RuntimeError("Eval Error: type mismatch between "+left.kind.String()+" and "+right.kind.String(), b.operator)
		return Literal{}, errors.New("type mismatch in binary operation")
	}

	// String concatenation
	if right.kind == tokentype.STRING && b.operator.GetType() == tokentype.PLUS {
		return *NewString(left.str + right.str), nil
	}

	// Report invalid non-numeric operations
	if right.kind != tokentype.NUMBER {
		fault.RuntimeError("Eval Error: invalid binary non-numeric operation", b.operator)
		return Literal{}, errors.New("invalid binary operation on non-numeric")
	}

	switch b.operator.GetType() {
	case tokentype.PLUS:
		return *NewNumber(left.num + right.num), nil
	case tokentype.MINUS:
		return *NewNumber(left.num - right.num), nil
	case tokentype.SLASH:
		return *NewNumber(left.num / right.num), nil
	case tokentype.STAR:
		return *NewNumber(left.num * right.num), nil
	case tokentype.GREATER:
		if left.num > right.num {
			return *NewLiteral(tokentype.TRUE), nil
		}
		return *NewLiteral(tokentype.FALSE), nil
	case tokentype.GREATER_EQUAL:
		if left.num >= right.num {
			return *NewLiteral(tokentype.TRUE), nil
		}
		return *NewLiteral(tokentype.FALSE), nil
	case tokentype.LESS:
		if left.num < right.num {
			return *NewLiteral(tokentype.TRUE), nil
		}
		return *NewLiteral(tokentype.FALSE), nil
	case tokentype.LESS_EQUAL:
		if left.num <= right.num {
			return *NewLiteral(tokentype.TRUE), nil
		}
		return *NewLiteral(tokentype.FALSE), nil
	}

	fault.RuntimeError("Eval Error: illegal binary operator", b.operator)
	return Literal{}, errors.New("illegal binary operator")
}

// Grouping
type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression}
}
func (g Grouping) accept() string {
	return "( group " + g.expression.accept() + ")"
}
func (g Grouping) evaluate() (Literal, error) {
	return g.expression.evaluate()
}

// Print the Abstract Syntax Tree
func ASTPrinter(e Expr) string {
	return e.accept()
}

func Interpret(e Expr) {
	literal, err := e.evaluate()

	if err != nil {
		fmt.Println(err)
		return
	}

	if literal.kind == tokentype.NUMBER {
		fmt.Println(literal.num)
		return
	}

	if literal.kind == tokentype.STRING {
		fmt.Println(literal.str)
		return
	}

	fmt.Println(literal.kind.String())
}
