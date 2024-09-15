package main

import (
	"fmt"
	"net"
	"os"
	"time"
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
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()

	fmt.Println("Connection accepted")

	for {
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	br, _ := conn.Read(buf)

	req := requestFromBytes(buf)
	if req.hdr.len != int32(br-4) {
		fmt.Printf("header length %d does not match buffer length %d\n", req.hdr.len, br-4)
		errResp := makeErrorResponce(req.hdr.correlationId, CORRUPT_MESSAGE)
		conn.Write(errResp.bytes())
		return
	}
	resp := makeResponse(req)

	fmt.Printf("Responding with: %v\n", resp.bytes())
	_, err := conn.Write(resp.bytes())
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}
}
