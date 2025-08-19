package lexer

import (
	"bufio"
	"bytes"

	"github.com/LxrdShadow/live.md/internal/token"
)

type Lexer struct {
	input string
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) Lex() []token.Token {
	tokens := []token.Token{}

	r := bytes.NewReader([]byte(l.input))
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		tokens = append(tokens, l.LexLine(scanner.Text()))
	}

	return tokens
}

func (l *Lexer) LexLine(line string) token.Token {
	tok := token.Token{}
	pos := 0

	if len(line) > 0 && line[pos] == '#' {
		tok.Type = token.TokenHeader
	} else {
		tok.Type = token.TokenParagraph
	}

	return tok
}
