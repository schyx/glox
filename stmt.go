package main

type Stmt interface {
	accept(StmtVisitor)
}

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
	visitExpression(Expression)
	visitPrint(Print)
	visitVar(Var)
}
