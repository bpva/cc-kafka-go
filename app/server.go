package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

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
	defer conn.Close()

	fmt.Println("Connection accepted")

	for {
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	br, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}

	if err == io.EOF {
		fmt.Println("Connection closed")
		return
	}

	req := requestFromBytes(buf)
	if req.hdr.len != int32(br-4) {
		fmt.Printf("header length %d does not match buffer length %d\n", req.hdr.len, br)
		errResp := makeErrorResponce(req.hdr.correlationId, CORRUPT_MESSAGE)
		conn.Write(errResp.bytes())
		return
	}
	resp := makeResponse(req)

	fmt.Printf("Responding with: %v\n", resp)
	_, err = conn.Write(resp.bytes())
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}
