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
	defer f.Close()


	current_line := ""

	for {
		buf := make([]uint8, 8)
		bytes_read, err := f.Read(buf)

		if err == io.EOF || bytes_read == 0 {
			fmt.Printf("read: %s\n", current_line)
			break
		}

		parts := strings.Split(string(buf), "\n")

		current_line += parts[0]
		if len(parts) > 1 {
			fmt.Printf("read: %s\n", current_line)
			current_line = parts[1]
		}
	}
}
