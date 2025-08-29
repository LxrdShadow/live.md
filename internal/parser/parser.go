package parser

import (
	"github.com/LxrdShadow/live.md/internal/ast"
	"github.com/LxrdShadow/live.md/internal/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) consume() token.Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) Parse() *ast.Node {
	root := &ast.Node{Type: ast.DOCUMENT}

	for p.current().Type != token.EOF {
		root.Children = append(root.Children, p.parseBlock())
	}

	return root
}
