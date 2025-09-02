package lexer

import (
	"bufio"
	"bytes"
	"unicode/utf8"

	"github.com/LxrdShadow/live.md/internal/token"
)

type Lexer struct {
	input string
	pos   int
	buf   []rune
}

func New(input string) *Lexer {
	return &Lexer{input: input, buf: []rune{}}
}

func (l *Lexer) Lex() []token.Token {
	tokens := []token.Token{}

	reader := bytes.NewReader([]byte(l.input))
	scanner := bufio.NewScanner(reader)

	firstLine := true
	for scanner.Scan() {
		if !firstLine {
			// add newline token before the next line
			tokens = append(tokens, token.Token{Type: token.NEWLINE})
		} else {
			firstLine = false
		}

		tokens = append(tokens, l.lexLine(scanner.Text())...)
	}
	tokens = append(tokens, token.Token{Type: token.EOF})

	return tokens
}

func (l *Lexer) remaining() string {
	return l.input[l.pos:]
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
				if i+2 < len(line) && line[i+2] == '*' {
					tok := token.Token{Type: token.BOLDITALIC, Value: "***"}
					tokens = append(tokens, tok)
					i += len("***")
				} else {
					tok := token.Token{Type: token.BOLD, Value: "**"}
					tokens = append(tokens, tok)
					i += len("**")
				}
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
