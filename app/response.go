package main

import (
	"encoding/binary"
)

type responseHeader struct {
	CorrelationId int32
}

type response struct {
	len  int32
	hdr  responseHeader
	body body
}

type body interface {
	bytes() []byte
}

type apiVersionsResponseBody struct {
	ErrorCode      int16
	ApiKeys        []apiKey
	ThrottleTimeMs int32
	TagBuffer      tagBuffer
}

type apiKey struct {
	ApiKey     int16
	MinVersion int16
	MaxVersion int16
	TagBuffer  tagBuffer
}

type tagBuffer struct {
	buf []byte
}

func (tb *tagBuffer) bytes() []byte {
	tb.buf = binary.AppendUvarint(tb.buf, uint64(len(tb.buf)))
	return tb.buf
}

func (avrb *apiVersionsResponseBody) bytes() []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(avrb.ErrorCode))
	buf = binary.AppendUvarint(buf, uint64(len(avrb.ApiKeys)+1))
	for _, apiKey := range avrb.ApiKeys {
		buf = append(buf, apiKey.bytes()...)
	}
	buf = append(buf, make([]byte, 4)...)
	binary.BigEndian.PutUint32(buf[len(buf)-4:], uint32(avrb.ThrottleTimeMs))
	buf = append(buf, avrb.TagBuffer.bytes()...)
	return buf
}

func (ak *apiKey) bytes() []byte {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf, uint16(ak.ApiKey))
	binary.BigEndian.PutUint16(buf[2:], uint16(ak.MinVersion))
	binary.BigEndian.PutUint16(buf[4:], uint16(ak.MaxVersion))
	buf = append(buf, ak.TagBuffer.bytes()...)
	return buf
}

func (r *response) bytes() []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf, uint32(r.len))
	binary.BigEndian.PutUint32(buf[4:], uint32(r.hdr.CorrelationId))
	buf = append(buf, r.body.bytes()...)
	return buf
}

func (r *response) setLen() {
	r.len = 4 + int32(len(r.body.bytes()))
}
