package faunadb

import (
	"encoding/json"
)

/*
Expr is the base type for FaunaDB query language expressions.

Expressions are created by using the query language functions in query.go. Query functions are designed to compose with each other, as well as with
custom data structures. For example:

	type User struct {
		Name string
	}

	_, _ := client.Query(
		Create(
			Collection("users"),
			Obj{"data": User{"John"}},
		),
	)

*/
type Expr interface {
	expr() // Make sure only internal structures can be marked as valid expressions
	String() string
}

type unescapedObj map[string]Expr
type unescapedArr []Expr
type invalidExpr struct{ err error }

func (obj unescapedObj) expr()          {}
func (obj unescapedObj) String() string { byte, _ := json.Marshal(obj); return string(byte) }

func (arr unescapedArr) expr()          {}
func (arr unescapedArr) String() string { byte, _ := json.Marshal(arr); return string(byte) }

func (inv invalidExpr) expr()          {}
func (inv invalidExpr) String() string { byte, _ := inv.MarshalJSON(); return string(byte) }

func (inv invalidExpr) MarshalJSON() ([]byte, error) {
	return nil, inv.err
}

// Obj is a expression shortcut to represent any valid JSON object
type Obj map[string]interface{}

// Arr is a expression shortcut to represent any valid JSON array
type Arr []interface{}

func (obj Obj) expr()          {}
func (obj Obj) String() string { byte, _ := obj.MarshalJSON(); return string(byte) }

func (arr Arr) expr()          {}
func (arr Arr) String() string { byte, _ := arr.MarshalJSON(); return string(byte) }

// MarshalJSON implements json.Marshaler for Obj expression
func (obj Obj) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(obj)) }

// MarshalJSON implements json.Marshaler for Arr expression
func (arr Arr) MarshalJSON() ([]byte, error) { return json.Marshal(wrap(arr)) }

// OptionalParameter describes optional parameters for query language functions
type OptionalParameter func(unescapedObj)

func applyOptionals(options []OptionalParameter, fn unescapedObj) Expr {
	for _, option := range options {
		option(fn)
	}
	return fn
}

func fn1(k1 string, v1 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
	})
}

func fn2(k1 string, v1 interface{}, k2 string, v2 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
		k2: wrap(v2),
	})
}

func fn3(k1 string, v1 interface{}, k2 string, v2 interface{}, k3 string, v3 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
		k2: wrap(v2),
		k3: wrap(v3),
	})
}

func fn4(k1 string, v1 interface{}, k2 string, v2 interface{}, k3 string, v3 interface{}, k4 string, v4 interface{}, options ...OptionalParameter) Expr {
	return applyOptionals(options, unescapedObj{
		k1: wrap(v1),
		k2: wrap(v2),
		k3: wrap(v3),
		k4: wrap(v4),
	})
}
