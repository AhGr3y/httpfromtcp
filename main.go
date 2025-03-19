package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const inputFilePath = "messages.txt"
const port = ":42069"

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to setup listener: %s\n", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("error accepting incoming connection: %s\n", err)
		}
		connSrc := conn.RemoteAddr().String()
		fmt.Printf("Incoming connection from %s has been accepted!\n", connSrc)

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		conn.Close()
		fmt.Printf("Connection from %s has been closed.\n", connSrc)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		defer f.Close()
		currentLine := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error reading from %s: %s\n", inputFilePath, err)
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLine += parts[i]
				lines <- currentLine
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return lines
}
