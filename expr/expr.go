package expr

import (
	"Interpreter/token"
	"Interpreter/tokentype"
)

type Expr interface {
	accept() string
}

// Literal
type Literal struct {
	kind  tokentype.TokenType
	value string
}

func NewLiteral(kind tokentype.TokenType, value string) *Literal {
	return &Literal{kind, value}
}
func (l Literal) accept() string {
	return l.kind.String() + "-" + l.value
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

// Binary
type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) *Binary {
	return &Binary{left, operator, right}
}
func (u Binary) accept() string {
	return "(" + u.operator.GetLexeme() + u.left.accept() + u.right.accept() + ")"
}

// Grouping
type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression}
}
func (u Grouping) accept() string {
	return "( group " + u.expression.accept() + ")"
}

// Print the Abstract Syntax Tree
func ASTPrinter(e Expr) string {
	return e.accept()
}
