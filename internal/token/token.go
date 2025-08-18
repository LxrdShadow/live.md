package token

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenHeader
	TokenParagraph
	TokenText
	TokenBold
	TokenItalic
	TokenCodeSpan
)

type Token struct {
	Type  TokenType
	Value string // The text content
	Level int    // Used only for headers (#=1, ##=2, ...)
}

func (t Token) String() string {
	switch t.Type {
	case TokenHeader:
		return "HEADER(" + string(rune('0'+t.Level)) + ", " + t.Value + ")"
	case TokenParagraph:
		return "PARAGRAPH(" + t.Value + ")"
	case TokenText:
		return "TEXT(" + t.Value + ")"
	case TokenBold:
		return "BOLD(" + t.Value + ")"
	case TokenItalic:
		return "ITALIC(" + t.Value + ")"
	case TokenCodeSpan:
		return "CODESPAN(" + t.Value + ")"
	case TokenEOF:
		return "EOF"
	}

	return "UNKNOWN(" + t.Value + ")"
}
