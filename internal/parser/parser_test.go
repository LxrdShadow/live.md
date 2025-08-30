package parser

import (
	"fmt"
	"testing"

	"github.com/LxrdShadow/live.md/internal/ast"
	"github.com/LxrdShadow/live.md/internal/lexer"
)

func TestParsingHeader(t *testing.T) {
	input := `# **Hello cruel
*World* of ***mine***`
	l := lexer.New(input)
	tokens := l.Lex()
	fmt.Println(tokens)
	p := New(tokens)
	document := p.Parse()

	fmt.Println(document)
	if len(document.Children) != 2 {
		t.Fatalf("len(document.Children) is not %d. got=%d", 2, len(document.Children))
	}

	header := document.Children[0]
	if header.Type != ast.HEADER {
		t.Fatalf("header.Type is not %s. got=%s", ast.HEADER, header.Type)
	}
	if header.Level != 1 {
		t.Fatalf("header.Level is not %d. got=%d", 1, header.Level)
	}

	if document.Children[1].Type != ast.PARAGRAPH {
		t.Fatalf("document.Children[1].Type is not %s. got=%s", ast.PARAGRAPH, document.Children[1].Type)
	}
}
