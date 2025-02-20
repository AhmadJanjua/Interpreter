package almond

import (
	"fmt"
	"os"
	"strconv"
)

type Object struct {
	kind    TokenType
	literal any
}

// Ctor
func NewObject(k TokenType, v any) *Object {
	switch k {
	case STRING:
		val, ok := v.(string)
		if !ok {
			fmt.Println("Implmentation Error: Created a string object and passed non-string value.")
			os.Exit(1)
		}
		return &Object{k, val}

	case NUMBER:
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

// -- Truth value helpers
func (o *Object) Bool() bool {
	// 0, null, false
	switch o.kind {
	case NULL:
		return false
	case FALSE:
		return false
	case NUMBER:
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
	case NUMBER:
	case STRING:
		if right.literal != left.literal {
			return false
		}
	}

	// otherwise if the kind matches (true, false, null etc)
	return true
}

// -- Data retrieval helpers
func (o *Object) GetKind() TokenType {
	return o.kind
}

func (o *Object) GetKindStr() string {
	return o.kind.String()
}

func (o *Object) GetLiteral() any {
	return o.literal
}

// -- Formatting helpers
func (o *Object) String() string {
	switch o.kind {
	case STRING:
		s, ok := o.literal.(string)

		if !ok {
			fmt.Println("Implementation Error: failed to extract string from type String Object.")
			os.Exit(8)
		}
		return s
	case NUMBER:
		f, ok := o.literal.(float64)

		if !ok {
			fmt.Println("Implementation Error: failed to extract number from type Number Object.")
			os.Exit(9)
		}
		return strconv.FormatFloat(f, 'f', -1, 64)
	default:
		return o.kind.String()
	}
}
