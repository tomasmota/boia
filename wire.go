package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type PacketType uint8

const (
	RequestVotePackeType = iota + 1
	AppendEntriesPacketType
)

// Header for every packet
type PacketHeader struct {
	Version       byte
	PacketType    PacketType
	PayloadLength uint16
}

const PROTOCOL_VERSION = 1

func NewPacketHeader() *PacketHeader {
	return &PacketHeader{
		Version: byte(PROTOCOL_VERSION),
	}
}

func (h *PacketHeader) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if h.Version != PROTOCOL_VERSION {
		if h.Version == 0 {
			return nil, fmt.Errorf("version is not set")
		}
		return nil, fmt.Errorf("unknown version: %v", h.Version)
	}
	buf.WriteByte(h.Version)

	if h.PacketType == 0 {
		return nil, fmt.Errorf("packet type not set")
	}
	buf.WriteByte(byte(h.PacketType))

	if h.PacketType == 0 {
		return nil, fmt.Errorf("packet length not set")
	}
	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, h.PayloadLength)
	buf.Write(lenBuf)

	return buf.Bytes(), nil
}

func UnmarshalPacketHeader(reader *bytes.Reader) (*PacketHeader, error) {
	header := &PacketHeader{}
	version, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if version != PROTOCOL_VERSION {
		return nil, fmt.Errorf("unknown version: %v", version)
	}
	header.Version = version

	packetType, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	header.PacketType = PacketType(packetType)

	err = binary.Read(reader, binary.BigEndian, &header.PayloadLength)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// Request made by candidates to gather votes
type RequestVoteRequest struct {
	PacketHeader

	// Candidate's Term
	Term uint64

	// Candidate's ID
	CandidateId ServerId

	// Index of Candidate's last Log Entry
	LastLogIndex uint64

	// Term of Candidate's last Log Entry
	LastLogTerm uint64
}

func (r *RequestVoteRequest) Marshal() ([]byte, error) {
	var payload bytes.Buffer

	var termBuf bytes.Buffer
	err := binary.Write(&termBuf, binary.BigEndian, r.Term)
	if err != nil {
		log.Fatalf("error parsing term into binary: %v", err)
	}

	var idBuf bytes.Buffer
	err = binary.Write(&idBuf, binary.BigEndian, r.CandidateId)
	if err != nil {
		log.Fatalf("error parsing candidate into binary: %v", err)
	}

	var lastLogIndexBuf bytes.Buffer
	err = binary.Write(&lastLogIndexBuf, binary.BigEndian, r.LastLogIndex)
	if err != nil {
		log.Fatalf("error parsing lastLogIndex into binary: %v", err)
	}

	var lastLogTermBuf bytes.Buffer
	err = binary.Write(&lastLogTermBuf, binary.BigEndian, r.LastLogTerm)
	if err != nil {
		log.Fatalf("error parsing lastLogTerm into binary: %v", err)
	}

	r.PacketHeader = *NewPacketHeader()
	r.PacketType = RequestVotePackeType
	r.PayloadLength = uint16(termBuf.Len() + idBuf.Len() + lastLogIndexBuf.Len() + lastLogTermBuf.Len())

	header, err := r.PacketHeader.Marshal()
	if err != nil {
		return nil, err
	}
	payload.Write(header)
	payload.Write(termBuf.Bytes())
	payload.Write(idBuf.Bytes())
	payload.Write(lastLogIndexBuf.Bytes())
	payload.Write(lastLogTermBuf.Bytes())

	return payload.Bytes(), nil
}

// func UnmarshalRequestVoteRequest(reader bufio.Reader) (*RequestVoteRequest, error) {
// 	version, err := reader.ReadByte()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if version != PROTOCOL_VERSION {
// 		return nil, fmt.Errorf("unknown version: %v", version)
// 	}
//
// 	return nil, nil
// }
