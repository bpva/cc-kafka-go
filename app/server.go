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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var eof bool
	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()
	for !eof {
		buf := make([]byte, 1024)
		br, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				eof = true
				break
			}
			fmt.Println("Error reading:", err.Error())
			return
		}

		req := requestFromBytes(buf)
		if req.hdr.len != int32(br-4) {
			fmt.Printf("header length %d does not match buffer length %d\n", req.hdr.len, br-4)
			errResp := makeErrorResponce(req.hdr.correlationId, CORRUPT_MESSAGE)
			conn.Write(errResp.bytes())
			return
		}
		resp := makeResponse(req)

		fmt.Printf("Responding with: %v\n", resp.bytes())
		_, err = conn.Write(resp.bytes())
		if err != nil {
			fmt.Println("Error writing:", err.Error())
		}
	}
}
