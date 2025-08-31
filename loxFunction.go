package main

import "fmt"

type LoxFunction struct {
	declaration   Function
	env           *Environment
	isInitializer bool
}

func (lf LoxFunction) bind(li LoxInstance) LoxFunction {
	env := Environment{values: make(map[string]any), enclosing: lf.env}
	env.define("this", li)
	return LoxFunction{declaration: lf.declaration, env: &env, isInitializer: lf.isInitializer}
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
	if lf.isInitializer {
		output, _ := lf.env.getAt(0, "this")
		return output
	} else {
		return interp.returnVal
	}
}

func (lf LoxFunction) arity() int {
	return len(lf.declaration.params)
}

func (lf LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", lf.declaration.name.lexeme)
}
