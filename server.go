package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type ServerState string

const (
	Follower  ServerState = "Follower"
	Leader    ServerState = "Leader"
	Candidate ServerState = "Candidate"
)

type ServerConfig struct {
	PeerPort   string
	ClientPort string

	PeerAddresses []string
}

type Server struct {
	config ServerConfig
	id     uuid.UUID

	clientLn *net.TCPListener
	peerLn   *net.TCPListener

	currentTerm int64
	votedFor    uuid.UUID
	state       ServerState
}

func NewServer(config ServerConfig) *Server {
	return &Server{
		config: config,
		id:     uuid.New(),

		// WARN: remove after testing client conn
		state: Follower,
		// WARN: remove after testing client conn
		currentTerm: 4,
	}
}

func (s *Server) Start() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+s.config.PeerPort)
	if err != nil {
		return fmt.Errorf("error resolving tcp address: %w", err)
	}
	pln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("error creating peer listener: %w", err)
	}
	s.peerLn = pln
	go s.serverLoop()

	// cln, err := net.Listen("tcp", fmt.Sprintf(":%s", s.config.ClientPort))
	// if err != nil {
	// 	return fmt.Errorf("error creating client listener: %w", err)
	// }
	// s.clientLn = cln
	// go s.clientLoop()

	slog.Info("raft node started")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down server...")
	// TODO: do graceful shutdown

	return nil
}

func (s *Server) serverLoop() error {
	slog.Info("accepting peer connections", "port", s.peerLn.Addr())
	for {
		timeout := rand.Intn(150) + 150 // timeout between 150 and 300 ms
		s.peerLn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
		conn, err := s.peerLn.Accept()
		if err != nil {
			if net, ok := err.(net.Error); ok && net.Timeout() {
				// TODO: we only want to timeout if we didnt receive AppendEntries, if we receive RequestVote timer should keep ticking
				slog.Warn("timed out waiting for peer connection", "timeout", timeout)
				continue
			}
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

// func (s *Server) clientLoop() error {
// 	slog.Info("accepting client connections", "port", s.clientLn.Addr())
// 	for {
// 		conn, err := s.clientLn.Accept()
// 		if err != nil {
// 			slog.Error("error accepting client connection", "err", err)
// 			continue
// 		}
//
// 		// fmt.Println(string(conn.Read))
// 		go s.handleClientConn(conn)
// 	}
// }
//
//
//
//
// type Command byte
//
// const (
// 	Read  Command = 'r'
// 	Write Command = 'w'
// )
//
// type ClientRequest struct {
// 	command Command
// 	key     string
// 	value   int
// 	serial  uint64
// }
//
// func handleClientRequest(r io.Reader) error {
// 	reader := bufio.NewReader(r)
//
// 	command, err := reader.ReadByte()
// 	if err != nil {
// 		log.Fatalf("error parsing command: %v", err)
// 	}
//
// 	switch Command(command) {
// 	case Read:
// 		fmt.Println("read command")
// 	case Write:
// 		fmt.Println("write command")
// 		handleWriteRequest(reader)
// 	default:
// 		return fmt.Errorf("unrecognized command: %c", command)
// 	}
// 	return nil
// }
//
// func handleWriteRequest(reader *bufio.Reader) error {
// 	lengthByte, err := reader.ReadByte()
// 	if err != nil {
// 		log.Fatalf("error reading payload length: %v", err)
// 	}
// 	length := uint8(lengthByte)
// 	fmt.Printf("length %d\n", length)
//
// 	key := make([]byte, length)
// 	n, err := io.ReadFull(reader, key)
// 	if err != nil {
// 		log.Fatalf("error reading key: %v", err)
// 	}
// 	if n != int(length) {
// 		log.Fatalf("expected to read %b bytes. got=%d", length, n)
// 	}
// 	fmt.Printf("key %s\n", key)
//
// 	var value int64
// 	err = binary.Read(reader, binary.BigEndian, &value)
// 	if err != nil {
// 		log.Fatalf("error reading value: %v", err)
// 	}
// 	// fmt.Printf("bytes client:%b", value.Bytes())
//
// 	newline, err := reader.ReadByte()
// 	if err != nil {
// 		log.Fatalf("error reading newline: %v", err)
// 	}
// 	if newline != '\n' {
// 		log.Fatalf("expected \\n. got=%s", string(newline))
// 	}
//
// 	return nil
// }
//
// func (s *Server) handleClientConn(conn net.Conn) {
// 	slog.Info("new client connection", "address", conn.RemoteAddr())
// 	conn.Write([]byte("I am not your leader"))
// }
