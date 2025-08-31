package main

import (
	"errors"
	"fmt"
)

// --------------- INTERPRETER ---------------

type Interpreter struct {
	output      any
	err         error
	badToken    Token
	checkReturn bool
	returnVal   any
	locals      map[Expr]int
	env         *Environment
	lx          *Lox
}

func (interp *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		_, _, err, _ := execStmt(statement, interp.env, interp.locals, interp.lx)
		if err != nil {
			return
		}
	}
}

func (interp *Interpreter) resolve(expr Expr, depth int) {
	interp.locals[expr] = depth
}

// --------------- STATEMENTS ---------------

func (interp *Interpreter) execute(stmt Stmt) {
	stmt.accept(interp)
}

func execStmt(stmt Stmt, env *Environment, locals map[Expr]int, lx *Lox) (any, bool, error, Token) {
	dummyInterp := &Interpreter{env: env, locals: locals, lx: lx}
	dummyInterp.execute(stmt)
	return dummyInterp.returnVal, dummyInterp.checkReturn, dummyInterp.err, dummyInterp.badToken
}

func (interp *Interpreter) visitBlock(stmt Block) {
	interp.executeBlock(stmt.statments, &Environment{values: make(map[string]any), enclosing: interp.env})
}

func (interp *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	for _, statement := range statements {
		returnVal, checkReturn, err, badToken := execStmt(statement, env, interp.locals, interp.lx)
		if err != nil {
			interp.err = err
			interp.badToken = badToken
			return
		}
		interp.returnVal = returnVal
		interp.checkReturn = checkReturn
		if interp.checkReturn {
			return
		}
	}
}

func (interp *Interpreter) visitClass(stmt Class) {
	var superclass any
	if stmt.superclass.id > 0 {
		var superclassErr error
		var superclassBadToken Token
		superclass, superclassErr, superclassBadToken = evalExpr(stmt.superclass, interp.env, interp.locals, interp.lx)
		if superclassErr != nil {
			interp.err = superclassErr
			interp.badToken = superclassBadToken
			return
		}
		switch superclass.(type) {
		case LoxClass:
			break
		default:
			err := fmt.Errorf("Superclass must be a class")
			interp.lx.RuntimeError(stmt.superclass.name, err)
			interp.err = err
			interp.badToken = stmt.superclass.name
			return
		}
	}
	super, ok := superclass.(LoxClass)
	interp.env.define(stmt.name.lexeme, nil)
	env := interp.env
	if stmt.superclass.id > 0 {
		env = &Environment{values: make(map[string]any), enclosing: interp.env}
		env.define("super", super)
	}
	methods := make(map[string]LoxFunction)
	for _, method := range stmt.methods {
		function := LoxFunction{declaration: method, env: env, isInitializer: method.name.lexeme == "init"}
		methods[method.name.lexeme] = function
	}
	var klass LoxClass
	if ok {
		klass = LoxClass{name: stmt.name.lexeme, methods: methods, superclass: &super}
	} else {
		klass = LoxClass{name: stmt.name.lexeme, methods: methods, superclass: nil}
	}
	interp.env.assign(stmt.name, klass)
}

func (interp *Interpreter) visitExpression(stmt Expression) {
	_, err, badToken := evalExpr(stmt.expr, interp.env, interp.locals, interp.lx)
	if err != nil {
		interp.err = err
		interp.badToken = badToken
		return
	}
}

func (interp *Interpreter) visitFunction(stmt Function) {
	function := LoxFunction{declaration: stmt, env: interp.env, isInitializer: false}
	interp.env.define(stmt.name.lexeme, function)
}

func (interp *Interpreter) visitIf(stmt If) {
	conditionVal, conditionErr, conditionBadToken := evalExpr(stmt.condition, interp.env, interp.locals, interp.lx)
	if conditionErr != nil {
		interp.err = conditionErr
		interp.badToken = conditionBadToken
		return
	}
	if isTruthy(conditionVal) {
		returnVal, checkReturn, execErr, execBadToken := execStmt(stmt.thenBranch, interp.env, interp.locals, interp.lx)
		if execErr != nil {
			interp.err = execErr
			interp.badToken = execBadToken
			return
		}
		interp.returnVal = returnVal
		interp.checkReturn = checkReturn
	} else if stmt.elseBranch != nil {
		returnVal, checkReturn, execErr, execBadToken := execStmt(stmt.elseBranch, interp.env, interp.locals, interp.lx)
		if execErr != nil {
			interp.err = execErr
			interp.badToken = execBadToken
			return
		}
		interp.returnVal = returnVal
		interp.checkReturn = checkReturn
	}
}

