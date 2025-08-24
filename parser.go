package main

import (
	"errors"
	"slices"
)

type Parser struct {
	tokens  []Token
	current int
	lx      *Lox
}

func (p *Parser) Parse() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return make([]Stmt, 0), err
		}
		statements = append(statements, statement)
	}
	return statements, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match([]TokenType{VAR}) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, identifierConsumeErr := p.consume(IDENTIFIER, "Expect variable name.") // TODO: refactor `consume` to throw error
	if identifierConsumeErr != nil {
		p.lx.ParseError(p.peek(), identifierConsumeErr.Error())
		return nil, errors.New(identifierConsumeErr.Error())
	}
	var initializer Expr = nil
	if p.match([]TokenType{EQUAL}) {
		var exprErr error
		initializer, exprErr = p.expression()
		if exprErr != nil {
			return nil, exprErr
		}
	}

	_, semicolonConsumeErr := p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if semicolonConsumeErr != nil {
		p.lx.ParseError(p.peek(), semicolonConsumeErr.Error())
		return nil, errors.New(semicolonConsumeErr.Error())
	}
	return Var{name: name, initializer: initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match([]TokenType{PRINT}) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, exprErr := p.expression()
	if exprErr != nil {
		return nil, exprErr
	}
	_, consumeErr := p.consume(SEMICOLON, "Expect ';' after vale.")
	if consumeErr != nil {
		p.lx.ParseError(p.peek(), consumeErr.Error())
		return nil, errors.New(consumeErr.Error())
	}
	return Print{value}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	value, exprErr := p.expression()
	if exprErr != nil {
		return nil, exprErr
	}
	_, consumeErr := p.consume(SEMICOLON, "Expect ';' after vale.")
	if consumeErr != nil {
		p.lx.ParseError(p.peek(), consumeErr.Error())
		return nil, errors.New(consumeErr.Error())
	}
	return Expression{value}, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return expr, err
	}

	for p.match([]TokenType{BANG_EQUAL, EQUAL_EQUAL}) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return right, err
		}
		expr = Binary{left: expr, operator: operator, right: right}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return expr, err
	}

	for p.match([]TokenType{GREATER_EQUAL, GREATER, LESS, LESS_EQUAL}) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return right, err
		}
		expr = Binary{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return expr, err
	}

	for p.match([]TokenType{MINUS, PLUS}) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return right, err
		}
		expr = Binary{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return expr, err
	}

	for p.match([]TokenType{SLASH, STAR}) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return right, err
		}
		expr = Binary{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match([]TokenType{BANG, MINUS}) {
		operator := p.previous()
		expr, err := p.unary()
		if err != nil {
			return expr, err
		}
		return Unary{operator: operator, right: expr}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match([]TokenType{FALSE}) {
		return Literal{value: false}, nil
	}
	if p.match([]TokenType{TRUE}) {
		return Literal{value: true}, nil
	}
	if p.match([]TokenType{NIL}) {
		return Literal{value: nil}, nil
	}
	if p.match([]TokenType{NUMBER, STRING}) {
		return Literal{p.previous().literal}, nil
	}
	if p.match([]TokenType{IDENTIFIER}) {
		return Variable{name: p.previous()}, nil
	}

	if p.match([]TokenType{LEFT_PAREN}) {
		expr, errExpression := p.expression()
		if errExpression != nil {
			return expr, errExpression
		}
		_, err := p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			p.lx.ParseError(p.peek(), err.Error())
			return Unary{}, errors.New(err.Error())
		}
		return Grouping{expression: expr}, nil
	}

	p.lx.ParseError(p.peek(), "Expect expression.")
	return Unary{}, errors.New("Expect expression.")
}

func (p *Parser) match(tokenTypes []TokenType) bool {
	if slices.ContainsFunc(tokenTypes, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return Token{}, errors.New(message)
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}
		switch p.peek().tokenType {
		case CLASS:
		case FOR:
		case FUN:
		case IF:
		case PRINT:
		case RETURN:
		case VAR:
		case WHILE:
			return
		}
	}
	p.advance()
}
