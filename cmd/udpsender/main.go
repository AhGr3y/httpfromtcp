package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const serverAddr = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to resolve udp addr: %s\n", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to dial up a connection: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s. Type a message and press Enter to send a message. Type Ctrl+C to exit the program.\n", serverAddr)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from input: %s\n", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to udp stream: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Message sent!")
	}
}
