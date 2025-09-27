package body

import (
	"bytes"
	"fmt"
)
type Body struct {
	Body string
	ContentLength int
	CurrentCL int
}

func NewBody() *Body {
	return &Body {
		Body: string(""),
		ContentLength: 0,
		CurrentCL: 0,
	}
}

var CRLF = []byte("\n")
var ErrExtractContentLengthFailed = fmt.Errorf("failed to extract content length")
// Parse will be called multiple time, each time should consume at least one line.
// Return the bytes consumed
// If the parse found contentLength, update and return immedietly
func (b *Body) Parse(data []byte) (int, bool, error) {
	idx := bytes.Index(data, CRLF)
	if idx == -1 { // not enough data
		return 0, false, nil
	}

	if idx == 0 { // last CRFL at the end
		return 0, true, nil
	}

	consumedN := idx + len(CRLF)
	currLine := data[:idx]
	// Todo refractor this
	b.Body += string(currLine) + "\n" // add back the /n

	b.CurrentCL += consumedN

	return consumedN, b.CurrentCL == consumedN, nil
}
