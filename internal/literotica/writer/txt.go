// Package writer 写
package writer

import (
	"bufio"
	"fmt"
	"os"

	"github.com/r0oot/tools/internal/literotica/protocol"
)

type TxtWriter struct {
	b      *protocol.BasicInfo
	writer *bufio.Writer
	file   *os.File
}

func (t *TxtWriter) Create(b *protocol.BasicInfo) error {
	filePath := b.Title + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}
	t.file = file
	t.writer = bufio.NewWriter(file)
	t.b = b
	return nil
}

func (t *TxtWriter) AddSection(s *protocol.SectionInfo) error {
	t.writer.WriteString(fmt.Sprintf("\n=======%v=======\n\n", s.Title))
	t.writer.WriteString(s.Content)
	return nil
}

func (t *TxtWriter) Save() error {
	t.writer.Flush()
	t.file.Close()
	return nil
}
