package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Lox struct {
	hadError bool
}

func main() {
	expr := Binary{
		left: Unary{
			operator: Token{tokenType: MINUS, lexeme: "-", literal: nil, line: 1},
			right: Literal{value: 123},
		},
		operator: Token{tokenType: STAR, lexeme: "*", literal: nil, line: 1},
		right:    Grouping{expression: Literal{value: 45.67}},
	}
	expr2 := Binary{
		left: Binary{
			left: Literal{value: 1},
			operator: Token{tokenType: PLUS, lexeme: "+", literal: nil, line: 1},
			right: Literal{value: 2},
		},
		operator: Token{tokenType: STAR, lexeme: "*", literal: nil, line: 1},
		right: Binary{
			left: Literal{value: 3},
			operator: Token{tokenType: MINUS, lexeme: "-", literal: nil, line: 1},
			right: Literal{value: 4},
		},
	}
	astp := AstPrinter{}
	fmt.Println(astp.Print(expr))
	rpn := RPN{}
	fmt.Println(rpn.Print(expr2))
	lx := Lox{}
	lx.main()
}

func (lx *Lox) main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		lx.runFile(args[0])
	} else {
		lx.runPrompt()
	}
}

func (lx *Lox) runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	lx.run(string(content))
	if lx.hadError {
		os.Exit(65)
	}
}

func (lx *Lox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		lx.run(line)
		lx.hadError = false
	}
}

func (lx *Lox) run(source string) {
	scanner := Scanner{source: source, tokens: make([]Token, 0), start: 0, current: 0, line: 1, lox: lx}
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Printf("%+v\n", token)
	}
}

func (lx *Lox) Error(line int, message string) {
	lx.report(line, "", message)
}

func (lx *Lox) report(line int, where string, message string) {
	lx.hadError = true
	fmt.Printf("[line %d] Error%v: %v\n", line, where, message)
}
