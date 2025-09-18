package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	f, err := os.OpenFile("messages.txt", os.O_RDONLY, 0644)

	if err != nil {
		panic(err)
	}

	ch := getLinesChannel(f)
	
	for newLine := range ch {
		fmt.Printf("read: %s\n", newLine)
	}


}

func getLinesChannel(f io.ReadCloser) <-chan string {
	currentLine := ""

	ch := make(chan string)

	go func() {
		for {
			buf := make([]uint8, 8)
			bytes_read, err := f.Read(buf)

			if err == io.EOF || bytes_read == 0 {
				ch <- currentLine
				defer f.Close()
				close(ch)
				break
			}

			parts := strings.Split(string(buf), "\n")

			currentLine += parts[0]
			if len(parts) > 1 {
				ch <- currentLine
				currentLine = parts[1]
			}
		}
	}()

	return ch
}
