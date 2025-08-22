package main

type Expr interface {
	accept(Visitor)
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) accept(v Visitor) { v.visitBinary(b) }

type Grouping struct {
	expression Expr
}

func (g Grouping) accept(v Visitor) { v.visitGrouping(g) }

type Literal struct {
	value any
}

func (l Literal) accept(v Visitor) { v.visitLiteral(l) }

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) accept(v Visitor) { v.visitUnary(u) }

type Visitor interface {
	visitBinary(Binary)
	visitGrouping(Grouping)
	visitLiteral(Literal)
	visitUnary(Unary)
}
