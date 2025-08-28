package main

import "fmt"

type LoxFunction struct {
	declaration Function
	env         *Environment
}

func (lf LoxFunction) call(interp *Interpreter, args []any) any {
	env := Environment{values: make(map[string]any), enclosing: lf.env}
	for i := range len(lf.declaration.params) {
		env.define(lf.declaration.params[i].lexeme, args[i])
	}
	interp.executeBlock(lf.declaration.body, &env)
	defer func() {
		interp.returnVal = nil
		interp.checkReturn = false
	}()
	return interp.returnVal
}

func (lf LoxFunction) arity() int {
	return len(lf.declaration.params)
}

func (lf LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", lf.declaration.name.lexeme)
}
