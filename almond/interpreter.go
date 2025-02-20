package almond

type Interpreter struct {
	env Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{*NewEnv()}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		err := statement.Evaluate(&i.env)

		if err != nil {
			break
		}
	}
}
