package interpreter

import (
	"Interpreter/expr"
	"fmt"
)

func Interpret(e expr.Expr) {
	obj, err := e.Evaluate()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(obj.String())
}
