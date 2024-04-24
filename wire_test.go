package main

import (
	"bytes"
	"slices"
	"testing"
)

func TestMarshalHeader(t *testing.T) {
	header := NewPacketHeader()
	header.PacketType = AppendEntriesPacketType
	header.PayloadLength = 32

	bytes, err := header.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) != 4 {
		t.Fatalf("expected header to be 4 bytes long, got=%d", len(bytes))
	}

	expected := []byte{0x01, 0x02, 0x00, 0x20}
	if !slices.Equal(bytes, expected) {
		t.Fatalf("marshaled bytes don't match, expected=%v got=%v", expected, bytes)
	}
}

func TestUnmarshalHeader(t *testing.T) {
	b := []byte{0x01, 0x02, 0x00, 0x20}
	expected := &PacketHeader{
		Version:       0x01,
		PacketType:    AppendEntriesPacketType,
		PayloadLength: 32,
	}

	header, err := UnmarshalPacketHeader(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	if header.Version != expected.Version {
		t.Fatalf("version mismatch. expected=%x got=%x", expected.Version, header.Version)
	}
	if header.PacketType != expected.PacketType {
		t.Fatalf("packet type mismatch. expected=%x got=%x", expected.PacketType, header.PacketType)
	}
	if header.PayloadLength != expected.PayloadLength {
		t.Fatalf("payload length mismatch. expected=%d got=%d", expected.PayloadLength, header.PayloadLength)
	}
}

func TestMarshalVoteRequest(t *testing.T) {
	r := RequestVoteRequest{
		Term:         54,
		CandidateId:  3,
		LastLogIndex: 42631,
		LastLogTerm:  54,
	}

	bytes, err := r.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if len(bytes) != 4+25 {
		t.Fatalf("expected packet to be 31 bytes long, got=%d", len(bytes))
	}
}
