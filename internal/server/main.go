package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
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
	s.close.Store(true)
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
			if s.close.Load() {
				return
			}
			log.Fatal(err)
			continue
		}
		fmt.Println("New connection accepted")
		go s.handle(conn)
	}
}

/*
Handle single conection then close
*/
func (s *Server) handle(conn net.Conn) {
	defer conn.Close() // DOC: why we defer instead of putting it in the end
	fmt.Println("Handling the new connection")
	err := WriteStatusLine(conn, Ok)
	if err != nil {
    log.Printf("write status line: %v", err)
    return
	}
	body := ""
	h := GetDefaultHeaders(len(body)) // 0 if ""
	err = WriteHeaders(conn, h)
	if err != nil {
		log.Fatal("error when writing header")
	}
}
