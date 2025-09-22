package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "1.1"
} 

func (r *RequestLine) ValidMethod() bool {
	isAllUpper := strings.ToUpper(r.Method) == r.Method 
	isAlphaBet := true

	for _, r := range r.Method {
		if !unicode.IsLetter(r) {
			isAlphaBet = false
		}
	}

	return isAllUpper && isAlphaBet
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine // Ex: GET /coffee HTTP/1.1
}

var ErrBadRequestLine = fmt.Errorf("bad request-line")
var ErrIncompleteRequestLine = fmt.Errorf("incomplete start line")
var ErrUnsupportedVersion = fmt.Errorf("unsupported HTTP version")
var ErrInvalidMethod = fmt.Errorf("invalid method")
var SEPARATOR = "\r\n"



func parseRequestLine(s string) (*RequestLine, string, error) {
	idx := strings.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, s, nil
	}

	startLine := s[:idx]
	restOfString := s[idx+len(SEPARATOR):] // START_LINE\r\n ->the rest

	requestLineParts := strings.Split(startLine, " ")	

	fmt.Println(requestLineParts)

	if len(requestLineParts) != 3 {
		return nil, restOfString, ErrIncompleteRequestLine// Empty
	}

	versionParts := strings.Split(requestLineParts[2], "/")

	requestLine := &RequestLine{
		Method: requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion: versionParts[1],
	}

	if !requestLine.ValidHTTP() {
		return nil, restOfString, ErrUnsupportedVersion
	}

	if !requestLine.ValidMethod() {
		return nil, restOfString, ErrInvalidMethod
	}

	return  requestLine, restOfString, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestBytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf("failed to io.ReadAll")
	}

	requestStr := string(requestBytes)

	requestLine, _, err := parseRequestLine(requestStr)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine, // <--------- panic if requestLine is nil 
	}, err
}