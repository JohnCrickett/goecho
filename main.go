// The Echo Protocol is a formally defined in RFC862
// https://datatracker.ietf.org/doc/html/rfc862
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	address = ":7" //This is a Well Known Port assigned for Echo Protocol
	//https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Well-known_ports
	bufferSize = 2048
)

func main() {
	udp := flag.Bool("udp", false, "Use UDP instead of TCP")
	flag.Parse()

	if *udp {
		fmt.Printf("UDP Echo server listening on %s\n", address)
		UdpEchoServer()
	} else {
		fmt.Printf("TCP Echo server listening on %s\n", address)
		TcpEchoServer()
	}
}

func UdpEchoServer() {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, bufferSize)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}
		go udpReply(conn, addr, buf[:n])
	}
}

func udpReply(conn net.PacketConn, addr net.Addr, buf []byte) {
	conn.WriteTo(buf, addr)
}

func TcpEchoServer() {
	// Go provides a high level abstraction instead of the Berkley Sockets API
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
		fmt.Printf("TCP connection from: %s\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, bufferSize)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Println(err)
		}
		_, err = conn.Write(buf)
		if err != nil {
			return
		}
		clear(buf)
	}
}
