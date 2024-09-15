package main

import "encoding/binary"

type errorCode int16
type kafkaError string

const (
	NO_ERROR        errorCode = 0
	CORRUPT_MESSAGE errorCode = 2
	UNKNOWN_VERSION errorCode = 35
)

const (
	UnknownVersionErr kafkaError = "Unknown version"
)

type errorBody struct {
	error errorCode
}

func makeErrorResponce(correlationId int32, error errorCode) *response {
	res := &response{
		hdr: responseHeader{
			CorrelationId: correlationId,
		},
		body: &errorBody{
			error: error,
		},
	}
	res.setLen()

	return res
}

func (e *errorBody) bytes() []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(e.error))
	return buf
}

func (k kafkaError) Error() string {
	return string(k)
}
