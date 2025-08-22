package main

import "fmt"

type RPN struct {
	output string
}

func (rpn *RPN) Print(expr Expr) string {
	expr.accept(rpn)
	return rpn.output
}

func (rpn *RPN) visitBinary(expr Binary) {
	left := RPN{}
	expr.left.accept(&left)
	right := RPN{}
	expr.right.accept(&right)
	rpn.output = fmt.Sprintf("%s %s %s", left.output, right.output, expr.operator.lexeme)
}

func (rpn *RPN) visitGrouping(expr Grouping) {
	intermediate := RPN{}
	expr.expression.accept(&intermediate)
	rpn.output = fmt.Sprintf("%s group ", intermediate.output)
}

func (rpn *RPN) visitLiteral(expr Literal) {
	if expr.value == nil {
		rpn.output = "nil"
	} else {
		rpn.output = fmt.Sprintf("%v", expr.value)
	}
}

func (rpn *RPN) visitUnary(expr Unary) {
	intermediate := RPN{}
	expr.right.accept(&intermediate)
	rpn.output = fmt.Sprintf("%s %s", intermediate.output, expr.operator.lexeme)
}
