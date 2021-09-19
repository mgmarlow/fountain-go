package lexer

import (
	"fmt"
	"strings"
	"unicode"
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
	TTilde
	TSlugline
	TDialogue
	TTransition
	TCenteredText
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
	"tilde",
	"slugline",
	"dialogue",
	"transition",
	"centered_text",
}

func (t T) String() string {
	return printmap[t]
}

type Lexer struct {
	Token   T
	Line    uint
	Value   string
	input   string
	current int
	start   int
	end     int
	w       int
	cp      rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.step()
	l.Next()
	return l
}

func (l *Lexer) Next() {
	for {
		l.start = l.end
		l.Token = 0

		switch l.cp {
		case eof:
			l.Token = TEndOfFile

		// \r\n for Windows
		case '\r':
			l.step()
			l.step()
			l.Line++
			continue

		case '\n':
			l.step()
			l.Line++
			continue

		case '\t':
			l.step()
			continue

		case '_':
			l.step()
			l.Token = TUnderscore
			l.Value = ""

		case '*':
			l.step()

			if l.cp == '/' {
				l.step()
				l.Token = TBoneyardClose
				l.Value = ""
				break
			}

			l.Token = TAsterisk
			l.Value = ""

		case '/':
			l.step()
			if l.cp == '*' {
				l.step()
				l.Token = TBoneyardOpen
				l.Value = ""
			}

		case '[':
			l.step()
			if l.cp == '[' {
				l.step()
				l.Token = TNoteOpen
				l.Value = ""
			}

		case ']':
			l.step()
			if l.cp == ']' {
				l.step()
				l.Token = TNoteClose
				l.Value = ""
			}

		case '=':
			l.step()
			l.Token = TEquals
			l.Value = ""

		case '~':
			l.step()
			l.Token = TTilde
			l.Value = ""

		case '(':
			l.step()
			l.Token = TParenOpen
			l.Value = ""

		case ')':
			l.step()
			l.Token = TParenClose
			l.Value = ""

		default:
			l.step()

		text:
			for {
				switch l.cp {
				case '\n', '\r', eof, '_', '*', '~', '(', ')':
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

			// These matches are line-based. Handling in the lexer for now,
			// but may be better job for the parser.
			if isUpper(contents) {
				if beginsWith(contents, "EXT.") {
					l.Token = TSlugline
					l.Value = contents
					return
				}

				if beginsWith(contents, "INT.") {
					l.Token = TSlugline
					l.Value = contents
					return
				}

				if beginsWith(contents, ".") {
					l.Token = TSlugline
					l.Value = contents
					return
				}

				if beginsWith(contents, ">") && endsWith(contents, "<") {
					l.Token = TCenteredText
					l.Value = contents
					return
				}

				if beginsWith(contents, ">") {
					l.Token = TTransition
					l.Value = contents
					return
				}

				if endsWith(contents, "TO:") {
					l.Token = TTransition
					l.Value = contents
					return
				}

				l.Token = TDialogue
				l.Value = contents
				return
			}

			l.Token = TText
			l.Value = contents
		}

		return
	}
}

func (l *Lexer) step() {
	cp, w := utf8.DecodeRuneInString(l.input[l.current:])

	if w == 0 {
		cp = eof
	}

	l.w = w
	l.cp = cp
	l.end = l.current
	l.current += w
}

func (l *Lexer) peek() rune {
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

func (l *Lexer) raw() string {
	return l.input[l.start:l.end]
}

func (l *Lexer) String() string {
	var val string
	if l.Value != "" {
		val = fmt.Sprintf(" value=\"%s\"", l.Value)
	}

	return fmt.Sprintf("<%v%s>", l.Token, val)
}

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func beginsWith(contents, s string) bool {
	return strings.Index(contents, s) == 0
}

func endsWith(contents, s string) bool {
	return strings.Index(contents, s) == len(contents)-len(s)
}
