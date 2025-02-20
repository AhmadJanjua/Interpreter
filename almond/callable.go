package almond

import (
	"errors"
	"fmt"
	"time"
)

type Callable interface {
	Arity() int
	Call(env Environment, args []Object) (Object, error)
	ToString() string
}

// Create a native clock for the program
type NativeClock struct {
	start time.Time
}

func NewNativeClock() *NativeClock {
	return &NativeClock{time.Now()}
}

func (n *NativeClock) Arity() int {
	return 0
}

func (n *NativeClock) Call(env Environment, args []Object) (Object, error) {
	elapsed := float64(time.Since(n.start).Nanoseconds()) / 1e9
	timeObj := NewObject(NUMBER, elapsed)

	return *timeObj, nil
}

func (n *NativeClock) ToString() string {
	return "<NATIVE CLK FN>"
}

// Create Native sleep function
type NativeSleep struct{}

func NewNativeSleep() *NativeSleep      { return &NativeSleep{} }
func (n *NativeSleep) Arity() int       { return 1 }
func (n *NativeSleep) ToString() string { return "<NATIVE SLEEP FN>" }
func (n *NativeSleep) Call(env Environment, args []Object) (Object, error) {
	arg := args[0]

	dt_ms, ok := arg.literal.(float64)

	if !ok {
		fmt.Print("sleepMS usage error: must supply a number!")
		return Object{}, errors.New("sleep argument must be type number")
	}

	time.Sleep(time.Duration(dt_ms) * time.Millisecond)
	return *NewObject(NULL, nil), nil
}

// Create a user function callable
type FunctionCall struct {
	declaration FnStmt
}

func NewFunctionCall(d FnStmt) *FunctionCall {
	return &FunctionCall{d}
}

func (f *FunctionCall) Call(env Environment, args []Object) (Object, error) {
	// copy environment
	copy := env

	fnEnv := NewEnclosedEnv(&copy)
	var err error = nil

	for idx, param := range f.declaration.params {
		fnEnv.Define(param.lexeme, args[idx])
	}

	for _, statement := range f.declaration.body {
		err = statement.Evaluate(fnEnv)

		if err != nil {
			break
		}
	}

	if err != nil {
		val, ok := err.(*Object)

		if ok {
			return *val, nil
		}
		fmt.Print(err)
	}

	return *NewObject(NULL, nil), nil
}

func (f *FunctionCall) Arity() int {
	return len(f.declaration.params)
}

func (f *FunctionCall) ToString() string {
	return "<fn " + f.declaration.name.lexeme + ">"
}
