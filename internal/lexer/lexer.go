package lexer

import (
	"bufio"
	"bytes"
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

	reader := bytes.NewReader([]byte(l.input))
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		tokens = append(tokens, l.lexLine(scanner.Text())...)
	}
	tokens = append(tokens, token.Token{Type: token.EOF})

	return tokens
}

func (l *Lexer) lexLine(line string) []token.Token {
	tokens := []token.Token{}
	i := 0

	isHeader, level := l.treatHeader(line)
	if isHeader {
		for range level {
			tokens = append(tokens, token.Token{Type: token.HEADER, Value: "#"})
		}

		i += level
	} else {
		tokens = append(tokens, token.Token{Type: token.PARAGRAPH})
	}

	tokens = append(tokens, l.lexInline(line, level)...)

	return tokens
}

func (l *Lexer) treatHeader(line string) (bool, int) {
	level := 0
	pos := 0

	for pos < len(line) && line[pos] == '#' {
		level++
		pos++
	}

	if pos < len(line) && line[pos] == ' ' {
		return true, level
	}

	return false, 0
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

		return token.Token{Type: token.HEADER, Children: l.lexInline(line[pos:]), Level: min(level, 6)}
	} else {
		return token.Token{Type: token.TEXT, Value: line}
	}
}

func (l *Lexer) lexInline(line string, start int) []token.Token {
	tokens := []token.Token{}
	buf := []rune{}

	flushBuf := func() {
		if len(buf) > 0 {
			tokens = append(tokens, token.Token{
				Type:  token.TEXT,
				Value: string(buf),
			})
			buf = []rune{}
		}
	}

	for i := start; i < len(line); {
		r, size := utf8.DecodeRuneInString(line[i:])

		switch r {
		case '*':
			flushBuf()
			if i+1 < len(line) && line[i+1] == '*' {
				tok := token.Token{Type: token.BOLD, Value: "**"}
				tokens = append(tokens, tok)
				i += len("**")
			} else {
				tok := token.Token{Type: token.ITALIC, Value: "*"}
				tokens = append(tokens, tok)
				i += size
			}
		case '`':
			flushBuf()
			tok := token.Token{Type: token.CODESPAN, Value: "*"}
			tokens = append(tokens, tok)
			i += size
		default:
			buf = append(buf, r)
			i += size
		}
	}

	flushBuf()
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
