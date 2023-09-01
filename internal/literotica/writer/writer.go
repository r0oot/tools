package writer

import "github.com/r0oot/tools/internal/literotica/processor"
import "github.com/r0oot/tools/internal/literotica/protocol"

func NewWriter(t protocol.BookType) processor.Writer {
	switch t {
	case protocol.BookTypeTxt:
		return &TxtWriter{}
	case protocol.BookTypeEpub:
		return &EpubWriter{}
	default:
		panic("book type error")
	}
}
