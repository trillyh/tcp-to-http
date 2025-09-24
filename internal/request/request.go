package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"https/internal/headers"
	"bytes"
	"errors"
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

type parserState string
const (
    initialized   parserState = "init"
    parsingRl     parserState = "parsingRl"
    parsingHeader parserState = "parsingHeader"
    done          parserState = "done"
)


type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine // Ex: GET /coffee HTTP/1.1
	Headers headers.Headers
	state parserState
}

func NewRequest() *Request {
	return &Request {
		state: initialized,
	}
}

var ErrBadRequestLine = fmt.Errorf("bad request-line")
var ErrIncompleteRequestLine = fmt.Errorf("incomplete start line")
var ErrUnsupportedVersion = fmt.Errorf("unsupported HTTP version")
var ErrInvalidMethod = fmt.Errorf("invalid method")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 { // not enough data
		return nil, 0, nil
	}

	startLine := b[:idx]
	consumedN := idx + len(SEPARATOR) // START_LINE\r\n ->the rest

	requestLineParts := bytes.Split(startLine, []byte(" "))	

	fmt.Println(requestLineParts)

	if len(requestLineParts) != 3 {
		return nil, consumedN, ErrIncompleteRequestLine// Empty
	}

	versionParts := bytes.Split(requestLineParts[2], []byte("/"))

	requestLine := &RequestLine{
		Method: string(requestLineParts[0]),
		RequestTarget: string(requestLineParts[1]),
		HttpVersion: string(versionParts[1]),
	}

	if !requestLine.ValidHTTP() {
		return nil, consumedN, ErrUnsupportedVersion
	}

	if !requestLine.ValidMethod() {
		return nil, consumedN, ErrInvalidMethod
	}

	return  requestLine, consumedN, nil
}

// TODO: use a switch and enum for this
func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case parsingRl:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 { // need more data
			return 0, nil
		}

		r.RequestLine = *rl
		r.state = parsingHeader

		return n, nil

	case parsingHeader:
		if r.Headers == nil {
			r.Headers = headers.NewHeaders()
		}
		n, isHeaderDone, err := r.Headers.Parse(data)
		if err != nil {
			return n, err
		}

		if isHeaderDone {
			r.state = done
			return n, err
		}
	}
	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()

	// buf could overrun when header + body > 1k 
	buf := make([]byte, 1024)
	bufLen := 0 // valid bytes currently in the buffer

	request.state = parsingRl
	for request.state != done {
		// Read from bufLen(start from 0) to 1024
		// n is the number of bytes it has read (element in the buf right now)
		n, err := reader.Read(buf[bufLen:]) 

		if err != nil {
			// Read returns n > 0, it may return err == nil or err == io.EOF (subsequent call after data stop comming in)
			if errors.Is(err, io.EOF) {  // <----------------- last read
				request.state = done
				break
			}
			fmt.Println("error when read from Reader")
			return nil, err
		}

		if n > 0 {
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
	}

	return request, nil
}