package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Lox struct {
	hadError        bool
	hadRuntimeError bool
}

func main() {
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
	if lx.hadRuntimeError {
		os.Exit(70)
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
	parser := Parser{tokens: tokens, current: 0, lx: lx}
	statements, err := parser.Parse()
	if err != nil {
		return
	}
	interpreter := Interpreter{env: &Environment{values: make(map[string]any)}, lx: lx}
	interpreter.Interpret(statements)
}

func (lx *Lox) Error(line int, message string) {
	lx.report(line, "", message)
}

func (lx *Lox) report(line int, where string, message string) {
	lx.hadError = true
	fmt.Printf("[line %d] Error%v: %v\n", line, where, message)
}

func (lx *Lox) ParseError(token Token, message string) {
	if token.tokenType == EOF {
		lx.report(token.line, " at end", message)
	} else {
		lx.report(token.line, fmt.Sprintf(" at '%s'", token.lexeme), message)
	}
}

func (lx *Lox) RuntimeError(token Token, err error) {
	lx.hadRuntimeError = true
	fmt.Printf("%s\n[line %d]\n", err.Error(), token.line)
}
