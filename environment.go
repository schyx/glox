package main

import (
	"fmt"
)

type Environment struct {
	values map[string]any
	lx     *Lox
}

func (env *Environment) define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) get(name Token) (any, error) {
	val, ok := env.values[name.lexeme]
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf("Undefined variable '%s'.", name.lexeme)
}
