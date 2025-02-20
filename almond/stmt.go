package almond

import (
	"fmt"
)

type Stmt interface {
	Evaluate(e *Environment) error
}

type BlockStmt struct {
	statements []Stmt
}

func NewBlockStmt(s []Stmt) *BlockStmt {
	return &BlockStmt{s}
}

func (b BlockStmt) Evaluate(e *Environment) error {
	blockEnv := NewEnclosedEnv(e)

	for _, statement := range b.statements {
		err := statement.Evaluate(blockEnv)
		if err != nil {
			return err
		}
	}

	return nil
}

type ExprStmt struct {
	expression Expr
}

func NewExprStmt(e Expr) *ExprStmt {
	return &ExprStmt{e}
}

func (x ExprStmt) Evaluate(e *Environment) error {
	_, err := x.expression.Evaluate(e)

	return err
}

type PrintStmt struct {
	expression Expr
}

func NewPrintStmt(e Expr) *PrintStmt {
	return &PrintStmt{e}
}

func (p PrintStmt) Evaluate(e *Environment) error {
	value, err := p.expression.Evaluate(e)

	if err != nil {
		return err
	}
	fmt.Println(value.String())

	return nil
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func NewVarStmt(n Token, i Expr) *VarStmt {
	return &VarStmt{n, i}
}

func (v VarStmt) Evaluate(e *Environment) error {
	if v.initializer == nil {
		e.Define(v.name.GetLexeme(), *NewObject(NULL, nil))
		return nil
	}

	value, err := v.initializer.Evaluate(e)

	if err != nil {
		return err
	}

	e.Define(v.name.GetLexeme(), value)

	return nil
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIfStmt(c Expr, t Stmt, e Stmt) *IfStmt {
	return &IfStmt{c, t, e}
}

func (i IfStmt) Evaluate(e *Environment) error {
	cond, err := i.condition.Evaluate(e)

	if err != nil {
		return err
	}

	if cond.Bool() {
		err = i.thenBranch.Evaluate(e)

		if err != nil {
			return err
		}

	} else if i.elseBranch != nil {
		err = i.elseBranch.Evaluate(e)

		if err != nil {
			return err
		}
	}
	return nil
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func NewWhileStmt(c Expr, b Stmt) *WhileStmt {
	return &WhileStmt{c, b}
}

func (w WhileStmt) Evaluate(e *Environment) error {
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
