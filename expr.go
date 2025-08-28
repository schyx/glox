package main

type Expr interface {
	accept(ExprVisitor)
}

type Assign struct {
	name  Token
	value Expr
}

func (a Assign) accept(v ExprVisitor) { v.visitAssign(a) }

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) accept(v ExprVisitor) { v.visitBinary(b) }

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (c Call) accept(v ExprVisitor) { v.visitCall(c) }

type Grouping struct {
	expression Expr
}

func (g Grouping) accept(v ExprVisitor) { v.visitGrouping(g) }

type Literal struct {
	value any
}

func (l Literal) accept(v ExprVisitor) { v.visitLiteral(l) }

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func (l Logical) accept(v ExprVisitor) { v.visitLogical(l) }

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) accept(v ExprVisitor) { v.visitUnary(u) }

type Variable struct {
	name Token
}

func (variable Variable) accept(v ExprVisitor) { v.visitVariable(variable) }

type ExprVisitor interface {
	visitAssign(Assign)
	visitBinary(Binary)
	visitCall(Call)
	visitGrouping(Grouping)
	visitLiteral(Literal)
	visitLogical(Logical)
	visitUnary(Unary)
	visitVariable(Variable)
}
