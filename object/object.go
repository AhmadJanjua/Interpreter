package object

import (
	"Interpreter/tokentype"
	"fmt"
	"os"
)

type Object struct {
	kind    tokentype.TokenType
	literal any
}

// Create object and ensure types are managed correctly
func NewObject(k tokentype.TokenType, v any) *Object {
	switch k {
	case tokentype.STRING:
		val, ok := v.(string)
		if !ok {
			fmt.Println("Implmentation Error: Created a string object and passed non-string value.")
			os.Exit(1)
		}
		return &Object{k, val}

	case tokentype.NUMBER:
		val, ok := v.(float64)
		if !ok {
			fmt.Println("Implmentation Error: Created a number object and passed non-float64 value.")
			os.Exit(2)
		}
		return &Object{k, val}
	default:
		return &Object{k, nil}
	}
}

// Get the truth value of an object
func (o *Object) Bool() bool {
	// 0, null, false
	switch o.kind {
	case tokentype.NULL:
		return false
	case tokentype.FALSE:
		return false
	case tokentype.NUMBER:
		num, ok := o.literal.(float64)

		if !ok {
			fmt.Println("Implementation Error: Created an invalid number object. Found while evaluating truth value in Object.")
		}

		if num == 0 {
			return false
		}
	}

	// all other values are true
	return true
}

func (left *Object) Equal(right *Object) bool {
	if left.kind != right.kind {
		return false
	}

	switch left.kind {
	case tokentype.NUMBER:
	case tokentype.STRING:
		if right.literal != left.literal {
			return false
		}
	}

	// otherwise if the kind matches (true, false, null etc)
	return true
}

func (o *Object) GetKind() tokentype.TokenType {
	return o.kind
}

func (o *Object) GetKindStr() string {
	return o.kind.String()
}

func (o *Object) GetLiteral() any {
	return o.literal
}

func (o *Object) String() string {
	switch o.kind {
	case tokentype.STRING:
		s, ok := o.literal.(string)

		if !ok {
			fmt.Println("Implementation Error: failed to extract string from type String Object.")
			os.Exit(8)
		}
		return s
	case tokentype.NUMBER:
		f, ok := o.literal.(float64)

		if !ok {
			fmt.Println("Implementation Error: failed to extract number from type Number Object.")
			os.Exit(9)
		}
		return fmt.Sprintf("%f", f)
	default:
		return o.kind.String()
	}
}
