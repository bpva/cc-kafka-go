package main

import (
	"fmt"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Connection accepted")

	buff := make([]byte, 1024)

	_, err = conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte{0, 0, 0, 0, 0, 0, 0, 7})
	if err != nil {
		fmt.Println("Error writing:", err.Error())
		os.Exit(1)
	}

	conn.Close()
}
