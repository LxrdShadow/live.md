package main

import (
	"fmt"
	"os"

	"github.com/LxrdShadow/live.md/internal/lexer"
	"github.com/LxrdShadow/live.md/internal/parser"
	"github.com/LxrdShadow/live.md/internal/render"
)

func main() {
	input := "# Hello **world**\nThis *is* `code` ."
	l := lexer.New(input)
	tokens := l.Lex()

	p := parser.New(tokens)
	root := p.Parse()
	root.Normalize()

	r := render.NewHTMLRenderer()
	html := r.Render(root)

	fmt.Println(html)

	var file *os.File
	filePath := "output.html"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("failed to create %s: %s\n", filePath, err)
			return
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_WRONLY, 0755) // 0755 is the file permission in octal
		if err != nil {
			fmt.Printf("failed to open file: %r\n", err)
			return
		}
	}

	file.Write([]byte(html))
	fmt.Println("\nFile written")
}
