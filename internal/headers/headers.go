package headers

import (
	"strings"
	"fmt"
)

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


type Headers map[string]string
func NewHeaders() Headers {
	return Headers{}
}

var ErrFieldNameContainsSpace = fmt.Errorf("field name contains space")
var ErrBadFieldName = fmt.Errorf("bad field-name")
var SEPARATOR = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	in := string(data)
	idx := strings.Index(in, SEPARATOR)
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 0, true, nil
	}

	consumedN := idx + len(SEPARATOR)

	field := in[:idx]
	parts := strings.SplitN(field, ":", 2)
	fieldName := parts[0]
	fieldValue := parts[1]
	
	if strings.HasSuffix(fieldName, " ") || len(fieldName) < 1 {
		return 0, false, ErrFieldNameContainsSpace
	}

	// constraint here
	for _, r := range fieldName {
		if !isTokenChar(r) {
			return 0, false, ErrBadFieldName
		}
	}

	fieldName = strings.ToLower(fieldName)

	fieldName, fieldValue = strings.TrimSpace(fieldName), strings.TrimSpace(fieldValue)
	h[fieldName] = fieldValue

	fmt.Printf("Consumed %d", idx)
	return consumedN, false, nil
}