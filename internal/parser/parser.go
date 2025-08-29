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

func (p *Parser) parseBlock() *ast.Node {
	tok := p.current()
	switch tok.Type {
	case token.HEADER:
		level := p.getHeaderLevel(p.pos)
		for range level {
			p.consume()
		}
		return &ast.Node{Type: ast.HEADER, Level: level}
	default:
		p.consume()
		return &ast.Node{Type: ast.PARAGRAPH}
	}
}

func (p *Parser) getHeaderLevel(start int) int {
	pos := 0

	for pos < len(p.tokens[start:]) && p.tokens[pos].Type == token.HEADER {
		pos++
	}

	return pos
}
