package main

type Resolver struct {
	interp          *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	lx              *Lox
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

func (r *Resolver) visitBlock(stmt Block) {
	r.beginScope()
	r.resolveStatements(stmt.statments)
	r.endScope()
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.accept(r)
}

func (r *Resolver) visitExpression(stmt Expression) {
	r.resolveExpression(stmt.expr)
}

func (r *Resolver) visitFunction(stmt Function) {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) visitIf(stmt If) {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStatement(stmt.elseBranch)
	}
}

func (r *Resolver) visitPrint(stmt Print) {
	r.resolveExpression(stmt.expr)
}

func (r *Resolver) visitReturn(stmt Return) {
	if r.currentFunction == NONE {
		r.lx.ResolveError(stmt.keyword, "Can't return from top-level code.")
	}
	if stmt.value == nil {
		return
	}
	r.resolveExpression(stmt.value)
}

func (r *Resolver) visitWhile(stmt While) {
	r.resolveExpression(stmt.condition)
	r.resolveStatement(stmt.body)
}

func (r *Resolver) resolveFunction(function Function, functionType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) visitVar(stmt Var) {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpression(stmt.initializer)
	}
	r.define(stmt.name)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.accept(r)
}

func (r *Resolver) visitAssign(expr Assign) {
	r.resolveExpression(expr.value)
	r.resolveLocal(expr, expr.name)
}

func (r *Resolver) visitBinary(expr Binary) {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
}

func (r *Resolver) visitCall(expr Call) {
	r.resolveExpression(expr.callee)
	for _, argument := range expr.arguments {
		r.resolveExpression(argument)
	}
}

func (r *Resolver) visitGrouping(expr Grouping) {
	r.resolveExpression(expr.expression)
}

func (r *Resolver) visitLiteral(expr Literal) {}

func (r *Resolver) visitLogical(expr Logical) {
	r.resolveExpression(expr.left)
	r.resolveExpression(expr.right)
}

func (r *Resolver) visitUnary(expr Unary) {
	r.resolveExpression(expr.right)
}

func (r *Resolver) visitVariable(expr Variable) {
	if (len(r.scopes) == 0) {
		return
	}
	if val, ok := r.scopes[len(r.scopes)-1][expr.name.lexeme]; (ok == true) && (val == false) {
		r.lx.ResolveError(expr.name, "Can't read local variable in its own initializer.")
	}
	r.resolveLocal(expr, expr.name)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.lexeme]; ok {
		r.lx.ResolveError(name, "Already had a variable with this name in this scope.")
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.interp.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}
