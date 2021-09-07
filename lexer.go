package main

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

var sceneheadings = [10]string{
	"INT ",
	"EXT ",
	"INT.",
	"EXT.",
	"EST ",
	"EST.",
	"INT./EXT.",
	"INT/EXT ",
	"I/E ",
	"I/E.",
}

type lexer struct {
	input   string
	current int
	start   int
	end     int
	w       int
	value   string
	cp      rune
	token   T
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
	l.step()
	r := l.cp
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.current -= l.w
	cp, w := utf8.DecodeRuneInString(l.input[l.current:])
	l.end = l.current - l.w
	l.cp = cp
	l.w = w
}

func (l *lexer) Next() {
	for {
		l.start = l.end
		l.token = 0

		switch l.cp {
		case eof:
			l.token = TEndOfFile

		case '\n', '\t', ' ':
			l.step()
			continue

		// Italics
		case '_':
			l.step()
			l.token = TUnderscore
			l.value = ""

		// Boneyard/underline/bold
		case '*':
			l.step()

			if l.cp == '/' {
				l.step()
				l.token = TBoneyardEnd
				l.value = ""
				break
			}

			l.token = TAsterisk
			l.value = ""

		// Boneyard
		case '/':
			if l.peek() == '*' {
				l.step()
				l.step()
				l.token = TBoneyardOpen
				l.value = ""
			}

		// Note
		case '[':
			if l.peek() == '[' {
				l.step()
				l.step()
				l.token = TNoteOpen
				l.value = ""
			}

		case ']':
			if l.peek() == ']' {
				l.step()
				l.step()
				l.token = TNoteEnd
				l.value = ""
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
				l.token = TPagebreak
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
					l.token = TCenteredAction
					break center_or_transition

				case '\n', eof:
					l.token = TTransition
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
			l.token = TLyric
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
			l.token = TParenthetical
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

				case ']':
					// Note found
					if l.peek() == ']' {
						break text
					}

				case '/':
					// Boneyard found
					if l.peek() == '*' {
						break text
					}
				}

				l.step()
			}

			contents := l.raw()

			// Heading/Transition/Character
			if isUpper(contents) {
				// Ends in TO:
				if x := strings.Index(contents, "TO:"); x == len(contents)-3 {
					l.token = TTransition
					l.value = contents
					break
				}

				if validSceneHeading(contents) {
					l.token = THeading
					l.value = contents
					break
				}

				if isAlphaNumeric(contents) {
					l.token = TCharacter
					l.value = contents
					break
				}
			}

			l.token = TText
			l.value = contents
		}

		return
	}
}

func validSceneHeading(contents string) bool {
	if !isUpper(contents) {
		return false
	}

	for _, heading := range sceneheadings {
		if x := strings.Index(contents, heading); x == 0 {
			return true
		}
	}

	// Length check to avoid "." false positives
	if x := strings.Index(contents, "."); x == 0 && len(contents) > 1 {
		return true
	}

	return false
}

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func (l *lexer) raw() string {
	return l.input[l.start:l.end]
}

func isAlphaNumeric(s string) bool {
	for _, r := range s {
		alphanumeric := unicode.IsLetter(r) || unicode.IsDigit(r)
		if !alphanumeric && r != ' ' {
			return false
		}
	}

	return true
}

func (l *lexer) String() string {
	return fmt.Sprintf("%v: %s\n", l.token, l.value)
}
