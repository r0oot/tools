package main

import (
	"fmt"

	epub "github.com/bmaupin/go-epub"
)

func main() {
	fmt.Println("vim-go")
	// Create a new EPUB
	e := epub.NewEpub("My title")

	// Set the author
	e.SetAuthor("Hingle McCringleberry")

	// Add a section
	for i:=0; i < 10; i ++ {
		section1Body := fmt.Sprintf(`<h1>Section %d</h1>
		<p>This is a paragraph %d.</p>范德萨范德萨范德萨`, i, i)
		e.AddSection(section1Body, fmt.Sprintf("Section %+v: 中文", i), "", "")
	}

	// Write the EPUB
	err := e.Write("My EPUB.epub")
	if err != nil {
		fmt.Println(err)
	}
}

