package main

import "time"

type Clock struct{}

func (c Clock) arity() int {
	return 0
}

func (c Clock) call(_interp *Interpreter, _args []any) any {
	return time.Now().Unix()
}

func (c Clock) String() string {
	return "<native fn>"
}
