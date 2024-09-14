package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

type responseHeader struct {
	CorrelationId int32
}

type response struct {
	len int32
	hdr responseHeader
}

func (r *response) bytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf, uint32(r.len))
	binary.BigEndian.PutUint32(buf[4:], uint32(r.hdr.CorrelationId))
	return buf
}

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

	lengthBuf := make([]byte, 4)

	_, err = io.ReadFull(conn, lengthBuf)
	if err != nil {
		fmt.Println("Error reading length:", err.Error())
		os.Exit(1)
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	fmt.Println("Length:", length)

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(conn, messageBuf)
	if err != nil {
		fmt.Println("Error reading message:", err.Error())
		os.Exit(1)
	}

	correlationId := binary.BigEndian.Uint32(messageBuf[4:8])

	resp := response{
		len: 0,
		hdr: responseHeader{
			CorrelationId: int32(correlationId),
		},
	}

	_, err = conn.Write(resp.bytes())

	if err != nil {
		fmt.Println("Error writing:", err.Error())
		os.Exit(1)
	}
}
