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
		return &ast.Node{Type: ast.HEADER, Level: level, Children: p.parseInlineUntil(token.NEWLINE)}
	default:
		p.consume()
		return &ast.Node{Type: ast.PARAGRAPH, Children: p.parseInlineUntil(token.NEWLINE)}
	}
}

func (p *Parser) getHeaderLevel(start int) int {
	pos := 0

	for pos < len(p.tokens[start:]) && p.tokens[pos].Type == token.HEADER {
		pos++
	}

	return pos
}

func (p *Parser) findClosing(t token.TokenType) int {
	for pos := p.pos; pos < len(p.tokens); pos++ {
		if p.tokens[pos].Type == token.NEWLINE {
			return -1
		}

		if p.tokens[pos].Type == t {
			return pos
		}
	}

	return -1
}

func (p *Parser) parseInlineUntil(stop token.TokenType) []*ast.Node {
	nodes := []*ast.Node{}

	for p.current().Type != stop && p.current().Type != token.EOF {
		tok := p.consume()
		switch tok.Type {
		case token.BOLD:
			end := p.findClosing(token.BOLD)

			if end != -1 {
				nodes = append(nodes, &ast.Node{
					Type:     ast.BOLD,
					Children: p.parseInlineUntil(token.BOLD),
				})
			} else {
				nodes = append(nodes, &ast.Node{
					Type:  ast.TEXT,
					Value: tok.Value,
				})
			}
		case token.ITALIC:
			end := p.findClosing(token.ITALIC)

			if end != -1 {
				nodes = append(nodes, &ast.Node{
					Type:     ast.ITALIC,
					Children: p.parseInlineUntil(token.ITALIC),
				})
			} else {
				nodes = append(nodes, &ast.Node{
					Type:  ast.TEXT,
					Value: tok.Value,
				})
			}
		case token.BOLDITALIC:
			end := p.findClosing(token.BOLDITALIC)

			if end != -1 {
				nodes = append(nodes, &ast.Node{
					Type:     ast.BOLDITALIC,
					Children: p.parseInlineUntil(token.BOLDITALIC),
				})
			} else {
				nodes = append(nodes, &ast.Node{
					Type:  ast.TEXT,
					Value: tok.Value,
				})
			}
		default:
			nodes = append(nodes, &ast.Node{
				Type:  ast.TEXT,
				Value: tok.Value,
			})
		}
	}

	if p.current().Type == stop {
		p.consume()
	}

	return nodes
}
