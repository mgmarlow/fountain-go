package lexer

import (
	"fmt"
	"unicode/utf8"
)

type T uint

const eof = -1

const (
	TEndOfFile T = iota
	TText
	TEquals
	TBoneyardOpen
	TBoneyardClose
	TUnderscore
	TAsterisk
	TParenOpen
	TParenClose
	TNoteOpen
	TNoteClose
	TCaretOpen
	TCaretClose
	TTilde
)

var printmap = [...]string{
	"eof",
	"text",
	"equals",
	"boneyard_open",
	"boneyard_end",
	"underscore",
	"asterisk",
	"paren_open",
	"paren_close",
	"note_open",
	"note_close",
	"caret_open",
	"caret_close",
	"tilde",
}

func (t T) String() string {
	return printmap[t]
}

// TODO: Collect newline and column.
type lexer struct {
	Token   T
	input   string
	current int
	start   int
	end     int
	w       int
	value   string
	cp      rune
}

func NewLexer(input string) *lexer {
	l := &lexer{
		input: input,
	}
	l.step()
	l.Next()
	return l
}

func (l *lexer) Next() {
	for {
		l.start = l.end
		l.Token = 0

		switch l.cp {
		case eof:
			l.Token = TEndOfFile

		case '\n', '\t':
			l.step()
			continue

		case '_':
			l.step()
			l.Token = TUnderscore
			l.value = ""

		case '*':
			l.step()

			if l.cp == '/' {
				l.step()
				l.Token = TBoneyardClose
				l.value = ""
				break
			}

			l.Token = TAsterisk
			l.value = ""

		case '/':
			l.step()
			if l.cp == '*' {
				l.step()
				l.Token = TBoneyardOpen
				l.value = ""
			}

		case '[':
			l.step()
			if l.cp == '[' {
				l.step()
				l.Token = TNoteOpen
				l.value = ""
			}

		case ']':
			l.step()
			if l.cp == ']' {
				l.step()
				l.Token = TNoteClose
				l.value = ""
			}

		case '=':
			l.step()
			l.Token = TEquals
			l.value = ""

		case '>':
			l.step()
			l.Token = TCaretOpen
			l.value = ""

		case '<':
			l.step()
			l.Token = TCaretClose
			l.value = ""

		case '~':
			l.step()
			l.Token = TTilde
			l.value = ""

		case '(':
			l.step()
			l.Token = TParenOpen
			l.value = ""

		case ')':
			l.step()
			l.Token = TParenClose
			l.value = ""

		default:
			l.step()

		text:
			for {
				switch l.cp {
				case '\n', eof, '_', '*', '~', '<', '>', '(', ')':
					break text

				case '[':
					// Note found
					if l.peek() == '[' {
						break text
					}
					l.step()

				case ']':
					// Note found
					if l.peek() == ']' {
						break text
					}
					l.step()

				case '/':
					// Boneyard found
					if l.peek() == '*' {
						break text
					}
					l.step()

				default:
					l.step()
				}
			}

			contents := l.raw()
			l.Token = TText
			l.value = contents
		}

		return
	}
}

func (l *lexer) step() {
	cp, w := utf8.DecodeRuneInString(l.input[l.current:])

	if w == 0 {
		cp = eof
	}

	l.w = w
	l.cp = cp
	l.end = l.current
	l.current += w
}

func (l *lexer) peek() rune {
	pw := l.w
	pcp := l.cp
	pend := l.end
	pcurrent := l.current

	l.step()
	r := l.cp

	l.w = pw
	l.cp = pcp
	l.end = pend
	l.current = pcurrent

	return r
}

func (l *lexer) raw() string {
	return l.input[l.start:l.end]
}

func (l *lexer) String() string {
	var val string
	if l.value != "" {
		val = fmt.Sprintf(" value=\"%s\"", l.value)
	}

	return fmt.Sprintf("<%v%s>", l.Token, val)
}
