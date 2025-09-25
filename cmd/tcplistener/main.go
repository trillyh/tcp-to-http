package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"https/internal/request"
)





func readFromFiles(path string) (f io.ReadCloser) {
	f, err := os.OpenFile("messages.txt", os.O_RDONLY, 0644)

	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close() 

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("New connection accepted")	
		go func(c net.Conn) {
			r, err := request.RequestFromReader(c)
			if err != nil {
				log.Fatal(err)
				return
			}

			rl := r.RequestLine
			fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", rl.Method, rl.RequestTarget, rl.HttpVersion)
			fmt.Printf("\nHeaders:\n")
			for k, v := range r.Headers.All() {
				fmt.Printf("- %s: %s\n", k, v)
			}
		}(conn)
	}
}