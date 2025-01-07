package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	address = ":7"
)

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Accepting connection from: %s\n", conn.RemoteAddr())
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Println(err)
		}
		conn.Write(buf)
	}
}
