package almond

import (
	"errors"
)

type Environment struct {
	enclosing *Environment
	lut       map[string]Object
}

func NewEnv() *Environment {
	lut := map[string]Object{}
	thisEnv := Environment{nil, lut}
	return &thisEnv
}

func NewEnclosedEnv(e *Environment) *Environment {
	lut := map[string]Object{}
	return &Environment{e, lut}
}

func (e *Environment) Define(name string, value Object) {
	e.lut[name] = value
}

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
