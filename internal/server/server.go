package server

import (
	"fmt"
	"https/internal/request"
	"https/internal/response"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

// The handler writes a success response body to w if everything goes well and returns nil
type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	close atomic.Bool
	handler Handler
}

func (s *Server) Close() error {
	err := s.listener.Close()
	s.close.Store(true)
	return err
}

/*
Handle single conection then close
*/
func (s *Server) handleConnection(conn net.Conn, handler Handler) {
	defer conn.Close() // DOC: why we defer instead of putting it in the end
	
	fmt.Println("Handling the new connection")
	
	responseWriter := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(response.GetDefaultHeaders(0))
		return
	}
	s.handler(responseWriter, r)
}

/* 
Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine. 
I used an atomic.Bool to track whether the server is closed or not so that I can ignore connection errors after the server is closed.
*/
func (s *Server) runServer() {
	 for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.close.Load() {
				return
			}
      log.Printf("accept error: %v", err)
			continue
		}
		fmt.Println("New connection accepted")
		go s.handleConnection(conn, s.handler)
	}
}

// Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
func Serve(port int, handler Handler) (*Server, error) {
	portStr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("error when creating listener for port %d", port)
	}
	server := &Server {
		listener: listener,
		handler: handler,
	}
	go server.runServer()
	return server, nil
}
