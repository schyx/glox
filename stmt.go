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

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i If) accept(v StmtVisitor) { v.visitIf(i) }

func (e Expression) accept(v StmtVisitor) { v.visitExpression(e) }

type Print struct {
	expr Expr
}

func (p Print) accept(v StmtVisitor) { v.visitPrint(p) }

type Var struct {
	name        Token
	initializer Expr
}

func (variable Var) accept(v StmtVisitor) { v.visitVar(variable) }

type StmtVisitor interface {
	visitBlock(Block)
	visitExpression(Expression)
	visitIf(If)
	visitPrint(Print)
	visitVar(Var)
}
