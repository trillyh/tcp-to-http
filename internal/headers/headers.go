package headers

import (
	"strings"
	"fmt"
	"bytes"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name string, value string) {
	name = strings.ToLower(name)

	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) All() map[string]string {
	if h == nil {
		return nil
	}
	return h.headers
}

func isTokenChar(r rune) bool {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return true
	}

	if r >= '0' && r <= '9' {
		return true
	}

	switch r {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~' :
		return true
	}		
	return false
}

var ErrFieldNameContainsSpace = fmt.Errorf("field name contains space")
var ErrBadFieldName = fmt.Errorf("bad field-name")
var CRLF = []byte("\r\n")

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, CRLF)
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil // consume 2 bytes (CRLF)
	}

	consumedN := idx + len(CRLF)

	field := data[:idx]
	parts := bytes.SplitN(field, []byte(":"), 2)
	fieldName := parts[0]
	fieldValue := parts[1]
	
	if bytes.HasSuffix(fieldName, []byte(" ")) || len(fieldName) < 1 {
		return 0, false, ErrFieldNameContainsSpace
	}

	// constraint here
	for _, r := range fieldName {
		if !isTokenChar(rune(r)) {
			return 0, false, ErrBadFieldName
		}
	}

	fieldName, fieldValue = bytes.TrimSpace(fieldName), bytes.TrimSpace(fieldValue)

	h.Set(string(fieldName), string(fieldValue))

	return consumedN, false, nil
}
