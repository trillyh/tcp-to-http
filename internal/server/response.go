package server

/* 
TODO: Consider these as well
    Content-Encoding: Is the response content encoded/compressed? If so, then this should be included to tell the client how to decode it. (Remember, encoded != encrypted)
    Date: The date and time that the message was sent. This is useful for caching and other things.
    Cache-Control: Directives for caching mechanisms in both requests and responses. This is useful for telling the client or any intermediaries how to cache the response.
*/
import (
	"fmt"
	"https/internal/headers"
	"io"
	"strconv"
)


type StatusCode int
const (
	StatusOk StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	// RFC 9112 status-line = HTTP-version SP status-code SP [ reason-phrase ]
	var statusLine []byte 
	switch statusCode {
	case StatusOk:
		statusLine = []byte("HTTP/1.1 200 OK \r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request \r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error \r\n")
	default:
		statusLine = []byte(" \r\n")
	}
	_, err := w.Write(statusLine)
	if err != nil {
		return fmt.Errorf("error when writing statusCode: %w", err)
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return *h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers.All() {
		hStr := fmt.Sprintf("%s: %s\r\n", k, v)
		fmt.Println(hStr)
		_, err := w.Write([]byte(hStr))
		if err != nil {
			return fmt.Errorf("error when writing header %w", err)
		}
	}

	if _, err := w.Write([]byte("\r\n")); err != nil {
		return fmt.Errorf("error when writing header terminator: %w", err)
	}

	return nil
}

type Writer struct{}
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	return 0, nil
}
