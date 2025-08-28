package main

type Stmt interface {
	accept(StmtVisitor)
}

type Block struct {
	statments []Stmt
}

func (b Block) accept(v StmtVisitor) { v.visitBlock(b) }

type Expression struct {
	expr Expr
}

func (e Expression) accept(v StmtVisitor) { v.visitExpression(e) }

type Function struct {
	name   Token
	params []Token
	body   []Stmt
}

func (f Function) accept(v StmtVisitor) { v.visitFunction(f) }

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i If) accept(v StmtVisitor) { v.visitIf(i) }

type Print struct {
	expr Expr
}

func (p Print) accept(v StmtVisitor) { v.visitPrint(p) }

type Return struct {
	keyword Token
	value   Expr
}

func (r Return) accept(v StmtVisitor) { v.visitReturn(r) }

type Var struct {
	name        Token
	initializer Expr
}

func (variable Var) accept(v StmtVisitor) { v.visitVar(variable) }

type While struct {
	condition Expr
	body      Stmt
}

func (w While) accept(v StmtVisitor) { v.visitWhile(w) }

type StmtVisitor interface {
	visitBlock(Block)
	visitExpression(Expression)
	visitFunction(Function)
	visitIf(If)
	visitPrint(Print)
	visitReturn(Return)
	visitVar(Var)
	visitWhile(While)
}
