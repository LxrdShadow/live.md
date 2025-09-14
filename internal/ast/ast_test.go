package ast_test

import (
	"testing"

	"github.com/LxrdShadow/live.md/internal/ast"
	"github.com/LxrdShadow/live.md/internal/lexer"
	"github.com/LxrdShadow/live.md/internal/parser"
)

func TestNormalization(t *testing.T) {
	input := " ** Hello world"
	l := lexer.New(input)
	p := parser.New(l.Lex())
	root := p.Parse()
	root.Normalize()
	paragraph := root.Children[0]

	if len(paragraph.Children) != 1 {
		t.Fatalf("len(paragraph.Children) is not %d. got=%d", 1, len(paragraph.Children))
	}

	if paragraph.Children[0].Type != ast.TEXT {
		t.Fatalf("paragraph.Children[0].Type is not %s. got=%s", ast.TEXT, paragraph.Children[0].Type)
	}
}
