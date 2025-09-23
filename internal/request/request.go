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

type parserState string
const (
	initialzed = "init"
	done = "done"
)

type Request struct {
	RequestLine RequestLine // Ex: GET /coffee HTTP/1.1
	state parserState
}

func newRequest() *Request {
	return &Request {
		state: initialzed,
	}
}

var ErrBadRequestLine = fmt.Errorf("bad request-line")
var ErrIncompleteRequestLine = fmt.Errorf("incomplete start line")
var ErrUnsupportedVersion = fmt.Errorf("unsupported HTTP version")
var ErrInvalidMethod = fmt.Errorf("invalid method")
var SEPARATOR = "\r\n"

func (r *Request) parse(data []byte) (int, error) {

}

func parseRequestLine(s string) (*RequestLine, int, error) {
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
	request := newRequest()

	// buf could overrun when header + body > 1k 
	buf := make([]byte, 1024)
	bufLen := 0 // valid bytes currently in the buffer
	for request.state != "done" {
		// Read from bufLen(start from 0) to 1024
		// n is the number of bytes it has read (element in the buf right now)
		n, err := reader.Read(buf[bufLen:]) 
		if err != nil {
			return nil, err
		}

		bufLen += n 
		// pass in the buf from 0->bufLen to the request's parser
		// parseN is the number of bytes the request parsed (read)
		parsedN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		// data from bufLen up to parsedN is already consumed, so no longer need 0<----bufLen and update bufLen
		// buf[parsedN:bufLen] are still unparsed leftovers
		// shift left (use big brain)
		copy(buf, buf[parsedN:bufLen])
		bufLen -= parsedN 	
	}

	return request, nil
}