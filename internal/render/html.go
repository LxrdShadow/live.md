package render

import (
	"strings"

	"github.com/LxrdShadow/live.md/internal/ast"
)

type HTMLRenderer struct{}

func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{}
}

func (r *HTMLRenderer) renderNode(sb *strings.Builder, n *ast.Node) {
	switch n.Type {
	case ast.DOCUMENT:
		for _, c := range n.Children {
			r.renderNode(sb, c)
		}
	case ast.PARAGRAPH:
		sb.WriteString("<p>")
		for _, c := range n.Children {
			r.renderNode(sb, c)
		}
		sb.WriteString("</p>\n")
	case ast.HEADER:
		opening, closing := getHeaderTags(n.Level)
		sb.WriteString(opening)
		for _, c := range n.Children {
			r.renderNode(sb, c)
		}
		sb.WriteString(closing)
	case ast.BOLD:
		sb.WriteString("<strong>")
		for _, c := range n.Children {
			r.renderNode(sb, c)
		}
		sb.WriteString("</strong>")
	case ast.ITALIC:
		sb.WriteString("<em>")
		for _, c := range n.Children {
			r.renderNode(sb, c)
		}
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
	closing := "</h" + string(rune('0'+level)) + ">"
	return opening, closing
}
