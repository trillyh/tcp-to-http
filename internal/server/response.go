package server

import (
	"io"
)
type StatusCode int

const (
	Ok StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine []byte 
	switch statusCode {
	case Ok:
		statusLine = []byte("HTTP/1.1 200 OK")
	case BadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request")
	case InternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error")
	default:
		statusLine = []byte("")
	}
	_, err := w.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers
