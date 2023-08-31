package writer

import (
	"strings"

	"github.com/bmaupin/go-epub"
	"github.com/r0oot/tools/internal/literotica/protocol"
)

type EpubWriter struct {
	basicInfo *protocol.BasicInfo
	writer    *epub.Epub
}

func (e *EpubWriter) Create(basicInfo *protocol.BasicInfo) error {
	e.basicInfo = basicInfo
	e.writer = epub.NewEpub(basicInfo.Title)
	e.writer.SetAuthor("Dirty Soul")
	return nil
}

func (e *EpubWriter) AddSection(sectionInfo *protocol.SectionInfo) error {
	content := strings.ReplaceAll(sectionInfo.Content, "\n", "</p><p>")
	content = "<p>" + sectionInfo.Content + "</p>"
	e.writer.AddSection(`<h1>`+sectionInfo.Title+`</h1>`+content, sectionInfo.Title, "", "")
	return nil
}

func (e *EpubWriter) Save() error {
	return e.writer.Write(e.basicInfo.Title + ".epub")
}
