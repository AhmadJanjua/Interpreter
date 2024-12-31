package interpreter

import (
	"Interpreter/environment"
	"Interpreter/stmt"
)

type Interpreter struct {
	env environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{*environment.NewEnv()}
}

func (i *Interpreter) Interpret(statements []stmt.Stmt) {
	for _, statement := range statements {
		err := statement.Evaluate(&i.env)

		if err != nil {
			break
		}
	}
}
