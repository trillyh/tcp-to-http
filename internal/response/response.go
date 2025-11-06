package response

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

type writerState int
const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type StatusCode int
const (
	StatusOk StatusCode = 200
	StatusBadRequest StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Writer struct {
	writer io.Writer	
	writerState writerState
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/html")
	return *h
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	// RFC 9112 status-line = HTTP-version SP status-code SP [ reason-phrase ]

	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write statusline in state %d", w.writerState)
	}

	defer func() {w.writerState = writerStateHeaders}()

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
	_, err := w.writer.Write(statusLine)
	if err != nil {
		return fmt.Errorf("error when writing statusCode: %w", err)
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write header in state %d", w.writerState)
	}

	defer func() {w.writerState = writerStateBody}()
	for k, v := range headers.All() {
		hStr := fmt.Sprintf("%s: %s\r\n", k, v)
		fmt.Println(hStr)
		_, err := w.writer.Write([]byte(hStr))
		if err != nil {
			return fmt.Errorf("error when writing header %w", err)
		}
	}

	if _, err := w.writer.Write([]byte("\r\n")); err != nil {
		return fmt.Errorf("error when writing header terminator: %w", err)
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", writerStateBody)
	}
	n, err := w.writer.Write(p)

	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	// Assume body len can be represenet by ui32
	// Get and convert len(p) to hex
	pLen := len(p)
	outLen := fmt.Sprintf("%x\r\n",pLen)
	n, err := w.writer.Write([]byte(outLen)) 
	if err != nil {
		return n, err
	}

	nConsumedFromWritingP, err := w.writer.Write(p[:pLen])
	n += nConsumedFromWritingP
	if err != nil {
		return n, err
	}
	nConsumedFromWritingRN, err := w.writer.Write([]byte("\r\n"))
	n += nConsumedFromWritingRN
	if err != nil {
		return n, err
	}
	return n, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	return n, err
}
