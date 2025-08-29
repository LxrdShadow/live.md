package ast

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
