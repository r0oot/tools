package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	epub "github.com/bmaupin/go-epub"
)

type writer interface {
	Create(title string) error
	AddSection(title, content string) error
	Save() error
}

type txtWriter struct {
	title  string
	writer *bufio.Writer
	file   *os.File
}

func (t *txtWriter) Create(title string) error {
	filePath := title + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}
	t.file = file
	t.writer = bufio.NewWriter(file)
	t.title = title
	return nil
}

func (t *txtWriter) AddSection(title, content string) error {
	t.writer.WriteString(fmt.Sprintf("\n=======%v=======\n\n", title))
	t.writer.WriteString(content)
	return nil
}

func (t *txtWriter) Save() error {
	t.writer.Flush()
	t.file.Close()
	return nil
}

type epubWriter struct {
	title  string
	writer *epub.Epub
}

func (e *epubWriter) Create(title string) error {
	e.title = title
	e.writer = epub.NewEpub(title)
	e.writer.SetAuthor("Dirty Soul")
	return nil
}

func (e *epubWriter) AddSection(title, content string) error {
	content = strings.ReplaceAll(content, "\n", "</p><p>")
	content = "<p>" + content + "</p>"
	e.writer.AddSection(`<h1>`+title+`</h1>`+content, title, "", "")
	return nil
}

func (e *epubWriter) Save() error {
	return e.writer.Write(e.title + ".epub")
}
