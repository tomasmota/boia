package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

const (
	PeerPort   = "3838"
	ClientPort = "3839"
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
	cln, err := net.Listen("tcp", fmt.Sprintf(":%s", ClientPort))
	if err != nil {
		return fmt.Errorf("error creating client listener: %w", err)
	}
	s.clientLn = cln
	go s.clientLoop()

	pln, err := net.Listen("tcp", fmt.Sprintf(":%s", PeerPort))
	if err != nil {
		return fmt.Errorf("error creating peer listener: %w", err)
	}
	s.peerLn = pln
	go s.peerLoop()

	slog.Info("raft node started")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down server...")
	// TODO: do graceful shutdown

	return nil
}

func (s *Server) clientLoop() error {
	slog.Info("accepting client connections", "port", s.clientLn.Addr())
	for {
		conn, err := s.clientLn.Accept()
		if err != nil {
			slog.Error("error accepting client connection", "err", err)
			continue
		}

		go s.handleClientConn(conn)
	}
}

func (s *Server) peerLoop() error {
	slog.Info("accepting peer connections", "port", s.peerLn.Addr())
	for {
		conn, err := s.peerLn.Accept()
		if err != nil {
			slog.Error("error accepting peer connection", "err", err)
			continue
		}

		go s.handlePeerConn(conn)
	}
}

func (s *Server) handlePeerConn(conn net.Conn) {
	slog.Info("new peer connection", "address", conn.RemoteAddr())
	conn.Write([]byte("Handshake and stuff"))
}

func (s *Server) handleClientConn(conn net.Conn) {
	slog.Info("new client connection", "address", conn.RemoteAddr())
	conn.Write([]byte("I am not your leader"))
}

func main() {
	s := NewServer()
	log.Fatal(s.Start())
}
