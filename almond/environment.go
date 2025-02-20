package almond

import (
	"errors"
)

type Environment struct {
	enclosing *Environment
	lut       map[string]Object
}

// Ctor
func NewEnv() *Environment {
	lut := map[string]Object{}
	thisEnv := Environment{nil, lut}
	return &thisEnv
}

// Ctor with existing env
func NewEnclosedEnv(e *Environment) *Environment {
	lut := map[string]Object{}
	return &Environment{e, lut}
}

// ---- Functions

// store variables or functions
func (e *Environment) Define(name string, value Object) {
	e.lut[name] = value
}

// retrieve variables or functions
func (e *Environment) Get(tok Token) (Object, error) {
	name := tok.GetLexeme()
	value, ok := e.lut[name]

	if ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(tok)
	}

	RuntimeError("Undefined variable '"+name+"'.", tok)
	return Object{}, errors.New("undefined variable error")
}

// update variables or function
func (e *Environment) Assign(tok Token, value Object) error {
	name := tok.GetLexeme()
	_, ok := e.lut[name]

	if ok {
		e.lut[name] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(tok, value)
	}

	RuntimeError("Undefined variable '"+name+"'.", tok)
	return errors.New("undefined variable error")
}
