package lexer

import (
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

	buf := []rune{}

	flushBuf := func() {
		if len(buf) > 0 {
			tokens = append(tokens, token.Token{
				Type:  token.TEXT,
				Value: string(buf),
			})
			clear(buf)
		}
	}

	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.remaining())

		switch r {
		case '#':
			isHeader, level := l.treatHeaderToken(tokens)
			if isHeader {
				flushBuf()

				for range level {
					tokens = append(tokens, token.Token{Type: token.HEADER, Value: "#"})
				}

				l.pos += level
			} else {
				buf = append(buf, r)
				l.pos += size
			}
		case '*':
			flushBuf()

		case '`':
			flushBuf()
			tok := token.Token{Type: token.CODESPAN, Value: "`"}
			tokens = append(tokens, tok)
			l.pos += size
		default:
			buf = append(buf, r)
			l.pos += size
		}
	}

	flushBuf()
	return tokens
}

func (l *Lexer) treatHeaderToken(tokens []token.Token) (bool, int) {
	if len(tokens) > 0 && tokens[len(tokens)-1].Type != token.NEWLINE {
		return false, 0
	}

	level := 0
	pos := 0
	remaining := l.remaining()

	for pos < len(remaining) && remaining[pos] == '#' {
		level++
		pos++
	}

	if pos < len(remaining) && remaining[pos] == ' ' {
		return true, level
	}

	return false, 0
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
