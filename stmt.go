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
	visitPrint(Print)
	visitVar(Var)
}
