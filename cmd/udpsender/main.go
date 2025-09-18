package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("W")

	raddr, err := net.ResolveUDPAddr("udp","localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n>")
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
			continue
		}

		if _, err := conn.Write(line); err != nil {
			log.Fatal(err)
			continue
		}
	}
}