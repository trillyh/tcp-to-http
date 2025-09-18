package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)


func getLinesChannel(f io.ReadCloser) <-chan string {
	currentLine := ""

	out := make(chan string)

	go func() {
		defer f.Close()
		defer close(out)

		for {
			buf := make([]uint8, 8)
			bytes_read, err := f.Read(buf)

			if err == io.EOF || bytes_read == 0 {
				out <- currentLine
				break
			}

			parts := strings.Split(string(buf), "\n")

			currentLine += parts[0]
			if len(parts) > 1 {
				out <- currentLine
				currentLine = parts[1]
			}
		}
	}()

	return out
}

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
			lines := getLinesChannel(conn)
			for newLine := range lines {
				fmt.Printf("read: %s\n", newLine)
			}
		}(conn)
	}
}