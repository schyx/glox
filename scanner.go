package main

type Scanner struct {
	source string
}

func (sc Scanner) ScanTokens() []Token {
	return make([]Token, 0)
}
