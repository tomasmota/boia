package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	PeerPort   = "3838"
	ClientPort = "3839"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log := NewLog()
	sm := NewStateMachine()

	log.Append(LogEntry{
		Term: 233123,
		Command: &SetCommand{
			Key:   "abc",
			Value: 32,
		},
		Committed: false,
	})

	log.Entries[0].Command.Apply(sm)

	fmt.Println(sm.state["abc"])
}
