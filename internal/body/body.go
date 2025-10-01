package body

import (
//	"bytes"
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
	// ----------------> CHECK IF LEN(DATA) + CURRENTCL == EXPECTED CONTENT LENGTH ELSE RETURN
	fmt.Printf("Body.Parse with %s" , string(data))
	//idx := bytes.Index(data, CRLF)
	idx := len(data)
	fmt.Printf("HERE----> %d", idx)
	if idx == -1 { // not enough data
		return 0, false, nil
	}

	if idx == 0 { // last CRFL at the end
		return 0, true, nil
	}

	//consumedN := idx + len(CRLF)
	consumedN := idx
	currLine := data[:idx]
	// Todo refractor this
	//b.Body += string(currLine) + "\n" // add back the /n
	b.Body += string(currLine) // add back the /n

	b.CurrentCL += consumedN
	fmt.Printf("b.CurrentCL %d contentLength %d", b.CurrentCL, b.ContentLength)
	return consumedN, b.CurrentCL == b.ContentLength, nil
}
