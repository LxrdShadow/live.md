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