func (interp *Interpreter) visitPrint(stmt Print) {
	val, err, badToken := evalExpr(stmt.expr, interp.env, interp.locals, interp.lx)
	if err != nil {
		interp.err = err
		interp.badToken = badToken
		return
	}
	if val == nil {
		fmt.Print("nil\n")
	} else {
		fmt.Printf("%v\n", val)
	}
}

func (interp *Interpreter) visitReturn(stmt Return) {
	var value any
	if stmt.value != nil {
		var valueErr error
		var valueBadToken Token
		value, valueErr, valueBadToken = evalExpr(stmt.value, interp.env, interp.locals, interp.lx)
		if valueErr != nil {
			interp.err = valueErr
			interp.badToken = valueBadToken
			return
		}
	}
	interp.returnVal = value
	interp.checkReturn = true
}

func (interp *Interpreter) visitVar(stmt Var) {
	var value any
	if stmt.initializer != nil {
		val, err, badToken := evalExpr(stmt.initializer, interp.env, interp.locals, interp.lx)
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
		conditionVal, conditionErr, conditionBadToken := evalExpr(stmt.condition, interp.env, interp.locals, interp.lx)
		if conditionErr != nil {
			interp.err = conditionErr
			interp.badToken = conditionBadToken
			return
		}
		if !isTruthy(conditionVal) {
			return
		}
		returnVal, checkReturn, bodyErr, bodyBadToken := execStmt(stmt.body, interp.env, interp.locals, interp.lx)
		if bodyErr != nil {
			interp.err = bodyErr
			interp.badToken = bodyBadToken
			return
		}
		if checkReturn {
			interp.returnVal = returnVal
			interp.checkReturn = checkReturn
			break
		}
	}
}

// --------------- EXPRESSIONS ---------------

func (interp *Interpreter) evaluate(expr Expr) {
	expr.accept(interp)
}

func evalExpr(expr Expr, env *Environment, local map[Expr]int, lx *Lox) (any, error, Token) {
	dummyInterp := Interpreter{env: env, locals: local, lx: lx}
	dummyInterp.evaluate(expr)
	return dummyInterp.output, dummyInterp.err, dummyInterp.badToken
}

func (interp *Interpreter) visitAssign(expr Assign) {
	value, err, badToken := evalExpr(expr.value, interp.env, interp.locals, interp.lx)
	if err != nil {
		interp.output = NIL
		interp.err = err
		interp.badToken = badToken
		return
	}
	distance, ok := interp.locals[expr]
	if ok {
		interp.env.assignAt(distance, expr.name, value)
	} else {
		globals := interp.getGlobals()
		assignErr := globals.assign(expr.name, value)
		if assignErr != nil {
			interp.lx.RuntimeError(expr.name, assignErr)
			interp.err = assignErr
			interp.badToken = expr.name
			return
		}
	}
	interp.output = value
}

