package stmt

import (
	"Interpreter/environment"
	"Interpreter/expr"
	"Interpreter/object"
	"Interpreter/token"
	"Interpreter/tokentype"
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
	block_env := environment.NewEnclosedEnv(e)

	for _, statement := range b.statements {
		err := statement.Evaluate(block_env)
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
	_, err := x.expression.Evaluate(e)

	return err
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

type If struct {
	condition   expr.Expr
	then_branch Stmt
	else_branch Stmt
}

func NewIf(c expr.Expr, t Stmt, e Stmt) *If {
	return &If{c, t, e}
}

func (i If) Evaluate(e *environment.Environment) error {
	cond, err := i.condition.Evaluate(e)

	if err != nil {
		return err
	}

	if cond.Bool() {
		err = i.then_branch.Evaluate(e)

		if err != nil {
			return err
		}

	} else if i.else_branch != nil {
		err = i.else_branch.Evaluate(e)

		if err != nil {
			return err
		}
	}
	return nil
}

type While struct {
	condition expr.Expr
	body      Stmt
}

func NewWhile(c expr.Expr, b Stmt) *While {
	return &While{c, b}
}

func (w While) Evaluate(e *environment.Environment) error {
	val, err := w.condition.Evaluate(e)

	if err != nil {
		return err
	}

	for val.Bool() {
		err = w.body.Evaluate(e)

		if err != nil {
			return err
		}

		val, err = w.condition.Evaluate(e)

		if err != nil {
			return err
		}
	}

	return nil
}
