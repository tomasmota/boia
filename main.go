package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

const (
	COMMAND_READ  = "COMMAND_READ"
	COMMAND_WRITE = "COMMAND_WRITE"
)

type Command struct {
	Operation    string
	SerialNumber int64
}

type Server struct {
	id uuid.UUID

	clientLn net.Listener
	peerLn   net.Listener
}

func NewServer() *Server {
	return &Server{
		id: uuid.New(),
	}
}

func (s *Server) Start() error {
	// Set up listener for client connections
	cln, err := net.Listen("tcp", ":3838")
	if err != nil {
		return fmt.Errorf("error creating client listener: %w", err)
	}
	s.clientLn = cln

	// Set up listening for peer connections
	pln, err := net.Listen("tcp", ":3839")
	if err != nil {
		return fmt.Errorf("error creating peer listener: %w", err)
	}
	s.peerLn = pln

	return nil
}

func (s *Server) loop() {
	// for now only listen for client connections
}

func main() {
	s := NewServer()
	log.Fatal(s.Start())
}
