package main

import (
	"fmt"
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

	if len(bytes) != 6 {
		t.Fatalf("expected header to be 6 bytes long, got=%d", len(bytes))
	}

	expected := []byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x20}
	if !slices.Equal(bytes, expected) {
		t.Fatalf("marshaled bytes don't match, expected=%v got=%v", expected, bytes)
	}
	fmt.Printf("% x", bytes)
}
