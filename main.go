package main

import "log"

const (
	PeerPort   = "3838"
	ClientPort = "3839"
)

func main() {
	s := NewServer(ServerConfig{
		PeerPort:   PeerPort,
		ClientPort: ClientPort,
	})
	log.Fatal(s.Start())
}
