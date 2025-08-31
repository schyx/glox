package main

type Stmt interface {
	accept(StmtVisitor)
}

type Block struct {
	statments []Stmt
	id        int
}

func (b Block) accept(v StmtVisitor) { v.visitBlock(b) }

type Class struct {
	name       Token
	superclass Variable
	methods    []Function
	id         int
}

func (c Class) accept(v StmtVisitor) { v.visitClass(c) }

type Expression struct {
	expr Expr
	id   int
}

func (e Expression) accept(v StmtVisitor) { v.visitExpression(e) }

type Function struct {
	name   Token
	params []Token
	body   []Stmt
	id     int
}

func (f Function) accept(v StmtVisitor) { v.visitFunction(f) }

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
	id         int
}

func (i If) accept(v StmtVisitor) { v.visitIf(i) }

type Print struct {
	expr Expr
	id   int
}

func (p Print) accept(v StmtVisitor) { v.visitPrint(p) }

type Return struct {
	keyword Token
	value   Expr
	id      int
}

func (r Return) accept(v StmtVisitor) { v.visitReturn(r) }

type Var struct {
	name        Token
	initializer Expr
	id          int
}

func (variable Var) accept(v StmtVisitor) { v.visitVar(variable) }

type While struct {
	condition Expr
	body      Stmt
	id        int
}

func (w While) accept(v StmtVisitor) { v.visitWhile(w) }

type StmtVisitor interface {
	visitBlock(Block)
	visitClass(Class)
	visitExpression(Expression)
	visitFunction(Function)
	visitIf(If)
	visitPrint(Print)
	visitReturn(Return)
	visitVar(Var)
	visitWhile(While)
}
