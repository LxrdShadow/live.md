package parser

import "github.com/LxrdShadow/live.md/internal/token"

type Parser struct {
	tokens []token.Token
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}
