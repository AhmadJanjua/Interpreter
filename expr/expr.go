package expr

import (
	"Interpreter/environment"
	"Interpreter/fault"
	"Interpreter/object"
	"Interpreter/token"
	"Interpreter/tokentype"
	"errors"
	"fmt"
	"os"
)

type Expr interface {
	Evaluate(e *environment.Environment) (object.Object, error)
}

// Literal
type Literal struct {
	value object.Object
}

func NewLiteral(t tokentype.TokenType) *Literal {
	return &Literal{*object.NewObject(t, nil)}
}

func NewNumber(f float64) *Literal {
	return &Literal{*object.NewObject(tokentype.NUMBER, f)}
}

func NewString(s string) *Literal {
	return &Literal{*object.NewObject(tokentype.STRING, s)}
}
func (l Literal) Evaluate(e *environment.Environment) (object.Object, error) {
	return l.value, nil
}

// Unary
type Unary struct {
	operator token.Token
	right    Expr
}

func NewUnary(operator token.Token, right Expr) *Unary {
	return &Unary{operator, right}
}

func (u Unary) Evaluate(e *environment.Environment) (object.Object, error) {
	right, err := u.right.Evaluate(e)

	if err != nil {
		return right, err
	}

	switch u.operator.GetType() {
	case tokentype.BANG:
		// if the value is true -> return false
		if right.Bool() {
			return *object.NewObject(tokentype.FALSE, nil), nil
		} else {
			// true -> false
			return *object.NewObject(tokentype.TRUE, nil), nil
		}
	case tokentype.MINUS:
		// negate number
		if right.GetKind() == tokentype.NUMBER {
			val, ok := right.GetLiteral().(float64)

			if !ok {
				fmt.Println("Implementation Error: Failed to parse float form number in expr -> unary -> Evaluate")
				os.Exit(6)
			}
			return *object.NewObject(tokentype.NUMBER, -1*val), nil
		}
	}

	// Report error
	fault.RuntimeError("Eval Error: illegal unary operator", u.operator)
	return object.Object{}, errors.New("illegal unary operator")
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
func (b Binary) Evaluate(e *environment.Environment) (object.Object, error) {
	right, err := b.right.Evaluate(e)

	// Cascade error
	if err != nil {
		return right, err
	}

	left, err := b.left.Evaluate(e)

	// Cascade error
	if err != nil {
		return left, err
	}

	// check equality
	if b.operator.GetType() == tokentype.EQUALS {
		if left.Equal(&right) {
			return *object.NewObject(tokentype.TRUE, nil), nil
		}
		return *object.NewObject(tokentype.FALSE, nil), nil
	}
	if b.operator.GetType() == tokentype.NOT_EQUALS {
		if !left.Equal(&right) {
			return *object.NewObject(tokentype.TRUE, nil), nil
		}
		return *object.NewObject(tokentype.FALSE, nil), nil
	}

	// Report type missmatch for non equality tests
	if right.GetKind() != left.GetKind() {
		fault.RuntimeError("Eval Error: type mismatch between "+left.GetKindStr()+" and "+right.GetKindStr(), b.operator)
		return object.Object{}, errors.New("type mismatch in binary operation")
	}

	// String concatenation
	if right.GetKind() == tokentype.STRING && b.operator.GetType() == tokentype.PLUS {
		l_str, l_ok := left.GetLiteral().(string)
		r_str, r_ok := right.GetLiteral().(string)

		if l_ok && r_ok {
			return *object.NewObject(tokentype.STRING, l_str+r_str), nil
		}

		fmt.Println("Implementation Error: Could not parse string passed in expr -> binary -> Evaluate")
		os.Exit(7)
	}

	// Report invalid non-numeric operations
	if right.GetKind() != tokentype.NUMBER {
		fault.RuntimeError("Eval Error: invalid binary non-numeric operation", b.operator)
		return object.Object{}, errors.New("invalid binary operation on non-numeric")
	}
	l_num, l_ok := left.GetLiteral().(float64)
	r_num, r_ok := right.GetLiteral().(float64)

	if !l_ok || !r_ok {
		fmt.Println("Implementation Error: Could not parse numbers passed in expr -> binary -> Evaluate")
		os.Exit(8)
	}

	switch b.operator.GetType() {
	case tokentype.PLUS:
		return *object.NewObject(tokentype.NUMBER, l_num+r_num), nil
	case tokentype.MINUS:
		return *object.NewObject(tokentype.NUMBER, l_num-r_num), nil
	case tokentype.SLASH:
		return *object.NewObject(tokentype.NUMBER, l_num/r_num), nil
	case tokentype.STAR:
		return *object.NewObject(tokentype.NUMBER, l_num*r_num), nil
	case tokentype.GREATER:
		if l_num > r_num {
			return *object.NewObject(tokentype.TRUE, nil), nil
		}
		return *object.NewObject(tokentype.FALSE, nil), nil
	case tokentype.GREATER_EQUAL:
		if l_num >= r_num {
			return *object.NewObject(tokentype.TRUE, nil), nil
		}
		return *object.NewObject(tokentype.FALSE, nil), nil
	case tokentype.LESS:
		if l_num < r_num {
			return *object.NewObject(tokentype.TRUE, nil), nil
		}
		return *object.NewObject(tokentype.FALSE, nil), nil
	case tokentype.LESS_EQUAL:
		if l_num <= r_num {
			return *object.NewObject(tokentype.TRUE, nil), nil
		}
		return *object.NewObject(tokentype.FALSE, nil), nil
	}

	fault.RuntimeError("Eval Error: illegal binary operator", b.operator)
	return object.Object{}, errors.New("illegal binary operator")
}

// Grouping
type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression}
}
func (g Grouping) Evaluate(e *environment.Environment) (object.Object, error) {
	return g.expression.Evaluate(e)
}

// Variable
type Var struct {
	name token.Token
}

func NewVar(n token.Token) *Var {
	return &Var{n}
}

func (v *Var) GetToken() token.Token {
	return v.name
}

func (v Var) Evaluate(e *environment.Environment) (object.Object, error) {
	return e.Get(v.name)
}

// Assignment
type Assign struct {
	name  token.Token
	value Expr
}

func NewAssign(name token.Token, value Expr) *Assign {
	return &Assign{name, value}
}

func (a Assign) Evaluate(e *environment.Environment) (object.Object, error) {
	value, err := a.value.Evaluate(e)

	if err != nil {
		return object.Object{}, err
	}

	e.Assign(a.name, value)
	return value, nil
}
