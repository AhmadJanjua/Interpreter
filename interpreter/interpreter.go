package interpreter

import (
	"Interpreter/environment"
	"Interpreter/fault"
	"Interpreter/stmt"
	"Interpreter/token"
)

type Interpreter struct {
	env environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{*environment.NewEnv()}
}

func (i *Interpreter) Interpret(statements []stmt.Stmt) {
	for _, statement := range statements {
		// TODO: Error handling
		err := statement.Evaluate(&i.env)

		if err != nil {
			fault.RuntimeError("", token.Token{})
		}
	}
}
