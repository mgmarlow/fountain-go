package main

import (
	"unicode"
)

type Token struct {
	kind  string
	value string
}

type Lexer struct {
	current int
	input   string
	char    rune
}

const eof = -1

func Tokenize(input string) []Token {
	l := NewLexer(input)
	return l.BuildTokens()
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		current: 0,
		// Append newline for easier end-of-string handling.
		input: input + "\n",
		char:  []rune(input)[0],
	}
}

func (l *Lexer) BuildTokens() []Token {
	tokens := []Token{}

	for l.current < len([]rune(l.input)) {
		switch l.char {
		case eof:
			l.current++
			continue

		case '\n':
			l.next()
			tokens = append(tokens, Token{
				kind:  "newline",
				value: "",
			})
			continue

		case '@':
			// Skip leading @ symbol
			l.next()
			value := l.collect()
			tokens = append(tokens, Token{
				kind:  "character",
				value: value,
			})
			continue

		case 'E', 'I', '.':
			if l.matches("EXT.") || l.matches("INT.") || l.peek() != '.' {
				value := l.collect()
				tokens = append(tokens, Token{
					kind:  "scene_heading",
					value: value,
				})
				continue
			}
			fallthrough

		default:
			value := l.collect()
			if isUpper(value) {
				tokens = append(tokens, Token{
					kind:  "character",
					value: value,
				})
			} else {
				tokens = append(tokens, Token{
					kind:  "text",
					value: value,
				})
			}
			continue
		}
	}

	return tokens
}

func (l *Lexer) collect() string {
	value := ""
	for l.char != '\n' && l.char != eof {
		value += string(l.char)
		l.next()
	}
	return value
}

func (l *Lexer) next() {
	l.current++
	if l.current < len([]rune(l.input)) {
		l.char = []rune(l.input)[l.current]
	} else {
		l.char = eof
	}
}

func (l *Lexer) peek() rune {
	return []rune(l.input)[l.current+1]
}

func (l *Lexer) matches(expected string) bool {
	for i, c := range expected {
		char := []rune(l.input)[l.current+i]
		if char != c {
			return false
		}
	}

	return true
}

func isUpper(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) && unicode.IsLetter(c) {
			return false
		}
	}

	return true
}
