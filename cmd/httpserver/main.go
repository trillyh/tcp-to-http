package main

import (
	"crypto/sha256"
	"fmt"
	"https/internal/headers"
	"https/internal/request"
	"https/internal/response"
	"https/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)
const port = 42069

func getBodyResponse400() []byte { 
	return []byte(`<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
	</body>
	</html>`)
}

func getBodyResponse500() []byte {
	return []byte(`<html>
	<head>
	<title>500 Internal Server Error</title>
	</head>
	<body>
	<h1>Internal Server Error</h1>
	<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>`)
}

func getBodyResponse200() []byte {
	return []byte(`<html>
	<head>
	<title>200 OK</title>
	</head>
	<body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
	</body>
	</html>`)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	h := response.GetDefaultHeaders(0)

	target := req.RequestLine.RequestTarget 
	url := "https://httpbin.org/" + target[len("/httpbin/"):]
	fmt.Println("Proxying to", url)
	res, err := http.Get(url)
	if err != nil {
		body := getBodyResponse500()
		w.WriteStatusLine(response.StatusInternalServerError)
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	} 

	w.WriteStatusLine(response.StatusOk)
	h.Delete("Content-Length")	
	h.Set("Transfer-Encoding", "chunked")
	h.Replace("Content-Type", "text/plain")
	h.Replace("Trailer", "X-Content-Length, X-Content-Sha256")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	data := make([]byte, maxChunkSize)
	xContentLength := 0
	for {
		n, err := res.Body.Read(data)
		if n > 0 {
			_, err = w.WriteChunkedBody(data[:n])
			if err != nil {
				fmt.Println("Error when writing chunked body")
			}
			xContentLength += n
		}
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error responding to body:", err)
		}
	}

	_, err = w.WriteChunkedBodyDone()
	fmt.Println("Done with chunk")
	if err != nil {
		fmt.Println("Error writing BodyDone", err)
	}
	trailers := headers.NewHeaders()
	checkSum := sha256.Sum256(w.BodyResponse)
	hash := checkSum[:]
	trailers.Set("X-Content-Sha256", fmt.Sprintf("%x", hash))
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", xContentLength))
	fmt.Println("X-Content-Sha256: ", trailers.Get("X-Content-Sha256"))
	w.WriteTrailers(trailers)
}

func handler(w *response.Writer, req *request.Request) {
	h := response.GetDefaultHeaders(0)
	body := getBodyResponse200()
	status := response.StatusOk

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)

	} else if req.RequestLine.RequestTarget == "/yourproblem" {	
		body = getBodyResponse400()
		status = response.StatusBadRequest

	} else if req.RequestLine.RequestTarget == "/myproblem" {
		body = getBodyResponse500()
		status = response.StatusInternalServerError
	} 

	h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	w.WriteStatusLine(status)
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
