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
