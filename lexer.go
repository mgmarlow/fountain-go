package main

import (
	"strings"
	"unicode"
)

type Token struct {
	kind  string
	value string
}

type Lexer struct {
	current int
	col     int
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
		col:     0,
		// Append newline for easier end-of-string handling.
		input: input + "\n",
		char:  []rune(input)[0],
	}
}

func (l *Lexer) BuildTokens() []Token {
	tokens := []Token{}

	// Conditionals using l.col == 0 apply extra semantic analysis for certain
	// newline lexemes. These lexemes have no significance when nested in text,
	// so we only parse them when they begin at the first column.
	for l.current < len([]rune(l.input)) {
		// These characters are only significant if they fall at the beginning of a line.
		if l.col == 0 {
			switch l.char {
			case '>':
				// Skip leading >
				l.next()

				value := l.collect()

				// Might need to rethink this and just use "gt"/"lt" tokens, depending on
				// support for italicized/underlined centered text.
				if strings.HasSuffix(value, "<") {
					tokens = append(tokens, Token{
						kind: "centered_text",
						// Sans ending '<' char
						value: strings.TrimSpace(value[:len(value)-1]),
					})
					continue
				}

				tokens = append(tokens, Token{
					kind:  "transition",
					value: strings.TrimSpace(value),
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

			case '.':
				if l.peek() != '.' {
					value := l.collect()
					tokens = append(tokens, Token{
						kind:  "scene_heading",
						value: value,
					})
					continue
				}

			case 'E', 'I':
				if l.matches("EXT.") || l.matches("INT.") {
					value := l.collect()
					tokens = append(tokens, Token{
						kind:  "scene_heading",
						value: value,
					})
					continue
				}

			}
		}

		switch l.char {
		case eof:
			l.current++
			continue

		case '\n', '\r':
			l.next()
			l.col = 0
			tokens = append(tokens, Token{
				kind:  "newline",
				value: "",
			})
			continue

		case '(':
			l.next()
			tokens = append(tokens, Token{
				kind:  "oparen",
				value: "(",
			})
			continue

		case ')':
			l.next()
			tokens = append(tokens, Token{
				kind:  "cparen",
				value: ")",
			})
			continue

		case '^':
			l.next()
			tokens = append(tokens, Token{
				kind:  "caret",
				value: "^",
			})
			continue

		case '~':
			l.next()
			tokens = append(tokens, Token{
				kind:  "tilde",
				value: "~",
			})
			continue

		case '*':
			l.next()
			tokens = append(tokens, Token{
				kind:  "asterisk",
				value: "*",
			})
			continue

		case '_':
			l.next()
			tokens = append(tokens, Token{
				kind:  "underscore",
				value: "_",
			})
			continue
		}

		// In all other cases, we want the lexer to fall into some text analysis.
		// A default case is avoided, since it's difficult to manage fallthroughts.
		//
		// When collecting, make sure to punt the lexer back into analysis to catch
		// specific lexemes. Otherwise, nested expressions like emphasis will not be
		// parsed appropriately.
		value := l.collect(func(r rune) bool { return !isNestable(r) })
		if isUpper(value) && containsAlphanumeric(value) {
			if strings.HasSuffix(value, "TO:") {
				tokens = append(tokens, Token{
					kind:  "transition",
					value: value,
				})
				continue
			}

			tokens = append(tokens, Token{
				kind: "character",
				// TODO: Should probably always trim values for tokens
				// other than text.
				value: strings.TrimSpace(value),
			})
			continue
		}

		tokens = append(tokens, Token{
			kind:  "text",
			value: value,
		})
		continue
	}

	return tokens
}

// Use variadic args for default predicate.
func (l *Lexer) collect(p ...func(r rune) bool) string {
	predicate := func(_ rune) bool {
		return true
	}

	if len(p) > 0 {
		predicate = p[0]
	}

	value := ""
	for !isEOL(l.char) && predicate(l.char) {
		value += string(l.char)
		l.next()
	}
	return value
}

func (l *Lexer) next() {
	l.current++
	l.col++
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

func isEOL(r rune) bool {
	return r == '\n' || r == eof || r == '\r'
}

// Expressions that can appear in normal text that we want to make sure get tokenized
// into specific lexemes.
//
// Lexemes that check l.col == 0 cannot appear in this list.
func isNestable(r rune) bool {
	return r == '*' || r == '_' || r == '(' || r == ')' || r == '^'
}

// Contains any numbers or letters. Useful for excluding punctuation lexemes.
func containsAlphanumeric(s string) bool {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			count++
		}
	}
	return count > 0
}
