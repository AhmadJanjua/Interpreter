package almond

import "time"

type Callable interface {
	Arity() int
	Call(args []Object) (Object, error)
	ToString() string
}

// Create a native clock for the program
type NativeClockFn struct{}

func NewNativeClockFn() *NativeClockFn {
	return &NativeClockFn{}
}

func (n *NativeClockFn) Arity() int {
	return 0
}

func (n *NativeClockFn) Call(args []Object) (Object, error) {
	time := float64(time.Now().Unix())
	timeObj := NewObject(NUMBER, time)

	return *timeObj, nil
}

func (n *NativeClockFn) ToString() string {
	return "<NATIVE CLK FN>"
}

// Create a user function callable
type FunctionCall struct {
	declaration FnStmt
}

func NewFunctionCall(d FnStmt) *FunctionCall {
	return &FunctionCall{d}
}

func (f *FunctionCall) Call(args []Object) (Object, error) {
	fnEnv := NewEnclosedEnv(GlobalEnv)

	for idx, param := range f.declaration.params {
		fnEnv.Define(param.lexeme, args[idx])
	}

	for _, statement := range f.declaration.body {
		err := statement.Evaluate(fnEnv)

		if err != nil {
			break
		}
	}

	return *NewObject(NULL, nil), nil
}

func (f *FunctionCall) Arity() int {
	return len(f.declaration.params)
}

func (f *FunctionCall) ToString() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}
