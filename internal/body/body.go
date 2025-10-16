// Package body provides a small helper for incrementally assembling a fixed-
// length request body from arbitrary byte chunks (e.g., as they arrive from a
// network connection). The caller sets the expected byte length up front and
// then feeds chunks via Parse until the body is complete.
package body

// Body holds state for assembling a fixed-length payload.
//
// Invariants:
//   - ContentLength is the total number of bytes expected.
//   - len(Body) is the number of bytes accumulated so far.
//   - When len(Body) == ContentLength, the body is complete.
//
// This type is NOT safe for concurrent use without external synchronization.
type Body struct {
	Body string
	ContentLength int // Must be >= 0
}

func NewBody() *Body {
	return &Body { Body: string(""),
		ContentLength: 0,
	}
}

func (b *Body) SetLength(cl int) {
	b.ContentLength = cl
}

// Parse appends up to the remaining number of bytes from data into the Body,
// returning (consumed, done, err):
//   - consumed: how many bytes from data were used
//   - done: whether the body is now complete (len(Body) == ContentLength)
//   - err: non-nil only on error (e.g., ContentLength < current size)
//
// Parse never over-reads beyond ContentLength. Extra bytes in data are not
// consumed and should be supplied to the next higher-level parser in request
// min b/c we want to make sure that we only parse the Body
// prevent parsing next request ["BODY" + some of REQUEST2]
func (b *Body) Parse(data []byte) (consumed int, done bool, err error) {	
	remaining := min(b.ContentLength - len(b.Body), len(data))
	b.Body += string(data[:remaining])

	if len(b.Body) == b.ContentLength {
		return remaining, true, nil
	}
	return remaining, false, nil
}
