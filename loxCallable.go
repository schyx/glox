package main

type LoxCallable interface {
	arity() int
	call(*Interpreter, []any) any
}
