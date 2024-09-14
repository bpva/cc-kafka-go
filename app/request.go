package main

import "encoding/binary"

type requestHeader struct {
	apiKey        int16
	apiVersion    int16
	correlationId int32
}

type request struct {
	hdr requestHeader
}

func (r *request) bytes() []byte {
	buf := make([]byte, 4+8)
	binary.BigEndian.PutUint32(buf, uint32(0))
	binary.BigEndian.PutUint16(buf[4:], uint16(r.hdr.apiKey))
	binary.BigEndian.PutUint16(buf[6:], uint16(r.hdr.apiVersion))
	binary.BigEndian.PutUint32(buf[8:], uint32(r.hdr.correlationId))
	return buf
}

func requestFromBytes(b []byte) *request {
	return &request{
		hdr: requestHeader{
			apiKey:        int16(binary.BigEndian.Uint16(b)),
			apiVersion:    int16(binary.BigEndian.Uint16(b[2:])),
			correlationId: int32(binary.BigEndian.Uint32(b[4:])),
		},
	}
}

func (r *request) validate() error {
	apiVersion := r.hdr.apiVersion
	if apiVersion != 0 && apiVersion != 1 && apiVersion != 2 &&
		apiVersion != 3 && apiVersion != 4 {
		return UnknownVersionErr
	}
	return nil
}
