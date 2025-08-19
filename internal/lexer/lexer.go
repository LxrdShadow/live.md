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
		tokens = append(tokens, l.lexLine(scanner.Text()))
	}

	return tokens
}

func (l *Lexer) lexLine(line string) token.Token {
	var tok token.Token

	if len(line) > 0 && line[0] == '#' {
		tok = l.lexHeader(line)
	} else {
		tok.Type = token.TokenParagraph
	}

	return tok
}

func (l *Lexer) lexHeader(line string) token.Token {
	pos := 0
	level := 0

	for pos < len(line) && line[pos] == '#' {
		level++
		pos++
	}

	if pos < len(line) && line[pos] == ' ' {
		pos++ // Skip the space

		return token.Token{Type: token.TokenHeader, Value: line[pos:], Level: min(level, 6)}
	} else {
		return token.Token{Type: token.TokenText, Value: line}
	}
}
