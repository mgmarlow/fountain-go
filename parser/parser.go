package parser

import (
	"fmt"

	"github.com/mgmarlow/fountain/lexer"
)

type Component interface{}

type Composite struct {
	ElementType string
	Children    []Component
}

func NewComposite(el string) Composite {
	return Composite{
		ElementType: el,
		Children:    []Component{},
	}
}

type Dialogue struct {
	Composite
	Character string
}

func NewDialogue(el, c string) Dialogue {
	return Dialogue{
		Composite: Composite{
			ElementType: el,
			Children:    []Component{},
		},
		Character: c,
	}
}

func (c *Composite) add(cmp Component) {
	c.Children = append(c.Children, cmp)
}

type Leaf struct {
	ElementType string
	Value       string
}

func NewLeaf(el, v string) Leaf {
	return Leaf{
		ElementType: el,
		Value:       v,
	}
}

// TODO: Have better-named element types
func getElementType(t lexer.T) string {
	return fmt.Sprint(t)
}

type Parser struct {
	Root Composite
	l    *lexer.Lexer
}

func NewParser(contents string) *Parser {
	parser := &Parser{
		Root: NewComposite("root"),
		l:    lexer.NewLexer(contents),
	}
	parser.parse()
	return parser
}

type Emitter interface {
	Emit(root Composite) []byte
}

func (p *Parser) Emit(emitter Emitter) []byte {
	return emitter.Emit(p.Root)
}

func (p *Parser) parse() {
	for p.l.Token != lexer.TEndOfFile {
		var node Component

		switch p.l.Token {
		case lexer.TDialogue:
			node = p.assembleDialogue()
		default:
			// Unimplemented, should default to text
			node = NewLeaf(getElementType(p.l.Token), p.l.Value)
		}

		p.Root.add(node)
		p.l.Next()
	}
}

func (p *Parser) assembleDialogue() Dialogue {
	node := NewDialogue("dialogue", p.l.Value)

	for {
		p.l.Next()

		switch p.l.Token {
		case lexer.TText, lexer.TAsterisk, lexer.TUnderscore, lexer.TParenOpen, lexer.TParenClose:
			node.add(NewLeaf("text", p.l.Value))
		default:
			return node
		}
	}
}
