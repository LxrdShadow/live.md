package ast

import (
	"fmt"
	"strings"
)

type NodeType string

const (
	DOCUMENT   = "AstDOCUMENT"
	HEADER     = "AstHEADER"
	PARAGRAPH  = "AstPARAGRAPH"
	BOLD       = "AstBOLD"
	ITALIC     = "AstITALIC"
	BOLDITALIC = "AstBOLDITALIC"
	TEXT       = "AstTEXT"
	CODESPAN   = "AstCODESPAN"
)

type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
	Level    int
}

func (n *Node) String() string {
	if len(n.Children) > 0 {
		childrenStr := []string{}

		for _, child := range n.Children {
			childrenStr = append(childrenStr, child.String())
		}
		return fmt.Sprintf("%s[%s]", n.Type, strings.Join(childrenStr, ", "))
	}

	if n.Value != "" {
		return fmt.Sprintf("%s(\"%s\")", n.Type, n.Value)
	}

	return string(n.Type)
}
