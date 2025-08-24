package main

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	output   any
	err      error
	badToken Token
	lx       *Lox
}

func (interp *Interpreter) Interpret(expr Expr) {
	interp.evaluate(expr)
	if interp.err != nil {
		interp.lx.RuntimeError(interp.badToken, interp.err)
		return
	}
	fmt.Printf("%v\n", interp.output)
}

func (interp *Interpreter) visitBinary(expr Binary) {
	left, err, token := evalExpr(expr.left)
	if err != nil {
		interp.output = 0
		interp.err = err
		interp.badToken = token
		return
	}
	right, err, token := evalExpr(expr.right)
	if err != nil {
		interp.output = 0
		interp.err = err
		interp.badToken = token
		return
	}

	switch expr.operator.tokenType {
	case BANG_EQUAL:
		interp.output = !isEqual(left, right)
		interp.err = nil
	case EQUAL_EQUAL:
		interp.output = isEqual(left, right)
		interp.err = nil
	case GREATER:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal > rightVal, err, expr.operator)
	case GREATER_EQUAL:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal >= rightVal, err, expr.operator)
	case LESS:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal < rightVal, err, expr.operator)
	case LESS_EQUAL:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal <= rightVal, err, expr.operator)
	case MINUS:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal - rightVal, err, expr.operator)
	case SLASH:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal / rightVal, err, expr.operator)
	case STAR:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal * rightVal, err, expr.operator)
	case PLUS:
		// numeric case
		leftVal, rightVal, err := toFloatPair(left, right)
		if err == nil {
			interp.getReturnVal(leftVal + rightVal, err, expr.operator)
			return
		}
		// string case
		leftString, rightString, err := toStringPair(left, right)
		interp.getReturnVal(leftString + rightString, err, expr.operator)
	}
}

func (interp *Interpreter) getReturnVal(okVal any, err error, badToken Token) {
	interp.err = err
	interp.badToken = badToken
	interp.output = okVal
}

func isEqual(left any, right any) bool {
	if left == nil && right == nil {
		return true
	} else if left == nil {
		return false
	}
	return left == right
}

func toStringPair(left any, right any) (string, string, error) {
	leftString, leftErr := toString(left)
	rightString, rightErr := toString(right)
	if leftErr != nil {
		return "", "", leftErr
	} else if rightErr != nil {
		return "", "", rightErr
	} else {
		return leftString, rightString, nil
	}
}

func toFloatPair(left any, right any) (float64, float64, error) {
	leftVal, leftErr := toFloat(left)
	rightVal, rightErr := toFloat(right)
	if leftErr != nil {
		return 0, 0, leftErr
	} else if rightErr != nil {
		return 0, 0, rightErr
	} else {
		return leftVal, rightVal, nil
	}
}

func (interp *Interpreter) visitGrouping(expr Grouping) {
	interp.output, interp.err, interp.badToken = evalExpr(expr.expression)
}

func (interp *Interpreter) visitLiteral(expr Literal) {
	interp.output = expr.value
}

func (interp *Interpreter) visitUnary(expr Unary) {
	right, err, token := evalExpr(expr.right)
	if err != nil {
		interp.output = 0
		interp.err = err
		interp.badToken = token
		return
	}
	switch expr.operator.tokenType {
	case BANG:
		interp.output = !isTruthy(right)
	case MINUS:
		val, err := toFloat(right)
		interp.output = -val
		interp.err = err
		interp.badToken = expr.operator
	}
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	switch t := object.(type) {
	case bool:
		return t
	default:
		return true
	}
}

func (interp *Interpreter) evaluate(expr Expr) {
	expr.accept(interp)
}

func evalExpr(expr Expr) (any, error, Token) {
	dummyInterp := Interpreter{}
	dummyInterp.evaluate(expr)
	return dummyInterp.output, dummyInterp.err, dummyInterp.badToken
}

func toFloat(val any) (float64, error) {
	switch v := val.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	}
	return 0, errors.New("Converting non-numeric value to float")
}

func toString(val any) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	default:
		return "", errors.New("Converting non-string value to string")
	}
}
