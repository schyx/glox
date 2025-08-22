package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
	output string
}

func (astp *AstPrinter) Print(expr Expr) string {
	expr.accept(astp)
	return astp.output
}

func (astp *AstPrinter) visitBinary(expr Binary) {
	astp.parenthesize(expr.operator.lexeme, []Expr{expr.left, expr.right})
}

func (astp *AstPrinter) visitLiteral(expr Literal) {
	if expr.value == nil {
		astp.output = "nil"
	} else {
		astp.output = fmt.Sprintf("%v", expr.value)
	}
}

func (astp *AstPrinter) visitGrouping(expr Grouping) {
	astp.parenthesize("group", []Expr{expr.expression})
}

func (astp *AstPrinter) visitUnary(expr Unary) {
	astp.parenthesize(expr.operator.lexeme, []Expr{expr.right})
}

func (astp *AstPrinter) parenthesize(name string, exprs []Expr) {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		innerAstp := AstPrinter{}
		builder.WriteString(" ")
		expr.accept(&innerAstp)
		builder.WriteString(innerAstp.output)
	}
	builder.WriteString(")")
	astp.output = builder.String()
}
