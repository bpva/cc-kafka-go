package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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

	var resp response

	req := requestFromBytes(messageBuf)
	err = req.validate()
	if err != nil {
		if errors.Is(err, UnknownVersionErr) {
			respondWithError(
				conn,
				correlationId,
				UNKNOWN_VERSION,
			)
		}
	}

	resp = response{
		hdr: responseHeader{
			CorrelationId: int32(correlationId),
		},
		body: &apiVersionsResponseBody{
			ErrorCode: 0,
			ApiKeys: []apiKey{
				{
					ApiKey:     18,
					MinVersion: 0,
					MaxVersion: 4,
				},
			},
			ThrottleTimeMs: 0,
		},
	}

	resp.setLen()

	fmt.Printf("Response: %+v\n", resp.bytes())

	_, err = conn.Write(resp.bytes())

	if err != nil {
		fmt.Println("Error writing:", err.Error())
		os.Exit(1)
	}
}

func respondWithError(conn net.Conn, correlationId uint32, err errorCode) {
	resp := make([]byte, 10)
	binary.BigEndian.PutUint32(resp, 10)
	binary.BigEndian.PutUint32(resp[4:], uint32(correlationId))
	binary.BigEndian.PutUint16(resp[8:], uint16(err))
	_, er := conn.Write(resp)

	if er != nil {
		fmt.Println("Error writing:", er.Error())
		os.Exit(1)
	}
}
