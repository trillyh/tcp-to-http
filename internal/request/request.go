package request

import (
	"bytes"
	"errors"
	"fmt"
	"https/internal/body"
	"https/internal/headers"
	"io"
	"strings"
	"unicode"
	"strconv"
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
		parsingBody 	parserState = "parsingBody"
    done          parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine // Ex: GET /coffee HTTP/1.1
	Headers *headers.Headers
	Body *body.Body
	state parserState
}

func NewRequest() *Request {
	return &Request {
		state: initialized,
		Body: body.NewBody(),
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
 
func (r *Request) parseSingle(data []byte) (int, error) {
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
			return 0, err
		}
		if isHeaderDone {
				fmt.Println("header done")
				r.state = parsingBody
				return n, nil
		}
		return n, nil
	case parsingBody:
		cl := r.Headers.Get("content-length") 
		length := 0
		if cl == "" { // nothing in body to parse
			r.state = done
			return 0, nil
		} 
		length, err := strconv.Atoi(cl)
		if err != nil {
			return 0, fmt.Errorf("erorr when trying to convert contentlength to int")
		}
		r.Body.ContentLength = length
		if length == 0 {
			r.state = done
			return 0, nil
		}
		// remaining in body awaiting to be parsed
		remaining := min(length - len(r.Body.Body), len(data))
		r.Body.Body += string(data[:remaining])
		//n, isBodyDone, err := r.Body.Parse(data)
		//if err != nil {
			//return n, err
		//}	
		if len(r.Body.Body) == length {
			r.state = done
			return remaining, nil
		}
		return remaining, nil
	
	case done:
		return 0, fmt.Errorf("error trying to read in done state")
	}
	return 0, fmt.Errorf("unknown state")
}

func (r *Request) parse(data []byte) (int, error) {
	parsedN := 0
	for r.state != done {
		n, err := r.parseSingle(data[parsedN:])
		if err != nil {
			return 0, err
		}
		parsedN += n
		if n == 0 {
			break
		}
	}
	return parsedN, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := NewRequest()

	// buf could overrun when header + body > 1k 
	buf := make([]byte, 1024)
	bufLen := 0 // valid bytes currently in the buffer

	r.state = parsingRl
	for r.state != done {
		// Read and append to buf at right side of buf[bufLen] 
		// buf[bufLen:] is the remaining free space of the buffer
		n, err := reader.Read(buf[bufLen:]) 
		if n > 0 {
			bufLen += n 
			// pass the bufLen of valid bytes in buf to parse.
			// parseN is the number of bytes the request parsed (read)
			parsedN, parseErr := r.parse(buf[:bufLen])
			if parseErr != nil {
				return nil, parseErr
			}
			copy(buf, buf[parsedN:bufLen])
			bufLen -= parsedN 	
			// data from bufLen up to parsedN is already consumed, so no longer need 0<----bufLen and update bufLen
			// buf[parsedN:bufLen] are still unparsed leftovers
			// shift left (use big brain)
		}
		if err == nil {
			continue
		}

		switch {
		case errors.Is(err, io.EOF):
				// drain buffer
				if (bufLen > 0 && r.state != done) {
					parseN, pErr := drainAndParse(r, buf[:bufLen])
					if pErr != nil {
						return nil, err
					}
					copy(buf, buf[parseN:bufLen])
					bufLen -= parseN
					// after this r.state should be done, else err
				}

				if r.state != done {
					return nil, fmt.Errorf("incomplete request, in state: %s, read n bytes on EOF: %d", r.state, n)
				}	
				
		default:
			return nil, err		
			}
		}	
	return r, nil
}

func drainAndParse(r *Request, data []byte) (int, error) {
	parseN, pErr := r.parse(data)
	if pErr != nil {
		return 0, pErr
	}

	if r.Body.CurrentCL != r.Body.ContentLength {
		return 0, fmt.Errorf("body's len does not match content length %d != %d",
			r.Body.CurrentCL, r.Body.ContentLength)
	}
	r.state = done
	return parseN, nil
}
