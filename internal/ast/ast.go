package ast

import (
	"fmt"
	"strings"
)

const (
	DOCUMENT  = "DOCUMENT"
	HEADER    = "HEADER"
	PARAGRAPH = "PARAGRAPH"
	BOLD      = "BOLD"
	ITALIC    = "ITALIC"
	TEXT      = "TEXT"
)

type Node struct {
	Type     string
	Value    string
	Children []Node
	Level    int
}

func (n Node) String() string {
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

	return n.Type
}
