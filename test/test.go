package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Data represents a test case, generic type T is the input type, generic type Y is the want type.
type Data[T, Y any] struct {
	Name  string
	Input T
	Want  Y
	E     error
}

// Option represents an option for Executor, generic type T is the input type, generic type Y is the want type.
type Option[T, Y any] func(*Executor[T, Y])
type assertFn[Y any] func(t *testing.T, expected, actual Y)

// WithComparison is an option to set the comparison function,
// generic type T is the input type, generic type Y is the output type.
func WithComparison[T, Y any](comparisonFn assertFn[Y]) Option[T, Y] {
	return func(e *Executor[T, Y]) {
		e.equalFn = comparisonFn
	}
}

// Executor manages and executes test cases, generic type T is the input type, generic type Y is the want type.
type Executor[T, Y any] struct {
	list    []Data[T, Y]
	equalFn assertFn[Y]
}

// NewExecutor creates an Executor, generic type T is the input type, generic type Y is the want type.
func NewExecutor[T, Y any](opt ...Option[T, Y]) *Executor[T, Y] {
	e := &Executor[T, Y]{}
	options := []Option[T, Y]{WithComparison[T, Y](func(t *testing.T, expected, actual Y) {
		gotBytes, err := json.Marshal(actual)
		if err != nil {
			t.Fatal(err)
		}
		wantBytes, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}
		assert.JSONEq(t, string(wantBytes), string(gotBytes))
	})}

	options = append(options, opt...)

	for _, o := range options {
		o(e)
	}
	return e
}

// Add adds test cases to the Executor, generic type T is the input type, generic type Y is the want type.
func (e *Executor[T, Y]) Add(data ...Data[T, Y]) {
	e.list = append(e.list, data...)
}

// Run executes the test cases without error response, generic type T is the input type, generic type Y is the want type.
func (e *Executor[T, Y]) Run(t *testing.T, do func(T) Y) {
	if do == nil {
		panic("execution body is nil")
		return
	}
	for _, v := range e.list {
		t.Run(v.Name, func(t *testing.T) {
			inner := do
			e.equalFn(t, v.Want, inner(v.Input))
		})
	}
}

// RunE executes the test cases with error response, generic type T is the input type, generic type Y is the want type.
func (e *Executor[T, Y]) RunE(t *testing.T, do func(T) (Y, error)) {
	if do == nil {
		panic("execution body is nil")
		return
	}
	for _, v := range e.list {
		t.Run(v.Name, func(t *testing.T) {
			inner := do
			got, err := inner(v.Input)
			if err != nil {
				assert.Equal(t, v.E, err)
				return
			}
			e.equalFn(t, v.Want, got)
		})
	}
}
