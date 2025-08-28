package main

import (
	"errors"
	"fmt"
	"slices"
)

// --------------- PARSER ---------------

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


// --------------- STATEMENTS ---------------

func (p *Parser) declaration() (Stmt, error) {
	if p.match([]TokenType{FUN}) {
		return p.function("function")
	}
	if p.match([]TokenType{VAR}) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) function(kind string) (Stmt, error) {
	name, identifierConsumeErr := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if identifierConsumeErr != nil {
		p.lx.ParseError(p.peek(), identifierConsumeErr.Error())
		return nil, identifierConsumeErr
	}
	_, leftParenConsumeErr := p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	if leftParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), leftParenConsumeErr.Error())
		return nil, identifierConsumeErr
	}
	parameters := make([]Token, 0)
	if !p.check(RIGHT_PAREN) {
		for isComma := true; isComma; isComma = p.match([]TokenType{COMMA}) {
			if len(parameters) >= 255 {
				p.lx.ParseError(p.peek(), "Can't have more than 255 parameters.")
			}
			param, paramConsumeErr := p.consume(IDENTIFIER, "Expect parameter name.")
			if paramConsumeErr != nil {
				p.lx.ParseError(p.peek(), paramConsumeErr.Error())
				return nil, paramConsumeErr
			}
			parameters = append(parameters, param)
		}
	}
	_, rightParenConsumeErr := p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if rightParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), rightParenConsumeErr.Error())
		return nil, rightParenConsumeErr
	}
	_, leftBraceConsumeErr := p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	if leftBraceConsumeErr != nil {
		p.lx.ParseError(p.peek(), leftBraceConsumeErr.Error())
		return nil, leftBraceConsumeErr
	}
	body, bodyErr := p.block()
	if bodyErr != nil {
		return nil, bodyErr
	}
	return Function{name: name, params: parameters, body: body}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, identifierConsumeErr := p.consume(IDENTIFIER, "Expect variable name.")
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
	if p.match([]TokenType{FOR}) {
		return p.forStatement()
	}
	if p.match([]TokenType{IF}) {
		return p.ifStatement()
	}
	if p.match([]TokenType{PRINT}) {
		return p.printStatement()
	}
	if p.match([]TokenType{RETURN}) {
		return p.returnStatement()
	}
	if p.match([]TokenType{WHILE}) {
		return p.whileStatement()
	}
	if p.match([]TokenType{LEFT_BRACE}) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return Block{statments: statements}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	_, leftParenConsumeErr := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if leftParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), leftParenConsumeErr.Error())
		return nil, leftParenConsumeErr
	}
	// Handle initializer part of for loop
	var initializer Stmt
	var initializerError error
	if p.match([]TokenType{SEMICOLON}) {
		initializer = nil
	} else if (p.match([]TokenType{VAR})) {
		initializer, initializerError = p.varDeclaration()
		if initializerError != nil {
			return nil, initializerError
		}
	} else {
		initializer, initializerError = p.expressionStatement()
		if initializerError != nil {
			return nil, initializerError
		}
	}
	// Handle condition check part of for loop
	var condition Expr
	var conditionError error
	if !p.check(SEMICOLON) {
		condition, conditionError = p.expression()
		if conditionError != nil {
			return nil, conditionError
		}
	}
	_, secondSemicolonConsumeErr := p.consume(SEMICOLON, "Expect ';' after a loop condition.")
	if secondSemicolonConsumeErr != nil {
		p.lx.ParseError(p.peek(), secondSemicolonConsumeErr.Error())
		return nil, secondSemicolonConsumeErr
	}
	// Handle increment
	var increment Expr
	var incrementErr error
	if !p.check(RIGHT_PAREN) {
		increment, incrementErr = p.expression()
		if incrementErr != nil {
			return nil, incrementErr
		}
	}
	_, rightParenConsumeErr := p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if rightParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), rightParenConsumeErr.Error())
		return nil, rightParenConsumeErr
	}
	// Handle body of for loop
	body, bodyErr := p.statement()
	if bodyErr != nil {
		return nil, bodyErr
	}
	// Desugar
	if increment != nil {
		body = Block{[]Stmt{body, Expression{expr: increment}}}
	}
	if condition == nil {
		condition = Literal{value: true}
	}
	body = While{condition: condition, body: body}
	if initializer != nil {
		body = Block{statments: []Stmt{initializer, body}}
	}
	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, leftParenConsumeErr := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if leftParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), leftParenConsumeErr.Error())
		return nil, leftParenConsumeErr
	}
	condition, conditionErr := p.expression()
	if conditionErr != nil {
		return nil, conditionErr
	}
	_, rightParenConsumeErr := p.consume(RIGHT_PAREN, "Expect '(' after 'if'.")
	if rightParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), rightParenConsumeErr.Error())
		return nil, rightParenConsumeErr
	}
	thenBranch, thenError := p.statement()
	if thenError != nil {
		return nil, thenError
	}
	var elseBranch Stmt
	if p.match([]TokenType{ELSE}) {
		var elseErr error
		elseBranch, elseErr = p.statement()
		if elseErr != nil {
			return nil, elseErr
		}
	}
	return If{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}, nil
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

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr
	var valueErr error
	if !p.check(SEMICOLON) {
		value, valueErr = p.expression()
		if valueErr != nil {
			return nil, valueErr
		}
	}
	_, semicolonConsumeErr := p.consume(SEMICOLON, "Expect ';' after return value.")
	if semicolonConsumeErr != nil {
		p.lx.ParseError(keyword, semicolonConsumeErr.Error())
		return nil, semicolonConsumeErr
	}
	return Return{keyword: keyword, value: value}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	_, leftParenConsumeErr := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if leftParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), leftParenConsumeErr.Error())
		return nil, leftParenConsumeErr
	}
	condition, conditionErr := p.expression()
	if conditionErr != nil {
		return nil, conditionErr
	}
	_, rightParenConsumeErr := p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if rightParenConsumeErr != nil {
		p.lx.ParseError(p.peek(), rightParenConsumeErr.Error())
		return nil, rightParenConsumeErr
	}
	body, stmtErr := p.statement()
	if stmtErr != nil {
		return nil, stmtErr
	}
	return While{condition: condition, body: body}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	_, consumeErr := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if consumeErr != nil {
		p.lx.ParseError(p.peek(), consumeErr.Error())
		return nil, errors.New(consumeErr.Error())
	}
	return statements, nil
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

// --------------- EXPRESSIONS ---------------

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match([]TokenType{EQUAL}) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		switch t := expr.(type) {
		case Variable:
			name := t.name
			return Assign{name: name, value: value}, nil
		default:
			p.lx.ParseError(equals, "Invalid assignment target.")
		}
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match([]TokenType{OR}) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = Logical{left: expr, operator: operator, right: right}
	}
	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match([]TokenType{AND}) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = Logical{left: expr, operator: operator, right: right}
	}
	return expr, nil
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
	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match([]TokenType{LEFT_PAREN}) {
			var finishCallErr error
			expr, finishCallErr = p.finishCall(expr)
			if finishCallErr != nil {
				return nil, finishCallErr
			}
		} else {
			break;
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)
	if !p.check(RIGHT_PAREN) {
		for next := true; next; next = p.match([]TokenType{COMMA}) {
			if len(arguments) >= 255 {
				p.lx.ParseError(p.peek(), "Can't have more than 255 arguments.")
			}
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
		}
	}
	paren, consumeErr := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if consumeErr != nil {
		p.lx.ParseError(paren, consumeErr.Error())
		return nil, consumeErr
	}
	return Call{callee: callee, paren: paren, arguments: arguments}, nil
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

// --------------- HELPERS ---------------

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
