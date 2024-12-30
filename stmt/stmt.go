package stmt

import (
	"Interpreter/environment"
	"Interpreter/expr"
	"Interpreter/object"
	"Interpreter/token"
	"Interpreter/tokentype"
	"errors"
	"fmt"
)

type Stmt interface {
	Evaluate(e *environment.Environment) error
}

type Block struct {
	statements []Stmt
}

func NewBlock(s []Stmt) *Block {
	return &Block{s}
}

func (b Block) Evaluate(e *environment.Environment) error {
	previous := e.Copy()
	defer func() {
		e = previous
	}()

	env := environment.NewEnclosedEnv(e)

	for _, statement := range b.statements {
		err := statement.Evaluate(env)
		if err != nil {
			return err
		}
	}

	return nil
}

type Expression struct {
	expression expr.Expr
}

func NewExpression(e expr.Expr) *Expression {
	return &Expression{e}
}

func (x Expression) Evaluate(e *environment.Environment) error {
	return errors.New("err expression statement")
}

type Print struct {
	expression expr.Expr
}

func NewPrint(e expr.Expr) *Print {
	return &Print{e}
}

func (p Print) Evaluate(e *environment.Environment) error {
	value, err := p.expression.Evaluate(e)

	if err != nil {
		return err
	}
	fmt.Println(value.String())

	return nil
}

type Var struct {
	name        token.Token
	initializer expr.Expr
}

func NewVar(n token.Token, i expr.Expr) *Var {
	return &Var{n, i}
}

func (v Var) Evaluate(e *environment.Environment) error {
	if v.initializer == nil {
		e.Define(v.name.GetLexeme(), *object.NewObject(tokentype.NULL, nil))
		return nil
	}

	value, err := v.initializer.Evaluate(e)

	if err != nil {
		return err
	}

	e.Define(v.name.GetLexeme(), value)

	return nil
}
