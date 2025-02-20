package almond

var GlobalEnv *Environment = NewEnv()

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{NewEnv()}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	DefineGlob()

	for _, statement := range statements {
		err := statement.Evaluate(i.env)

		if err != nil {
			break
		}
	}
}

// -- Helper function
func DefineGlob() {
	clockObj := NewObject(CALLABLE, NewNativeClockFn())
	GlobalEnv.Define("clock", *clockObj)
}
