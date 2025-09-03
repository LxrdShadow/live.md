package lexer

import (
	"fmt"
	"testing"
)

func TestLexingSimpleTokens(t *testing.T) {
	input := `## Hello world
*italic* **bold** ***bolditalic***`

	l := New(input)
	tokens := l.Lex()
	fmt.Println(tokens)

	if len(tokens) != 15 {
		t.Fatalf("len(tokens) is not %d. got=%d", 15, len(tokens))
	}
}

func TestLexingCodeSpan(t *testing.T) {
	input := "`\nThis\nis a\ncodespan\n`"

	l := New(input)
	tokens := l.Lex()
	fmt.Println(tokens)

	if len(tokens) != 1 {
		t.Fatalf("len(tokens) is not %d. got=%d", 1, len(tokens))
	}
}
