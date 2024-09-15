package main

import "encoding/binary"

type requestHeader struct {
	len           int32
	apiKey        int16
	apiVersion    int16
	correlationId int32
	clientId      string
	tagBuffer     tagBuffer
}

type request struct {
	hdr requestHeader
}

func requestFromBytes(b []byte) *request {
	return &request{
		hdr: requestHeader{
			len:           int32(binary.BigEndian.Uint32(b)),
			apiKey:        int16(binary.BigEndian.Uint16(b[4:])),
			apiVersion:    int16(binary.BigEndian.Uint16(b[6:])),
			correlationId: int32(binary.BigEndian.Uint32(b[8:])),
			clientId:      string(b[12:14]),
			tagBuffer:     tagBuffer{b[14:]},
		},
	}
}

func (r *request) validate() error {
	switch r.hdr.apiVersion {
	case 0, 1, 2, 3, 4:
	default:
		return UnknownVersionErr
	}
	return nil
}
