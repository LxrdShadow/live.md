package lexer

import (
	"fmt"
	"testing"

	"github.com/LxrdShadow/live.md/internal/token"
)

func TestLexingHeader(t *testing.T) {
	input := "# Hello world"
	l := New(input)
	tokens := l.Lex()

	if len(tokens) != 1 {
		t.Fatalf("len(tokens) is not %d. got=%d", 1, len(tokens))
	}

	if tokens[0].Type != token.TokenHeader {
		t.Fatalf("token is not HEADER. got=%s", tokens[0])
	}

	if tokens[0].Children[0].Type != token.TokenText {
		t.Fatalf("token child is not TEXT. got=%s", tokens[0].Children[0])
	}

	if tokens[0].Children[0].Value != "Hello world" {
		t.Fatalf("value of token child is not \"Hello world\". got=%s", tokens[0].Children[0].Value)
	}
}

func TestLexingBold(t *testing.T) {
	input := "# Hello **world**!"
	fmt.Println(len(input))

	l := New(input)
	tokens := l.Lex()

	if len(tokens) != 1 {
		t.Fatalf("len(tokens) is not %d. got=%d", 1, len(tokens))
	}

	headerTokenChildren := tokens[0].Children
	fmt.Println(headerTokenChildren)
	if len(headerTokenChildren) != 3 {
		t.Fatalf("len(headerTokenChildren) is not %d. got=%d", 3, len(headerTokenChildren))
	}
}
