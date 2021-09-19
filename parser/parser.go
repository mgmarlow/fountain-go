package parser

import (
	"fmt"

	"github.com/mgmarlow/fountain/lexer"
)

type Component interface{}

type Composite struct {
	ElementType string
	Children    []Component
	Value       string
}

func NewComposite(el, v string) Composite {
	return Composite{
		ElementType: el,
		Children:    []Component{},
		Value:       v,
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
}

func NewParser(contents string) *Parser {
	return &Parser{
		Root: parse(contents),
	}
}

type Emitter interface {
	Emit(root Composite) []byte
}

func (p *Parser) Emit(emitter Emitter) []byte {
	return emitter.Emit(p.Root)
}

func parse(contents string) Composite {
	root := NewComposite("root", "")

	l := lexer.NewLexer(contents)
	for l.Token != lexer.TEndOfFile {
		var node Component

		switch l.Token {
		case lexer.TDialogue:
			node = assembleDialogue(l)
		default:
			node = NewLeaf(getElementType(l.Token), l.Value)
		}

		root.add(node)
		l.Next()
	}

	return root
}

func assembleDialogue(l *lexer.Lexer) Composite {
	node := NewComposite("dialogue", l.Value)

	for {
		l.Next()

		switch l.Token {
		case lexer.TText, lexer.TAsterisk, lexer.TUnderscore, lexer.TParenOpen, lexer.TParenClose:
			node.add(NewLeaf(getElementType(l.Token), l.Value))
		default:
			return node
		}
	}
}
