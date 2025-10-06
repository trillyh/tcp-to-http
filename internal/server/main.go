package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"atomic"
)

type Server struct {
	listener net.Listener
	close atomic.Bool
}

// Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
func Serve(port int) (*Server, error)  {
	portStr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, fmt.Errorf("error when creating listener for port %d", port)
	}
	server := &Server {
		listener: listener,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	return err
}

/* 
Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine. 
I used an atomic.Bool to track whether the server is closed or not so that I can ignore connection errors after the server is closed.
*/
func (s *Server) listen() {
	 for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New connection accepted")
		go s.handle(conn)
	}
}

/*
Handle single conection then close
*/
func (s *Server) handle(conn net.Conn) {
	fmt.Println("Handling the new connection")
	response := "HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/plain\r\n" +
	"Content-Length: 12\r\n" + // "Hello World!" is 12 bytes
	"\r\n" +                   // blank line separates headers from body
	"Hello World!"
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Fatal("Error ")
	}
	conn.Close()
}
