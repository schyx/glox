package main

type Expr interface {
	accept(ExprVisitor)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) accept(v ExprVisitor) { v.visitBinary(b) }

type Grouping struct {
	expression Expr
}

func (g Grouping) accept(v ExprVisitor) { v.visitGrouping(g) }

type Literal struct {
	value any
}

func (l Literal) accept(v ExprVisitor) { v.visitLiteral(l) }

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) accept(v ExprVisitor) { v.visitUnary(u) }

type Variable struct {
	name Token
}

func (variable Variable) accept(v ExprVisitor) {v.visitVariable(variable) }

type ExprVisitor interface {
	visitBinary(Binary)
	visitGrouping(Grouping)
	visitLiteral(Literal)
	visitUnary(Unary)
	visitVariable(Variable)
}
