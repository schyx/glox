package main

type Expr interface {
	accept(ExprVisitor)
}

type Assign struct {
	name  Token
	value Expr
	id    int
}

func (a Assign) accept(v ExprVisitor) { v.visitAssign(a) }

type Binary struct {
	left     Expr
	operator Token
	right    Expr
	id       int
}

func (b Binary) accept(v ExprVisitor) { v.visitBinary(b) }

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
	id        int
}

func (c Call) accept(v ExprVisitor) { v.visitCall(c) }

type Get struct {
	object Expr
	name   Token
	id     int
}

func (g Get) accept(v ExprVisitor) { v.visitGet(g) }

type Grouping struct {
	expression Expr
	id         int
}

func (g Grouping) accept(v ExprVisitor) { v.visitGrouping(g) }

type Literal struct {
	value any
	id    int
}

func (l Literal) accept(v ExprVisitor) { v.visitLiteral(l) }

type Logical struct {
	left     Expr
	operator Token
	right    Expr
	id       int
}

func (l Logical) accept(v ExprVisitor) { v.visitLogical(l) }

type Set struct {
	object Expr
	name   Token
	value  Expr
	id     int
}

func (s Set) accept(v ExprVisitor) { v.visitSet(s) }

type This struct {
	keyword Token
	id      int
}

func (t This) accept(v ExprVisitor) { v.visitThis(t) }

type Unary struct {
	operator Token
	right    Expr
	id       int
}

func (u Unary) accept(v ExprVisitor) { v.visitUnary(u) }

type Variable struct {
	name Token
	id   int
}

func (variable Variable) accept(v ExprVisitor) { v.visitVariable(variable) }

type ExprVisitor interface {
	visitAssign(Assign)
	visitBinary(Binary)
	visitCall(Call)
	visitGet(Get)
	visitGrouping(Grouping)
	visitLiteral(Literal)
	visitLogical(Logical)
	visitSet(Set)
	visitThis(This)
	visitUnary(Unary)
	visitVariable(Variable)
}
