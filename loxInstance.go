package main

import "fmt"

type LoxInstance struct {
	klass  LoxClass
	fields map[string]any
}

func (li LoxInstance) get(name Token) (any, error) {
	val, ok := li.fields[name.lexeme]
	if ok {
		return val, nil
	}
	method, methodErr := li.klass.findMethod(name.lexeme)
	if methodErr != nil {
		return nil, fmt.Errorf("Undefined property '%s'.",  name.lexeme)
	} else {
		return method.bind(li), nil
	}
}

func (li LoxInstance) set(name Token, value any) {
	li.fields[name.lexeme] = value
}

func (li LoxInstance) String() string {
	return li.klass.name + " instance"
}
