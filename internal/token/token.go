package token

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	EOF       = "TokenEOF"
	HEADER    = "TokenHEADER"
	PARAGRAPH = "TokenPARAGRAPH"
	TEXT      = "TokenTEXT"
	BOLD      = "TokenBOLD"
	ITALIC    = "TokenITALIC"
	CODESPAN  = "TokenCODESPAN"
)

type Token struct {
	Type     TokenType
	Value    string // The text content
	Children []Token
	Level    int // Used only for headers (#=1, ##=2, ...)
}

func (t Token) String() string {
	if len(t.Children) > 0 {
		childrenStr := []string{}

		for _, child := range t.Children {
			childrenStr = append(childrenStr, child.String())
		}
		return fmt.Sprintf("%s[%s]", t.Type, strings.Join(childrenStr, ", "))
	}

	if t.Value != "" {
		return fmt.Sprintf("%s(\"%s\")", t.Type, t.Value)
	}

	return string(t.Type) // Fallback
}
