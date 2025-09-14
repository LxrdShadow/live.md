package parser

import (
	"github.com/LxrdShadow/live.md/internal/ast"
	"github.com/LxrdShadow/live.md/internal/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

// Initialize new parser
func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) peek() token.Token {
	if p.pos+1 >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos+1]
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
		return &ast.Node{Type: ast.PARAGRAPH, Children: p.parseParagraph()}
	}
}

func (p *Parser) getHeaderLevel(start int) int {
	pos := 0

	// loop while there's still header tokens
	for pos < len(p.tokens[start:]) && p.tokens[pos].Type == token.HEADER {
		pos++
	}

	return pos
}

// search for the token passed as argument and return its index if found, return -1 if not
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

func (p *Parser) parseParagraph() []*ast.Node {
	nodes := []*ast.Node{}

	for {
		if p.current().Type == token.EOF {
			break
		}

		if p.current().Type == token.NEWLINE {
			// if there's 2 newlines, treat it as a break (hard break, paragraph end)
			if p.peek().Type == token.NEWLINE {
				// consume the 2 newlines
				p.consume()
				p.consume()
				break
			}

			// if there's just 1 newline, treat it as a space (soft break)
			nodes = append(nodes, &ast.Node{
				Type:  ast.TEXT,
				Value: " ",
			})
			p.consume()
			continue
		}

		// otherwise parse inline as usual
		tok := p.consume()
		switch tok.Type {
		case token.BOLD, token.ITALIC, token.BOLDITALIC:
			end := p.findClosing(tok.Type)

			if end != -1 {
				nodes = append(nodes, p.parseInline(tok.Type))
			} else {
				nodes = append(nodes, &ast.Node{
					Type:  ast.TEXT,
					Value: tok.Value,
				})
			}
		case token.CODESPAN:
			nodes = append(nodes, &ast.Node{
				Type:  ast.CODESPAN,
				Value: tok.Value,
			})
		default:
			nodes = append(nodes, &ast.Node{
				Type:  ast.TEXT,
				Value: tok.Value,
			})
		}
	}

	return nodes
}

func (p *Parser) parseInline(stop token.TokenType) *ast.Node {
	// consume the opening marker
	node := &ast.Node{Type: tokenToAstType(stop)}

	// collect children until we see the same marker or EOF/NEWLINE
	for p.current().Type != stop &&
		p.current().Type != token.NEWLINE &&
		p.current().Type != token.EOF {

		tok := p.consume()
		switch tok.Type {
		case token.BOLD, token.ITALIC, token.BOLDITALIC:
			end := p.findClosing(tok.Type)

			if end != -1 {
				node.Children = append(node.Children, p.parseInline(tok.Type))
			} else {
				node.Children = append(node.Children, &ast.Node{
					Type:  ast.TEXT,
					Value: tok.Value,
				})
			}
		case token.CODESPAN:
			node.Children = append(node.Children, &ast.Node{
				Type:  ast.CODESPAN,
				Value: tok.Value,
			})
		default:
			node.Children = append(node.Children, &ast.Node{
				Type:  ast.TEXT,
				Value: tok.Value,
			})
		}
	}

	// consume the closing marker if found
	if p.current().Type == stop {
		p.consume()
	}

	return node
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
		case token.CODESPAN:
			nodes = append(nodes, &ast.Node{
				Type:  ast.CODESPAN,
				Value: tok.Value,
			})
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

func tokenToAstType(tt token.TokenType) ast.NodeType {
	switch tt {
	case token.BOLD:
		return ast.BOLD
	case token.ITALIC:
		return ast.ITALIC
	case token.BOLDITALIC:
		return ast.BOLDITALIC
	case token.CODESPAN:
		return ast.CODESPAN
	case token.HEADER:
		return ast.HEADER
	case token.PARAGRAPH:
		return ast.PARAGRAPH
	default:
		return ast.TEXT
	}
}
