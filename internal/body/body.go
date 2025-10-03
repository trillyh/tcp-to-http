package body

type Body struct {
	Body string
	ContentLength int
	CurrentCL int
}

func NewBody() *Body {
	return &Body { Body: string(""),
		ContentLength: 0,
		CurrentCL: 0,
	}
}

func (b *Body) SetLength(cl int) {
	b.ContentLength = cl
}

func (b *Body) Parse(data []byte) (int, bool, error) {
	// remaining in body awaiting to be parsed
	// min b/c we want to make sure that we only parse the Body
	// prevent parsing next request ["BODY" + some of REQUEST2]
	remaining := min(b.ContentLength - len(b.Body), len(data))
	b.Body += string(data[:remaining])

	if len(b.Body) == b.ContentLength {
		return remaining, true, nil
	}
	return remaining, false, nil
}
