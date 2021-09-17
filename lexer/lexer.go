package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type T uint

const eof = -1

const (
	TEndOfFile T = iota
	THeading
	TText
	TCharacter
	TParenthetical
	TLyric
	TTransition
	TCenteredAction
	TPagebreak
	TBoneyardOpen
	TBoneyardEnd
	TNoteOpen
	TNoteEnd
	TUnderscore
	TAsterisk
	TSection1
	TSection2
	TSection3
)

var printmap = [...]string{
	"EOF",
	"Heading",
	"Text",
	"Character",
	"Parenthetical",
	"Lyric",
	"Transition",
	"Centered Action",
	"Pagebreak",
	"Boneyard Open",
	"Boneyard End",
	"Note Open",
	"Note End",
	"Underscore",
	"Asterisk",
	"Section 1",
	"Section 2",
	"Section 3",
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

func (l *lexer) Next() {
	for {
		l.start = l.end
		l.Token = 0

		switch l.cp {
		case eof:
			l.Token = TEndOfFile

		case '\n', '\t', ' ':
			l.step()
			continue

		// Italics
		case '_':
			l.step()
			l.Token = TUnderscore
			l.value = l.raw()

		// Boneyard/underline/bold
		case '*':
			l.step()

			if l.cp == '/' {
				l.step()
				l.Token = TBoneyardEnd
				l.value = l.raw()
				break
			}

			l.Token = TAsterisk
			l.value = "*"

		// Boneyard
		case '/':
			l.step()
			if l.cp == '*' {
				l.step()
				l.Token = TBoneyardOpen
				l.value = l.raw()
			}

		// Note
		case '[':
			l.step()
			if l.cp == '[' {
				l.step()
				l.Token = TNoteOpen
				l.value = l.raw()
			}

		case ']':
			l.step()
			if l.cp == ']' {
				l.step()
				l.Token = TNoteEnd
				l.value = l.raw()
			}

		// Pagebreak
		case '=':
			l.step()
			numEquals := 1

		pagebreak:
			for {
				switch l.cp {
				case '=':
					numEquals++
				default:
					break pagebreak
				}
				l.step()
			}

			if numEquals >= 3 {
				l.Token = TPagebreak
				l.value = l.raw()
			} else {
				panic("unterminated pagebreak")
			}

		// Center/transition
		case '>':
			l.step()

		center_or_transition:
			for {
				switch l.cp {
				case '<':
					l.step()
					l.Token = TCenteredAction
					break center_or_transition

				case '\n', eof:
					l.Token = TTransition
					break center_or_transition
				}

				l.step()
			}

			text := l.input[l.start+1 : l.end-1]
			// Extra spaces are normally preserved, but not for centered elements
			l.value = strings.TrimSpace(text)

		// Lyric
		case '~':
			l.step()

		lyric:
			for {
				switch l.cp {
				case '~':
					l.step()
					break lyric

				case '\n', eof:
					panic("unterminated lyric")
				}

				l.step()
			}

			text := l.input[l.start+1 : l.end-1]
			l.Token = TLyric
			l.value = text

		// Parenthetical
		case '(':
			l.step()

		parenthetical:
			for {
				switch l.cp {
				case ')':
					l.step()
					break parenthetical

				case '\n', eof:
					panic("unterminated parenthetical")
				}

				l.step()
			}

			text := l.raw()
			l.Token = TParenthetical
			l.value = text

		default:
			l.step()

		text:
			for {
				switch l.cp {
				case '\n', eof:
					break text

				case '_', '*', '~':
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

func (l *lexer) raw() string {
	return l.input[l.start:l.end]
}

// TODO: Format as such: <TToken value="...">
func (l *lexer) String() string {
	return fmt.Sprintf("%v: %s", l.Token, l.value)
}
