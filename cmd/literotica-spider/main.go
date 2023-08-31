package main

import (
	"fmt"

	"github.com/r0oot/tools/internal/literotica/processor"
	"github.com/r0oot/tools/internal/literotica/protocol"
	"github.com/r0oot/tools/internal/literotica/reader"
	"github.com/r0oot/tools/internal/literotica/writer"
)

func main() {
	reqBody := protocol.ParseReqBody()
	p := processor.New(
		&reader.Reader{},
		newWriter(reqBody.BookType))
	if err := p.Do(&reqBody); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("finish")
}

func newWriter(t protocol.BookType) processor.Writer {
	switch t {
	case protocol.BookTypeTxt:
		return &writer.TxtWriter{}
	case protocol.BookTypeEpub:
		return &writer.EpubWriter{}
	default:
		panic("book type error")
	}
}
