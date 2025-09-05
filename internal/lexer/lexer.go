package lexer

import (
	"strings"
	"unicode/utf8"

	"github.com/LxrdShadow/live.md/internal/token"
)

type Lexer struct {
	input string
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
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
			buf = []rune{}
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

			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '*' {
				if l.pos+2 < len(l.input) && l.input[l.pos+2] == '*' {
					tok := token.Token{Type: token.BOLDITALIC, Value: "***"}
					tokens = append(tokens, tok)
					l.pos += len("***")
				} else {
					tok := token.Token{Type: token.BOLD, Value: "**"}
					tokens = append(tokens, tok)
					l.pos += len("**")
				}
			} else {
				tok := token.Token{Type: token.ITALIC, Value: "*"}
				tokens = append(tokens, tok)
				l.pos += size
			}
		case '`':
			flushBuf()
			count := 1
			for l.pos+count < len(l.input) && l.input[l.pos+count] == '`' {
				count++
			}

			delim := strings.Repeat("`", count)
			l.pos += count

			start := l.pos
			for {
				if l.pos >= len(l.input) {
					tokens = append(tokens, token.Token{Type: token.TEXT, Value: delim + l.input[start:]})
					break
				}

				if strings.HasPrefix(l.input[l.pos:], delim) {
					tokens = append(tokens, token.Token{Type: token.CODESPAN, Value: l.input[start:l.pos]})
					l.pos += count
					break
				}
				l.pos++
			}
		case '\n':
			flushBuf()
			tokens = append(tokens, token.Token{Type: token.NEWLINE})
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
