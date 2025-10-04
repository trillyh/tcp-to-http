package server

import (
	"net"
)

type Server struct {

}

//Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
func Serve(port int) (*Server, error) 


func (s *Server) Close() error 

/* 
Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine. 
I used an atomic.Bool to track whether the server is closed or not so that I can ignore connection errors after the server is closed.
*/
func (s *Server) listen()

/*
Handle single conection then close
*/
func (s *Server) handle(conn net.Conn)
