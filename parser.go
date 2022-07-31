package main

import (
	"fmt"
)

// Influenced by mdast:
// https://github.com/remarkjs/remark#syntax-tree

//
// Node Overview
//
// FlowContent: Sections of a screenplay
// Transition | Section | SceneHeading | PageBreak | Lyric | Character | Note | Boneyard
//
// Content: Text that forms paragraphs
// Dialogue | Action | Parenthetical
//
// PhrasingContent: Text in a document & its markup
// Emphasis | Underline | Strong | Text
//

type node struct {
	kind     string
	value    string
	children []node
}

type parser struct {
	current int
	tokens  []Token
}

func Parse(tokens []Token) node {
	p := NewParser(tokens)
	ast := node{
		kind:     "Root",
		children: []node{},
	}

	for p.current < len(p.tokens) {
		ast.children = append(ast.children, p.walk())
	}

	return ast
}

func NewParser(tokens []Token) *parser {
	return &parser{
		current: 0,
		tokens:  tokens,
	}
}

func (p *parser) walk() node {
	token := p.tokens[p.current]

	fmt.Printf("%s: %s\n", token.kind, token.value)

	p.current++

	return node{}
}
