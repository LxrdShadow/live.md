package render

import (
	"strings"

	"github.com/LxrdShadow/live.md/internal/ast"
)

type HTMLRenderer struct{}

func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{}
}

func (r *HTMLRenderer) Render(node *ast.Node) string {
	var sb strings.Builder
	r.renderNode(&sb, node)
	return sb.String()
}

func (r *HTMLRenderer) renderNode(sb *strings.Builder, n *ast.Node) {
	switch n.Type {
	case ast.DOCUMENT:
		r.renderChildren(sb, n)
	case ast.PARAGRAPH:
		sb.WriteString("<p>")
		r.renderChildren(sb, n)
		sb.WriteString("</p>\n")
	case ast.HEADER:
		opening, closing := getHeaderTags(n.Level)
		sb.WriteString(opening)
		r.renderChildren(sb, n)
		sb.WriteString(closing)
	case ast.BOLD:
		sb.WriteString("<strong>")
		r.renderChildren(sb, n)
		sb.WriteString("</strong>")
	case ast.ITALIC:
		sb.WriteString("<em>")
		r.renderChildren(sb, n)
		sb.WriteString("</em>")
	case ast.CODESPAN:
		sb.WriteString("<code>")
		sb.WriteString(n.Value)
		sb.WriteString("</code>")
	case ast.TEXT:
		sb.WriteString(n.Value)
	}
}

func getHeaderTags(level int) (string, string) {
	opening := "<h" + string(rune('0'+level)) + ">"
	closing := "</h" + string(rune('0'+level)) + ">\n"
	return opening, closing
}

func (r *HTMLRenderer) renderChildren(sb *strings.Builder, n *ast.Node) {
	for _, c := range n.Children {
		r.renderNode(sb, c)
	}
}
