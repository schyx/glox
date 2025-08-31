package main

import "fmt"

type LoxClass struct {
	name       string
	methods    map[string]LoxFunction
	superclass *LoxClass
}

func (lc LoxClass) findMethod(name string) (LoxFunction, error) {
	method, ok := lc.methods[name]
	if !ok {
		if lc.superclass != nil {
			return lc.superclass.findMethod(name)
		} else {
			return LoxFunction{}, fmt.Errorf("No method found.")
		}
	} else {
		return method, nil
	}
}

func (lc LoxClass) String() string {
	return lc.name
}

func (lc LoxClass) call(interp *Interpreter, arguments []any) any {
	instance := LoxInstance{klass: lc, fields: make(map[string]any)}
	initializer, err := lc.findMethod("init")
	if err == nil { // user provided constructor
		initializer.bind(instance).call(interp, arguments)
	}
	return instance
}

func (lc LoxClass) arity() int {
	intializer, err := lc.findMethod("init")
	if err == nil {
		return intializer.arity()
	} else {
		return 0
	}
}
