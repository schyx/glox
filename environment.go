package main

import (
	"fmt"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func (env *Environment) assign(name Token, value any) error {
	_, ok := env.values[name.lexeme]
	if ok {
		env.values[name.lexeme] = value
		return nil
	}
	if env.enclosing != nil {
		err := env.enclosing.assign(name, value)
		return err
	}
	return fmt.Errorf("Undefined variable '%s'.", name.lexeme)
}

func (env *Environment) define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) get(name Token) (any, error) {
	val, ok := env.values[name.lexeme]
	if ok {
		return val, nil
	}
	if env.enclosing != nil {
		val, err := env.enclosing.get(name)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return nil, fmt.Errorf("Undefined variable '%s'.", name.lexeme)
}

func (env *Environment) getAt(distance int, name string) (any, error) {
	ancestor := env.ancestor(distance)
	val, ok := ancestor.values[name]
	if !ok {
		return nil, fmt.Errorf("Undefined variable '%s'.", name)
	}
	return val, nil
}

func (env *Environment) assignAt(distance int, name Token, value any) {
	env.ancestor(distance).values[name.lexeme] = value
}

func (env *Environment) ancestor(distance int) *Environment {
	output := env
	for range distance {
		output = output.enclosing
	}
	return output
}
