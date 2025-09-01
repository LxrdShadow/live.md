package lexer

import (
	"testing"
)

func TestLexingHeader(t *testing.T) {
	input := `# Hello *world*
**Good** morning`

	l := New(input)
	tokens := l.Lex()

	if len(tokens) != 12 {
		t.Fatalf("len(tokens) is not %d. got=%d", 12, len(tokens))
	}
}
