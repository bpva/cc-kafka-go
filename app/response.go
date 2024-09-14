package main

import "encoding/binary"

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
