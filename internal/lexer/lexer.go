package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

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
	fmt.Println("lex line:", len(line))
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

		return token.Token{Type: token.TokenHeader, Children: l.lexInline(line[pos:]), Level: min(level, 6)}
	} else {
		return token.Token{Type: token.TokenText, Value: line}
	}
}

func (l *Lexer) lexInline(line string) []token.Token {
	tokens := []token.Token{}
	buf := []rune{}

	flushBuf := func(tt token.TokenType) {
		if len(buf) > 0 {
			tokens = append(tokens, token.Token{
				Type:  tt,
				Value: string(buf),
			})
			buf = []rune{}
		}
	}

	for i := 0; i < len(line); {
		r, size := utf8.DecodeRuneInString(line[i:])

		switch r {
		case '*':
			// Case of bold text
			flushBuf(token.TokenText)
			if i+1 < len(line) && line[i+1] == '*' {
				end := l.findClosing(line, i+2, "**")

				if end != -1 {
					tok := token.Token{
						Type:     token.TokenBold,
						Children: l.lexInline(line[i+2 : end]),
					}
					tokens = append(tokens, tok)
					i = end + len("**")
				} else {
					buf = append(buf, '*', '*')
					i += 2
				}
			} else { // Case of italic text
				end := l.findClosing(line, i+2, "*")

				if end != -1 {
					tok := token.Token{
						Type:     token.TokenItalic,
						Children: l.lexInline(line[i+1 : end]),
					}
					tokens = append(tokens, tok)
					i = end + len("*")
				} else {
					buf = append(buf, '*')
					i += 2
				}
			}
		default:
			buf = append(buf, r)
			i += size
		}
	}

	flushBuf(token.TokenText)
	return tokens
}

func (l *Lexer) findClosing(s string, start int, delim string) int {
	if start > len(s) {
		return -1
	}

	idx := strings.Index(s[start:], delim)
	if idx == -1 {
		return idx
	}
	return start + idx
}
