package environment

import (
	"Interpreter/fault"
	"Interpreter/object"
	"Interpreter/token"
	"errors"
)

type Environment struct {
	enclosing *Environment
	lut       map[string]object.Object
}

func NewEnv() *Environment {
	lut := map[string]object.Object{}
	return &Environment{nil, lut}
}

func NewEnclosedEnv(e *Environment) *Environment {
	lut := map[string]object.Object{}
	return &Environment{e, lut}
}

func (e *Environment) Define(name string, value object.Object) {
	e.lut[name] = value
}

func (e *Environment) Get(tok token.Token) (object.Object, error) {
	name := tok.GetLexeme()
	value, ok := e.lut[name]

	if ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(tok)
	}

	fault.RuntimeError("Undefined variable '"+name+"'.", tok)
	return object.Object{}, errors.New("undefined variable error")
}

func (e *Environment) Assign(tok token.Token, value object.Object) error {
	name := tok.GetLexeme()
	_, ok := e.lut[name]

	if ok {
		e.lut[name] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(tok, value)
	}

	fault.RuntimeError("Undefined variable '"+name+"'.", tok)
	return errors.New("undefined variable error")
}

func (e *Environment) Copy() *Environment {
	new_lut := make(map[string]object.Object)
	for k, v := range e.lut {
		new_lut[k] = v
	}
	var enclosing *Environment
	if e.enclosing != nil {
		enclosing = e.enclosing.Copy()
	}
	return &Environment{
		enclosing: enclosing,
		lut:       new_lut,
	}
}
