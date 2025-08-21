package main

import "strconv"

var keywords = map[string]TokenType{
	"and": AND,
	"class": CLASS,
	"else": ELSE,
	"false": FALSE,
	"for": FOR,
	"fun": FUN,
	"if": IF,
	"nil": NIL,
	"or": OR,
	"print": PRINT,
	"return": RETURN,
	"super": SUPER,
	"this": THIS,
	"true": TRUE,
	"var": VAR,
	"while": WHILE,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
	lox     *Lox
}

func (sc *Scanner) ScanTokens() []Token {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.scanToken()
	}
	sc.tokens = append(
		sc.tokens,
		Token{tokenType: EOF, lexeme: "", literal: struct{}{}, line: sc.line},
	)
	return sc.tokens
}

func (sc *Scanner) scanToken() {
	c := sc.advance()
	switch c {
	case '(':
		sc.addShortToken(LEFT_PAREN)
	case ')':
		sc.addShortToken(RIGHT_PAREN)
	case '{':
		sc.addShortToken(LEFT_BRACE)
	case '}':
		sc.addShortToken(RIGHT_BRACE)
	case ',':
		sc.addShortToken(COMMA)
	case '.':
		sc.addShortToken(DOT)
	case '-':
		sc.addShortToken(MINUS)
	case '+':
		sc.addShortToken(PLUS)
	case ';':
		sc.addShortToken(SEMICOLON)
	case '*':
		sc.addShortToken(STAR)
	case '!':
		if sc.match('=') {
			sc.addShortToken(BANG_EQUAL)
		} else {
			sc.addShortToken(BANG)
		}
	case '=':
		if sc.match('=') {
			sc.addShortToken(EQUAL_EQUAL)
		} else {
			sc.addShortToken(EQUAL)
		}
	case '<':
		if sc.match('=') {
			sc.addShortToken(LESS_EQUAL)
		} else {
			sc.addShortToken(LESS)
		}
	case '>':
		if sc.match('=') {
			sc.addShortToken(GREATER_EQUAL)
		} else {
			sc.addShortToken(GREATER)
		}
	case '/':
		if sc.match('/') {
			// A comment goes until the end of the string
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
		} else {
			sc.addShortToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace
		break
	case '\n':
		sc.line += 1
	case '"':
		sc.string()
	default:
		if isDigit(c) {
			sc.number()
		} else if isAlpha(c) {
			sc.identifier()
		} else {
			sc.lox.Error(sc.line, "Unexpected character.")
		}
	}
}

func (sc *Scanner) identifier() {
	for isAlphaNumeric(sc.peek()) {
		sc.advance()
	}
	text := sc.source[sc.start:sc.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	sc.addShortToken(tokenType)
}

func (sc *Scanner) number() {
	for isDigit(sc.peek()) {
		sc.advance()
	}

	// Look for a fractional part
	if sc.peek() == '.' && isDigit(sc.peekNext()) {
		// Consume the "."
		sc.advance()
		for isDigit(sc.peek()) {
			sc.advance()
		}
	}

	num, _ := strconv.ParseFloat(sc.source[sc.start:sc.current], 64)
	sc.addToken(NUMBER, num)
}

func (sc *Scanner) string() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line += 1
		}
		sc.advance()
	}

	if sc.isAtEnd() {
		sc.lox.Error(sc.line, "Unterminated string.")
		return
	}

	sc.advance() // Advance past the closing "

	// Trim surrounding quotes
	value := sc.source[sc.start+1 : sc.current-1]
	sc.addToken(STRING, value)
}

func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != expected {
		return false
	}
	sc.current += 1
	return true
}

func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return byte(0)
	}
	return sc.source[sc.current]
}

func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.source) {
		return byte(0)
	}
	return sc.source[sc.current+1]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) advance() byte {
	defer func() { sc.current += 1 }()
	return sc.source[sc.current]
}

func (sc *Scanner) addShortToken(tokenType TokenType) {
	sc.addToken(tokenType, struct{}{})
}

func (sc *Scanner) addToken(tokenType TokenType, literal any) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, Token{tokenType: tokenType, lexeme: text, literal: literal, line: sc.line})
}
