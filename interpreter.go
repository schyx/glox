package main

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	output   any
	err      error
	badToken Token
	env      *Environment
	lx       *Lox
}

func (interp *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		_, err, badToken := execStmt(statement, interp.env)
		if err != nil {
			interp.lx.RuntimeError(badToken, err)
			return
		}
	}
}

func (interp *Interpreter) visitBinary(expr Binary) {
	left, err, token := evalExpr(expr.left, interp.env)
	if err != nil {
		interp.output = nil
		interp.err = err
		interp.badToken = token
		return
	}
	right, err, token := evalExpr(expr.right, interp.env)
	if err != nil {
		interp.output = nil
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
		interp.getReturnVal(leftVal-rightVal, err, expr.operator)
	case SLASH:
		leftVal, rightVal, err := toFloatPair(left, right)
		if err == nil && rightVal == 0 {
			interp.getReturnVal(0, errors.New("Dividing by zero"), expr.operator)
			return
		}
		interp.getReturnVal(leftVal/rightVal, err, expr.operator)
	case STAR:
		leftVal, rightVal, err := toFloatPair(left, right)
		interp.getReturnVal(leftVal*rightVal, err, expr.operator)
	case PLUS:
		// numeric case
		leftVal, rightVal, err := toFloatPair(left, right)
		if err == nil {
			interp.getReturnVal(leftVal+rightVal, err, expr.operator)
			return
		}
		// string case
		leftString, rightString, err := toStringPair(left, right)
		interp.getReturnVal(leftString+rightString, err, expr.operator)
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
	interp.output, interp.err, interp.badToken = evalExpr(expr.expression, interp.env)
}

func (interp *Interpreter) visitLiteral(expr Literal) {
	interp.output = expr.value
}

func (interp *Interpreter) visitLogical(expr Logical) {
	left, leftErr, leftBadToken := evalExpr(expr.left, interp.env)
	if leftErr != nil {
		interp.err = leftErr
		interp.badToken = leftBadToken
		return
	}
	if expr.operator.tokenType == OR {
		if isTruthy(left) {
			interp.output = left
			return
		}
	} else {
		if !isTruthy(left) {
			interp.output = left
			return
		}
	}

	right, rightErr, rightBadToken := evalExpr(expr.right, interp.env)
	if rightErr != nil {
		interp.err = rightErr
		interp.badToken = rightBadToken
		return
	}
	interp.output = right
}

func (interp *Interpreter) visitUnary(expr Unary) {
	right, err, token := evalExpr(expr.right, interp.env)
	if err != nil {
		interp.output = nil
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

func evalExpr(expr Expr, env *Environment) (any, error, Token) {
	dummyInterp := Interpreter{env: env}
	dummyInterp.evaluate(expr)
	return dummyInterp.output, dummyInterp.err, dummyInterp.badToken
}

func (interp *Interpreter) execute(stmt Stmt) {
	stmt.accept(interp)
}

func execStmt(stmt Stmt, env *Environment) (any, error, Token) {
	dummyInterp := Interpreter{env: env}
	dummyInterp.execute(stmt)
	return dummyInterp.output, dummyInterp.err, dummyInterp.badToken
}

func (interp *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	for _, statement := range statements {
		_, err, badToken := execStmt(statement, env)
		if interp.err != nil {
			interp.err = err
			interp.badToken = badToken
			return
		}
	}
}

func (interp *Interpreter) visitBlock(stmt Block) {
	interp.executeBlock(stmt.statments, &Environment{values: make(map[string]any), enclosing: interp.env})
}

func (interp *Interpreter) visitExpression(stmt Expression) {
	_, err, badToken := evalExpr(stmt.expr, interp.env)
	if err != nil {
		interp.err = err
		interp.badToken = badToken
		return
	}
}

func (interp *Interpreter) visitIf(stmt If) {
	conditionVal, conditionErr, conditionBadToken := evalExpr(stmt.condition, interp.env)
	if conditionErr != nil {
		interp.err = conditionErr
		interp.badToken = conditionBadToken
		return
	}
	if isTruthy(conditionVal) {
		_, execErr, execBadToken := execStmt(stmt.thenBranch, interp.env)
		if execErr != nil {
			interp.err = execErr
			interp.badToken = execBadToken
			return
		}
	} else {
		_, execErr, execBadToken := execStmt(stmt.elseBranch, interp.env)
		if execErr != nil {
			interp.err = execErr
			interp.badToken = execBadToken
			return
		}
	}
}

func (interp *Interpreter) visitPrint(stmt Print) {
	val, err, badToken := evalExpr(stmt.expr, interp.env)
	if err != nil {
		interp.err = err
		interp.badToken = badToken
		return
	}
	fmt.Printf("%v\n", val)
}

func (interp *Interpreter) visitVar(stmt Var) {
	var value any
	if stmt.initializer != nil {
		val, err, badToken := evalExpr(stmt.initializer, interp.env)
		if err != nil {
			interp.err = err
			interp.badToken = badToken
			return
		}
		value = val
	}
	interp.env.define(stmt.name.lexeme, value)
}

func (interp *Interpreter) visitWhile(stmt While) {
	for {
		conditionVal, conditionErr, conditionBadToken := evalExpr(stmt.condition, interp.env)
		if conditionErr != nil {
			interp.err = conditionErr
			interp.badToken = conditionBadToken
			return
		}
		if !isTruthy(conditionVal) {
			return
		}
		_, bodyErr, bodyBadToken := execStmt(stmt.body, interp.env)
		if bodyErr != nil {
			interp.err = bodyErr
			interp.badToken = bodyBadToken
			return
		}
	}
}

func (interp *Interpreter) visitAssign(expr Assign) {
	value, err, badToken := evalExpr(expr.value, interp.env)
	if err != nil {
		interp.output = NIL
		interp.err = err
		interp.badToken = badToken
		return
	}
	assignErr := interp.env.assign(expr.name, value)
	if assignErr != nil {
		interp.err = err
		interp.badToken = expr.name
		return
	}
	interp.output = value
}

func (interp *Interpreter) visitVariable(expr Variable) {
	val, err := interp.env.get(expr.name)
	if err != nil {
		interp.output = nil
		interp.err = err
		interp.badToken = expr.name
		return
	}
	interp.output = val
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
