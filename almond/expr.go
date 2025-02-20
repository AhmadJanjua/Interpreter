package almond

import (
	"errors"
	"fmt"
	"os"
)

type Expr interface {
	Evaluate(e *Environment) (Object, error)
}

// Literals
type Literal struct {
	value Object
}

func NewLiteral(t TokenType) *Literal {
	return &Literal{*NewObject(t, nil)}
}

func NewNumber(f float64) *Literal {
	return &Literal{*NewObject(NUMBER, f)}
}

func NewString(s string) *Literal {
	return &Literal{*NewObject(STRING, s)}
}
func (l Literal) Evaluate(e *Environment) (Object, error) {
	return l.value, nil
}

// UnaryExpr
type UnaryExpr struct {
	operator Token
	right    Expr
}

func NewUnaryExpr(operator Token, right Expr) *UnaryExpr {
	return &UnaryExpr{operator, right}
}

func (u UnaryExpr) Evaluate(e *Environment) (Object, error) {
	right, err := u.right.Evaluate(e)

	if err != nil {
		return right, err
	}

	switch u.operator.GetType() {
	case BANG:
		// if the value is true -> return false
		if right.Bool() {
			return *NewObject(FALSE, nil), nil
		} else {
			// true -> false
			return *NewObject(TRUE, nil), nil
		}
	case MINUS:
		// negate number
		if right.GetKind() == NUMBER {
			val, ok := right.GetLiteral().(float64)

			if !ok {
				fmt.Println("Implementation Error: Failed to parse float form number in expr -> unary -> Evaluate")
				os.Exit(6)
			}
			return *NewObject(NUMBER, -1*val), nil
		}
	}

	// Report error
	RuntimeError("Eval Error: illegal unary operator", u.operator)
	return Object{}, errors.New("illegal unary operator")
}

// BinaryExpr
type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func NewBinaryExpr(left Expr, operator Token, right Expr) *BinaryExpr {
	return &BinaryExpr{left, operator, right}
}
func (b BinaryExpr) Evaluate(e *Environment) (Object, error) {
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
	if b.operator.GetType() == EQUALS {
		if left.Equal(&right) {
			return *NewObject(TRUE, nil), nil
		}
		return *NewObject(FALSE, nil), nil
	}
	if b.operator.GetType() == NOT_EQUALS {
		if !left.Equal(&right) {
			return *NewObject(TRUE, nil), nil
		}
		return *NewObject(FALSE, nil), nil
	}

	// Report type missmatch for non equality tests
	if right.GetKind() != left.GetKind() {
		RuntimeError("Eval Error: type mismatch between "+left.GetKindStr()+" and "+right.GetKindStr(), b.operator)
		return Object{}, errors.New("type mismatch in binary operation")
	}

	// String concatenation
	if right.GetKind() == STRING && b.operator.GetType() == PLUS {
		lStr, lOk := left.GetLiteral().(string)
		rStr, rOk := right.GetLiteral().(string)

		if lOk && rOk {
			return *NewObject(STRING, lStr+rStr), nil
		}

		fmt.Println("Implementation Error: Could not parse string passed in expr -> binary -> Evaluate")
		os.Exit(7)
	}

	// Report invalid non-numeric operations
	if right.GetKind() != NUMBER {
		RuntimeError("Eval Error: invalid binary non-numeric operation", b.operator)
		return Object{}, errors.New("invalid binary operation on non-numeric")
	}
	lNum, lOk := left.GetLiteral().(float64)
	rNum, rOk := right.GetLiteral().(float64)

	if !lOk || !rOk {
		fmt.Println("Implementation Error: Could not parse numbers passed in expr -> binary -> Evaluate")
		os.Exit(8)
	}

	switch b.operator.GetType() {
	case PLUS:
		return *NewObject(NUMBER, lNum+rNum), nil
	case MINUS:
		return *NewObject(NUMBER, lNum-rNum), nil
	case SLASH:
		return *NewObject(NUMBER, lNum/rNum), nil
	case STAR:
		return *NewObject(NUMBER, lNum*rNum), nil
	case GREATER:
		if lNum > rNum {
			return *NewObject(TRUE, nil), nil
		}
		return *NewObject(FALSE, nil), nil
	case GREATER_EQUAL:
		if lNum >= rNum {
			return *NewObject(TRUE, nil), nil
		}
		return *NewObject(FALSE, nil), nil
	case LESS:
		if lNum < rNum {
			return *NewObject(TRUE, nil), nil
		}
		return *NewObject(FALSE, nil), nil
	case LESS_EQUAL:
		if lNum <= rNum {
			return *NewObject(TRUE, nil), nil
		}
		return *NewObject(FALSE, nil), nil
	}

	RuntimeError("Eval Error: illegal binary operator", b.operator)
	return Object{}, errors.New("illegal binary operator")
}

// GroupingExpr
type GroupingExpr struct {
	expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{expression}
}
func (g GroupingExpr) Evaluate(e *Environment) (Object, error) {
	return g.expression.Evaluate(e)
}

// Variable
type VarExpr struct {
	name Token
}

func NewVarExpr(n Token) *VarExpr {
	return &VarExpr{n}
}

func (v *VarExpr) GetToken() Token {
	return v.name
}

func (v VarExpr) Evaluate(e *Environment) (Object, error) {
	return e.Get(v.name)
}

// Assignment
type AssignExpr struct {
	name  Token
	value Expr
}

func NewAssignExpr(name Token, value Expr) *AssignExpr {
	return &AssignExpr{name, value}
}

func (a AssignExpr) Evaluate(e *Environment) (Object, error) {
	value, err := a.value.Evaluate(e)

	if err != nil {
		return value, err
	}

	err = e.Assign(a.name, value)

	if err != nil {
		return Object{}, err
	}

	return value, nil
}

type LogicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func NewLogicalExpr(left Expr, operator Token, right Expr) *LogicalExpr {
	return &LogicalExpr{left, operator, right}
}

func (a LogicalExpr) Evaluate(e *Environment) (Object, error) {
	left, err := a.left.Evaluate(e)

	if err != nil {
		return left, err
	}

	if a.operator.GetType() == OR {
		if left.Bool() {
			return left, nil
		}
	} else {
		if !left.Bool() {
			return left, nil
		}
	}
	return a.right.Evaluate(e)
}
