package main

import "encoding/binary"

type errorCode int16
type kafkaError string

const (
	NO_ERROR        errorCode = 0
	UNKNOWN_VERSION errorCode = 35
)

const (
	UnknownVersionErr kafkaError = "Unknown version"
)

type errorBody struct {
	error errorCode
}

func (e *errorBody) bytes() []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(e.error))
	return buf
}

func (k kafkaError) Error() string {
	return string(k)
}