func (interp *Interpreter) visitBinary(expr Binary) {
	left, err, token := evalExpr(expr.left, interp.env, interp.locals, interp.lx)
	if err != nil {
		interp.output = nil
		interp.err = err
		interp.badToken = token
		return
	}
	right, err, token := evalExpr(expr.right, interp.env, interp.locals, interp.lx)
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

func (interp *Interpreter) visitCall(expr Call) {
	callee, calleeErr, calleeBadToken := evalExpr(expr.callee, interp.env, interp.locals, interp.lx)
	if calleeErr != nil {
		interp.err = calleeErr
		interp.badToken = calleeBadToken
		return
	}
	arguments := make([]any, 0)
	for _, argument := range expr.arguments {
		arg, argErr, argBadToken := evalExpr(argument, interp.env, interp.locals, interp.lx)
		if argErr != nil {
			interp.err = argErr
			interp.badToken = argBadToken
			return
		}
		arguments = append(arguments, arg)
	}
	switch function := callee.(type) {
	case LoxCallable:
		if len(arguments) != function.arity() {
			err := fmt.Errorf("Expected %d arguments but got %d.", function.arity(), len(arguments))
			interp.lx.RuntimeError(expr.paren, err)
			interp.err = err
			interp.badToken = expr.paren
			return
		}
		interp.output = function.call(interp, arguments)
	default:
		err := fmt.Errorf("Can only call functions and classes.")
		interp.lx.RuntimeError(expr.paren, err)
		interp.err = err
		interp.badToken = expr.paren
	}
}

func (interp *Interpreter) visitGet(expr Get) {
	object, objectErr, objectBadToken := evalExpr(expr.object, interp.env, interp.locals, interp.lx)
	if objectErr != nil {
		interp.err = objectErr
		interp.badToken = objectBadToken
		return
	}
	switch li := object.(type) {
	case LoxInstance:
		val, getErr := li.get(expr.name)
		if getErr != nil {
			interp.lx.RuntimeError(expr.name, getErr)
			interp.err = getErr
			interp.badToken = expr.name
			return
		} else {
			interp.output = val
			return
		}
	default:
		err := fmt.Errorf("Only instances have properties.")
		interp.lx.RuntimeError(expr.name, err)
		interp.err = err
		interp.badToken = expr.name
		return
	}
}

func (interp *Interpreter) visitGrouping(expr Grouping) {
	interp.output, interp.err, interp.badToken = evalExpr(expr.expression, interp.env, interp.locals, interp.lx)
}

func (interp *Interpreter) visitLiteral(expr Literal) {
	interp.output = expr.value
}

func (interp *Interpreter) visitLogical(expr Logical) {
	left, leftErr, leftBadToken := evalExpr(expr.left, interp.env, interp.locals, interp.lx)
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
	right, rightErr, rightBadToken := evalExpr(expr.right, interp.env, interp.locals, interp.lx)
	if rightErr != nil {
		interp.err = rightErr
		interp.badToken = rightBadToken
		return
	}
	interp.output = right
}

func (interp *Interpreter) visitSet(expr Set) {
	object, objectErr, objectBadToken := evalExpr(expr.object, interp.env, interp.locals, interp.lx)
	if objectErr != nil {
		interp.err = objectErr
		interp.badToken = objectBadToken
		return
	}
	switch li := object.(type) {
	case LoxInstance:
		value, valueErr, valueBadToken := evalExpr(expr.value, interp.env, interp.locals, interp.lx)
		if valueErr != nil {
			interp.err = valueErr
			interp.badToken = valueBadToken
			return
		}
		li.set(expr.name, value)
		interp.output = value
	default:
		err := fmt.Errorf("Only instances have fields.")
		interp.lx.RuntimeError(expr.name, err)
		interp.err = err
		interp.badToken = expr.name
	}
}

func (interp *Interpreter) visitSuper(expr Super) {
	distance := interp.locals[expr]
	super, _ := interp.env.getAt(distance, "super")
	superclass, _ := super.(LoxClass)
	obj, _ := interp.env.getAt(distance-1, "this")
	object, _ := obj.(LoxInstance)
	method, findMethodErr := superclass.findMethod(expr.method.lexeme)
	if findMethodErr != nil {
		interp.lx.RuntimeError(expr.method, findMethodErr)
		interp.err = findMethodErr
		interp.badToken = expr.method
		return
	}
	interp.output = method.bind(object)
}

func (interp *Interpreter) visitThis(expr This) {
	value, err := interp.lookUpVariable(expr.keyword, expr)
	if err != nil {
		interp.err = err
		interp.badToken = expr.keyword
		return
	}
	interp.output = value
}

func (interp *Interpreter) visitUnary(expr Unary) {
	right, err, token := evalExpr(expr.right, interp.env, interp.locals, interp.lx)
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

func (interp *Interpreter) visitVariable(expr Variable) {
	val, err := interp.lookUpVariable(expr.name, expr)
	if err != nil {
		interp.err = err
		interp.badToken = expr.name
		interp.lx.RuntimeError(expr.name, err)
		return
	}
	interp.output = val
}

func (interp *Interpreter) lookUpVariable(name Token, expr Expr) (any, error) {
	distance, ok := interp.locals[expr]
	if ok {
		return interp.env.getAt(distance, name.lexeme)
	} else {
		globals := interp.getGlobals()
		return globals.get(name)
	}
}

// --------------- HELPERS ---------------

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

func (interp *Interpreter) getGlobals() *Environment {
	env := interp.env
	for env.enclosing != nil {
		env = env.enclosing
	}
	return env
}
