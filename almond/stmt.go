package almond

import (
	"fmt"
)

type Stmt interface {
	Evaluate(e *Environment) error
}

// EXPRESSION STATEMENTS
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

// BLOCK STATMENTS
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

// PRINT STATEMENTS
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

// CONDITION STATEMENTS
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

// LOOPS FOR|WHILE
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

// RETURN STATEMENTS
type ReturnStmt struct {
	keyword Token
	value   Expr
}

func NewReturnStmt(key Token, val Expr) *ReturnStmt {
	return &ReturnStmt{key, val}
}

func (r ReturnStmt) Evaluate(e *Environment) error {
	value := *NewObject(NULL, nil)
	var err error = nil

	if r.value != nil {
		value, err = r.value.Evaluate(e)

		if err != nil {
			return err
		}
	}

	return &value
}

// FUNCTION STATEMENTS
type FnStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

func NewFnStmt(n Token, p []Token, b []Stmt) *FnStmt {
	return &FnStmt{n, p, b}
}

func (f FnStmt) Evaluate(e *Environment) error {
	function := NewFunctionCall(f)
	e.Define(f.name.lexeme, *NewObject(CALLABLE, function))
	return nil
}

// VARIABLE STATEMENTS
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
